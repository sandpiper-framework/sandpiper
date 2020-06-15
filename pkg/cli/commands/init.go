// Copyright The Sandpiper Authors. All rights reserved.
// This file is licensed under the Artistic License 2.0.
// License text can be found in the project's LICENSE file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

// sandpiper init

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
	args "github.com/urfave/cli/v2" // conflicts with our package name

	"github.com/sandpiper-framework/sandpiper/pkg/shared/config"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/database"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/model"
	"github.com/sandpiper-framework/sandpiper/pkg/shared/secure"
)

// Conn creates a receiver struct for a database connection
type Conn struct {
	*pg.DB
	conf       config.Database
	serverRole string
	httpURL    string
}

var debug bool

// Init seeds the database for initial use
func Init(c *args.Context) error {
	id := c.String("id")
	debug = c.Bool("debug")
	fmt.Printf("INITIALIZE A SANDPIPER DATABASE\n\n")

	// connect to master (template) database on server
	mc := masterConfig()
	mdb, err := connectDB(mc)
	if err != nil {
		connectionHelp(mc)
		return err
	}
	fmt.Printf("connected to host\n\n")

	// create the sandpiper database
	sc := sandpiperConfig(mc)
	debugMessage("Connection: %s", sc.DSN())
	if err := mdb.createDatabase(sc); err != nil {
		connectionHelp(sc)
		return err
	}

	// Update the database if necessary
	fmt.Printf("\napplying migrations...\n")
	msg, err := database.Migrate(sc.DSN())
	if err != nil {
		return err
	}
	fmt.Printf("Database: \"%s\"\n%s\n\n", sc.Database, msg)

	// connect to the new sandpiper database
	sdb, err := connectDB(sc)
	if err != nil {
		connectionHelp(sc)
		return err
	}

	// seed the database
	if err := sdb.seedDatabase(id); err != nil {
		return err
	}

	fmt.Printf("\ninitialization complete for \"%s\"\n\n", sc.Database)

	api, cli, err := sdb.createConfigFiles()
	if err != nil {
		return err
	}
	dir := cwd()
	fmt.Printf("Server config file \"%s\" created in %s\n", api, dir)
	fmt.Printf("Command config file \"%s\" created in %s\n\n", cli, dir)

	return nil
}

func masterConfig() config.Database {
	host := "localhost"
	if runtime.GOOS == "linux" {
		// use unix domain sockets (best guess default)
		host = "/var/run/postgresql"
	}
	conf := config.Database{
		Dialect:  "postgres",
		Database: "template1", // every postgres sever has this database
		Host:     Prompt("PostgreSQL Address ("+host+"): ", host),
		Port:     Prompt("PostgreSQL Port (5432): ", "5432"),
		User:     Prompt("PostgreSQL Superuser (postgres): ", "postgres"),
		Password: GetPassword("PostgreSQL Superuser Password: "),
		SSLMode:  Prompt("SSL Mode (disable): ", "disable"),
	}
	conf.Network = network(conf.Host)
	return conf
}

func sandpiperConfig(c config.Database) config.Database {
	return config.Database{
		Dialect:  "postgres",
		Database: strings.ToLower(Prompt("New Database Name (sandpiper): ", "sandpiper")),
		User:     Prompt("Database Owner (sandpiper): ", "sandpiper"),
		Password: Prompt("Database Owner Password: ", ""),
		Network:  c.Network,
		Host:     c.Host,
		Port:     c.Port,
		SSLMode:  c.SSLMode,
	}
}

func network(host string) string {
	if host[0:1] == "/" {
		return "unix"
	}
	return "tcp"
}

func connectDB(conf config.Database) (*Conn, error) {
	// connect to the database from a config
	opts := database.ConnectOptions(&conf)
	debugMessage("Connect Options: %v", opts)
	db := pg.Connect(opts)

	// test connectivity
	if _, err := db.Exec("SELECT 1"); err != nil {
		return nil, err
	}

	return &Conn{DB: db, conf: conf}, nil
}

func (db *Conn) createDatabase(c config.Database) error {
	var s string

	s = fmt.Sprintf("CREATE DATABASE %s;", c.Database)
	fmt.Println(s)
	if _, err := db.Exec(s); err != nil {
		pgErr, ok := err.(pg.Error)
		// allow duplicate role errors
		if ok && pgErr.Field('C') != "42P04" {
			return err
		}
		fmt.Printf("database \"%s\" already exists\n", c.Database)
	}

	s = fmt.Sprintf("CREATE USER %s WITH ENCRYPTED PASSWORD '%s';", c.User, c.Password)
	fmt.Println(s)
	if _, err := db.Exec(s); err != nil {
		pgErr, ok := err.(pg.Error)
		// allow duplicate role errors
		if ok && pgErr.Field('C') != "42710" {
			return err
		}
		fmt.Printf("user \"%s\" already exists\n", c.User)
	}

	s = fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;", c.Database, c.User)
	fmt.Println(s)
	if _, err := db.Exec(s); err != nil {
		return err
	}

	return nil
}

func (db *Conn) seedDatabase(id string) error {
	settings, err := db.getSettings()
	if err != nil && err != pg.ErrNoRows {
		return err
	}
	if settings != nil {
		fmt.Printf("%s\n", settings.Display())
		return errors.New("ERROR: database is already initialized")
	}
	company, err := db.addCompany(id)
	if err != nil {
		return err
	}
	if err := db.addUser(company.ID); err != nil {
		return err
	}
	return nil
}

func (db *Conn) addCompany(companyID string) (*sandpiper.Company, error) {
	var syncAddr string

	companyName := Prompt("Company Name: ", "")
	for syncAddr == "" {
		db.serverRole = Prompt("Server-Role (primary*/secondary): ", "primary")
		switch db.serverRole {
		case sandpiper.PrimaryServer:
			syncAddr = Prompt("Public Sync URL: ", "")
		case sandpiper.SecondaryServer:
			syncAddr = "(none)"
		default:
			fmt.Println("error: expected \"primary\" or \"secondary\"")
		}
	}
	db.httpURL = Prompt("Server http URL (http://localhost): ", "http://localhost")

	id, err := uuid.Parse(companyID)
	if err != nil {
		id = uuid.New()
	}

	company := sandpiper.Company{
		ID:       id,
		Name:     companyName,
		SyncAddr: syncAddr,
		Active:   true,
	}
	if err := db.Insert(&company); err != nil {
		return nil, err
	}

	if err := db.addSettings(company.ID, db.serverRole); err != nil {
		return nil, err
	}

	fmt.Printf("Added Company \"%s\"\n\n", company.Name)
	return &company, nil
}

func (db *Conn) getSettings() (*sandpiper.Setting, error) {
	setting := sandpiper.Setting{ID: true}
	if err := db.Select(&setting); err != nil {
		return nil, err
	}
	return &setting, nil
}

func (db *Conn) addSettings(companyID uuid.UUID, role string) error {
	setting := sandpiper.Setting{
		ID:         true,
		ServerRole: role,
		ServerID:   companyID,
	}
	return db.Insert(&setting)
}

func (db *Conn) addUser(companyID uuid.UUID) error {
	sec := secure.New(1, "")

	pw := Prompt("Sandpiper Admin Password: ", "")
	u := sandpiper.User{
		FirstName: "Sandpiper",
		LastName:  "Admin",
		Username:  "admin",
		Password:  sec.Hash(pw),
		Email:     "admin@mail.com",
		CompanyID: companyID,
		Role:      sandpiper.SuperAdminRole,
		Active:    true,
	}
	if err := db.Insert(&u); err != nil {
		return err
	}
	fmt.Printf("Added User \"%s\"\n", u.Username)
	return nil
}

func (db *Conn) createConfigFiles() (string, string, error) {
	var protect = func(name string) string {
		if fileExists(name) {
			_ = os.Rename(name, name+".bak")
		}
		return name
	}

	nameAPI := protect("api-" + db.serverRole + ".yaml")
	api := config.Configuration{
		Server: configServer(db.serverRole),
		DB:     &db.conf,
		JWT:    configJWT(),
		App:    configApp(),
	}
	if err := config.Save(&api, nameAPI); err != nil {
		return "", "", err
	}

	nameCLI := protect("cli-" + db.serverRole + ".yaml")
	cmd := config.Command{
		URL:          db.httpURL,
		Port:         api.Server.Port,
		MaxSyncProcs: 5,
	}
	if err := config.Save(&config.Configuration{Command: &cmd}, nameCLI); err != nil {
		return "", "", err
	}

	return nameAPI, nameCLI, nil
}

func configServer(role string) *config.Server {
	port := "8080"
	if role != sandpiper.PrimaryServer {
		port = "8081"
	}
	key, err := APISecret()
	if err != nil {
		key = "(generate a suitable key using the `sandpiper secret` command)"
	}
	return &config.Server{
		Port:         port,
		ReadTimeout:  10,
		WriteTimeout: 5,
		MaxSyncProcs: 5,
		Debug:        false,
		APIKeySecret: key,
	}
}

func configJWT() *config.JWT {
	key, err := JWTSecret()
	if err != nil {
		key = "(generate a suitable key using the `sandpiper secret` command)"
	}
	return &config.JWT{
		Secret:           key,
		Duration:         15,
		RefreshDuration:  15,
		MaxRefresh:       1440,
		SigningAlgorithm: "HS256",
		MinSecretLength:  64,
	}
}

func configApp() *config.Application {
	return &config.Application{
		MinPasswordStr: 1,
		ServiceLogging: true,
	}
}

func fileExists(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func connectionHelp(c config.Database) {
	fmt.Printf("\nConnection failed for: %s\n\n", c.DSN())
	fmt.Printf("Check \"pg_hba.conf\" for entry: \"%s\", database: \"%s\", user: \"%s\", addr: \"%s\", auth: \"md5\".\n",
		connectType(c.Network), c.Database, currentUser(), configAddr(c),
	)
	fmt.Printf("See https://www.postgresql.org/docs/current/client-authentication.html\n\n")
	if c.Network == "tcp" {
		fmt.Printf("NOTE: \"Remote TCP/IP connections will not be possible unless the server is started with an appropriate value\n" +
			"for the 'listen_addresses' configuration parameter, since the default behavior is to listen for TCP/IP\n" +
			"connections only on the local loopback address localhost.\"\n\n")
	}
}

func currentUser() string {
	if u, err := user.Current(); err == nil {
		return u.Username
	}
	return "(unknown)"
}

func configAddr(c config.Database) string {
	if c.Network == "unix" {
		return "(blank)"
	}
	return c.Host
}

func connectType(c string) string {
	if c == "unix" {
		return "local"
	}
	return "host"
}

// debugMessage prints a help message to the console
func debugMessage(format string, a ...interface{}) {
	if debug {
		fmt.Printf("\n"+format+"\n", a...)
	}
}

func cwd() string {
	d, err := os.Getwd()
	if err != nil {
		return "this folder"
	}
	return d
}
