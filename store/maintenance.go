package store

import (
	"github.com/evoL/gif/config"
	"github.com/evoL/gif/store/migrations"
	"github.com/rubenv/sql-migrate"
)

func DefaultMigrationSource() migrate.MigrationSource {
	migrate.SetTable("gif_migrations")
	return &migrate.AssetMigrationSource{
		Asset:    migrations.Asset,
		AssetDir: migrations.AssetDir,
		Dir:      "db_migrations",
	}
}

func (s *Store) Migrate(migrations migrate.MigrationSource) (err error) {
	_, err = migrate.Exec(s.db, config.Global.Db.Driver, migrations, migrate.Up)
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

func (s *Store) plannedMigrations(migrations migrate.MigrationSource) (planned []*migrate.PlannedMigration, err error) {
	planned, _, err = migrate.PlanMigration(s.db, config.Global.Db.Driver, migrations, migrate.Up, 0)
	return
}
