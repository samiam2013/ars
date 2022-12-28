package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type config struct {
	DiscordBotToken string `env:"DISCORD_BOT_TOKEN"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Fatalf("%+v\n", err)
	}
	logrus.Infof("config %+v", cfg)

	dg, err := discordgo.New("Bot " + cfg.DiscordBotToken)
	if err != nil {
		logrus.Fatalf("%+v\n", err)
	}
	defer dg.Close()

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	if err = dg.Open(); err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	logrus.Warn("Caught interrupt signal, exiting...")
}

const SirenMessage = "Siren is active"

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	// In this example, we only care about messages that are "ping".
	logrus.Info("message content:", m.Content)
	switch {
	case m.Content == "siren":
		go siren()

		_, err := s.ChannelMessageSend(m.ChannelID, SirenMessage)
		if err != nil {
			fmt.Println("error sending siren affirmation message:", err)
			failed(s, m)
		}
	case m.Content == "help":
		_, err := s.ChannelMessageSend(m.ChannelID,
			"```siren ........... a loud sound to wake sam\n"+
				"help ........................ this message```\n")
		if err != nil {
			fmt.Println("error sending help message:", err)
			failed(s, m)
		}
	}

}

func failed(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Failed to send you a respone")
	if err != nil {
		fmt.Println("error sending message:", err)
	}
}

func siren() error {
	if err := exec.Command("afplay", "./siren.mp3").Run(); err != nil {
		return fmt.Errorf("error playing siren: %v", err)
	}
	return nil
}
