package core

import (
	"SeedBot/helper"
	"fmt"
	"net/http"
	"time"

	"github.com/gookit/config/v2"
)

func launchBot(account *Account, config *config.Config, walletAddress string, proxy string, isBindWallet bool) {
	var isClaimFirstEgg, isClaimFarmingSeed bool
	var speedSeedLevel, storageSeedLevel, holyWaterLevel int

	isAutoHatchEgg := config.Bool("AUTO_HATCH_EGG")
	isAutoFeedBird := config.Bool("AUTO_FEED_BIRD")
	isAutoBirdHunt := config.Bool("AUTO_BIRD_HUNT")
	isAutoPlaySpinEgg := config.Bool("AUTO_PLAY_SPIN_EGG")
	isAutoUpgradeSpeed := config.Bool("AUTO_UPGRADE.SPEED")
	isAutoUpgradeStorage := config.Bool("AUTO_UPGRADE.STORAGE")
	isAutoUpgradeHolyWater := config.Bool("AUTO_UPGRADE.HOLY_WATER")

	claimFarmingSeedAfter := config.Int("CLAIM_FARMING_SEED_AFTER")

	client := &Client{
		username:   account.Username,
		proxy:      proxy,
		authToken:  account.QueryData,
		httpClient: &http.Client{},
	}

	ip := client.checkIp()
	if len(ip) > 0 {
		helper.PrettyLog("success", fmt.Sprintf("| %s | Check Ip Successfully | Ip Address: %s", client.username, ip))
	} else {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Check Ip Failed | Please Check Your Proxy...", client.username))
		return
	}

	profile := client.getProfile()
	if profile == nil {
		return
	}
	if data, exits := profile["data"].(map[string]interface{}); exits && data != nil {
		var isWalletConnected bool

		if !data["give_first_egg"].(bool) {
			isClaimFirstEgg = true
		}

		parsedTime, err := time.Parse(time.RFC3339Nano, data["last_claim"].(string))
		if err != nil {
			helper.PrettyLog("error", fmt.Sprintf("| %s | Check Last Claim Time Seed Failed | Failed to parse time: %v", client.username, err))
		} else {
			if (parsedTime.Unix() + int64(claimFarmingSeedAfter)) < time.Now().Unix() {
				isClaimFarmingSeed = true
			}
		}

		if upgradesInfo, exits := data["upgrades"].([]interface{}); exits {
			for _, info := range upgradesInfo {
				infoMap := info.(map[string]interface{})
				if infoMap["upgrade_type"] == "mining-speed" {
					speedSeedLevel = int(infoMap["upgrade_level"].(float64))
				}
				if infoMap["upgrade_type"] == "storage-size" {
					storageSeedLevel = int(infoMap["upgrade_level"].(float64))
				}
				if infoMap["upgrade_type"] == "holy-water" {
					holyWaterLevel = int(infoMap["upgrade_level"].(float64))
				}
			}
		}

		if len(data["wallet_address_ton"].(string)) > 0 {
			isWalletConnected = true
		}

		helper.PrettyLog("success", fmt.Sprintf("| %s | Claimed First Egg: %v | Speed Seed Level: %v | Storage Seed Level: %v | Holy Water Level: %v | Wallet Connected: %v", client.username, data["give_first_egg"].(bool), speedSeedLevel, storageSeedLevel, holyWaterLevel, isWalletConnected))
	}

	balance := client.getBalance()
	if data, exits := balance["data"].(float64); exits && balance != nil {
		helper.PrettyLog("success", fmt.Sprintf("| %s | Balance: %.9f", client.username, (data/1e9)))
	}

	checkGuild := client.checkGuild()
	if data, exits := checkGuild["data"].(map[string]interface{}); exits && data != nil {
		if data["guild_id"].(string) != "9e02254f-d921-43d3-839f-903706dedeb5" {
			client.leaveGuild()
			time.Sleep(3 * time.Second)
			client.joinGuild()
			time.Sleep(3 * time.Second)
		}
	} else {
		client.joinGuild()
	}

	guildInfo := client.getGuildInfo()
	if data, exits := guildInfo["data"].(map[string]interface{}); exits && data != nil {
		helper.PrettyLog("success", fmt.Sprintf("| %s | %s | Members: %v | Hunted: %.9f | Reward: %.9f | Rank: %v", client.username, data["name"].(string), int(data["number_member"].(float64)), (data["hunted"].(float64)/1e9), (data["reward"].(float64)/1e9), int(data["rank_index"].(float64))))
	}

	birdInventory := client.getBirdInventory()
	wormInventory := client.getWormInventory()
	eggInventory := client.getEggInventory()
	if birdInventory != nil && wormInventory != nil && eggInventory != nil {
		if bird, exists := birdInventory["data"].(map[string]interface{}); exists {
			if worm, exists := wormInventory["data"].(map[string]interface{}); exists {
				if egg, exists := eggInventory["data"].(map[string]interface{}); exists {
					helper.PrettyLog("success", fmt.Sprintf("| %s | Bird Inventory: %v | Worm Inventory: %v | Egg Inventory: %v", client.username, int(bird["total"].(float64)), int(worm["total"].(float64)), int(egg["total"].(float64))))
				}
			}
		}
	}

	checkLoginBonus := client.checkLoginBonus()
	if data, exits := checkLoginBonus["data"].([]interface{}); exits && data != nil {
		if len(data) > 0 {
			for _, detail := range data {
				detailMap := detail.(map[string]interface{})
				parsedTime, err := time.Parse(time.RFC3339Nano, detailMap["timestamp"].(string))
				if err != nil {
					helper.PrettyLog("error", fmt.Sprintf("| %s | Claim Login Bonus Failed | Failed to parse time: %v", client.username, err))
					continue
				}

				if parsedTime.Format("2006/01/02") != time.Now().Format("2006/01/02") {
					claimLoginBonus := client.claimLoginBonus()

					if data, exits := claimLoginBonus["data"].(map[string]interface{}); exits && claimLoginBonus != nil {
						parsedTime, err := time.Parse(time.RFC3339Nano, data["timestamp"].(string))
						if err != nil {
							helper.PrettyLog("error", fmt.Sprintf("| %s | Claim Login Bonus Failed | Failed to parse time: %v", client.username, err))
							continue
						}

						if parsedTime.Format("2006/01/02") == time.Now().Format("2006/01/02") {
							helper.PrettyLog("success", fmt.Sprintf("| %s | Claim Login Bonus Successfully | Amount: %.9f", client.username, (data["amount"].(float64)/1e9)))
						}
					}
				}
			}
		}
	}

	if isClaimFarmingSeed {
		claimSeed := client.claimFarmingSeed()
		if data, exits := claimSeed["data"].(map[string]interface{}); exits && data != nil {
			helper.PrettyLog("success", fmt.Sprintf("| %s | Claim Seed Successfully | Amount: %.9f", client.username, (data["amount"].(float64)/1e9)))
		}
	}

	if isClaimFirstEgg {
		claimFirstEgg := client.claimFirstEgg()
		if data, exits := claimFirstEgg["data"].(map[string]interface{}); exits && data != nil {
			helper.PrettyLog("success", fmt.Sprintf("| %s | Claim First Egg Successfully | Egg Type: %s | Status: %s", client.username, data["type"].(string), data["status"].(string)))
		}
	}

	if isAutoHatchEgg {
		eggInventory = client.getEggInventory()
		if eggInventory != nil {
			if data, exits := eggInventory["data"].(map[string]interface{}); exits && data != nil {
				if int(data["total"].(float64)) > 0 {
					egg := data["items"].([]interface{})
					for _, item := range egg {
						itemMap := item.(map[string]interface{})
						hatchEgg := client.hatchEgg(itemMap["id"].(string))
						if hatchEgg != nil {
							if data, exits := hatchEgg["data"].(map[string]interface{}); exits && data != nil {
								helper.PrettyLog("success", fmt.Sprintf("| %s | Hatch Egg Successfully | Bird Type: %s", client.username, data["type"].(string)))
							}
						}
					}
				}
			}
		}
	}

	birdStatus := client.getBirdStatus()
	if data, exits := birdStatus["data"].(map[string]interface{}); exits && data != nil {
		helper.PrettyLog("success", fmt.Sprintf("| %s | Bird Status: %s | Type: %s | Energy Level: %v | Energy Max: %v | Happiness: %v | Task Level: %v", client.username, data["status"].(string), data["type"].(string), int(data["energy_level"].(float64)/1e9), int(data["energy_max"].(float64)/1e9), int(data["happiness_level"].(float64)/1e9), int(data["task_level"].(float64)/1e9)))
	}

	wormStatus := client.getWormStatus()
	if data, exits := wormStatus["data"].(map[string]interface{}); exits && data != nil {
		parsedTime, err := time.Parse(time.RFC3339, data["ended_at"].(string))
		if err != nil {
			helper.PrettyLog("error", fmt.Sprintf("| %s | Check Worm Failed | Failed to parse time: %v", client.username, err))
		}

		if parsedTime.Unix() < time.Now().Unix() {
			catchWorm := client.catchWorm()
			if data, exits := catchWorm["data"].(map[string]interface{}); exits && data != nil {
				helper.PrettyLog("success", fmt.Sprintf("| %s | Catch Worm Successfully | Worm Type: %s | Catch Status: %s", client.username, data["type"].(string), data["status"].(string)))
			}
		}
	}

	if isAutoFeedBird {
		birdInventory = client.getBirdInventory()
		wormInventory = client.getWormInventory()
		if birdInventory != nil && wormInventory != nil {
			if birdsData, exists := birdInventory["data"].(map[string]interface{}); exists {
				if int(birdsData["total"].(float64)) > 0 {
					if wormsData, exists := wormInventory["data"].(map[string]interface{}); exists {
						if int(wormsData["total"].(float64)) > 0 {
							birds := birdsData["items"].([]interface{})
							for _, bird := range birds {
								birdMap := bird.(map[string]interface{})
								if birdMap["is_leader"].(bool) {
									worms := wormsData["items"].([]interface{})
									for _, worm := range worms {
										wormMap := worm.(map[string]interface{})
										feedBird := client.feedBird(birdMap["id"].(string), wormMap["id"].(string))
										if feedBird != nil {
											if data, exits := feedBird["data"].(map[string]interface{}); exits && data != nil {
												helper.PrettyLog("success", fmt.Sprintf("| %s | Feed Bird Successfully | Current Energy: %v | Max Energy: %v", client.username, (data["energy_level"].(float64)/100), (data["energy_max"].(float64)/100)))
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	mainTask := client.getTasks()
	if mainTask != nil {
		if data, exits := mainTask["data"].([]interface{}); exits && data != nil {
			for _, task := range data {
				taskMap := task.(map[string]interface{})
				if taskMap["task_user"] != nil {
					taskUser := taskMap["task_user"].(map[string]interface{})
					if !taskUser["completed"].(bool) {
						claimTask := client.claimTask(taskUser["id"].(string), taskMap["name"].(string))
						if claimTask != nil {
							if data, exits := claimTask["data"].(map[string]interface{}); exits && data != nil {
								if status, exits := data["data"].(map[string]interface{}); exits {
									if status["completed"].(bool) {
										helper.PrettyLog("success", fmt.Sprintf("| %s | Claim Task %s Successfully | Sleep 15s Before Next Task...", client.username, taskMap["name"].(string)))
									}
								} else {
									helper.PrettyLog("error", fmt.Sprintf("| %s | Claim Task %s Failed | Status: %s | You Can Try Manual | Sleep 15s Before Next Task...", client.username, taskMap["name"].(string), data["error"].(string)))
								}
							}
						}
					}
				} else {
					startTask := client.startTask(taskMap["id"].(string), taskMap["name"].(string))
					if startTask != nil {
						if taskId, exits := startTask["data"].(string); exits && len(taskId) > 0 {
							helper.PrettyLog("success", fmt.Sprintf("| %s | Start Task %s Successfully | Sleep 5s Before Claim Task...", client.username, taskMap["name"].(string)))

							time.Sleep(5 * time.Second)

							claimTask := client.claimTask(taskId, taskMap["name"].(string))
							if claimTask != nil {
								if data, exits := claimTask["data"].(map[string]interface{}); exits && data != nil {
									if status, exits := data["data"].(map[string]interface{}); exits {
										if status["completed"].(bool) {
											helper.PrettyLog("success", fmt.Sprintf("| %s | Claim Task %s Successfully | Sleep 15s Before Next Task...", client.username, taskMap["name"].(string)))
										}
									} else {
										helper.PrettyLog("error", fmt.Sprintf("| %s | Claim Task %s Failed | Status: %s | You Can Try Manual | Sleep 15s Before Next Task...", client.username, taskMap["name"].(string), data["error"].(string)))
									}
								}
							}
						}
					}
				}

				time.Sleep(15 * time.Second)
			}
		}
	}

	holyWaterTask := client.getTaskHolyWater()
	if holyWaterTask != nil {
		if data, exits := holyWaterTask["data"].([]interface{}); exits && data != nil {
			for _, task := range data {
				taskMap := task.(map[string]interface{})
				if taskMap["task_user"] == nil {
					var isCompletingReferTask bool
					friendsInfo := client.getFriendsInfo()
					if friendsInfo != nil {
						if data, exits := friendsInfo["data"].(map[string]interface{}); exits && data != nil {
							if len(data["referees"].([]interface{})) > 0 {
								isCompletingReferTask = true
							}
						}
					}

					if !isCompletingReferTask && taskMap["type"].(string) == "refer" {
						continue
					}

					startTask := client.startTaskHolyWater(taskMap["id"].(string), taskMap["name"].(string))
					if startTask != nil {
						if taskId, exits := startTask["data"].(string); exits && len(taskId) > 0 {
							helper.PrettyLog("success", fmt.Sprintf("| %s | Start Holy Water Task %s Successfully | Sleep 15s Before Claim Task...", client.username, taskMap["name"].(string)))

							time.Sleep(15 * time.Second)

							claimTask := client.claimTask(taskId, taskMap["name"].(string))
							if claimTask != nil {
								if data, exits := claimTask["data"].(map[string]interface{}); exits && data != nil {
									status := data["data"].(map[string]interface{})
									if status["id"].(string) == taskId {
										helper.PrettyLog("success", fmt.Sprintf("| %s | Claim Task %s Successfully | Sleep 15s Before Next Task...", client.username, taskMap["name"].(string)))
									}
								}
							}
						}
					}
				}

				time.Sleep(15 * time.Second)
			}
		}
	}

	if isAutoBirdHunt {
		birdStatus = client.getBirdStatus()
		if data, exits := birdStatus["data"].(map[string]interface{}); exits && data != nil {
			if (int(data["happiness_level"].(float64)) / 100) < 100 {
				birdHappiness := client.birdHappiness(data["id"].(string), 100)
				if data, exits := birdHappiness["data"].(map[string]interface{}); exits && data != nil {
					helper.PrettyLog("success", fmt.Sprintf("| %s | Bird Happiness Successfully | Current Happiness: %v", client.username, (int(data["happiness_level"].(float64))/100)))
				}
			}

			parsedTime, err := time.Parse(time.RFC3339, data["hunt_end_at"].(string))
			if err != nil {
				helper.PrettyLog("error", fmt.Sprintf("| %s | Bird Hunt Failed | Failed to parse time: %v", client.username, err))
			} else {
				if parsedTime.Unix() < time.Now().Unix() {
					if data["hunt_end_at"].(string) != "0001-01-01T00:00:00Z" {
						claimBirdHunt := client.claimBirdHunt(data["id"].(string))
						if data, exits := claimBirdHunt["data"].(map[string]interface{}); exits && data != nil {
							helper.PrettyLog("success", fmt.Sprintf("| %s | Claim Bird Hunt Successfully | Sleep 3s Before Start Hunt", client.username))
						}

						time.Sleep(3 * time.Second)
					}

					if int(data["energy_level"].(float64)/1e9) > 0 {
						startBirdHunt := client.startBirdHunt(data["id"].(string), int(data["task_level"].(float64)))
						if data, exits := startBirdHunt["data"].(map[string]interface{}); exits && data != nil {
							parsedTime, err := time.Parse(time.RFC3339, data["hunt_end_at"].(string))
							if err != nil {
								helper.PrettyLog("success", fmt.Sprintf("| %s | Bird Hunt Successfully | Claim Reward After: %v", client.username, data["hunt_end_at"].(string)))
							} else {
								helper.PrettyLog("success", fmt.Sprintf("| %s | Bird Hunt Successfully | Claim Reward After: %vs", client.username, (parsedTime.Unix()-time.Now().Unix())))
							}
						}
					} else {
						helper.PrettyLog("info", fmt.Sprintf("| %s | Failed Start Bird Hunt | Energy Bird Not Enough", client.username))
					}
				} else {
					helper.PrettyLog("info", fmt.Sprintf("| %s | Bird Still Hunting | Claim After %vs", client.username, (parsedTime.Unix()-time.Now().Unix())))
				}
			}
		}
	}

	streakReward := client.getStreakReward()
	if data, exits := streakReward["data"].([]interface{}); exits && data != nil {
		if len(data) > 0 {
			for _, streak := range data {
				streakMap := streak.(map[string]interface{})
				claimStreakReward := client.claimStreakReward(streakMap["id"].(string))
				if data, exits := claimStreakReward["data"].([]interface{}); exits && data != nil {
					for _, status := range data {
						statusMap := status.(map[string]interface{})
						if statusMap["status"].(string) == "received" {
							helper.PrettyLog("success", fmt.Sprintf("| %s | Claim Streak Reward Successfully | Streak Reward: %s", client.username, statusMap["status"].(string)))
						}
					}
				}
			}
		}
	}

	if isAutoPlaySpinEgg {
		spinEgg := client.getSpinEggTicket()
		if data, exits := spinEgg["data"].([]interface{}); exits && data != nil {
			if len(data) > 0 {
				for _, ticket := range data {
					ticketMap := ticket.(map[string]interface{})
					playSpinEgg := client.claimSpinEgg(ticketMap["id"].(string))
					if data, exits := playSpinEgg["data"].(map[string]interface{}); exits && data != nil {
						helper.PrettyLog("success", fmt.Sprintf("| %s | Play Spin Egg Successfully | Reward Status: %s | Type: %s", client.username, data["status"].(string), data["type"].(string)))
					}
				}
			}
		}
	}

	if isAutoUpgradeSpeed || isAutoUpgradeStorage || isAutoUpgradeHolyWater {
		settings := client.getSettings()
		if data, exits := settings["data"].(map[string]interface{}); exits && data != nil {
			speedSeedCosts := data["mining-speed-costs"].([]interface{})
			storageSeedCosts := data["mining-speed-costs"].([]interface{})

			if isAutoUpgradeSpeed {
				var currentBalance float64
				maxLevel := config.Int("AUTO_UPGRADE.MAX_LEVEL.SPEED")
				balance = client.getBalance()
				if data, exits := balance["data"].(float64); exits && balance != nil {
					currentBalance = data / 1e9
				}

				if speedSeedLevel < maxLevel {
					if currentBalance > (float64(speedSeedCosts[(speedSeedLevel-1)].(float64)) / 1e9) {
						upgradeSpeed := client.upgradeSpeedSeed()
						if upgradeSpeed == "{}" {
							helper.PrettyLog("success", fmt.Sprintf("| %s | Upgrade Speed Seed Successfully | Current Speed Level: %v", client.username, (speedSeedLevel+1)))
						}
					} else {
						helper.PrettyLog("error", fmt.Sprintf("| %s | Upgrade Speed Seed Failed | Not Enough Balance...", client.username))
					}
				}
			}

			if isAutoUpgradeStorage {
				var currentBalance float64
				maxLevel := config.Int("AUTO_UPGRADE.MAX_LEVEL.STORAGE")
				balance = client.getBalance()
				if data, exits := balance["data"].(float64); exits && balance != nil {
					currentBalance = data / 1e9
				}

				if storageSeedLevel < maxLevel {
					if currentBalance > (float64(storageSeedCosts[(storageSeedLevel-1)].(float64)) / 1e9) {
						upgradeStorage := client.upgradeStorageSeed()
						if upgradeStorage == "{}" {
							helper.PrettyLog("success", fmt.Sprintf("| %s | Upgrade Storage Seed Successfully | Current Storage Level: %v", client.username, (storageSeedLevel+1)))
						}
					} else {
						helper.PrettyLog("error", fmt.Sprintf("| %s | Upgrade Storage Seed Failed | Not Enough Balance...", client.username))
					}
				}
			}

			if isAutoUpgradeHolyWater {
				var taskCompleted int
				maxLevel := config.Int("AUTO_UPGRADE.MAX_LEVEL.HOLY_WATER")
				holyWaterTask = client.getTaskHolyWater()
				if holyWaterTask != nil {
					if data, exits := holyWaterTask["data"].([]interface{}); exits && data != nil {
						for _, task := range data {
							taskMap := task.(map[string]interface{})
							if taskMap["task_user"] == nil {
								taskCompleted++
							}
						}
					}
				}

				if holyWaterLevel < maxLevel {
					if taskCompleted > holyWaterLevel {
						upgradeHolyWater := client.upgradeHolyWater()
						if upgradeHolyWater == "{}" {
							helper.PrettyLog("success", fmt.Sprintf("| %s | Upgrade Holy Water Successfully | Current Holy Water Level: %v", client.username, (holyWaterLevel+1)))
						}
					} else {
						helper.PrettyLog("error", fmt.Sprintf("| %s | Upgrade Holy Water Failed | Not Enough Balance...", client.username))
					}
				}
			}
		}
	}
}
