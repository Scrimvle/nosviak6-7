package database

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/glebarez/go-sqlite"
	_ "github.com/go-sql-driver/mysql"
)

// Config is polled from server.toml.
type Config struct {
	Sqlite   bool   `json:"sqlite"`
	Schemas  string `json:"schemas"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

// Database is the main stream.
type Database struct {
	db        *sql.DB
	Config    *Config
	Connected time.Time
}

// DB is the mainstream provider for Database.
var DB *Database = new(Database)

// Connect will push for the connection with the database.
func (database *Database) Connect() error {
	database.Connected = time.Now()
	err := source.OPTIONS.MarshalFromPath(&database.Config, "database")
	if err != nil {
		return err
	}

	switch database.Config.Sqlite {

	case true:
		database.db, err = sql.Open("sqlite", database.Config.Database)
		if err != nil {
			return err
		}

	case false:
		database.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", database.Config.Username, database.Config.Password, database.Config.Host, database.Config.Database))
		if err != nil {
			return err
		}

		if err := database.db.Ping(); err != nil {
			return err
		}
	}

	return database.schemaBuild()
}

// execute wraps the (sql.DB).Exec function and removes the sql.Result
func (database *Database) execute(str string, args ...any) error {
	if _, err := database.db.Exec(str, args...); err != nil {
		return err
	} else {
		return nil
	}
}

// schemaBuild will execute the schemas if required.
func (database *Database) schemaBuild() error {
	return filepath.Walk(database.Config.Schemas, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(info.Name()) != ".mysql" && !database.Config.Sqlite || filepath.Ext(info.Name()) != ".sqlite" && database.Config.Sqlite {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil || content == nil || strings.Count(string(content), "`") <= 1 {
			return err
		}

		if query, err := database.db.Query(fmt.Sprintf("SELECT * FROM `%s`", strings.Split(string(content), "`")[1:2][0])); err == nil && query != nil {
			return query.Close()
		}

		if err := database.execute(string(content)); err != nil {
			return err
		}

		switch strings.Split(string(content), "`")[1:2][0] {

		default:
			return nil

		case "users":
			user := new(User)
			user.API = source.OPTIONS.Bool("default", "api")
			if user.API {
				user.APIKey = []byte(hex.EncodeToString(*NewSalt(8)))
				source.LOGGER.AggregateTerminal().WriteLog(gologr.DEFAULT, "[NEW-USER] APIKey: %s", string(user.APIKey))
			}

			/* configures metadata about the user  */
			user.Username = "root"
			user.Password = []byte(hex.EncodeToString(*NewSalt(4)))
			user.Expiry = time.Unix(int64(source.OPTIONS.Ints("default_user", "expiry_days")*86400), 0).Unix()
			user.Roles = append(user.Roles, "admin", "mod", "reseller", "api")

			/* configures plan information about the user */
			user.Conns = source.OPTIONS.Ints("default_user", "conns")
			user.Maxtime = source.OPTIONS.Ints("default_user", "maxtime")
			user.Cooldown = source.OPTIONS.Ints("default_user", "cooldown")
			user.Sessions = source.OPTIONS.Ints("default_user", "max_sessions")

			source.LOGGER.AggregateTerminal().WriteLog(gologr.DEFAULT, "[NEW-USER] Username: %s", user.Username)
			source.LOGGER.AggregateTerminal().WriteLog(gologr.DEFAULT, "[NEW-USER] Password: %s", string(user.Password))

			defaults, err := os.Create("default_user.txt")
			if err != nil {
				return err
			}

			defaults.WriteString(fmt.Sprintf("%s:%s", user.Username, string(user.Password)))
			return database.NewUser(user, SYSTEM, nil)
		}
	})
}
