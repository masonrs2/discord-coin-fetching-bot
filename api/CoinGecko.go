package api

import (
	. "DiscordBot/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetCoins() ([]Token, error) {
	var coins []Token

	resp, err := http.Get("https://api.coingecko.com/api/v3/coins/list")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = json.Unmarshal(body, &coins)
	if err != nil {
		log.Fatal("error unmarshalling json: ", err)
		return nil, err
	}

	return coins, nil
}

func GetCoinPrice(coin string) (float64, error) {
	if coin == "" {
		fmt.Println("No coin provided")
		return 0, nil
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s", coin)
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal("Error calling api using given coin", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body", err)
	}

	token := Response{}
	err = json.Unmarshal(body, &token)

	fmt.Println("body of JSON: ", (body))
	if err != nil {
		fmt.Println("error unmarshalling json: ", err)
	}

	fmt.Println(token)

	return token.MarketData.CurrentPrice.Usd, nil
}
