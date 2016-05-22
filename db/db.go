package db

// Replace this with a proper database abstraction when the database needs
// to be, well, abstracted away.

import (
	"database/sql"
	"github.com/evoL/gif/config"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
)

func New(driver, dataSource string) (db *sql.DB, err error) {
	if driver == "sqlite3" {
		os.MkdirAll(path.Dir(dataSource), 0755)
	}

	db, err = sql.Open(driver, dataSource)
	return
}

func Default() (*sql.DB, error) {
	driver := config.Global.Db.Driver
	dataSource := config.Global.Db.DataSource

	return New(driver, dataSource)
}
