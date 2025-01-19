# crypto-notify-telegram (Binance Price Alert Bot)

This is a Go application that monitors the price of a specified cryptocurrency pair (e.g., BTC/USDT) using the Binance API and sends price alerts to a Telegram group or user. The bot also allows users to set custom price thresholds and conditions via Telegram commands.

## Features
- Fetches the latest cryptocurrency price from Binance.
- Sends price alerts to a specified Telegram chat when the price crosses a defined threshold.
- Supports dynamic configuration of threshold prices via Telegram commands.
- Allows users to choose whether to notify when the price is greater than or less than the threshold.

## Prerequisites
- Go 1.18 or later installed on your system.
- A Binance API key is **not required**, as this bot uses public endpoints.
- A Telegram bot token from [BotFather](https://t.me/BotFather).
- The chat ID of the Telegram group or user you want to send alerts to.

## Installation
1. Clone this repository:
   ```bash
   git clone <repository-url>
   cd <repository-folder>
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Update the configuration constants in the code:
   - `symbol`: The cryptocurrency pair to monitor (e.g., `BTCUSDT`).
   - `telegramToken`: Your Telegram bot token.
   - `telegramChatID`: The chat ID of the group or user where the bot will send messages. Group chat IDs usually start with `-100`.

4. Build and run the application:
   ```bash
   go run main.go
   ```

## Usage
### Commands
The bot supports the following Telegram commands:

1. **/setprice `< or > <value>`**
   - Sets the price threshold and condition.
   - Example:
     - `/setprice > 50000`: Notify if the price goes above $50,000.
     - `/setprice < 30000`: Notify if the price drops below $30,000.

2. **/getprice**
   - Displays the current price threshold and condition.
   - Example response:
     - `📊 Current threshold: price greater than $50,000.`

### Running the Bot in a Group
1. Add the bot to your Telegram group.
2. Ensure the bot has permission to send messages.
3. Obtain the group chat ID using the `/getUpdates` method of the Telegram Bot API:
   ```bash
   curl -X GET "https://api.telegram.org/bot<your-bot-token>/getUpdates"
   ```
   Look for the `id` field under the `chat` section.

4. Update the `telegramChatID` constant in the code with the group chat ID.
5. Restart the bot.

## Example Output
### In the Group Chat
- When the price crosses the threshold:
  ```
  🚨 Price Alert! BTCUSDT is now above $50,000 (Threshold: $50,000)
  ```

- In response to `/getprice`:
  ```
  📊 Current threshold: price greater than $50,000.
  ```

## Configuration
The following constants can be updated in the code to customize the bot:

| Constant            | Description                                           | Default Value     |
|---------------------|-------------------------------------------------------|-------------------|
| `binanceURL`        | URL for the Binance API.                              | Binance API URL   |
| `symbol`            | The cryptocurrency pair to monitor (e.g., `BTCUSDT`).| `BTCUSDT`         |
| `telegramToken`     | Telegram bot token.                                   | Replace with your token |
| `telegramChatID`    | Telegram chat ID to send alerts.                     | Replace with your chat ID |
| `checkInterval`     | Interval (in seconds) to fetch the price.             | `60s`             |

## Logs
The application logs the following information:
- Current price fetched from Binance.
- Price alerts sent to the Telegram chat.
- Errors encountered while fetching prices or sending messages.

## Notes
- Ensure the bot is not in privacy mode if you want it to process all group messages. You can disable privacy mode via BotFather with the `/setprivacy` command.
- For groups, ensure you use the correct chat ID (starting with `-100`) and that the bot has sufficient permissions.

## License
This project is open-source and available under the [MIT License](LICENSE).

## Acknowledgments
- [Binance API](https://binance-docs.github.io/apidocs/spot/en/)
- [go-resty](https://github.com/go-resty/resty) for HTTP requests.
- [tgbotapi](https://github.com/go-telegram-bot-api/telegram-bot-api) for Telegram bot integration.
