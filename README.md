# BubuKabum Discord Bot

Discord bot that clones and optimizes 7TV emotes for Discord servers.

## Features

- Clones emotes from 7TV
- Optimizes sizes for Discord limits
- Handles static and animated emotes  
- Optional 2-frame GIF conversion

## Prerequisites

- Docker Desktop (Windows/Mac) or Docker Engine (Linux)
- Discord Bot Token & Application ID ([Setup Instructions](https://postimg.cc/gallery/WqhXSfd))

## Installation

### Linux/Mac
1. Clone repo:
```bash
git clone https://github.com/Stefanistkuhl/bubuKabum.git
cd bubuKabum
```
2. Copy environment template:
```bash
cp .env.example .env
```

### Windows
1. Clone repo:
```bash
git clone https://github.com/Stefanistkuhl/bubuKabum.git
cd bubuKabum
```
2. Copy environment template:
```bash
copy .env.example .env
```

### All platforms

3. Edit `.env` with your tokens:
```
TOKEN=your_discord_bot_token
APPID=your_discord_application_id
```

4. Build the Docker image:
```bash
docker build -t bubukabum .
```

5. Start the bot:
```bash
docker compose up -d
```

## Usage

1. Invite bot to server (needs Manage Emojis permission)
2. Use `/emoteclone` command
3. Input 7TV emote URL and optional settings

That's it! The bot handles all image processing in the container.
