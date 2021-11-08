package main

import (
	"AleksandrCherepanov/go/tinkoff-invest/internal/tinkoff"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

const URL = "https://api-invest.tinkoff.ru/openapi/"

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = godotenv.Load(currentDir + "/tinkoff-invest/.env")
	if err != nil {
		fmt.Println(err)
		return
	}

	api := tinkoff.New(os.Getenv("TOKEN"), os.Getenv("ACCOUNT_ID"), URL)

	portfolio, err := api.Portfolio()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(portfolio)
}
