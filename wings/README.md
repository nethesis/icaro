# wings

Captive portal hosted along with Wax, it's static responsive web page for user login

## Build Setup

``` bash
# install dependencies
npm install

# serve with hot reload at localhost:8080
npm run dev

# build for production with minification
npm run build

# build for production and view the bundle analyzer report
npm run build --report
```

For a detailed explanation on how things work, check out the [guide](http://vuejs-templates.github.io/webpack/) and [docs for vue-loader](http://vuejs.github.io/vue-loader).

## Auth mode and Digest
Captive portal is opend by Dedalo when a new device is connected to its network and, to make the auth process valid, Dedalo shows the captive portal page with some params added in the URL.

The URL is composed by:

`http://<icaro_server_name>/wings/?digest=<MD5(unit_secret+unit_uuid)>&uuid=<unit_uuid>&sessionid=<unique_session_id>&uamip=<dedalo_ip_network>&uamport=<dedalo_ip_port>`

A full example is:

`http://my.icaro-cloud.org/wings/?digest=e8040859e1bdfd6cb9a1c2e57edd3d05&uuid=9f621375-7c29-4b47-aec7-f1021d219bce&sessionid=151816839300000002&uamip=10.1.0.1&uamport=3993`

When the page is open, Wings component extracts:
- `digest`: calculated with `MD5(unit_secret+unit_uuid)`
- `uuid`: is the id of the unit, created when unit is registered
- `sessionid`: created by dedalo when a new device is connected
- `uamip`: dedalo network range ip
- `uamport`: dedalo instance port

All of above params are used to authenticate the requests from Wings to Wax, the retrieve preferences or social authentications.
