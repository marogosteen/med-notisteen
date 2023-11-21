package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var DiscordToken string
var ChannelID string

func loadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DiscordToken = os.Getenv("TOKEN")
	ChannelID = os.Getenv("CHANNEL_ID")
}

// Function to schedule event and send notification
func scheduleEventAndNotify(discordSession *discordgo.Session) {
	location, _ := time.LoadLocation("Asia/Tokyo") // Load Japan timezone
	now := time.Now().In(location)

	// Calculate next notification time
	next := now.Add(time.Hour * 24)
	next = time.Date(next.Year(), next.Month(), next.Day(), 7, 0, 0, 0, location) // Set to next 7 AM

	if now.Hour() >= 7 { // If it's already past 7 AM, schedule for the next day
		next = next.Add(time.Hour * 24)
	}

	duration := next.Sub(now)

	// Schedule the notification
	time.AfterFunc(duration, func() {
		// Send a message to a specific channel
		_, err := discordSession.ChannelMessageSend(ChannelID, "7時になったぞ!野間薬飲めや!")
		if err != nil {
			fmt.Println("error sending message,", err)
			return
		}

		// Schedule the next notification
		scheduleEventAndNotify(discordSession)
	})
}

func main() {
	loadEnvVariables()
	discordSession, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Open a websocket connection to Discord and begin listening.
	err = discordSession.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Schedule event and send notification
	go scheduleEventAndNotify(discordSession)
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-signalChannel

	// Cleanly close down the Discord session.
	err = discordSession.Close()
	if err != nil {
		fmt.Println("error closing the connection,", err)
		return
	}
}