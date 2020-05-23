// Copyright The Sandpiper Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

// Package command implements `sandpiper` commands (add, pull, list, ...)
package command

// sandpiper init

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
	args "github.com/urfave/cli/v2" // conflicts with our package name

	"sandpiper/pkg/shared/config"
	"sandpiper/pkg/shared/database"
	"sandpiper/pkg/shared/model"
	"sandpiper/pkg/shared/secure"
)

// Conn creates a receiver struct for a database connection
type Conn struct {
	*pg.DB
	conf       config.Database
	serverRole string
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
		fmt.Println(mc.SafeURL())
		return err
	}
	fmt.Printf("connected to host\n\n")

	// create the sandpiper database
	sc := sandpiperConfig(mc)
	if err := mdb.createDatabase(sc); err != nil {
		fmt.Println(sc.SafeURL())
		return err
	}

	// Update the database if necessary (from bindata embedded files)
	fmt.Printf("\napplying migrations...\n")
	debugMessage("Connect Options: %s", sc.DSN())
	msg, err := database.Migrate(sc.DSN())
	if err != nil {
		return err
	}
	fmt.Printf("Database: \"%s\"\n%s\n\n", sc.Database, msg)

	// connect to the new sandpiper database
	sdb, err := connectDB(sc)
	if err != nil {
		fmt.Println(sc.SafeURL())
		return err
	}

	// seed the database
	if err := sdb.seedDatabase(id); err != nil {
		return err
	}

	fmt.Printf("\ninitialization complete for \"%s\"\n\n", sc.Database)

	filename, err := sdb.createConfigFile()
	if err != nil {
		return err
	}
	fmt.Printf("A server config file \"%s\" was created in this folder\n\n", filename)

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
	if host[0:1]=="/" {
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
		} else {
			fmt.Printf("database \"%s\" already exists\n", c.Database)
		}
	}

	s = fmt.Sprintf("CREATE USER %s WITH ENCRYPTED PASSWORD '%s';", c.User, c.Password)
	fmt.Println(s)
	if _, err := db.Exec(s); err != nil {
		pgErr, ok := err.(pg.Error)
		// allow duplicate role errors
		if ok && pgErr.Field('C') != "42710" {
			return err
		} else {
			fmt.Printf("user \"%s\" already exists\n", c.User)
		}
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
		case "primary":
			syncAddr = Prompt("Public Sync URL: ", "")
		case "secondary":
			syncAddr = "(none)"
		default:
			fmt.Println("error: expected \"primary\" or \"secondary\"")
		}
	}

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

	fmt.Printf("Added Company \"%s\"\n", company.Name)
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
	user := sandpiper.User{
		FirstName: "Sandpiper",
		LastName:  "Admin",
		Username:  "admin",
		Password:  sec.Hash(pw),
		Email:     "admin@mail.com",
		CompanyID: companyID,
		Role:      sandpiper.SuperAdminRole,
		Active:    true,
	}
	if err := db.Insert(&user); err != nil {
		return err
	}
	fmt.Printf("Added User \"%s\"\n", user.Username)
	return nil
}

func (db *Conn) createConfigFile() (string, error) {
	name := "api-" + db.serverRole + ".yaml"
	if fileExists(name) {
		if err := os.Rename(name, name+".bak"); err != nil {
			return "", err
		}
	}
	c := config.Configuration{
		Server: configServer(db.serverRole),
		DB:     &db.conf,
		JWT:    configJWT(),
		App:    configApp(),
	}
	if err := config.Save(&c, name); err != nil {
		return "", err
	}
	return name, nil
}

func configServer(role string) *config.Server {
	port := "8080"
	if role != "primary" {
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

// debugMessage prints a help message to the console
func debugMessage(format string, a ...interface{}) {
	if debug {
		fmt.Printf("\n"+format+"\n", a...)
	}
}
