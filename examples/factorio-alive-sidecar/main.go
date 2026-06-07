package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/baelish/alive/api"
	"github.com/baelish/alive/client"
	"github.com/gorcon/rcon"
	goflags "github.com/jessevdk/go-flags"
	"go.uber.org/zap"
)

type options struct {
	AliveURL           string `long:"alive-api-url" env:"ALIVE_API_URL" required:"true" description:"URL of the alive API"`
	BoxID              string `long:"box-id" env:"ALIVE_BOX_ID" required:"true" description:"Box ID to update"`
	BoxName            string `long:"box-name" env:"ALIVE_BOX_NAME" description:"Display name for the box (defaults to box-id)"`
	BoxSize            string `long:"box-size" env:"ALIVE_BOX_SIZE" default:"small" description:"Box size: e.g. small, dsmall, dot, medium"`
	MaxTBU             string `long:"max-tbu" env:"ALIVE_MAX_TBU" description:"Maximum time between updates (e.g. 5m)"`
	RconURL            string `long:"rcon-url" env:"RCON_URL" default:"127.0.0.1:27015" description:"Factorio rcon endpoint to scrape"`
	RconPass           string `long:"rcon-pass" env:"RCON_PASS" description:"Factorio rcon password"`
	RconPassFile       string `long:"rcon-pass-file" env:"RCON_PASS_FILE" default:"/factorio/config/rconpw" description:"location of the rcon password file"`
	ServerSettingsFile string `long:"server-settings-file" env:"SERVER_SETTINGS_FILE" default:"/factorio/config/server-settings.json" description:"location of the server settings file"`
	Interval           string `long:"interval" env:"CHECK_INTERVAL" default:"30s" description:"How often to scrape and update"`
	Debug              bool   `long:"debug" env:"DEBUG" description:"Enable debug logging"`
}

var (
	logger *zap.Logger
)

type serverSettings struct {
	MaxPlayers int `json:"max_players"`
}

func main() {
	var opts options
	parser := goflags.NewParser(&opts, goflags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*goflags.Error); ok && flagsErr.Type == goflags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	logger = zap.Must(zap.NewProduction())
	if opts.Debug {
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger = zap.Must(cfg.Build())
	}
	defer logger.Sync()

	interval, err := time.ParseDuration(opts.Interval)
	if err != nil {
		logger.Fatal("invalid interval", zap.String("interval", opts.Interval), zap.Error(err))
	}

	var boxSize api.BoxSize
	if err := json.Unmarshal([]byte(`"`+opts.BoxSize+`"`), &boxSize); err != nil {
		logger.Warn("invalid box size, defaulting to medium", zap.String("size", opts.BoxSize))
		boxSize = api.Medium
	}

	var maxTBU *api.Duration
	if opts.MaxTBU != "" {
		d, err := time.ParseDuration(opts.MaxTBU)
		if err != nil {
			logger.Fatal("invalid max-tbu", zap.String("value", opts.MaxTBU), zap.Error(err))
		}
		dur := api.Duration(d)
		maxTBU = &dur
	}

	boxName := opts.BoxName
	if boxName == "" {
		boxName = opts.BoxID
	}

	c := client.NewClient(opts.AliveURL)
	ensureBox(c, opts.BoxID, boxName, boxSize, maxTBU)

	rconPass := opts.RconPass
	if rconPass == "" {
		// get from file
		data, err := os.ReadFile(opts.RconPassFile)
		if err != nil {
			logger.Fatal("failed to read rcon password file", zap.String("path", opts.RconPassFile), zap.Error(err))
		}
		rconPass = strings.TrimSpace(string(data))
	}

	if rconPass == "" {
		logger.Fatal("rcon password must be set via --rcon-pass or readable from --rcon-pass-file")
	}

	var settings serverSettings
	jsonData, err := os.ReadFile(opts.ServerSettingsFile)
	if err != nil {
		logger.Error("unable to read server settings file", zap.String("path", opts.ServerSettingsFile), zap.Error(err))
	} else {
		if err := json.Unmarshal(jsonData, &settings); err != nil {
			logger.Error("unable to unmarshal data into server settings", zap.Error(err))
		}
	}

	logger.Info("starting sidecar", zap.String("rcon-url", opts.RconURL), zap.Duration("interval", interval))
	for {
		checkAndUpdate(c, opts.BoxID, opts.RconURL, rconPass, settings)
		time.Sleep(interval)
	}
}

func ensureBox(c *client.Client, id, name string, size api.BoxSize, maxTBU *api.Duration) {
	if _, err := c.GetBox(id); err == nil {
		logger.Info("box already exists", zap.String("id", id))
		return
	}

	box := api.Box{
		ID:          id,
		Name:        name,
		Description: fmt.Sprintf("Status of the %s Factorio server", name),
		Size:        size,
		MaxTBU:      maxTBU,
	}

	if _, err := c.CreateBox(box); err != nil {
		logger.Fatal("failed to create box", zap.String("id", id), zap.Error(err))
	}
	logger.Info("created box", zap.String("id", id))
}

func checkAndUpdate(c *client.Client, boxID, rconURL string, rconPass string, settings serverSettings) {
	status, message := getStatus(rconURL, rconPass, settings.MaxPlayers)
	if err := c.CreateEvent(api.Event{ID: boxID, Status: status, Message: message}); err != nil {
		logger.Error("failed to post event", zap.Error(err))
	}
}

func getStatus(rconURL, rconPass string, maxPlayers int) (api.Status, string) {
	conn, err := rcon.Dial(rconURL, rconPass)
	if err != nil {
		return api.Red, fmt.Sprintf("connect to rcon failed: %v", err)
	}
	defer conn.Close()

	response, err := conn.Execute("/c rcon.print(#game.connected_players)")
	if err != nil {
		return api.Red, fmt.Sprintf("failed to get player count: %v", err)
	}

	playerCount, err := strconv.Atoi(strings.TrimSpace(response))
	if err != nil {
		return api.Amber, fmt.Sprintf("invalid player count:%s (%v)", response, err)
	}

	if maxPlayers > 0 {
		return api.Green, fmt.Sprintf("%d/%d Online", playerCount, maxPlayers)
	}

	return api.Green, fmt.Sprintf("%d Online", playerCount)
}
