package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	AppName  string // APP_NAME
	LogLevel string // LOG_LEVEL

	FileSavePath string // FILE_SAVE_PATH

	ParserWorkerAmount int // PARSER_WORKER_AMOUNT

	// several conditions can be specified separated by comma (example: "new, used")
	// if empty string then no conditions are applied, by default env is empty string
	ConditionsToParse []string // CONDITIONS_TO_PARSE
	ParsingStartURL   string   // PARSING_START_URL
}

// Load func loads config from .env file located by envFilePath and initialize Config struct
// If env variable is not present, fallback will be used
func Load(log *log.Logger, envFilePath string) *Config {
	var config Config
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Infof("failed to load config file: %v; applying default values", err)
	}
	config.AppName = GetString("APP_NAME", "crawler")
	config.LogLevel = GetString("LOG_LEVEL", "debug")
	config.FileSavePath = GetString("FILE_SAVE_PATH", "data")
	config.ParserWorkerAmount = GetInt("PARSER_WORKER_AMOUNT", 5)
	config.FileSavePath = GetString("FILE_SAVE_PATH", "data")
	conditionsToParse := GetString("CONDITIONS_TO_PARSE", "")
	conditionsToParse = strings.ReplaceAll(conditionsToParse, "", "")
	if conditionsToParse == "" {
		config.ConditionsToParse = []string{}
	} else {
		config.ConditionsToParse = strings.Split(conditionsToParse, ",")
	}
	config.ParsingStartURL = GetString(
		"PARSING_START_URL",
		"https://kz.ebay.com/b/XFX-Computer-Graphics-Cards/27386/bn_2774033?Memory%2520Size=8%2520GB%7C24%2520GB&rt=nc&Connectors=HDMI%7CDisplayPort&Memory%2520Type=HBM%7CHBM2%7CSDRAM%7CDDR4%7CDDR5&mag=1",
	)

	return &config

}

// GetString func returns environment variable value as a string value,
// If variable doesn't exist or is not set, returns fallback value
func GetString(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

// GetInt func returns environment variable value as a integer value,
// If variable doesn't exist or is not set, returns fallback value
func GetInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	res, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return fallback
	}
	return int(res)
}
