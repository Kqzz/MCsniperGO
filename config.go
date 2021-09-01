package main

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type ConfigStruct struct {
	Sniper struct {
		TimingSystemPreference string  `toml:"timing_system_preference"`
		CycleTimingSystems     bool    `toml:"cycle_timing_systems"`
		AutoClaimNamemc        bool    `toml:"auto_claim_namemc"`
		SnipeRequests          int     `toml:"snipe_requests"`
		PrenameRequests        int     `toml:"prename_requests"`
		Spread                 float64 `toml:"spread"`
	} `toml:"sniper"`
	Accounts struct {
		MaxAccounts        int `toml:"max_accounts"`
		PrenameMaxAccounts int `toml:"prename_max_accounts"`
		StartAuth          int `toml:"start_auth"`
		AuthDelay          int `toml:"auth_delay"`
	} `toml:"accounts"`
	Skin struct {
		ChangeSkinOnSnipe bool   `toml:"change_skin_on_snipe"`
		SkinChangeType    string `toml:"skin_change_type"`
		Skin              string `toml:"skin"`
	} `toml:"skin"`
	Announce struct {
		McsnipergoAnnounceCode   string `toml:"mcsnipergo_announce_code"`
		WebhookURL               string `toml:"webhook_url"`
		WebhookFormat            string `toml:"webhook_format"`
		PrivateWebhookURL        string `toml:"private_webhook_url"`
		PrivateWebhookFormat     string `toml:"private_webhook_format"`
		PrivateWebhookIncludeAcc bool   `toml:"private_webhook_include_acc"`
	} `toml:"announce"`
}

const (
	defaultConfigString = `[sniper]
timing_system_preference = "star.shopping"
cycle_timing_systems = true # Go through each timing system until droptime is successfully grabbed
auto_claim_namemc = false

snipe_requests = 2 # requests to be sent per acc for normal sniping
prename_requests = 6 # requests to be sent per acc for prename sniping

spread = 5 # delay between requests in milliseconds

[accounts]
max_accounts = 1
prename_max_accounts = 10

start_auth = 720 # start auth 720 minutes before drop
auth_delay = 1 # time between acc auth

[skin]
change_skin_on_snipe = false
skin_change_type = "url"
skin = "" # this value depends on the skin_change_type value

[announce]
mcsnipergo_announce_code = "" # leave blank to not announce snipe


# discord webhook-related things
webhook_url = "" # public webhook url to announce snipe
webhook_format = "sniped {name} with {searches} searches using MCsniperGO!"

private_webhook_url = "" # webhook url where PRIVATE information will be sent, DO NOT MAKE THIS IN A PUBLIC DISCORD CHANNEL.
private_webhook_format = "sniped {name} with {searches} searches and {offset} offset!"
private_webhook_include_acc = false # include acc details in webhook`
)

func defaultConfig() error {

	cfgBytes := []byte(defaultConfigString)
	file, err := os.Create("config.toml")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(cfgBytes)
	if err != nil {
		return err
	}
	return nil

}

func getConfig() (ConfigStruct, error) {
	if !fileExists("config.toml") {
		return ConfigStruct{}, errors.New("config file does not exist")
	}

	confBytes, err := ioutil.ReadFile("config.toml")

	if err != nil {
		return ConfigStruct{}, err
	}

	var cfg ConfigStruct
	err = toml.Unmarshal(confBytes, &cfg)

	if err != nil {
		return ConfigStruct{}, err
	}

	return cfg, nil
}
