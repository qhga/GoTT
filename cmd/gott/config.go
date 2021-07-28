package main

import (
	"encoding/json"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	WebPort            string                          `json:"web_port"`
	WebDomain          string                          `json:"web_domain"`
	WebAdmin           string                          `json:"web_admin"`
	WebAdminPass       string                          `json:"web_admin_pass"`
	DbUser             string                          `json:"db_user"`
	DbPass             string                          `json:"db_pass"`
	DbHost             string                          `json:"db_host"`
	DbPort             string                          `json:"db_port"`
	MailUser           string                          `json:"mail_user"`
	MailPass           string                          `json:"mail_pass"`
	MailSmtpServer     string                          `json:"mail_smtp_server"`
	MailSmtpPort       string                          `json:"mail_smtp_port"`
	TestTime           int                             `json:"test_time"`
	TestAmount         int                             `json:"test_amount"`
	TestMaxWpm         int                             `json:"test_max_wpm"`
	TestShowTimer      bool                            `json:"test_show_timer"`
	TestShowResult     bool                            `json:"test_show_result"`
	Surveys            map[string][]primitive.ObjectID `json:"surveys"`
	AvailableKeyboards []Keyboard                      `json:"available_keyboards"`
}

var config Config

func init() {
	f, err := os.Open("../../configs/config.json")
	if err != nil {
		log.Fatalln(err, "Could not read config file")
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		log.Fatalln(err, "Could not decode config file")
	}
}
