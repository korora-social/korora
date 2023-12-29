package migrations

import (
	"context"
	"errors"
	"sort"

	"github.com/uptrace/bun"
)

var (
	ErrBadOrder = errors.New("migration: bad order")

	singleton *Migrator = &Migrator{}
)

func AddMigration(comment string, fn MigrationFunc) bool {
	singleton.Add(comment, fn)
	return true
}

func RunMigrations(db *bun.DB) error {
	return singleton.Migrate(db)
}

type MigrationFunc func(tx bun.Tx) error

type Migration struct {
	migration MigrationFunc
	comment   string
}

func (m *Migration) Run(tx bun.Tx) error {
	return m.migration(tx)
}

type Migrator struct {
	migrations []*Migration
}

func (m *Migrator) Add(comment string, fn MigrationFunc) {
	if m.migrations == nil {
		m.migrations = []*Migration{}
	}

	m.migrations = append(m.migrations, &Migration{migration: fn, comment: comment})
}

type DbMigration struct {
	bun.BaseModel `bun:"table:migrations"`
	Pk            int    `bun:"pk,pk"`
	Comment       string `bun:"comment,notnull"`
}

func (m *Migrator) Migrate(db *bun.DB) error {
	// ensure migrations are sorted, using their comment as the key
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].comment < m.migrations[j].comment
	})

	_, err := db.NewCreateTable().Model(&DbMigration{}).IfNotExists().Exec(context.Background())
	if err != nil {
		return err
	}

	runMigrations := []*DbMigration{}
	err = db.NewSelect().Model(&DbMigration{}).Order("pk asc").Scan(context.Background(), &runMigrations)
	if err != nil {
		return err
	}

	for idx, migration := range m.migrations {
		migrationNum := idx + 1

		if idx < len(runMigrations) {
			if runMigrations[idx].Comment != migration.comment {
				return ErrBadOrder
			}
			continue
		}

		tx, err := db.Begin()
		defer tx.Rollback()
		if err != nil {
			return err
		}

		err = migration.Run(tx)
		if err != nil {
			return err
		}

		record := &DbMigration{
			Pk:      migrationNum,
			Comment: migration.comment,
		}

		_, err = tx.NewInsert().Model(record).Exec(context.Background())
		if err != nil {
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}
