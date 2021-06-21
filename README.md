# migrate-mysql

## Usage

```
  -migrations-table string
     the table name to use for schema versioning (default "SchemaMigrations")
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
     use no lock with migrate tool
  -path string
     the folder with the migrations in
```