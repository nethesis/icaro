#!/usr/bin/env python3
"""
One-shot bootstrap: link existing Icaro company accounts to My (Logto)
organizations by writing accounts.logto_org_id.

Identity chain (no email heuristics at runtime):

    icaro accounts.uuid == old-my (noc) reseller.uuid
    noc reseller --(uuid / piva / company name)--> new-my reseller org

Inputs (CSV with header):

  --icaro-accounts  id,uuid,username,email
      mysql icaro:  SELECT id, uuid, username, email FROM accounts
                    WHERE type='reseller' AND logto_org_id IS NULL;

  --noc-resellers   uuid,piva,company,username,email
      mysql noc:    SELECT r.uuid, r.piva, r.company, u.username, u.email
                    FROM reseller r JOIN User u ON u.id = r.user_id
                    WHERE r.deleted IS NULL;

  --my-orgs         logto_id,name,custom_data
      psql my:      SELECT logto_id, name, custom_data FROM resellers
                    WHERE deleted_at IS NULL;
      (custom_data as JSON text)

Outputs (in --output-dir, default ./bootstrap-out):

  matched.sql             UPDATE statements to apply (review first!)
  report_matched.csv      account -> org with match method
  report_icaro_only.csv   nethspot-only accounts (no noc reseller): keep classic login
  report_no_my_org.csv    noc reseller found, but no org on new my
  report_ambiguous.csv    multiple candidate orgs, resolve by hand

Nothing is written to any database: apply matched.sql manually.
"""

import argparse
import csv
import json
import os
import re
import sys
from collections import defaultdict


def norm_name(s):
    """Normalize a company name for loose comparison."""
    s = (s or "").lower()
    s = re.sub(r"\b(s\.?r\.?l\.?s?|s\.?p\.?a\.?|s\.?n\.?c\.?|s\.?a\.?s\.?|srls|unipersonale)\b", "", s)
    return re.sub(r"[^a-z0-9]", "", s)


def load_csv(path):
    """Load a CSV or TSV file (mysql -B output works as-is)."""
    with open(path, newline="", encoding="utf-8", errors="replace") as f:
        header = f.readline()
        delimiter = "\t" if "\t" in header else ","
        f.seek(0)
        return list(csv.DictReader(f, delimiter=delimiter))


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--icaro-accounts", required=True)
    ap.add_argument("--noc-resellers", required=True)
    ap.add_argument("--my-orgs", required=True)
    ap.add_argument("--output-dir", default="bootstrap-out")
    args = ap.parse_args()

    accounts = load_csv(args.icaro_accounts)
    noc = load_csv(args.noc_resellers)
    orgs = load_csv(args.my_orgs)

    noc_by_uuid = {r["uuid"].strip(): r for r in noc if r.get("uuid", "").strip()}

    # Index my orgs by the legacy keys found in custom_data (raw scan: the
    # exact key names vary between migration batches) and by normalized name
    orgs_by_token = defaultdict(list)   # token (uuid/piva found in custom_data) -> orgs
    orgs_by_name = defaultdict(list)
    for org in orgs:
        raw = (org.get("custom_data") or "").lower()
        try:
            cd = json.loads(org.get("custom_data") or "{}")
        except ValueError:
            cd = {}
        tokens = set(re.findall(r"[0-9a-f]{8}-[0-9a-f-]{10,}[0-9a-f]|\b\d{9,11}\b", raw))
        for v in cd.values():
            if isinstance(v, str) and v.strip():
                tokens.add(v.strip().lower())
        for t in tokens:
            orgs_by_token[t].append(org)
        orgs_by_name[norm_name(org.get("name"))].append(org)

    matched, icaro_only, no_my_org, ambiguous = [], [], [], []
    used_orgs = defaultdict(list)  # logto_id -> accounts (detect N:1 collisions)

    for acc in accounts:
        uuid = (acc.get("uuid") or "").strip()
        noc_r = noc_by_uuid.get(uuid)
        if not noc_r:
            icaro_only.append(acc)
            continue

        candidates, method = [], None
        for key, m in ((uuid.lower(), "uuid"), ((noc_r.get("piva") or "").strip().lower(), "piva")):
            if key and orgs_by_token.get(key):
                candidates, method = orgs_by_token[key], m
                break
        if not candidates:
            nname = norm_name(noc_r.get("company") or acc.get("username"))
            if nname and orgs_by_name.get(nname):
                candidates, method = orgs_by_name[nname], "name"

        uniq = {o["logto_id"]: o for o in candidates}
        if len(uniq) == 1:
            org = next(iter(uniq.values()))
            matched.append((acc, org, method))
            used_orgs[org["logto_id"]].append(acc)
        elif len(uniq) > 1:
            ambiguous.append((acc, list(uniq.values()), method))
        else:
            no_my_org.append((acc, noc_r))

    # An org linked to more than one account is a conflict: move to ambiguous
    conflicts = {oid for oid, accs in used_orgs.items() if len(accs) > 1}
    matched, dropped = [m for m in matched if m[1]["logto_id"] not in conflicts], [m for m in matched if m[1]["logto_id"] in conflicts]
    for acc, org, method in dropped:
        ambiguous.append((acc, [org], method + " (org linked by multiple accounts)"))

    os.makedirs(args.output_dir, exist_ok=True)

    def out(name):
        return open(os.path.join(args.output_dir, name), "w", newline="", encoding="utf-8")

    with out("matched.sql") as f:
        f.write("-- Review before applying: mysql icaro < matched.sql\n")
        for acc, org, method in matched:
            f.write(
                "UPDATE accounts SET logto_org_id = '%s' WHERE id = %s AND logto_org_id IS NULL; -- %s -> %s [%s]\n"
                % (org["logto_id"], acc["id"], acc.get("username"), org.get("name"), method)
            )

    with out("report_matched.csv") as f:
        w = csv.writer(f)
        w.writerow(["account_id", "username", "email", "logto_org_id", "org_name", "method"])
        for acc, org, method in matched:
            w.writerow([acc["id"], acc.get("username"), acc.get("email"), org["logto_id"], org.get("name"), method])

    with out("report_icaro_only.csv") as f:
        w = csv.writer(f)
        w.writerow(["account_id", "username", "email"])
        for acc in icaro_only:
            w.writerow([acc["id"], acc.get("username"), acc.get("email")])

    with out("report_no_my_org.csv") as f:
        w = csv.writer(f)
        w.writerow(["account_id", "username", "noc_company", "noc_piva", "noc_email"])
        for acc, noc_r in no_my_org:
            w.writerow([acc["id"], acc.get("username"), noc_r.get("company"), noc_r.get("piva"), noc_r.get("email")])

    with out("report_ambiguous.csv") as f:
        w = csv.writer(f)
        w.writerow(["account_id", "username", "method", "candidate_orgs"])
        for acc, cands, method in ambiguous:
            w.writerow([acc["id"], acc.get("username"), method, "; ".join("%s (%s)" % (o["logto_id"], o.get("name")) for o in cands)])

    print("accounts:   %d" % len(accounts))
    print("matched:    %d  -> matched.sql" % len(matched))
    print("icaro-only: %d  (classic login, untouched)" % len(icaro_only))
    print("no my org:  %d  (on noc but not migrated to new my)" % len(no_my_org))
    print("ambiguous:  %d  (resolve by hand)" % len(ambiguous))
    return 0


if __name__ == "__main__":
    sys.exit(main())
