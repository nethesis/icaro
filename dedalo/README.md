# Dedalo

## Daemon
- dedalo (ad-hoc instance of chilli with dedalo configuration)
- check dedalo status: `systemctl status dedalo`

## Configuration
`# cat /opt/icaro/dedalo/config`
```ini
HS_INTERFACE="interface_eth"
HS_HOSTNAME="host:port"
HS_ID="hotspot-id"
HS_UNIT_NAME="unit-id"
HS_UUID="generated-uuid"
HS_SECRET="your-secret"
HS_ALLOW_ORIGINS="your-allow-origin-for-cors"
```

### Example
```ini
HS_INTERFACE="eth2"
HS_HOSTNAME="myicaro"
HS_ID="hotel-king"
HS_UNIT_NAME="unit-king"
HS_UUID="1610fe6d-8578-4247-b4a2-c40dced94ber"
HS_SECRET="My__Secret99"
HS_ALLOW_ORIGINS="*"
```

## Helper
`Usage: ./dedalo { query | config | start | stop | restart | status | info }`
- **query**: wrapper to chilli_query
- **config**: generate `chilli.conf` based on your `config` file
- **start**: equal to `systemctl start dedalo`
- **stop**: equal to `systemctl stop dedalo`
- **restart**: equal to `systemctl restart dedalo`
- **status**: equal to `systemctl status dedalo`
- **info**: show the hotspot unit info
  ```
  Dedalo: Network Access Controller
    HotSpot:    hs-test
    Unit name:  unit-test
    UUID:       1610fe6d-8578-4247-b4a2-c40dced94bdd
    Secret:     Nethesis,1234
    MAC:        00-0D-B9-41-7C-F8
  ```
