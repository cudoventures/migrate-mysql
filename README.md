# migrate-mysql

This was made because the formatting required to pass the DSN parameter to the golang-migrate tool is arduous for MySQL, especially if you want to connect using TLS. This tool uses a different convention of with environment variables and flags. Not all features of golang-migrate are available, only the basics.

## Usage

```
The mysql client connection configuration will be picked up from the following environment variables.

MYSQL_CA
MYSQL_CLIENT_CERT
MYSQL_CLIENT_KEY
MYSQL_DATABASE
MYSQL_HOST
MYSQL_PASS
MYSQL_PORT
MYSQL_USER

If the environment variables are not set, the following flags must be used instead:

  -mysql-ca string
     MySQL TLS Certificate Authority, provide if you want to connect to MySQL with TLS
  -mysql-client-cert string
     MySQL TLS Client Certificate, provide if you want to connect to MySQL with TLS
  -mysql-client-key string
     MySQL TLS Client Key, provide if you want to connect to MySQL with TLS
  -mysql-database string
     the database to use when connected to MySQL
  -mysql-host string
     the host to use when connecting to MySQL
  -mysql-pass string
     the password to use when connecting to MySQL
  -mysql-port int
     the port to use when connecting to MySQL
  -mysql-user string
     the user to use when connecting to MySQL
  -no-lock
     use no lock flag with migrate tool

The required flags are:

  -migrations-table string
     the table name to use for schema versioning (default "SchemaMigrations")
  -path string
     the folder with the migrations in

Not all of the features of the migrate tool are available yet.

The accepted commands are:

up - perform all up migrations
down - perform all down migrations
step [n] - perform n migration steps
version - show the db version

Sample usage:

export MYSQL_DATABASE=test
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_HOST=localhost

migrate-mysql -path ./migrations up
```