package main

import (
	"context"
	"flag"
	"fmt"
	"stock/cmd/cronjob"
	"stock/config"
	"time"
)

var (
	loc, _      = time.LoadLocation("Asia/Jakarta")
	currentTime = time.Now().In(loc)
	currentDate = currentTime.Format(time.DateOnly)
)

func main() {
	var (
		serverType string
		configPath string
		fromDate   string
		toDate     string
		symbols    string
	)
	flag.StringVar(
		&serverType,
		"type",
		"api",
		"provide server to run",
	)
	flag.StringVar(
		&configPath,
		"config",
		"config.yaml",
		"path to config file",
	)
	flag.StringVar(
		&fromDate,
		"from-date",
		currentDate,
		"from date",
	)
	flag.StringVar(
		&toDate,
		"to-date",
		currentDate,
		"to date",
	)
	flag.StringVar(
		&symbols,
		"symbols",
		"",
		"symbols",
	)
	flag.Parse()

	// Initialize config
	if err := config.Init(configPath); err != nil {
		panic(fmt.Sprintf("failed to initialize config: %v", err))
	}

	var (
		shutdown func(context.Context) error
		err      error
	)

	defer func() {
		_ = shutdown(context.Background())
	}()

	// Choose runner based on type.
	var runner Runner
	switch serverType {
	// case "cronjob-upsert-emitten-profile":
	// 	runner = cronjob.NewUpsertEmittenProfile()
	case "cronjob-upsert-broker-summary":
		runner = cronjob.NewUpsertBrokerSummary(fromDate, toDate, symbols)
	case "cronjob-upsert-price-feed":
		runner = cronjob.NewUpsertPriceFeed(fromDate, toDate, symbols)
	// case "cronjob-upsert-subsector":
	// 	runner = cronjob.NewUpsertSubsector()
	// case "cronjob-upsert-emitten-profile-info":
	// 	runner = cronjob.NewUpsertEmittenProfileInfo()
	// case "cronjob-update-emitten-underwriter-code":
	// 	runner = cronjob.NewUpdateEmittenUnderwriterCode()
	default:
		panic("invalid server type")
	}

	if err = RunApp(runner); err != nil {
		panic(fmt.Sprintf("cannot run app %s", serverType))
	}
}
