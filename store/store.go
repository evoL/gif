package store

import (
	"database/sql"
	"github.com/evoL/gif/config"
	"github.com/evoL/gif/db"
	. "github.com/evoL/gif/image"
	"io/ioutil"
	"os"
	"path"
)

type Store struct {
	path string
	db   *sql.DB
}

func Default() (*Store, error) {
	return New(config.StorePath())
}

func New(path string) (*Store, error) {
	defaultDb, err := db.Default()
	if err != nil {
		return nil, err
	}

	store := &Store{
		path: path,
		db:   defaultDb,
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}
	return store, nil
}

func (store *Store) Save(image *Image) error {
	if err := ioutil.WriteFile(store.PathFor(image), image.Data, 0644); err != nil {
		return err
	}

	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO images (id, url) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if image.Url == "" {
		_, err = stmt.Exec(image.Id, sql.NullString{})
	} else {
		_, err = stmt.Exec(image.Id, image.Url)
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (store *Store) Close() error {
	return store.db.Close()
}

func (store *Store) PathFor(image *Image) string {
	return path.Join(store.path, image.Id+".gif")
}

func (store *Store) Contains(image *Image) bool {
	var result bool
	store.db.QueryRow("SELECT (COUNT(*) > 0) AS result FROM images WHERE id = ? OR url = ?", image.Id, image.Url).Scan(&result)
	return result
}

func (store *Store) Purge() error {
	return os.RemoveAll(store.path)
}
