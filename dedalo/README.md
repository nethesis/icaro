# Dedalo

Network Access Controller, runs on the firewall and intercepts all guest connections, based on CoovaChilli

## Daemon
- dedalo (ad-hoc instance of chilli with dedalo configuration)
- check dedalo status: `systemctl status dedalo`

## Configuration

Available options:

- ``HS_INTERFACE``: dedicated network interface for CoovaChilli
- ``HS_NETWORK``: network for clients connected to Dedalo
- ``HS_SPLASH_PAGE_URL``: Sun (capitve portal) URL hosted on your Icaro installation, eg: ``http://icaro.mydomain.com``
- ``HS_AAA_URL``: Wax (Radius over HTTP) URL hosted on your Icaro installation, eg: ``http://icaro.mydomain.com/wax/aaa``
- ``HS_ID``: the name of the Hotspot already present inside Icaro, eg: ``MyHotelCompany``
- ``HS_UNIT_NAME``: a descriptive name of local installation, eg: ``MyHotelAtTheSea``
- ``HS_UUID``: a unique unit idenifier, usually a UUID, eg ``161fre6d-8578-4247-b4a2-c40dced94bdd``
- ``HS_SECRET``: a shared secret between this unit and Icaro installation, eg: ``My$uperS3cret``
- ``HS_ALLOW_ORIGIN``: hosts allowed to execute CORS requests to Dedalo, usually it corresponds to ``HS_SPLASH_PAGE_URL``, eg: ``http://icaro.mydomain.com``

### Example
```ini
HS_INTERFACE="eth2"
HS_NETWORK="192.168.182.0/24"
HS_SPLASH_PAGE_URL="http://myicaro"
HS_AAA_URL="https://myicaro"
HS_ID="hotel-king"
HS_UNIT_NAME="unit-king"
HS_UUID="1610fe6d-8578-4247-b4a2-c40dced94ber"
HS_SECRET="My__Secret99"
HS_ALLOW_ORIGINS="*"
```

## Helper
`Usage: ./dedalo { query | config | register| start | stop | restart | status | info }`
- **query**: wrapper to chilli_query
- **config**: generate `chilli.conf` based on your `config` file
- **register**: register an unit using SUN accounts credentials
  ```
  dedalo register -u sun-username -p sun-password
  ```
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
