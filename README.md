# TefCON 2023 Telegram Bot

This repository contains the source code for the Telegram bot used for TefCON 2023.

## Setup

### Prerequisites

- Golang 1.21 or higher

### Installation

1. Clone the repository
2. Run `go build` in the root directory
3. Run the executable

## Configuration

The bot can be configured using environment variables. The following variables are available:

| Variable             | Description                   | Default |
| -------------------- | ----------------------------- | ------- |
| `TELEGRAM_BOT_TOKEN` | The token of the Telegram bot | `""`    |

## Usage

First, start a converstation with [tefconbot](https://t.me/tefconbot). The bot can be used to query the events of the TEFCon 2023. The following commands are available:

- `/start`: Starts the bot and shows a welcome message
- `/help`: Shows a help message
- `/map`: Shows a map of the event
- `/current_events`: Shows the current events
- `/next_events`: Shows the next events


