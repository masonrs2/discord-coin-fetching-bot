package main

import (
	"DiscordBot/config"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Token struct {
	Id     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type Response struct {
	MarketData struct {
		CurrentPrice struct {
			Usd float64 `json:"usd"`
		} `json:"current_price"`
	} `json:"market_data"`
}

var (
	DiscordBotId string
	DiscordBot   *discordgo.Session
)

func RunBot() {
	DiscordBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	DiscordBot.AddHandler(OnMessageCreate)

	err = DiscordBot.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	DiscordBotId = DiscordBot.State.User.ID
	fmt.Println("DiscordBotId: ", DiscordBotId)
}

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == DiscordBotId {
		fmt.Println("Message from self, ignoring...")
		return
	}

	fmt.Println("Message received: ", m.Content)

	if string(m.Content[0]) == config.BotPrefix {
		var symbol = m.Content[1:]

		coins, err := GetCoins()
		if err != nil { log.Fatal("Error getting coins: ", err) }

		for _, coin := range coins {
			if coin.Symbol == strings.ToLower(symbol) {
				price, err := GetCoinPrice(coin.Id)
				if err != nil { log.Fatal("Error getting coin price: ", err) }

				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: $%.2f", coin.Name, price))
				if err != nil { log.Fatal("Error sending message: ", err) }
			}
		}

	}
}

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
	if err != nil {
		log.Fatal("error unmarshalling json: ", err)
	}

	fmt.Println(token)

	return token.MarketData.CurrentPrice.Usd, nil
}

func main() {
	err := config.ConfigureBot()
	if err != nil {
		fmt.Println("Error configuring bot: ", err)
		return
	}

	RunBot()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	DiscordBot.Close()

	GetCoins()
	GetCoinPrice("bitcoin")
}
