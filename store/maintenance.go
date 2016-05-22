package store

import (
	"bufio"
	"github.com/evoL/gif/config"
	"github.com/evoL/gif/store/migrations"
	"github.com/rubenv/sql-migrate"
	"io/ioutil"
	"os"
	"strings"
)

func DefaultMigrationSource() migrate.MigrationSource {
	migrate.SetTable("gif_migrations")
	return &migrate.AssetMigrationSource{
		Asset:    migrations.Asset,
		AssetDir: migrations.AssetDir,
		Dir:      "db_migrations",
	}
}

func (s *Store) plannedMigrations(migrations migrate.MigrationSource) (planned []*migrate.PlannedMigration, err error) {
	planned, _, err = migrate.PlanMigration(s.db, config.Global.Db.Driver, migrations, migrate.Up, 0)
	return
}

func (s *Store) Prepare(migrationSource migrate.MigrationSource) (err error) {
	migrate, err := s.ShouldMigrate(migrationSource)
	if err != nil || !migrate {
		return
	}

	recreate, err := s.ShouldRecreate(migrationSource)
	if err != nil {
		return
	}

	if recreate {
		if err = s.Recreate(migrationSource); err != nil {
			return
		}
	} else if err = s.Migrate(migrationSource); err != nil {
		return
	}

	return
}

func (s *Store) Migrate(migrationSource migrate.MigrationSource) (err error) {
	planned, err := s.plannedMigrations(migrationSource)
	if err != nil {
		return
	}

	// Perform DB migrations
	_, err = migrate.Exec(s.db, config.Global.Db.Driver, migrationSource, migrate.Up)
	if err != nil {
		return
	}

	// Perform data migrations
	for _, migration := range planned {
		if dataMigration, ok := migrations.DataMigrations[migration.VersionInt()]; ok {
			if err = dataMigration(); err != nil {
				return
			}
		}
	}

	return
}

func (s *Store) ShouldMigrate(migrations migrate.MigrationSource) (bool, error) {
	planned, err := s.plannedMigrations(migrations)
	if err != nil {
		return false, err
	}

	return len(planned) > 0, nil
}

func (s *Store) ShouldRecreate(migrations migrate.MigrationSource) (bool, error) {
	planned, err := s.plannedMigrations(migrations)
	if err != nil {
		return false, err
	}

	for _, migration := range planned {
		if migration.VersionInt() == 1 {
			return true, nil
		}
	}

	return false, nil
}

func (s *Store) Implode() (err error) {
	query := `
  PRAGMA writable_schema = 1;
  DELETE FROM sqlite_master WHERE type in ('table', 'index', 'trigger');
  PRAGMA writable_schema = 0;
  VACUUM;`

	_, err = s.db.Exec(query)
	return
}

func (s *Store) Recreate(migrations migrate.MigrationSource) error {
	emptyStore := false

	// Create backup
	backupFile, err := ioutil.TempFile("", "gif-backup")
	if err != nil {
		return err
	}

	// Export backup
	writer := bufio.NewWriter(backupFile)
	if err = s.Export(writer, NullFilter{}, false); err != nil {
		if strings.HasPrefix(err.Error(), "no such table") {
			emptyStore = true
		} else {
			return err
		}
	}
	writer.Flush()

	// Drop schema
	if err = s.Implode(); err != nil {
		return err
	}

	// Migrate
	if err = s.Migrate(migrations); err != nil {
		return err
	}

	// Import
	if !emptyStore {
		backupFile.Seek(0, 0)
		reader := bufio.NewReader(backupFile)
		images, err := ParseMetadata(reader)
		if err != nil {
			return err
		}

		if err = s.ImportMetadata(images); err != nil {
			return err
		}
	}

	// Cleanup
	backupFile.Close()
	_ = os.Remove(backupFile.Name())

	return nil
}
