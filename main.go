package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
	migrate "github.com/golang-migrate/migrate/v4"
	migrate_mysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		mysqlCA         string
		mysqlClientKey  string
		mysqlClientCert string
		mysqlDatabase   string
		mysqlHost       string
		mysqlPass       string
		mysqlPort       int
		mysqlUser       string
		migrationsPath  string
		migrationsTable string
		noLock          bool
	)

	fs := flag.NewFlagSet("migrate-mysql", flag.ExitOnError)
	{
		fs.StringVar(&migrationsPath, "path", "", "the folder with the migrations in")
		fs.StringVar(&migrationsTable, "migrations-table", "SchemaMigrations", "the table name to use for schema versioning")
		fs.StringVar(&mysqlCA, "mysql-ca", "", "MySQL TLS Certificate Authority, provide if you want to connect to MySQL with TLS")
		fs.StringVar(&mysqlClientCert, "mysql-client-cert", "", "MySQL TLS Client Certificate, provide if you want to connect to MySQL with TLS")
		fs.StringVar(&mysqlClientKey, "mysql-client-key", "", "MySQL TLS Client Key, provide if you want to connect to MySQL with TLS")
		fs.StringVar(&mysqlDatabase, "mysql-database", "", "the database to use when connected to MySQL")
		fs.StringVar(&mysqlHost, "mysql-host", "", "the host to use when connecting to MySQL")
		fs.StringVar(&mysqlPass, "mysql-pass", "", "the password to use when connecting to MySQL")
		fs.IntVar(&mysqlPort, "mysql-port", 0, "the port to use when connecting to MySQL")
		fs.StringVar(&mysqlUser, "mysql-user", "", "the user to use when connecting to MySQL")
		fs.BoolVar(&noLock, "no-lock", false, "use no lock with migrate tool")

		err := fs.Parse(os.Args[1:])
		if err != nil {
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not parse flags: %v\n", err)
				os.Exit(1)
			}
		}
		if mysqlCA == "" {
			mysqlCA = os.Getenv("MYSQL_CA")
		}
		if mysqlClientCert == "" {
			mysqlClientCert = os.Getenv("MYSQL_CLIENT_CERT")
		}
		if mysqlClientKey == "" {
			mysqlClientKey = os.Getenv("MYSQL_CLIENT_KEY")
		}
		if mysqlDatabase == "" {
			mysqlDatabase = os.Getenv("MYSQL_DATABASE")
		}
		if mysqlHost == "" {
			mysqlHost = os.Getenv("MYSQL_HOST")
		}
		if mysqlPass == "" {
			mysqlPass = os.Getenv("MYSQL_PASS")
		}
		if mysqlPort == 0 {
			portEnv := os.Getenv("MYSQL_PORT")
			port, err := strconv.Atoi(portEnv)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not set MYSQL_PORT to %s, expected integer\n", portEnv)
				os.Exit(1)
			}
			mysqlPort = port
		}
		if mysqlUser == "" {
			mysqlUser = os.Getenv("MYSQL_USER")
		}
	}
	mySQLClientConfig := &ClientConfig{
		User:       mysqlUser,
		Pass:       mysqlPass,
		Host:       mysqlHost,
		Port:       mysqlPort,
		Database:   mysqlDatabase,
		ClientKey:  mysqlClientKey,
		ClientCert: mysqlClientCert,
		CA:         mysqlCA,
	}

	db, err := NewClient(mySQLClientConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting to MySQL: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	driver, err := migrate_mysql.WithInstance(db, &migrate_mysql.Config{
		DatabaseName:    mysqlDatabase,
		MigrationsTable: migrationsTable,
		NoLock:          noLock,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error making migrator driver: %v\n", err)
		return
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%v", migrationsPath),
		mysqlDatabase,
		driver,
	)
	m.Log = &logger{}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error making migrator: %v\n", err)
		return
	}

	commandParameters := fs.Args()

	if len(commandParameters) < 1 {
		fmt.Fprintf(os.Stderr, "expected command up/down/step/version\n")
		return
	}
	var commandErr error
	switch commandParameters[0] {
	case "up":
		commandErr = m.Up()
	case "down":
		commandErr = m.Down()
	case "step":
		stepsParam, err := strconv.Atoi(commandParameters[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "step parameter must be an int\n")
		}
		commandErr = m.Steps(stepsParam)
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				fmt.Printf("database not versioned\n")
			} else {
				fmt.Fprintf(os.Stderr, "error getting version: %v\n", err)
			}
		}
		fmt.Printf("database: %s, version: %v, dirty %v\n", mysqlDatabase, version, dirty)
		return
	}
	if commandErr != nil {
		fmt.Fprintf(os.Stderr, "error performing migration: %v\n", commandErr)
	}
}

type logger struct{}

func (l *logger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *logger) Verbose() bool {
	return false
}

type ClientConfig struct {
	User       string
	Pass       string
	Host       string
	Port       int
	Database   string
	ClientKey  string
	ClientCert string
	CA         string
}

func (c *ClientConfig) SSLEnabled() bool {
	return c.ClientKey != "" && c.ClientCert != "" && c.CA != ""
}

func (c *ClientConfig) DSN() string {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Database,
	)

	dsnOpts := "?multiStatements=true"
	if c.SSLEnabled() {
		dsnOpts = dsnOpts + "&tls=custom"
	}
	dsn = dsn + dsnOpts

	return dsn
}

func NewClient(cfg *ClientConfig) (*sql.DB, error) {
	if cfg.SSLEnabled() {
		rootCertPool := x509.NewCertPool()

		if ok := rootCertPool.AppendCertsFromPEM([]byte(cfg.CA)); !ok {
			return nil, fmt.Errorf("failed to append PEM.")
		}

		clientCert := make([]tls.Certificate, 0, 1)
		certs, err := tls.X509KeyPair([]byte(cfg.ClientCert), []byte(cfg.ClientKey))
		if err != nil {
			return nil, fmt.Errorf("failed to produce x509 keypair: %w", err)
		}
		clientCert = append(clientCert, certs)

		err = mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:            rootCertPool,
			Certificates:       clientCert,
			InsecureSkipVerify: true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to configure mysql tls config: %w", err)
		}
	}

	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}

	return db, nil
}
