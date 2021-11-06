package flusher_test

import (
	"errors"
	"strconv"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozonva/ova-link-api/internal/flusher"
	"github.com/ozonva/ova-link-api/internal/link"
	"github.com/ozonva/ova-link-api/internal/mocks"
)

var _ = Describe("Flusher", func() {
	var ctrl *gomock.Controller
	var mockRepo *mocks.MockRepo

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Flush entities.", func() {
		var flusherImpl flusher.Flusher
		var entities []link.Link
		var sliceLenAndValuesMatcher1 gomock.Matcher
		var sliceLenAndValuesMatcher2 gomock.Matcher
		var sliceLenAndValuesMatcher3 gomock.Matcher

		BeforeEach(func() {
			flusherImpl = flusher.NewFlusher(2, mockRepo)
			entities = make([]link.Link, 0, 6)
			for i := 0; i < 6; i++ {
				link := link.New(uint64(i), "https://link"+strconv.Itoa(i))
				link.ID = uint64(i)
				entities = append(entities, *link)
			}

			sliceLenAndValuesMatcher1 = gomock.All(gomock.Len(2), gomock.Eq(entities[0:2]))
			sliceLenAndValuesMatcher2 = gomock.All(gomock.Len(2), gomock.Eq(entities[2:4]))
			sliceLenAndValuesMatcher3 = gomock.All(gomock.Len(2), gomock.Eq(entities[4:6]))
		})

		Context("Flush chunks successfully.", func() {
			It("Everything was saved. Should return nil.", func() {
				gomock.InOrder(
					mockRepo.EXPECT().AddEntities(sliceLenAndValuesMatcher1).Return(nil),
					mockRepo.EXPECT().AddEntities(sliceLenAndValuesMatcher2).Return(nil),
					mockRepo.EXPECT().AddEntities(sliceLenAndValuesMatcher3).Return(nil),
				)

				Expect(flusherImpl.Flush(entities)).Should(BeNil())
			})
		})
		Context("Flush chunks with errors.", func() {
			It("One chunk was not saved. Should return all unprocessed entities", func() {
				gomock.InOrder(
					mockRepo.EXPECT().AddEntities(sliceLenAndValuesMatcher1).Return(nil),
					mockRepo.EXPECT().AddEntities(sliceLenAndValuesMatcher2).Return(errors.New("something goes wrong")),
					mockRepo.EXPECT().AddEntities(sliceLenAndValuesMatcher3).Return(nil),
				)

				unprocessed := make([]link.Link, 0, 2)
				unprocessed = append(unprocessed, entities[2:4]...)

				Expect(flusherImpl.Flush(entities)).Should(BeEquivalentTo(unprocessed))
			})
			It("Several chunks were not saved. Should return all unprocessed entities", func() {
				gomock.InOrder(
					mockRepo.EXPECT().AddEntities(sliceLenAndValuesMatcher1).Return(errors.New("something goes wrong")),
					mockRepo.EXPECT().AddEntities(sliceLenAndValuesMatcher2).Return(nil),
					mockRepo.EXPECT().AddEntities(sliceLenAndValuesMatcher3).Return(errors.New("something goes wrong")),
				)

				unprocessed := make([]link.Link, 0, 4)
				unprocessed = append(unprocessed, entities[0:2]...)
				unprocessed = append(unprocessed, entities[4:6]...)

				Expect(flusherImpl.Flush(entities)).Should(BeEquivalentTo(unprocessed))
			})
		})
	})
})
