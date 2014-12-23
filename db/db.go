package db

import (
	"database/sql"
	"github.com/evoL/gif/config"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
)

func New(driver, dataSource string) (*sql.DB, error) {
	var needsInit bool
	if driver == "sqlite3" {
		os.MkdirAll(path.Dir(dataSource), 0755)

		_, err := os.Stat(dataSource)
		needsInit = err != nil && os.IsNotExist(err)
	}

	db, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, err
	}

	if needsInit {
		if err = Setup(db); err != nil {
			return db, err
		}
	}
	return db, nil
}

func Default() (*sql.DB, error) {
	driver := config.Global.Db.Driver
	dataSource := config.Global.Db.DataSource

	return New(driver, dataSource)
}

func Setup(db *sql.DB) error {
	schema := `
	CREATE TABLE images (
	  id VARCHAR(40) PRIMARY KEY,
	  url TEXT
	);

	CREATE TABLE tags (tag VARCHAR(255) PRIMARY KEY);

	CREATE TABLE images_tags (
	  image_id VARCHAR(40) NOT NULL,
	  tag VARCHAR(255) NOT NULL
	);

	CREATE INDEX images_tags_index ON images_tags (tag);
	CREATE UNIQUE INDEX images_tags_unique ON images_tags (image_id, tag);`

	_, err := db.Exec(schema)
	return err
}
