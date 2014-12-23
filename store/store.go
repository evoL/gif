package store

import (
	"database/sql"
	"github.com/evoL/gif/config"
	"github.com/evoL/gif/db"
	. "github.com/evoL/gif/image"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type Store struct {
	path string
	db   *sql.DB
}

func Default() (*Store, error) {
	return New(config.Global.StorePath)
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

func (store *Store) Add(image *Image) error {
	if err := ioutil.WriteFile(store.PathFor(image), image.Data, 0644); err != nil {
		return err
	}

	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO images (id, url, added_at) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()

	if image.Url == "" {
		_, err = stmt.Exec(image.Id, sql.NullString{}, now)
	} else {
		_, err = stmt.Exec(image.Id, image.Url, now)
	}
	if err != nil {
		return err
	}

	image.AddedAt = &now

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

func (store *Store) List() (result []Image, err error) {
	result = make([]Image, 0)

	tx, err := store.db.Begin()
	if err != nil {
		return
	}

	// Fetch all tags first - less queries!
	tagMap := make(map[string][]string)
	rows, err := tx.Query("SELECT image_id, tag FROM image_tags")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			imageId string
			tag     string
		)
		err = rows.Scan(&imageId, &tag)
		if err != nil {
			return
		}

		if tagList, keyExists := tagMap[imageId]; keyExists {
			tagMap[imageId] = append(tagList, tag)
		} else {
			tagMap[imageId] = []string{tag}
		}
	}

	err = rows.Err()
	if err != nil {
		return
	}

	// Fetch all images
	rows, err = tx.Query("SELECT id, url, added_at FROM images ORDER BY added_at DESC")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			img     = Image{}
			url     sql.NullString
			addedAt time.Time
		)
		err = rows.Scan(&img.Id, &url, &addedAt)
		if err != nil {
			return
		}

		img.AddedAt = &addedAt
		if url.Valid {
			img.Url = url.String
		}
		if tagList, hasTags := tagMap[img.Id]; hasTags {
			img.Tags = tagList
		}

		result = append(result, img)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (store *Store) UpdateTags(image *Image, tags []string) error {
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}

	// First, remove old tags
	_, err = tx.Exec("DELETE FROM image_tags WHERE image_id = ?", image.Id)
	if err != nil {
		return err
	}

	// Second, add new ones
	stmt, err := tx.Prepare("INSERT INTO image_tags (image_id, tag) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, tag := range tags {
		_, err = stmt.Exec(image.Id, tag)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	image.Tags = tags

	return nil
}
