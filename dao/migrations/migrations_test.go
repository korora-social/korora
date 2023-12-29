package migrations_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"database/sql"

	"github.com/korora-social/korora/dao/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func TestMigrations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database Migrations Test Suite")
}

func nilMigrator(_ bun.Tx) error {
	return nil
}

var _ = Describe("Migrator", func() {
	var migrator *migrations.Migrator
	var db *bun.DB

	BeforeEach(func() {
		migrator = &migrations.Migrator{}

		sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:")
		if err != nil {
			panic(err)
		}

		db = bun.NewDB(sqldb, sqlitedialect.New())
	})

	It("Calls the migration functions", func() {
		called := false
		migrator.Add("first", func(tx bun.Tx) error {
			called = true
			return nil
		})

		Expect(migrator.Migrate(db)).To(Succeed())
		Expect(called).To(BeTrue())
	})

	It("Bubbles up errors", func() {
		theErr := fmt.Errorf("boom")
		migrator.Add("first", func(tx bun.Tx) error {
			return theErr
		})

		err := migrator.Migrate(db)
		Expect(err).To(Equal(theErr))
	})

	It("Creates a migration table if none exists", func() {
		tables := []string{}
		Expect(db.NewRaw("SELECT name FROM sqlite_master WHERE type='table'").Scan(context.TODO(), &tables)).To(Succeed())
		Expect(tables).To(BeEmpty())

		Expect(migrator.Migrate(db)).To(Succeed())
		Expect(db.NewRaw("SELECT name FROM sqlite_master WHERE type='table'").Scan(context.TODO(), &tables)).To(Succeed())
		Expect(tables).To(ContainElement("migrations"))
	})

	It("Records migrations made", func() {
		migrator.Add("first", func(tx bun.Tx) error {
			return nil
		})
		Expect(migrator.Migrate(db)).To(Succeed())

		migration := &migrations.DbMigration{}
		Expect(db.NewSelect().Model(migration).Order("pk desc").Limit(1).Scan(context.TODO(), migration)).To(Succeed())
		Expect(migration).To(Equal(&migrations.DbMigration{
			Pk:      1,
			Comment: "first",
		}))
	})

	It("Doesn't re-run migrations", func() {
		migrator.Add("first", func(tx bun.Tx) error {
			return nil
		})
		Expect(migrator.Migrate(db)).To(Succeed())

		migrator.Add("second", func(tx bun.Tx) error {
			return nil
		})
		Expect(migrator.Migrate(db)).To(Succeed())

		migration := &migrations.DbMigration{}
		Expect(db.NewSelect().Model(migration).Order("pk desc").Limit(1).Scan(context.TODO(), migration)).To(Succeed())
		Expect(migration).To(Equal(&migrations.DbMigration{
			Pk:      2,
			Comment: "second",
		}))

		var count int
		Expect(db.NewRaw("SELECT count(*) FROM migrations").Scan(context.TODO(), &count)).To(Succeed())
		Expect(count).To(Equal(2))
	})

	It("Fails if the migration comments are different", func() {
		migrator.Add("first", func(tx bun.Tx) error {
			return nil
		})
		Expect(migrator.Migrate(db)).To(Succeed())

		newMigrator := &migrations.Migrator{}
		newMigrator.Add("second", func(tx bun.Tx) error {
			return nil
		})
		Expect(newMigrator.Migrate(db)).To(MatchError(migrations.ErrBadOrder))

		var count int
		Expect(db.NewRaw("SELECT count(*) FROM migrations").Scan(context.TODO(), &count)).To(Succeed())
		Expect(count).To(Equal(1))
	})

	It("Only records the migration if it was successful", func() {
		migrator.Add("first", func(tx bun.Tx) error {
			return nil
		})
		migrator.Add("second", func(tx bun.Tx) error {
			return fmt.Errorf("boom")
		})
		Expect(migrator.Migrate(db)).ToNot(Succeed())

		migration := &migrations.DbMigration{}
		Expect(db.NewSelect().Model(migration).Order("pk desc").Limit(1).Scan(context.TODO(), migration)).To(Succeed())
		Expect(migration).To(Equal(&migrations.DbMigration{
			Pk:      1,
			Comment: "first",
		}))

		var count int
		Expect(db.NewRaw("SELECT count(*) FROM migrations").Scan(context.TODO(), &count)).To(Succeed())
		Expect(count).To(Equal(1))
	})

	It("Runs migrations in lexical order of comments", func() {
		migrator.Add("99-last", nilMigrator)
		migrator.Add("01-first", nilMigrator)
		migrator.Add("50-middle", nilMigrator)
		Expect(migrator.Migrate(db)).To(Succeed())

		records := []*migrations.DbMigration{}
		Expect(db.NewSelect().Model(&migrations.DbMigration{}).Order("pk asc").Scan(context.Background(), &records)).To(Succeed())
		Expect(records).To(BeEquivalentTo([]*migrations.DbMigration{
			&migrations.DbMigration{Pk: 1, Comment: "01-first"},
			&migrations.DbMigration{Pk: 2, Comment: "50-middle"},
			&migrations.DbMigration{Pk: 3, Comment: "99-last"},
		}))
	})

	It("Runs each migration in a separate transaction", func() {
		type tmp struct {
			bun.BaseModel `bun:"table:testTable"`
			Pk            int    `bun:"pk,pk,autoincrement"`
			Note          string `bun:"note,notnull"`
		}

		migrator.Add("01-first", func(tx bun.Tx) error {
			_, err := tx.NewCreateTable().Model(&tmp{}).Exec(context.TODO())
			if err != nil {
				return err
			}

			record := &tmp{Note: "comment from first migration"}
			_, err = tx.NewInsert().Model(record).Exec(context.TODO())
			return err
		})

		boom := fmt.Errorf("boom")
		migrator.Add("02-boom", func(tx bun.Tx) error {
			record := &tmp{Note: "comment from second migration"}
			_, err := tx.NewInsert().Model(record).Exec(context.TODO())
			if err != nil {
				return err
			}

			return boom
		})

		Expect(migrator.Migrate(db)).To(MatchError(boom))

		records := []*tmp{}
		Expect(db.NewSelect().Model(&tmp{}).Order("pk asc").Scan(context.Background(), &records)).To(Succeed())
		Expect(records).To(BeEquivalentTo([]*tmp{
			&tmp{Pk: 1, Note: "comment from first migration"},
		}))
	})
})

var _ = Describe("Global migrator", func() {
	It("Has the global migrations", func() {
		sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:")
		if err != nil {
			panic(err)
		}
		db := bun.NewDB(sqldb, sqlitedialect.New())

		Expect(migrations.RunMigrations(db)).To(Succeed())

		// confirm the initial users table exists
		tables := []string{}
		Expect(db.NewRaw("SELECT name FROM sqlite_master WHERE type='table'").Scan(context.TODO(), &tables)).To(Succeed())
		Expect(tables).To(ContainElement("users"))
	})
})
