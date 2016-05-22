package migrations

//go:generate go-bindata -pkg migrations -prefix "../" -o assets.go ../db_migrations/

var DataMigrations = map[int64]func() error{}
