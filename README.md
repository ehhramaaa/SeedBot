[![Static Badge](https://img.shields.io/badge/Telegram-Bot%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](http://t.me/seed_coin_bot/app?startapp=5024522783)
[![Static Badge](https://img.shields.io/badge/Telegram-Channel%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/bansos_code)
[![Static Badge](https://img.shields.io/badge/Telegram-Chat%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/bansos_code_chat)

![demo](https://raw.githubusercontent.com/ehhramaaa/SeedBot/main/demo/Screenshot_4.png)

# 🔥🔥 Support Cross Platform 🔥🔥

- #### Why Golang? Cause Less Modules Error, Just Run 1 Command And It's Automatic Download Modules And Run Program

- #### If Your Device Is Android, You Can Use Termux For Run This Bot

## Recommendation before use

#### Go Version >= 1.23

## Features

|             Feature             | Supported |
| :-----------------------------: | :-------: |
|          Auto Farming           |    ✅     |
|         Multithreading          |    ✅     |
|         Auto Bird Hunt          |    ✅     |
|         Auto Feed Bird          |    ✅     |
|         Auto Hatch Egg          |    ✅     |
|         Use Query Data          |    ✅     |
|         Auto Catch Worm         |    ✅     |
|        Auto Upgrade Tree        |    ✅     |
|       Auto Play Spin Egg        |    ✅     |
|        Random User Agent        |    ✅     |
|      Auto Claim First Egg       |    ✅     |
|      Auto Check Inventory       |    ✅     |
|      Auto Upgrade Storage       |    ✅     |
|     Auto Claim Login Bonus      |    ✅     |
|     Auto Upgrade Holy Water     |    ✅     |
|    Auto Claim Streak Reward     |    ✅     |
|    Auto Completing Main Task    |    ✅     |
|  Auto Increase Bird Happiness   |    ✅     |
| Auto Completing Holy Water Task |    ✅     |
|        Auto Fusion Piece        |    ⏳     |
|       Auto Connect Wallet       |    ⏳     |
|       Proxy Socks5 / HTTP       |    ✅     |

## [Settings](https://github.com/ehhramaaa/SeedBot/blob/main/config.yml)

|               Settings                |                             Description                             |
| :-----------------------------------: | :-----------------------------------------------------------------: |
|          **AUTO_FEED_BIRD**           |                  Auto Feed Bird If Worm Available                   |
|          **AUTO_BIRD_HUNT**           |                 Auto Bird Hunt If Energy Available                  |
|          **AUTO_HATCH_EGG**           |              Auto Hatch Egg If Available In Inventory               |
|        **AUTO_PLAY_SPIN_EGG**         |               Auto Play Spin Egg If Ticket Available                |
|        **AUTO_UPGRADE.SPEED**         |                 Auto Upgrade Tree If Balance Enough                 |
|       **AUTO_UPGRADE.STORAGE**        |               Auto Upgrade Storage If Balance Enough                |
|     **CLAIM_FARMING_SEED_AFTER**      | Delay Before Claim Farming Seed (e.g. MIN:3600, MAX:7200) In Second |
|      **AUTO_UPGRADE.HOLY_WATER**      |        Auto Upgrade Holy Water If Task Holy Water Completed         |
|   **AUTO_UPGRADE.MAX_LEVEL.SPEED**    |                  Max Upgrade Level Amount For Tree                  |
|  **AUTO_UPGRADE.MAX_LEVEL.STORAGE**   |                Max Upgrade Level Amount For Storage                 |
| **AUTO_UPGRADE.MAX_LEVEL.HOLY_WATER** |               Max Upgrade Level Amount For Holy Water               |
|           **RANDOM_SLEEP**            |    Delay Before The Next Lap (e.g. MIN:3600, MAX:7200) In Second    |

## Prerequisites 📚

Before you begin, make sure you have the following installed:

- [Golang](https://go.dev/doc/install) **Go Version Tested 1.23.1**
- Fresh Query Data Because If The Query Data Can Expired In Seed
- Remove .example From Files With Have That Name And Insert With Your Data
- **If You Using Termux**
  ```shell
  apt update && apt upgrade -y
  apt install golang
  go version
  ```

## Installation

You can download the [**repository**](https://github.com/ehhramaaa/SeedBot.git) by cloning it to your system and installing the necessary dependencies:

```shell
git clone https://github.com/ehhramaaa/SeedBot.git
cd SeedBot
go run .
```

Then you can do build application by typing:

Windows:

```shell
go build -o SeedBot.exe
```

Linux:

```shell
go build -o SeedBot
```

## Usage

```shell
go run .
```

Or

```shell
go run main.go
```

**If You Want Auto Select Choice In Terminal**

For Option 1

```shell
go run . -c 1
```

For Option 2

```shell
go run . -c 2
```
