package repo_test

import (
	"errors"
	"log"
	"time"

	"github.com/ozonva/ova-link-api/internal/link"

	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozonva/ova-link-api/internal/repo"
)

var _ = Describe("Repo", func() {
	Context("Link", func() {
		var linkRepo repo.Repo
		var db *sqlx.DB
		var dbMock sqlxmock.Sqlmock
		var err error

		BeforeEach(func() {
			db, dbMock, err = sqlxmock.Newx()
			if err != nil {
				log.Fatalln("cannot create db mock")
			}
			linkRepo = repo.NewLinkRepo(db)
		})

		It("Describe success", func() {
			selectTime := time.Now()
			dbMock.ExpectQuery("SELECT id, user_id, url, description, tags, created_at FROM links WHERE id = \\$1").
				WithArgs(1).
				WillReturnRows(
					sqlxmock.
						NewRows([]string{"id", "user_id", "url", "description", "tags", "created_at"}).
						AddRow(1, 1, "https://test.com", "test description", "tag1#tag2", selectTime),
				)

			result, err := linkRepo.DescribeEntity(1)

			Expect(result).Should(BeEquivalentTo(&link.Link{
				ID:          1,
				UserID:      1,
				Url:         "https://test.com",
				Description: "test description",
				Tags:        "tag1#tag2",
				CreatedAt:   selectTime,
			}))
			Expect(err).Should(Succeed())
		})

		It("Describe error", func() {
			dbMock.ExpectQuery("SELECT id, user_id, url, description, tags, created_at FROM links WHERE id = \\$1").
				WithArgs(1).
				WillReturnError(errors.New("not found"))

			result, err := linkRepo.DescribeEntity(1)

			Expect(result).Should(BeNil())
			Expect(err).Should(HaveOccurred())
		})

		It("List success", func() {
			selectTime1 := time.Now()
			selectTime2 := time.Now()
			selectTime3 := time.Now()
			dbMock.ExpectQuery("SELECT id, user_id, url, description, tags, created_at FROM links LIMIT 2 OFFSET 2").
				WillReturnRows(
					sqlxmock.
						NewRows([]string{"id", "user_id", "url", "description", "tags", "created_at"}).
						AddRow(3, 1, "https://test.com3", "test description3", "tag3#tag6", selectTime1).
						AddRow(4, 1, "https://test.com4", "test description4", "tag4#tag7", selectTime2).
						AddRow(5, 3, "https://test.com5", "test description5", "tag5#tag8", selectTime3),
				)

			expected := []link.Link{
				{
					ID:          3,
					UserID:      1,
					Url:         "https://test.com3",
					Description: "test description3",
					Tags:        "tag3#tag6",
					CreatedAt:   selectTime1,
				},
				{
					ID:          4,
					UserID:      1,
					Url:         "https://test.com4",
					Description: "test description4",
					Tags:        "tag4#tag7",
					CreatedAt:   selectTime2,
				},
				{
					ID:          5,
					UserID:      3,
					Url:         "https://test.com5",
					Description: "test description5",
					Tags:        "tag5#tag8",
					CreatedAt:   selectTime3,
				},
			}
			result, err := linkRepo.ListEntities(2, 2)

			Expect(result).Should(BeEquivalentTo(expected))
			Expect(err).Should(Succeed())
		})

		It("List error", func() {
			dbMock.ExpectQuery("SELECT id, user_id, url, description, tags, created_at FROM links LIMIT 2 OFFSET 2").
				WillReturnError(errors.New("something goes wrong"))

			result, err := linkRepo.ListEntities(2, 2)

			Expect(result).Should(BeNil())
			Expect(err).Should(HaveOccurred())
		})

		It("Delete success", func() {
			dbMock.ExpectExec("DELETE FROM links WHERE id = \\$1").
				WithArgs(1).
				WillReturnResult(sqlxmock.NewResult(1, 1))

			err := linkRepo.DeleteEntity(1)

			Expect(err).Should(Succeed())
		})

		It("Delete error", func() {
			dbMock.ExpectExec("DELETE FROM links WHERE id = \\$1").
				WithArgs(1).
				WillReturnError(errors.New("something goes wrong"))

			err := linkRepo.DeleteEntity(1)
			Expect(err).Should(HaveOccurred())
		})

		It("Create success", func() {
			selectTime1 := time.Now()
			selectTime2 := time.Now()
			insert := []link.Link{
				{
					ID:          3,
					UserID:      1,
					Url:         "https://test.com3",
					Description: "test description3",
					Tags:        "tag3#tag6",
					CreatedAt:   selectTime1,
				},
				{
					ID:          4,
					UserID:      1,
					Url:         "https://test.com4",
					Description: "test description4",
					Tags:        "tag4#tag7",
					CreatedAt:   selectTime2,
				},
			}

			dbMock.ExpectExec("INSERT INTO links \\(user_id,url,description,tags\\) VALUES \\(\\$1,\\$2,\\$3,\\$4\\),\\(\\$5,\\$6,\\$7,\\$8\\)").
				WithArgs(1, "https://test.com3", "test description3", "tag3#tag6", 1, "https://test.com4", "test description4", "tag4#tag7").
				WillReturnResult(sqlxmock.NewResult(2, 2))

			err := linkRepo.AddEntities(insert)
			Expect(err).Should(Succeed())
		})

		It("Create error", func() {
			selectTime1 := time.Now()
			selectTime2 := time.Now()
			insert := []link.Link{
				{
					ID:          3,
					UserID:      1,
					Url:         "https://test.com3",
					Description: "test description3",
					Tags:        "tag3#tag6",
					CreatedAt:   selectTime1,
				},
				{
					ID:          4,
					UserID:      1,
					Url:         "https://test.com4",
					Description: "test description4",
					Tags:        "tag4#tag7",
					CreatedAt:   selectTime2,
				},
			}

			dbMock.ExpectExec("INSERT INTO links \\(user_id,url,description,tags\\) VALUES \\(\\$1,\\$2,\\$3,\\$4\\),\\(\\$5,\\$6,\\$7,\\$8\\)").
				WithArgs(1, "https://test.com3", "test description3", "tag3#tag6", 1, "https://test.com4", "test description4", "tag4#tag7").
				WillReturnError(errors.New("something goes wrong"))

			err := linkRepo.AddEntities(insert)
			Expect(err).Should(HaveOccurred())
		})
	})
})
