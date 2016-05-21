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

type ExportedImage struct {
	Id   string
	Url  string
	Tags []string
}

type ExportFormat struct {
	Creator string
	Images  []ExportedImage
}

type TagInformation struct {
	Tag   string
	Count int64
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

func (store *Store) Find(filter Filter) (image *Image, err error) {
	limitedFilter := Limiter{Filter: filter, Limit: 1}

	var imageSlice []Image
	imageSlice, err = store.List(limitedFilter)
	if err != nil {
		return
	}

	if len(imageSlice) == 0 {
		return
	}

	image = &imageSlice[0]
	return
}

func (store *Store) Get(imageId string) (*Image, error) {
	return store.Find(ExactIdFilter{Id: imageId})
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

	// fmt.Println(queryString)

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
			} else {
				img.Url = ""
			}
			if tag.Valid {
				img.Tags = []string{tag.String}
			} else {
				img.Tags = []string{}
			}

			// Fetch file size
			var info os.FileInfo
			info, err = os.Stat(store.PathFor(&img))
			if err != nil {
				return
			}
			img.Size = uint64(info.Size())
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

func (store *Store) ListTags(filter Filter) (result []TagInformation, err error) {
	result = make([]TagInformation, 0)

	queryString := fmt.Sprintf(`
		SELECT tag, COUNT(*)
		FROM image_tags
		WHERE %s
		GROUP BY tag
		ORDER BY tag ASC`, filter.Condition())

	rows, err := store.db.Query(queryString, filter.Values()...)
	if err != nil {
		return
	}
	defer rows.Close()

	tagInfo := TagInformation{}
	for rows.Next() {
		err = rows.Scan(&tagInfo.Tag, &tagInfo.Count)
		if err != nil {
			return
		}

		result = append(result, tagInfo)
	}

	err = rows.Err()
	return
}

func (store *Store) UpdateUrl(image *Image, url string) error {
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE images SET url = ? WHERE id = ?", url, image.Id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	image.Url = url

	return nil
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

func (s *Store) Remove(image *Image) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Remove the file
	err = os.Remove(s.PathFor(image))
	if err != nil {
		return err
	}

	// Remove from database
	_, err = tx.Exec("DELETE FROM images WHERE id = ?", image.Id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	image.AddedAt = nil

	return nil
}

func (s *Store) RemoveAll(images []Image) (err error) {
	for _, img := range images {
		if err = s.Remove(&img); err != nil {
			return
		}
	}
	return
}

func (s *Store) Hydrate(img *Image) (err error) {
	path := s.PathFor(img)
	img.Data, err = ioutil.ReadFile(path)
	return
}

func (s *Store) Version() (int, error) {
	var version int64 = 0

	err := s.db.QueryRow("PRAGMA user_version").Scan(&version)

	// If not specified, the version is 1
	if err == nil && version < 1 {
		version = 1
	}

	return int(version), err
}
