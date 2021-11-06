package saver_test

import (
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/ozonva/ova-link-api/internal/flusher"
	"github.com/ozonva/ova-link-api/internal/link"
	"github.com/ozonva/ova-link-api/internal/mocks"
	"github.com/ozonva/ova-link-api/internal/saver"
)

var _ = Describe("Saver", func() {
	var repo *mocks.MockRepo
	var ctrl *gomock.Controller
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		repo = mocks.NewMockRepo(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("Saving by timeout. Should be call when timeout is expired.", func() {
		defer GinkgoRecover()
		flusherImpl := flusher.NewFlusher(3, repo)
		timeoutSaver := saver.NewTimeOutSaver(5, flusherImpl, 1)

		repo.EXPECT().AddEntities(gomock.Any()).Times(2).Return(nil)

		timeoutSaver.Save(*link.New(1, "1"))
		timeoutSaver.Save(*link.New(1, "2"))
		timeoutSaver.Save(*link.New(2, "3"))
		timeoutSaver.Save(*link.New(2, "4"))

		time.Sleep(2 * time.Second)
	})

	It("Saving by close. Should be call immediately. Timeout should be stopped", func() {
		flusherImpl := flusher.NewFlusher(3, repo)
		timeoutSaver := saver.NewTimeOutSaver(5, flusherImpl, 2)

		repo.EXPECT().AddEntities(gomock.Any()).Times(2).Return(nil)

		timeoutSaver.Save(*link.New(1, "1"))
		timeoutSaver.Save(*link.New(1, "2"))
		timeoutSaver.Save(*link.New(2, "3"))
		timeoutSaver.Save(*link.New(2, "4"))

		timeoutSaver.Close()
		time.Sleep(1 * time.Second)
	})
})
