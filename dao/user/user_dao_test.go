package user_test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/korora-social/korora/dao/daoerrors"
	"github.com/korora-social/korora/dao/user"
	"github.com/korora-social/korora/models"

	sq "github.com/sleepdeprecation/squirrelly"
	_ "modernc.org/sqlite"
)

func TestUserDao(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User DAO Test Suite")
}

var migrationSql string
var _ = BeforeSuite(func() {
	// test assumes that this is running from the root directory of korora
	// TODO: actually setup database migrations
	wd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	schemaFile := filepath.Join(wd, "..", "..", "db", "schema.sql")
	rawSchema, err := os.ReadFile(schemaFile)
	Expect(err).NotTo(HaveOccurred())

	migrationSql = string(rawSchema)
})

var _ = Describe("User DAO", func() {
	var db *sq.Db
	var dao user.Dao

	BeforeEach(func() {
		var err error
		db, err = sq.Open("sqlite", "file::memory:")
		Expect(err).NotTo(HaveOccurred())

		_, err = db.DB.Exec(migrationSql)
		Expect(err).NotTo(HaveOccurred())

		dao = user.New(db)
	})

	Describe("GetByUsername", func() {
		It("Returns daoerrors.NotFound when a record doesn't exist", func() {
			_, err := dao.GetByUsername("not-found")
			Expect(err).To(Equal(daoerrors.NotFound))
		})
	})

	Describe("Save user", func() {
		It("Saves the user in the database", func() {
			u := models.NewUser("sleep", "https://korora.social")
			Expect(dao.Save(u)).To(Succeed())

			byUsername, err := dao.GetByUsername("sleep")
			Expect(err).To(BeNil())
			Expect(byUsername).To(BeEquivalentTo(u))
		})

		It("Updates the user", func() {
			u := models.NewUser("sleep", "https://korora.social")
			Expect(dao.Save(u)).To(Succeed())

			u.PrivateKey = []byte("fake")
			Expect(dao.Save(u)).To(Succeed())

			byUsername, err := dao.GetByUsername("sleep")
			Expect(err).To(BeNil())
			Expect(byUsername).To(BeEquivalentTo(u))
		})
	})
})
