package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const configPath = "./../../configs/"

func ConfigUpdaterMock(configPath string) error {
	for i := 0; i < 5; i++ {
		err := readConfig(configPath)
		if err != nil {
			return err
		}
	}

	return nil
}

var _ = Describe("Update config.", func() {
	It("Success update.", func() {
		Expect(
			UpdateConfig(configPath+"config.json", ConfigUpdaterMock),
		).Should(BeNil())
	})
	It("Failure update.", func() {
		Expect(
			UpdateConfig(configPath+"config_wrong.json", ConfigUpdaterMock),
		).Should(MatchError("open ./../../configs/config_wrong.json: no such file or directory"))
	})
})
