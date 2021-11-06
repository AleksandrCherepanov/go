package utils

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map utils.", func() {
	DescribeTable("Invert map.",
		func(input map[string]int, expected map[int]string, expectedErr error) {
			if expectedErr == nil {
				result, err := MapInvert(input)
				Expect(err).Should(Succeed())
				Expect(result).Should(BeEquivalentTo(expected))
			} else {
				result, err := MapInvert(input)
				Expect(err).Should(HaveOccurred())
				Expect(result).Should(BeNil())
			}
		},
		Entry(
			"Distinct values, sorted keys.",
			map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5},
			map[int]string{1: "a", 2: "b", 3: "c", 4: "d", 5: "e"},
			nil,
		),
		Entry(
			"Distinct values, unsorted keys.",
			map[string]int{"c": 2, "d": 4, "a": 1, "e": 3, "b": 5},
			map[int]string{1: "a", 5: "b", 2: "c", 4: "d", 3: "e"},
			nil,
		),
		Entry(
			"Duplicated values, error",
			map[string]int{"a": 1, "b": 2, "c": 2, "d": 4, "e": 4},
			nil,
			errors.New("duplicate key exists"),
		),
	)
})
