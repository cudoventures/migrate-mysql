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

var defaultMysqlSocketPath = "/var/run/mysqld/mysqld.sock"

func main() {
	var (
		mysqlCA         string
		mysqlClientKey  string
		mysqlClientCert string
		mysqlDatabase   string
		mysqlDsn        string
		mysqlDsnEmpty   bool
		mysqlHost       string
		mysqlPass       string
		mysqlPort       int
		mysqlPortEmpty  bool
		mysqlSocket     string
		mysqlUser       string
		migrationsPath  string
		migrationsTable string
		noLock          bool
	)

	fs := flag.NewFlagSet("migrate-mysql", flag.ExitOnError)
	{
		fs.StringVar(&migrationsPath, "path", "", "the folder with the migrations in")
		fs.StringVar(&migrationsTable, "migrations-table", "", "the table name to use for schema versioning")
		fs.StringVar(&mysqlCA, "server-ca", "", "MySQL TLS Certificate Authority, provide if you want to connect to MySQL with TLS")
		fs.StringVar(&mysqlClientCert, "client-cert", "", "MySQL TLS Client Certificate, provide if you want to connect to MySQL with TLS")
		fs.StringVar(&mysqlClientKey, "client-key", "", "MySQL TLS Client Key, provide if you want to connect to MySQL with TLS")
		fs.StringVar(&mysqlDatabase, "database", "", "the database to use when connected to MySQL")
		fs.StringVar(&mysqlDsn, "dsn", "", "the complete dsn to use when connecting to MySQL")
		fs.StringVar(&mysqlHost, "host", "", "the host to use when connecting to MySQL")
		fs.StringVar(&mysqlPass, "pass", "", "the password to use when connecting to MySQL")
		fs.IntVar(&mysqlPort, "port", 0, "the port to use when connecting to MySQL")
		fs.StringVar(&mysqlSocket, "sock", "", "the socket to use when connecting to MySQL")
		fs.StringVar(&mysqlUser, "user", "", "the user to use when connecting to MySQL")
		fs.BoolVar(&noLock, "no-lock", false, "use no lock with migrate tool")

		// already set to exit on error
		_ = fs.Parse(os.Args[1:])

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
		if mysqlDsn == "" {
			mysqlDsnEnv, set := os.LookupEnv("MYSQL_DSN")
			mysqlDsnEmpty = !set
			if set {
				mysqlDsn = mysqlDsnEnv
			}
		}
		if mysqlHost == "" {
			mysqlHostEnv, set := os.LookupEnv("MYSQL_HOST")
			if set {
				mysqlHost = mysqlHostEnv
			}
		}
		if mysqlPass == "" {
			mysqlPass = os.Getenv("MYSQL_PASS")
		}
		if mysqlPort == 0 {
			portEnv, set := os.LookupEnv("MYSQL_PORT")
			mysqlPortEmpty = !set
			if set {
				port, err := strconv.Atoi(portEnv)
				if err != nil {
					fmt.Fprintf(os.Stderr, "could not set MYSQL_PORT to %s, expected integer\n", portEnv)
					os.Exit(1)
				}
				mysqlPort = port
			}
		}
		if mysqlUser == "" {
			mysqlUser = os.Getenv("MYSQL_USER")
		}
		if mysqlSocket == "" {
			mysqlSocketEnv, set := os.LookupEnv("MYSQL_SOCK")
			if set {
				mysqlSocket = mysqlSocketEnv
			} else {
				mysqlSocket = defaultMysqlSocketPath
			}
		}
	}
	mySQLClientConfig := &ClientConfig{
		User:         mysqlUser,
		Pass:         mysqlPass,
		Host:         mysqlHost,
		Port:         mysqlPort,
		PortEmpty:    mysqlPortEmpty,
		SocketPath:   mysqlSocket,
		Database:     mysqlDatabase,
		FullDsn:      mysqlDsn,
		FullDsnEmpty: mysqlDsnEmpty,
		ClientKey:    mysqlClientKey,
		ClientCert:   mysqlClientCert,
		CA:           mysqlCA,
	}

	db, err := NewClient(mySQLClientConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()

	driver, err := migrate_mysql.WithInstance(db, &migrate_mysql.Config{
		DatabaseName:    mysqlDatabase,
		MigrationsTable: migrationsTable,
		NoLock:          noLock,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%v", migrationsPath),
		mysqlDatabase,
		driver,
	)
	m.Log = &logger{}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	commandParameters := fs.Args()

	if len(commandParameters) < 1 {
		fmt.Fprintln(os.Stderr, "expected command up/down/step/version")
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
			fmt.Fprintln(os.Stderr, "step parameter must be an int")
		}
		commandErr = m.Steps(stepsParam)
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				fmt.Fprintln(os.Stderr, "database not versioned")
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		fmt.Printf("database: %s, version: %v, dirty %v\n", mysqlDatabase, version, dirty)
		return
	}
	if commandErr != nil {
		fmt.Fprintln(os.Stderr, commandErr)
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
	User         string
	Pass         string
	Host         string
	Port         int
	PortEmpty    bool
	Database     string
	FullDsn      string
	FullDsnEmpty bool
	ServerName   string
	SocketPath   string
	ClientKey    string
	ClientCert   string
	CA           string
}

func (c *ClientConfig) MySQLConfig() (string, error) {
	if !c.FullDsnEmpty {
		dsn, err := mysql.ParseDSN(c.FullDsn)
		if err != nil {
			return "", err
		}
		return dsn.FormatDSN(), nil
	}
	cfg := mysql.NewConfig()
	cfg.User = c.User
	cfg.Passwd = c.Pass
	if !c.PortEmpty {
		cfg.Net = "tcp"
		cfg.Addr = fmt.Sprintf("%s:%d", c.Host, c.Port)
	} else {
		cfg.Net = "unix"
		cfg.Addr = c.SocketPath
	}
	cfg.DBName = c.Database
	if c.SSLEnabled() {
		rootCertPool := x509.NewCertPool()

		if ok := rootCertPool.AppendCertsFromPEM([]byte(c.CA)); !ok {
			return "", fmt.Errorf("failed to append PEM.")
		}

		clientCert := make([]tls.Certificate, 0, 1)
		certs, err := tls.X509KeyPair([]byte(c.ClientCert), []byte(c.ClientKey))
		if err != nil {
			return "", fmt.Errorf("failed to produce x509 keypair: %w", err)
		}
		clientCert = append(clientCert, certs)

		err = mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:            rootCertPool,
			Certificates:       clientCert,
			InsecureSkipVerify: true,
		})
		if err != nil {
			return "", fmt.Errorf("failed to register mysql tls config: %w", err)
		}
		cfg.TLSConfig = "custom"
	}
	cfg.MultiStatements = true
	return cfg.FormatDSN(), nil
}

func (c *ClientConfig) SSLEnabled() bool {
	return c.ClientKey != "" && c.ClientCert != "" && c.CA != ""
}

func NewClient(cfg *ClientConfig) (*sql.DB, error) {

	dsn, err := cfg.MySQLConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to configure mysql connection config: %w", err)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}

	return db, nil
}
