package utils

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/ozonva/ova-link-api/internal/link"
)

var _ = Describe("Slice utils", func() {
	Context("Split slice of int into chunks", func() {
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		DescribeTable("",
			func(input []int, size int, expected [][]int) {
				Expect(SliceChunk(input, size)).Should(BeEquivalentTo(expected))
			},
			Entry("empty slice", []int{}, 3, [][]int{}),
			Entry("one element, chunked equally", []int{1}, 1, [][]int{{1}}),
			Entry("several elements, size is more than length", slice[:3], 5, [][]int{slice[:3]}),
			Entry("several elements, size = 0", slice[:3], 0, [][]int{}),
			Entry("several elements, size < 0", slice[:3], -1, [][]int{}),
			Entry("odd amount of elements, chunked equally", slice, 3, [][]int{slice[:3], slice[3:6], slice[6:9]}),
			Entry("odd amount of elements, chunked unequally", slice[:7], 2, [][]int{slice[:2], slice[2:4], slice[4:6], slice[6:7]}),
			Entry("even amount of elements, chunked equally", slice[:6], 2, [][]int{slice[:2], slice[2:4], slice[4:6]}),
			Entry("even amount of elements, chunked unequally", slice[:6], 4, [][]int{slice[:4], slice[4:6]}),
		)
	})
	Context("Filter slice by another slice", func() {
		DescribeTable("",
			func(input []int, expected []int) {
				Expect(SliceFilterByList(input)).Should(BeEquivalentTo(expected))
			},
			Entry("empty slice", []int{}, []int{}),
			Entry("one element, not in filter", []int{4}, []int{4}),
			Entry("one element, in filter", []int{3}, []int{}),
			Entry("several elements, all elements are distinct, filtered", []int{1, 2, 3}, []int{2}),
			Entry(
				"several elements, duplicated elements, filtered",
				[]int{-5, -4, -3, -3, -2, -2, 0, 0, 1, 2, 3, 4, 5, 6, 6, 7, 7, 7},
				[]int{-4, -2, -2, 2, 4, 6, 6},
			),
			Entry("several elements, all not in list", []int{-6, -4, -2, 2, 4, 6, 8, 10}, []int{-6, -4, -2, 2, 4, 6, 8, 10}),
		)
	})
	Context("Split slice of link entities into chunks", func() {
		link1 := link.New(1, "https://test1.com")
		link1.ID = 1
		link2 := link.New(1, "https://test2.com")
		link2.ID = 2
		link3 := link.New(1, "https://test3.com")
		link3.ID = 3
		link4 := link.New(2, "https://test4.com")
		link4.ID = 4
		link5 := link.New(2, "https://test5.com")
		link5.ID = 5
		link6 := link.New(2, "https://test6.com")
		link6.ID = 6
		link7 := link.New(3, "https://test7.com")
		link7.ID = 7
		link8 := link.New(3, "https://test8.com")
		link8.ID = 8
		link9 := link.New(3, "https://test9.com")
		link9.ID = 9
		slice := []link.Link{*link1, *link2, *link3, *link4, *link5, *link6, *link7, *link8, *link9}

		DescribeTable("Convert slice of link entities into map",
			func(input []link.Link, size uint, expected [][]link.Link) {
				Expect(SliceChunkLink(input, size)).Should(BeEquivalentTo(expected))
			},
			Entry("empty slice", []link.Link{}, uint(3), [][]link.Link{}),
			Entry("one element, chunked equally", slice[:1], uint(1), [][]link.Link{slice[:1]}),
			Entry("several elements, size is more than length", slice[:3], uint(5), [][]link.Link{slice[:3]}),
			Entry("several elements, size = 0", slice[:3], uint(0), [][]link.Link{}),
			Entry("odd amount of elements, chunked equally", slice, uint(3), [][]link.Link{slice[:3], slice[3:6], slice[6:9]}),
			Entry("odd amount of elements, chunked unequally", slice[:7], uint(2), [][]link.Link{slice[:2], slice[2:4], slice[4:6], slice[6:7]}),
			Entry("even amount of elements, chunked equally", slice[:6], uint(2), [][]link.Link{slice[:2], slice[2:4], slice[4:6]}),
			Entry("even amount of elements, chunked unequally", slice[:6], uint(4), [][]link.Link{slice[:4], slice[4:6]}),
		)

		DescribeTable("Convert slice of link entities into map",
			func(input []link.Link, expected map[uint64]link.Link, expectedErr error) {
				if expectedErr != nil {
					result, err := SliceLinkToMapLink(input)
					Expect(err).Should(HaveOccurred())
					Expect(result).Should(BeNil())
				} else {
					result, err := SliceLinkToMapLink(input)
					Expect(err).Should(Succeed())
					Expect(result).Should(BeEquivalentTo(expected))
				}
			},
			Entry("empty slice", []link.Link{}, make(map[uint64]link.Link, 0), nil),
			Entry("slice without duplicates", slice[:3], map[uint64]link.Link{1: slice[0], 2: slice[1], 3: slice[2]}, nil),
			Entry("slice with duplicates", append(slice[:3], slice[0]), nil, errors.New("duplicate link id")),
		)
	})
})
