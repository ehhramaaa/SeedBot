package core

import (
	"SeedBot/helper"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

func (c *Client) setProxy() error {
	// Parse the proxy URL
	parsedURL, err := url.Parse(c.proxy)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %v", err)
	}

	// Extract username and password if available
	var proxyUser, proxyPass string
	if parsedURL.User != nil {
		proxyUser = parsedURL.User.Username()
		proxyPass, _ = parsedURL.User.Password()
	}

	// Handle based on scheme
	switch parsedURL.Scheme {
	case "http", "https":
		// HTTP Proxy
		transport := &http.Transport{
			Proxy: http.ProxyURL(parsedURL), // Handles HTTP Proxy with auth
		}
		c.httpClient = &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}

	case "socks5":
		// SOCKS5 Proxy
		var auth *proxy.Auth
		if proxyUser != "" && proxyPass != "" {
			auth = &proxy.Auth{
				User:     proxyUser,
				Password: proxyPass,
			}
		}

		// Create SOCKS5 dialer
		dialer, err := proxy.SOCKS5("tcp", parsedURL.Host, auth, proxy.Direct)
		if err != nil {
			return fmt.Errorf("error creating SOCKS5 dialer: %v", err)
		}

		// Set the transport to use the SOCKS5 dialer
		transport := &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}
		c.httpClient = &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}

	default:
		return fmt.Errorf("unsupported proxy scheme: %s", parsedURL.Scheme)
	}

	return nil
}

func (c *Client) makeRequest(method string, url string, jsonBody interface{}) ([]byte, error) {
	var err error
	// Set proxy if available
	err = c.setProxy()
	if err != nil {
		return nil, err
	}

	// Convert body to JSON
	var reqBody []byte
	if jsonBody != nil {
		reqBody, err = json.Marshal(jsonBody)
		if err != nil {
			return nil, err
		}
	}

	// Create new request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	setHeader(req, c.authToken)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 status code
	if resp.StatusCode >= 400 {
		// Read the response body to include in the error message
		bodyBytes, bodyErr := io.ReadAll(resp.Body)
		if bodyErr != nil {
			return nil, fmt.Errorf("error status: %v, and failed to read body: %v", resp.StatusCode, bodyErr)
		}
		return nil, fmt.Errorf("error status: %v, error message: %s", resp.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}

// Check Ip
func (c *Client) checkIp() string {
	res, err := c.makeRequest("GET", "https://ipinfo.io/ip", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to check ip: %v", c.username, err))
		return ""
	}

	return string(res)
}

// Get Profile
func (c *Client) getProfile() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/profile", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get profile: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Balance
func (c *Client) getBalance() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/profile/balance", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get balance: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Check Guild
func (c *Client) checkGuild() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/guild/member/detail", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to check guild: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Join Guild
func (c *Client) joinGuild() map[string]interface{} {
	payload := map[string]string{
		"guild_id": "9e02254f-d921-43d3-839f-903706dedeb5",
	}
	req, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/guild/join", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to join squad: %v", c.username, err))
		return nil
	}

	res, err := handleResponseMap(req)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return res
}

// Leave Guild
func (c *Client) leaveGuild() map[string]interface{} {
	payload := map[string]string{
		"guild_id": "9e02254f-d921-43d3-839f-903706dedeb5",
	}

	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/guild/leave", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to leave squad: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Guild Info
func (c *Client) getGuildInfo() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/guild/detail?guild_id=9e02254f-d921-43d3-839f-903706dedeb5&sort_by=total_hunted", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get guild info: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Bird Inventory
func (c *Client) getBirdInventory() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/bird/me?page=1", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get bird inventory: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Worm Inventory
func (c *Client) getWormInventory() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/worms/me?page=1", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get worm inventory: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Egg Inventory
func (c *Client) getEggInventory() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/egg/me?page=1", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get egg inventory: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Check Login Bonus
func (c *Client) checkLoginBonus() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/login-bonuses", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to check login bonus: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Claim Login Bonus
func (c *Client) claimLoginBonus() map[string]interface{} {
	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/login-bonuses", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to claim login bonus: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Claim Farming Seed
func (c *Client) claimFarmingSeed() map[string]interface{} {
	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/seed/claim", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to claim farming seed: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Claim First Egg
func (c *Client) claimFirstEgg() map[string]interface{} {
	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/give-first-egg", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get first egg: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Hatch Egg
func (c *Client) hatchEgg(eggId string) map[string]interface{} {
	payload := map[string]interface{}{
		"egg_id": eggId,
	}
	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/egg-hatch/complete", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to hatch egg: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Bird Status
func (c *Client) getBirdStatus() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/bird/is-leader", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get bird status: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Worm Status
func (c *Client) getWormStatus() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/worms", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get worm status: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Catch Worm
func (c *Client) catchWorm() map[string]interface{} {
	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/worms/catch", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to catch worm: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Feed Bird
func (c *Client) feedBird(birdId, wormId string) map[string]interface{} {
	wormIds := []string{wormId}
	payload := map[string]interface{}{
		"bird_id":  birdId,
		"worm_ids": wormIds,
	}

	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/bird-feed", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to feed bird: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Bird Happiness
func (c *Client) birdHappiness(birdId string, happinessRate int) map[string]interface{} {
	payload := map[string]interface{}{
		"bird_id":        birdId,
		"happiness_rate": happinessRate * 100,
	}

	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/bird-happiness", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to happiness bird: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Start Bird Hunt
func (c *Client) startBirdHunt(birdId string, taskLevel int) map[string]interface{} {
	payload := map[string]interface{}{
		"bird_id":    birdId,
		"task_level": taskLevel,
	}

	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/bird-hunt/start", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to start bird hunt: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Claim Bird Hunt
func (c *Client) claimBirdHunt(birdId string) map[string]interface{} {
	payload := map[string]string{
		"bird_id": birdId,
	}

	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/bird-hunt/complete", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to claim bird hunt: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Tasks
func (c *Client) getTasks() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/tasks/progresses", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get task: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Start Task
func (c *Client) startTask(taskId, taskName string) map[string]interface{} {
	res, err := c.makeRequest("POST", fmt.Sprintf("https://elb.seeddao.org/api/v1/tasks/%s", taskId), nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to start task %v: %v", c.username, taskName, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Friends Info
func (c *Client) getFriendsInfo() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/profile/recent-referees", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get friends info: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Task Holy Water
func (c *Client) getTaskHolyWater() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/upgrades/tasks/progresses", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get task holy water: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Start Task Holy Water
func (c *Client) startTaskHolyWater(taskId, taskName string) map[string]interface{} {
	res, err := c.makeRequest("GET", fmt.Sprintf("https://elb.seeddao.org/api/v1/upgrades/tasks/%s", taskId), nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to start task holy water %v: %v", c.username, taskName, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Claim Task
func (c *Client) claimTask(taskId, taskName string) map[string]interface{} {
	res, err := c.makeRequest("GET", fmt.Sprintf("https://elb.seeddao.org/api/v1/tasks/notification/%s", taskId), nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to claim task %v: %v", c.username, taskName, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Streak Reward
func (c *Client) getStreakReward() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/streak-reward", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get streak reward: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Claim Streak Reward
func (c *Client) claimStreakReward(streakId string) map[string]interface{} {
	streakIds := []string{streakId}
	payload := map[string]interface{}{
		"streak_reward_ids": streakIds,
	}

	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/streak-reward", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to claim streak reward: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Spin Egg Ticket
func (c *Client) getSpinEggTicket() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/spin-ticket", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get spin egg ticket: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Claim Spin Egg
func (c *Client) claimSpinEgg(ticketId string) map[string]interface{} {
	payload := map[string]string{
		"ticket_id": ticketId,
	}

	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/spin-reward", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to claim spin egg: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Get Settings
func (c *Client) getSettings() map[string]interface{} {
	res, err := c.makeRequest("GET", "https://elb.seeddao.org/api/v1/settings", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get settings: %v", c.username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Error handling response: %v", c.username, err))
		return nil
	}

	return result
}

// Upgrade Storage Seed
func (c *Client) upgradeStorageSeed() string {
	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/seed/storage-size/upgrade", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to upgrade storage seed: %v", c.username, err))
		return ""
	}

	return string(res)
}

// Upgrade Speed Seed
func (c *Client) upgradeSpeedSeed() string {
	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/seed/mining-speed/upgrade", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to upgrade speed seed: %v", c.username, err))
		return ""
	}

	return string(res)
}

// Upgrade Holy Water
func (c *Client) upgradeHolyWater() string {
	res, err := c.makeRequest("POST", "https://elb.seeddao.org/api/v1/upgrades/holy-water", nil)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to upgrade holy water: %v", c.username, err))
		return ""
	}

	return string(res)
}

// Todo Fusion Piece
