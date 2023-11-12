package main

import (
	"DiscordBot/config"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	. "DiscordBot/api"
	. "DiscordBot/structs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	DiscordBotId string
	DiscordBot   *discordgo.Session
	coins        []Token
)

func RunBot() {
	DiscordBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	coins, err = GetCoins()
	if err != nil {
		fmt.Println("Error getting coins: ", err)
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

		for _, coin := range coins {
			if coin.Symbol == strings.ToLower(symbol) {
				price, err := GetCoinPrice(coin.Id)
				if err != nil {
					log.Fatal("Error getting coin price: ", err)
				}

				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: $%.2f", coin.Name, price))
				if err != nil {
					log.Fatal("Error sending message: ", err)
				}
			}
		}

	}
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
