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

## APIs

The APIs are documented using [Postman](https://www.getpostman.com/).


