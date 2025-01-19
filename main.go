package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	binanceURL     = "https://api.binance.com/api/v3/ticker/price"
	symbol         = "BTCUSDT"                                        // Replace with your desired symbol (e.g., ETHUSDT)
	telegramToken  = "<your-bot-telegram-token>"                      // Replace with your Telegram bot token
	telegramChatID = "<your-bot-telegram-chatID>"                     // Replace with your Telegram chat ID
	checkInterval  = 60 * time.Second                                 // Interval to check the price
)

var (
	thresholdPrice float64 = 92000.0 // Default threshold price
	notifyGreater  bool    = false   // Notify if price > threshold (default: false)
)

func getPrice(symbol string) (float64, error) {
	client := resty.New()
	resp, err := client.R().
		SetQueryParam("symbol", symbol).
		Get(binanceURL)

	if err != nil {
		return 0, err
	}

	if resp.StatusCode() != 200 {
		return 0, errors.New("failed to fetch data from Binance")
	}

	// Parse JSON response
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return 0, err
	}

	priceStr, ok := result["price"].(string)
	if !ok {
		return 0, errors.New("price not found in Binance response")
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func sendTelegramMessage(bot *tgbotapi.BotAPI, chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	return err
}

func handleCommands(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Ensure the update contains a command
	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	// Log the received command
	log.Printf("Received command: %s", update.Message.Text)

	// Switch on the command type
	switch update.Message.Command() {
	case "setprice":
		args := strings.Fields(update.Message.CommandArguments())
		if len(args) < 2 {
			_ = sendTelegramMessage(bot, update.Message.Chat.ID, "❌ Invalid format. Use `/setprice < or > <value>` (e.g., `/setprice > 50000`).")
			return
		}

		// Parse the direction (< or >)
		direction := args[0]
		if direction != "<" && direction != ">" {
			_ = sendTelegramMessage(bot, update.Message.Chat.ID, "❌ Invalid direction. Use `/setprice <value>` with `<` or `>`.")
			return
		}

		// Parse the threshold price
		newPrice, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			_ = sendTelegramMessage(bot, update.Message.Chat.ID, "❌ Invalid price format. Use `/setprice <value>`.")
			return
		}

		// Update global variables
		thresholdPrice = newPrice
		notifyGreater = (direction == ">")

		response := fmt.Sprintf("✅ Threshold price set to $%.2f with condition: price %s %.2f.", thresholdPrice, direction, thresholdPrice)
		_ = sendTelegramMessage(bot, update.Message.Chat.ID, response)

	case "getprice":
		condition := "less than"
		if notifyGreater {
			condition = "greater than"
		}
		response := fmt.Sprintf("📊 Current threshold: price %s $%.2f.", condition, thresholdPrice)
		_ = sendTelegramMessage(bot, update.Message.Chat.ID, response)

	default:
		// Ignore unrecognized commands and do not spam the chat
		_ = sendTelegramMessage(bot, update.Message.Chat.ID, "❓ Unknown command. Use /setprice or /getprice.")
	}
}

func main() {
	// Initialize Telegram Bot
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}
	log.Printf("Authorized on Telegram bot: %s", bot.Self.UserName)

	// Parse Telegram Chat ID
	chatID, err := strconv.ParseInt(telegramChatID, 10, 64)
	if err != nil {
		log.Fatalf("Invalid chat ID: %v", err)
	}

	// Set up an update configuration
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Get updates channel
	updates := bot.GetUpdatesChan(u)

	// Main loop to monitor updates and price
	for {
		select {
		case update := <-updates:
			if update.Message != nil {
				handleCommands(bot, update) // Handle user commands like /setprice
			}
		default:
			// Monitor Binance price in the background
			price, err := getPrice(symbol)
			if err != nil {
				log.Printf("Failed to fetch price: %v", err)
				time.Sleep(checkInterval)
				continue
			}

			log.Printf("Current price of %s: %f", symbol, price)

			// Notify based on the condition (< or >)
			if (notifyGreater && price > thresholdPrice) || (!notifyGreater && price < thresholdPrice) {
				direction := "below"
				if notifyGreater {
					direction = "above"
				}
				message := fmt.Sprintf("🚨 Price Alert! %s is now %s $%.2f (Threshold: $%.2f)", symbol, direction, price, thresholdPrice)
				err := sendTelegramMessage(bot, chatID, message)
				if err != nil {
					log.Printf("Failed to send Telegram message: %v", err)
				} else {
					log.Printf("Notification sent: %s", message)
				}
			}

			time.Sleep(checkInterval)
		}
	}
}
