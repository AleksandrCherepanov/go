package link

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Link.", func() {
	id := uint64(1)
	userId := uint64(2)
	url := "https://test.com"
	expected := &Link{id, userId, url, "", "", time.Now()}

	Context("Creation.", func() {
		linkEntity := New(userId, url)
		linkEntity.ID = id

		It("ID should be the same", func() {
			Expect(linkEntity.ID).Should(BeIdenticalTo(expected.ID))
		})
		It("User should ID be the same", func() {
			Expect(linkEntity.UserID).Should(BeIdenticalTo(expected.UserID))
		})
		It("URL should be the same", func() {
			Expect(linkEntity.Url).Should(BeIdenticalTo(expected.Url))
		})
		It("Description should be the same", func() {
			Expect(linkEntity.Description).Should(BeIdenticalTo(expected.Description))
		})
		It("Tags should be the same", func() {
			Expect(linkEntity.Tags).Should(BeEquivalentTo(expected.Tags))
		})
		It("Time should be different", func() {
			Expect(linkEntity.CreatedAt.After(expected.CreatedAt)).Should(BeTrue())
		})
	})
	Context("Set values.", func() {
		linkEntity := New(userId, url)
		linkEntity.ID = id
		description := "some description"
		newUrl := "https://new_test.com"
		linkEntity.Description = description
		linkEntity.Url = newUrl
		linkEntity.SetTagsAsSlice([]string{"tag1", "tag2"})

		It("Description should be updated", func() {
			Expect(linkEntity.Description).Should(BeIdenticalTo(description))
			Expect(linkEntity.Description).ShouldNot(BeIdenticalTo(expected.Description))
		})
		It("Tags should be updated", func() {
			Expect(linkEntity.Tags).Should(BeEquivalentTo("tag1#tag2"))
			Expect(linkEntity.Tags).ShouldNot(BeEquivalentTo(expected.Tags))
		})
	})
	Context("Add/remove tags.", func() {
		linkEntity := New(userId, url)
		linkEntity.ID = id
		linkEntity.SetTagsAsSlice([]string{"tag1", "tag2"})

		It("Tags should be added.", func() {
			linkEntity.AddTag("tag3")
			linkEntity.AddTag("tag4")
			linkEntity.AddTag("tag1")
			linkEntity.AddTag("tag2")

			Expect(linkEntity.Tags).Should(BeIdenticalTo("tag1#tag2#tag3#tag4"))
		})
		It("Tags should be removed.", func() {
			linkEntity.RemoveTag("tag3")
			linkEntity.RemoveTag("tag1")
			linkEntity.RemoveTag("tag3")
			linkEntity.RemoveTag("tag1")

			Expect(linkEntity.Tags).Should(BeIdenticalTo("tag2#tag4"))
		})
	})
	Context("Compare with another link entity.", func() {
		linkEntity1 := New(userId, url)
		linkEntity1.ID = id
		linkCopy := *linkEntity1
		linkEntity2 := &linkCopy

		description := "some description"

		It("Two empty link should be equal", func() {
			Expect(linkEntity1.Equals(linkEntity2)).Should(BeTrue())
		})

		It("Add tags. Links should not be equal.", func() {
			linkEntity2.AddTag("tag1")
			linkEntity2.AddTag("tag2")
			Expect(linkEntity1.Equals(linkEntity2)).Should(BeFalse())

			linkEntity1.AddTag("tag2")
			Expect(linkEntity1.Equals(linkEntity2)).Should(BeFalse())

			linkEntity1.AddTag("tag3")
			Expect(linkEntity1.Equals(linkEntity2)).Should(BeFalse())
		})

		It("Make tags are the same. Links should be equal.", func() {
			linkEntity1.RemoveTag("tag3")
			linkEntity1.AddTag("tag1")
			Expect(linkEntity1.Equals(linkEntity2)).Should(BeTrue())
		})

		It("Change description. Links should not be equal.", func() {
			linkEntity1.Description = description
			Expect(linkEntity1.Equals(linkEntity2)).Should(BeFalse())
		})

		It("Make description the same. Links should be equal.", func() {
			linkEntity2.Description = description
			Expect(linkEntity1.Equals(linkEntity2)).Should(BeTrue())
		})
	})
	Context("Cast to string.", func() {
		defer GinkgoRecover()
		linkEntity := New(userId, url)
		linkEntity.ID = uint64(1)
		linkEntity.Description = "Ozon Go School. Project."
		linkEntity.SetTagsAsSlice([]string{"tag1", "tag2"})
		regexpString := `ID: 1,
UserID: 2,
URL: "https://test.com",
Description: "Ozon Go School. Project.",
Tags: "tag1#tag2",
DateCreated: [0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\+[0-9]{2}:[0-9]{2}`
		Expect(linkEntity.String()).Should(MatchRegexp(regexpString))
	})
})
