package api_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/ozonva/ova-link-api/internal/link"

	"github.com/golang/mock/gomock"
	"github.com/ozonva/ova-link-api/internal/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/ozonva/ova-link-api/internal/api"
	grpc "github.com/ozonva/ova-link-api/pkg/ova-link-api"
	ova_link_api "github.com/ozonva/ova-link-api/pkg/ova-link-api"
)

type LinkMatcher struct {
	expected []link.Link
}

func (lm LinkMatcher) Matches(x interface{}) bool {
	entities := x.([]link.Link)
	if len(entities) != len(lm.expected) {
		return false
	}

	for i := 0; i < len(entities); i++ {
		if entities[i].UserID != lm.expected[i].UserID {
			return false
		}
		if entities[i].Description != lm.expected[i].Description {
			return false
		}
		if entities[i].Tags != lm.expected[i].Tags {
			return false
		}
		if entities[i].Url != lm.expected[i].Url {
			return false
		}
	}

	return true
}

func (lm LinkMatcher) String() string {
	return fmt.Sprintf("is equal to %v (%T)", lm.expected, lm.expected)
}

var _ = Describe("Api", func() {
	Context("Database", func() {
		var API grpc.LinkAPIServer
		var ctrl *gomock.Controller
		var mockRepo *mocks.MockRepo

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockRepo = mocks.NewMockRepo(ctrl)
			API = api.NewLinkAPI(mockRepo, zerolog.Nop())
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("Describe success", func() {
			selectTime := time.Now()
			expected := &link.Link{
				ID:          1,
				UserID:      1,
				Url:         "https://test.com",
				Description: "test description",
				Tags:        "tag1#tag2",
				CreatedAt:   selectTime,
			}
			mockRepo.EXPECT().DescribeEntity(gomock.Eq(uint64(1))).Times(1).Return(expected, nil)

			_, err := API.DescribeLink(
				context.Background(),
				&ova_link_api.DescribeLinkRequest{
					Id: 1,
				},
			)

			Expect(err).Should(Succeed())
		})

		It("Describe error", func() {
			mockRepo.EXPECT().DescribeEntity(gomock.Eq(uint64(1))).Times(1).
				Return(nil, errors.New("something goes wrong"))

			_, err := API.DescribeLink(
				context.Background(),
				&ova_link_api.DescribeLinkRequest{
					Id: 1,
				},
			)

			Expect(err).Should(HaveOccurred())
		})

		It("List success", func() {
			selectTime1 := time.Now()
			selectTime2 := time.Now()
			selectTime3 := time.Now()

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

			mockRepo.EXPECT().ListEntities(gomock.Eq(uint64(2)), gomock.Eq(uint64(2))).
				Times(1).
				Return(expected, nil)

			limit := uint64(2)
			offset := uint64(2)
			_, err := API.ListLink(
				context.Background(),
				&ova_link_api.ListLinkRequest{
					Limit:  &limit,
					Offset: &offset,
				},
			)

			Expect(err).Should(Succeed())
		})

		It("List error", func() {
			mockRepo.EXPECT().ListEntities(gomock.Eq(uint64(2)), gomock.Eq(uint64(2))).
				Times(1).
				Return(nil, errors.New("something goes wrong"))

			limit := uint64(2)
			offset := uint64(2)
			_, err := API.ListLink(
				context.Background(),
				&ova_link_api.ListLinkRequest{
					Limit:  &limit,
					Offset: &offset,
				},
			)

			Expect(err).Should(HaveOccurred())
		})

		It("Delete success", func() {
			mockRepo.EXPECT().DeleteEntity(gomock.Eq(uint64(1))).Times(1).Return(nil)

			_, err := API.DeleteLink(
				context.Background(),
				&ova_link_api.DeleteLinkRequest{
					Id: 1,
				},
			)

			Expect(err).Should(Succeed())
		})

		It("Delete error", func() {
			mockRepo.EXPECT().DeleteEntity(gomock.Eq(uint64(1))).
				Times(1).Return(errors.New("something goes wrong"))

			_, err := API.DeleteLink(
				context.Background(),
				&ova_link_api.DeleteLinkRequest{
					Id: 1,
				},
			)

			Expect(err).Should(HaveOccurred())
		})

		It("Create success", func() {
			selectTime1 := time.Now()
			insert := []link.Link{
				{
					ID:          0,
					UserID:      1,
					Url:         "https://test.com3",
					Description: "test description3",
					Tags:        "tag3#tag6",
					CreatedAt:   selectTime1,
				},
			}

			linkMatcher := LinkMatcher{
				expected: insert,
			}
			mockRepo.EXPECT().AddEntities(linkMatcher).Times(1).Return(nil)

			_, err := API.CreateLink(
				context.Background(),
				&ova_link_api.CreateLinkRequest{
					UserId:      1,
					Url:         "https://test.com3",
					Description: "test description3",
					Tags:        []string{"tag3", "tag6"},
				},
			)
			time.Sleep(1 * time.Second)
			Expect(err).Should(Succeed())
		})
	})
})
