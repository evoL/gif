package store

import (
	"database/sql"
	"fmt"
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

func (store *Store) List(filter Filter) (result []Image, err error) {
	result = make([]Image, 0)

	// The code assumes that the database returns images that aren't mixed with each
	// other; that is, the following situation doesn't happen:
	//
	// id | tag
	// 0  | a
	// 1  | b
	// 0  | c
	//
	// This is extremely unlikely, because those images would be added at the exact
	// same time. However, to prevent this, "ORDER BY id" needs to be added.
	// I didn't include this now, because it might add an unnecessary
	// performance hit.

	queryString := fmt.Sprintf(`
	SELECT id, tag, url, added_at
	FROM (
		SELECT DISTINCT images.*
		FROM images
		LEFT JOIN image_tags
		ON images.id = image_tags.image_id
		WHERE %s
	) images
	LEFT JOIN image_tags
	ON images.id = image_tags.image_id`, filter.Condition())

	rows, err := store.db.Query(queryString, filter.Values()...)
	if err != nil {
		return
	}
	defer rows.Close()

	img := Image{}
	for rows.Next() {
		var (
			id      string
			url     sql.NullString
			addedAt time.Time
			tag     sql.NullString
		)

		err = rows.Scan(&id, &tag, &url, &addedAt)
		if err != nil {
			return
		}

		if img.Id == id {
			if tag.Valid {
				// Add another tag to existing image
				img.Tags = append(img.Tags, tag.String)
			}
			// The opposite case should not happen, but if it does, we're just ignoring it.
		} else {
			// Append the previously processed image if there's one
			if img.Id != "" {
				result = append(result, img)
			}

			// Create a new image
			img.Id = id
			img.AddedAt = &addedAt
			if url.Valid {
				img.Url = url.String
			}
			if tag.Valid {
				img.Tags = []string{tag.String}
			} else {
				img.Tags = []string{}
			}
		}
	}

	// As we're only appending when a different image turns up, we need to append
	// the last image separately.
	if img.Id != "" {
		result = append(result, img)
	}

	err = rows.Err()
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
