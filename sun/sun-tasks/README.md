# sun-tasks

Exec tasks immediatly or in a background jobs

# Use
`sun-tasks` has three options:
- `-a <action>`: Specify action to execute, current available actions:
  - `clean-tokens`: Removes useless access tokens older than their expiration dates
  - `store-sessions`: Stores stopped sessions into a separate table
- `-w`: Runs actions in background mode as a cron job
  - `clean-tokens`: runs `daily`
  - `store-sessions`: runs `daily`
- `-c`: Path to configuration file, default `/opt/icaro/sun-api/conf.json`