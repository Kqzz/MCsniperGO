package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type ConfigStruct struct {
	Sniper struct {
		TimingSystemPreference string `toml:"timing_system_preference"`
		CycleTimingSystems     bool   `toml:"cycle_timing_systems"`
		AutoClaimNamemc        bool   `toml:"auto_claim_namemc"`
		SnipeRequests          int    `toml:"snipe_requests"`
		PrenameRequests        int    `toml:"prename_requests"`
	} `toml:"sniper"`
	Accounts struct {
		MaxAccounts        int `toml:"max_accounts"`
		PrenameMaxAccounts int `toml:"prename_max_accounts"`
		StartAuth          int `toml:"start_auth"`
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

func defaultConfig() error {

	config := ConfigStruct{
		Sniper: struct {
			TimingSystemPreference string "toml:\"timing_system_preference\""
			CycleTimingSystems     bool   "toml:\"cycle_timing_systems\""
			AutoClaimNamemc        bool   "toml:\"auto_claim_namemc\""
			SnipeRequests          int    "toml:\"snipe_requests\""
			PrenameRequests        int    "toml:\"prename_requests\""
		}{
			TimingSystemPreference: "ckm",
			CycleTimingSystems:     true,
			AutoClaimNamemc:        false,
			SnipeRequests:          2,
			PrenameRequests:        6,
		},
		Accounts: struct {
			MaxAccounts        int "toml:\"max_accounts\""
			PrenameMaxAccounts int "toml:\"prename_max_accounts\""
			StartAuth          int "toml:\"start_auth\""
		}{
			MaxAccounts:        1,
			PrenameMaxAccounts: 10,
			StartAuth:          720,
		},
		Skin: struct {
			ChangeSkinOnSnipe bool   "toml:\"change_skin_on_snipe\""
			SkinChangeType    string "toml:\"skin_change_type\""
			Skin              string "toml:\"skin\""
		}{
			ChangeSkinOnSnipe: false,
			SkinChangeType:    "url",
			Skin:              "",
		},
		Announce: struct {
			McsnipergoAnnounceCode   string "toml:\"mcsnipergo_announce_code\""
			WebhookURL               string "toml:\"webhook_url\""
			WebhookFormat            string "toml:\"webhook_format\""
			PrivateWebhookURL        string "toml:\"private_webhook_url\""
			PrivateWebhookFormat     string "toml:\"private_webhook_format\""
			PrivateWebhookIncludeAcc bool   "toml:\"private_webhook_include_acc\""
		}{
			McsnipergoAnnounceCode:   "",
			WebhookURL:               "",
			WebhookFormat:            "",
			PrivateWebhookURL:        "",
			PrivateWebhookFormat:     "sniped {name} with {searches} searches and {offset} offset",
			PrivateWebhookIncludeAcc: false,
		},
	}

	cfgBytes, err := toml.Marshal(config)
	if err != nil {
		return err
	}
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
