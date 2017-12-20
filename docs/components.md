# Software components

Icaro is divided in 4 parts:

- **Dedalo**: network access controller, runs on the firewal and intercepts all guest connections, based on [CoovaChilli](http://coova.github.io/CoovaChilli/)
- **Wax**: accounting server, it speaks RADIUS AAA protocol via HTTP with Dedalo
- **Wings**: captive portal hosted along with Wax, it's static responsive web page for user login. It talks with Wax and Dedalo
- **Sun**: Hotspot manager, web applications which talks with Wax

Both Wax and Sun use the same MariaDB (or MySQL) database.

## Sun

Sun is written in Go.

Used tools:

- Web framework: https://github.com/gin-gonic/gin
- ORM: https://github.com/jinzhu/gorm

## Wax

Wax is a Go server listening on port 8181.
It handles AAA requests from CoovaChilli, only HTTP GET is supported.

Eeach AAA request is in the form `stage=<stage>&service=<service>&<param1>=<value1>...<paramN>=<valueN>`.

The `stage` parameter is mandatory; there are 3 available stages:

- `login`
- `counters`
- `register`


### login

This stage is invoked when a user tries to login from Wings.

Available parameters:

- `service`: service to restrict login access (unused)
- `user`: user name 
- `chap_chal: random challenge for CHAP authentication
- `chap_pass: password in CHAP format
- `chap_id`: ?  (unused)
- `ap`: access point MAC address (unused)
- `mac`: user MAC address (unused)
- `ip`: user local IP address given by Wings (unused)
- `sessionid`: unique session id (unused)
- `nasid`: unique name of the hotspot unit
- `md`: ? (unused)

Request example:
```
curl http://localhost:8181/?stage=login &service=login&user=aaa&chap_chal=57efc66af8df7ecb164a39758b18a407&chap_pass=47bce5c74f589f4867dbd57e9ca9f808&chap_id=0&ap=00-0D-B9-41-7C-F8&mac=84-3A-4B-11-44-D4&ip=10.1.0.3&sessionid=151335010500000001&nasid=hs-test&md=CCBA58A3D03CA1BC5D02B0E66A997914
```

### counters

This stage is invoked for session accounting.

Request examples:
```
curl http://localhost:8181/?stage=counters&status=stop&user=aaa&ap=00-0D-B9-41-7C-F8&mac=84-3A-4B-11-44-D4&ip=10.1.0.3&sessionid=151324185800000001&nasid=hs-test&duration=702&bytes_down=47487&pkts_down=135&bytes_up=5292&pkts_up=65&md=F2D7D9B3184E2890140C9B7FE28CC0FB
curl http://localhost:8181/?stage=counters&status=up&ap=00-0D-B9-41-7C-F8&mac=00-00-00-00-00-00&nasid=hs-test&md=F998859655817FE494BB9CEC3133BD2A
curl http://localhost:8181/?stage=counters&status=update&user=aaa&ap=00-0D-B9-41-7C-F8&mac=84-3A-4B-11-44-D4&ip=10.1.0.3&sessionid=151318184500000001&nasid=hs-test&duration=10205&bytes_down=656025&pkts_down=1655&bytes_up=164192&pkts_up=1794&md=F614C20150DCD5A8D867A977DEFF154C
```


## register

Not implemented.


## APIs

The APIs are documented using [Postman](https://www.getpostman.com/).


