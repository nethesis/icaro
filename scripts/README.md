# OIDC organization bootstrap

`oidc_org_bootstrap.py` links existing Icaro **reseller** accounts to their
My (Logto) organization by generating the `UPDATE` statements that populate
`accounts.logto_org_id`. That column is the **only** link the OIDC login
uses at runtime: no email/VAT matching ever happens during login. This
script is where the (reviewed, one-time) matching happens instead.

Identity chain used for matching:

```
icaro accounts.uuid == old-my (noc) reseller.uuid
noc reseller --(uuid / piva / normalized company name)--> new-my reseller org
```

The script is **offline** (CSV in, SQL + reports out — it never touches a
database) and **idempotent**: every generated `UPDATE` carries
`AND logto_org_id IS NULL`, so re-running and re-applying only adds new
links and never rewrites existing ones.

## Prerequisites

1. The `accounts.logto_org_id` column must exist (new installs get it from
   `deploy/ansible/roles/icarodb/files/icaro.sql`; existing databases need):

   ```sql
   ALTER TABLE accounts
     ADD COLUMN logto_org_id varchar(64) DEFAULT NULL AFTER uuid,
     ADD UNIQUE KEY logto_org_id (logto_org_id);
   ```

2. Read access to the three data sources: the Icaro MariaDB, the old My
   (noc) MariaDB, and the new My PostgreSQL.

## Step 1 — extract the input files

Both TSV (`mysql -B` output) and CSV are accepted; the delimiter is
auto-detected. Keep the headers exactly as below.

**Icaro accounts** (only `type='reseller'`; desk/customer accounts are
internal to Icaro and must never be linked):

```bash
mysql -B icaro -e "SELECT id, IFNULL(uuid,'') AS uuid, username, IFNULL(email,'') AS email
                   FROM accounts WHERE type='reseller';" > icaro_accounts.tsv
```

**Old My resellers** (noc):

```bash
mysql -B noc -e "SELECT r.uuid, IFNULL(r.piva,'') AS piva, IFNULL(r.company,'') AS company,
                        u.username, u.email
                 FROM reseller r JOIN User u ON u.id = r.user_id
                 WHERE r.deleted IS NULL;" > noc_resellers.tsv
```

**New My reseller organizations** (PostgreSQL; use `psql` so the JSON
`custom_data` is quoted correctly):

```bash
psql "$DATABASE_URL" -c "\copy (SELECT COALESCE(logto_id,'') AS logto_id, name,
                                COALESCE(custom_data::text,'{}') AS custom_data
                                FROM resellers WHERE deleted_at IS NULL)
                         TO 'my_orgs.csv' WITH CSV HEADER"
```

## Step 2 — dry run

```bash
python3 scripts/oidc_org_bootstrap.py \
  --icaro-accounts icaro_accounts.tsv \
  --noc-resellers  noc_resellers.tsv \
  --my-orgs        my_orgs.csv \
  --output-dir     bootstrap-out
```

Nothing is written to any database. Output files:

| File | Content | Action |
|---|---|---|
| `matched.sql` | one `UPDATE` per linked account, with account → org and match method (`uuid`/`piva`/`name`) as a comment | review, then apply |
| `report_matched.csv` | same matches in tabular form | review |
| `report_icaro_only.csv` | accounts with no old-My reseller (self-registered / Icaro-only) | none: they keep the classic password login |
| `report_no_my_org.csv` | old-My resellers whose organization is not on the new My yet | re-run after the next org import batch |
| `report_ambiguous.csv` | multiple candidate orgs, or one org claimed by multiple accounts | resolve by hand |

## Step 3 — review and apply

Read `matched.sql` (especially any `name`-method matches, which are the
loosest) and the ambiguous report, then:

```bash
mysql icaro < bootstrap-out/matched.sql
```

## Re-running (important)

Run steps 1–3 again **after every batch of organization imports on the new
My**. If a reseller's organization exists on My but the bootstrap has not
linked it yet, their first OIDC login triggers JIT provisioning and creates
a *new empty* account instead of linking the historical one (runtime
deliberately performs no heuristic matching). Keeping import and bootstrap
close together minimizes that window; the section below covers the cleanup
when a duplicate slips through.

Apply with `--force` so one conflict does not stop the rest, and read the
errors — they are the duplicate detector:

```bash
mysql --force icaro < bootstrap-out/matched.sql
```

A `Duplicate entry ... for key 'logto_org_id'` error means that org was
already claimed by a JIT-created duplicate: reconcile it as below.

## Reconciling a JIT duplicate

Symptom: a reseller logged in via My before the bootstrap linked their
historical account. JIT created a fresh empty account (recognizable:
`username = logto_org_id`, journal line `OIDC provisioning: created company
account ...`), while the historical account still has the hotspots.

```sql
-- 0. find the two accounts
SELECT id, username, name, logto_org_id, created FROM accounts
 WHERE logto_org_id = '<org_id>' OR (name LIKE '%<company>%' AND type = 'reseller');

-- 1. make sure the duplicate is empty (otherwise it needs a manual merge)
SELECT COUNT(*) FROM hotspots WHERE account_id = <dup_id>;

-- 2. move the link to the historical account
UPDATE accounts SET logto_org_id = NULL WHERE id = <dup_id>;
UPDATE accounts SET logto_org_id = '<org_id>' WHERE id = <legacy_id> AND logto_org_id IS NULL;

-- 3. drop the empty duplicate
DELETE FROM access_tokens      WHERE account_id = <dup_id>;
DELETE FROM subscriptions      WHERE account_id = <dup_id>;
DELETE FROM account_sms_counts WHERE account_id = <dup_id>;
DELETE FROM accounts           WHERE id = <dup_id>;
```

The reseller's next login lands on the historical account. If step 1 shows
hotspots in the duplicate, do **not** delete it: move its hotspots to the
historical account first (or escalate — that is a merge, not a cleanup).
