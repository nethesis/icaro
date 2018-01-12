## Wax

Listen port 8181.

# Build

go build

# Running

```
DB_USER=testuser DB_PASSWORD=testpassword DB_HOST=192.168.1.1 DB_NAME=icaro DB_PORT=3306 ./wax
```

# Launch tests

Execute:
```
DB_USER=<db_user> DB_PASSWORD=<db_password> DB_HOST=<db_host>6 DB_PORT=<db_port> DB_NAME=icaro go test
```

Example:
```
DB_USER=testuser DB_PASSWORD=testpassword DB_HOST=192.168.1.1 DB_NAME=icaro DB_PORT=3306 go test
```
