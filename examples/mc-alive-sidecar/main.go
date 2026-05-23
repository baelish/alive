package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/baelish/alive/api"
	"github.com/baelish/alive/client"
	goflags "github.com/jessevdk/go-flags"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"go.uber.org/zap"
)

const (
	statusMetric    string = "minecraft_status_healthy"
	onlineMetric    string = "minecraft_status_players_online_count"
	maxOnlineMetric string = "minecraft_status_players_max_count"
)

type options struct {
	AliveURL   string `long:"alive-api-url" env:"ALIVE_API_URL" required:"true" description:"URL of the alive API"`
	BoxID      string `long:"box-id" env:"ALIVE_BOX_ID" required:"true" description:"Box ID to update"`
	BoxName    string `long:"box-name" env:"ALIVE_BOX_NAME" description:"Display name for the box (defaults to box-id)"`
	BoxSize    string `long:"box-size" env:"ALIVE_BOX_SIZE" default:"small" description:"Box size: e.g. small, dsmall, dot, medium"`
	MaxTBU     string `long:"max-tbu" env:"ALIVE_MAX_TBU" description:"Maximum time between updates (e.g. 5m)"`
	MetricsURL string `long:"metrics-url" env:"METRICS_URL" default:"http://127.0.0.1:9090/metrics" description:"Prometheus metrics endpoint to scrape"`
	Interval   string `long:"interval" env:"CHECK_INTERVAL" default:"30s" description:"How often to scrape and update"`
	Debug      bool   `long:"debug" env:"DEBUG" description:"Enable debug logging"`
}

var (
	logger     *zap.Logger
	httpClient = &http.Client{Timeout: 10 * time.Second}
)

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

	logger.Info("starting sidecar", zap.String("metrics-url", opts.MetricsURL), zap.Duration("interval", interval))
	for {
		checkAndUpdate(c, opts.BoxID, opts.MetricsURL)
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
		Description: fmt.Sprintf("Status of the %s minecraft server", name),
		Size:        size,
		MaxTBU:      maxTBU,
	}

	if _, err := c.CreateBox(box); err != nil {
		logger.Fatal("failed to create box", zap.String("id", id), zap.Error(err))
	}
	logger.Info("created box", zap.String("id", id))
}

func checkAndUpdate(c *client.Client, boxID, metricsURL string) {
	metrics, scrapeErr := scrapeMetrics(metricsURL)

	var status api.Status
	var message string

	if scrapeErr != nil {
		status = api.Red
		message = fmt.Sprintf("scrape failed: %v", scrapeErr)
	} else {
		val, found := metrics[statusMetric]
		switch {
		case !found:
			status = api.Red
			message = fmt.Sprintf("metric %q not found", statusMetric)
		case val >= 1:
			status = api.Green
			online, onlineFound := metrics[onlineMetric]
			if !onlineFound {
				logger.Warn("metric not found", zap.String("onlineMetric", onlineMetric))
			}
			maxOnline, maxFound := metrics[maxOnlineMetric]
			if !maxFound {
				logger.Warn("metric not found", zap.String("maxOnlineMetric", maxOnlineMetric))
			}
			message = fmt.Sprintf("%d/%d Online", int(online), int(maxOnline))
		default:
			status = api.Red
			message = "Unhealthy"
		}
	}

	if err := c.CreateEvent(api.Event{ID: boxID, Status: status, Message: message}); err != nil {
		logger.Error("failed to post event", zap.Error(err))
	}
}

func scrapeMetrics(url string) (map[string]float64, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metrics endpoint returned %d", resp.StatusCode)
	}

	decoder := expfmt.NewDecoder(resp.Body, expfmt.ResponseFormat(resp.Header))
	metrics := make(map[string]float64)
	for {
		var mf dto.MetricFamily
		if err := decoder.Decode(&mf); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		for _, m := range mf.GetMetric() {
			var val float64
			switch mf.GetType() {
			case dto.MetricType_GAUGE:
				val = m.GetGauge().GetValue()
			case dto.MetricType_COUNTER:
				val = m.GetCounter().GetValue()
			default:
				val = m.GetUntyped().GetValue()
			}
			metrics[mf.GetName()] = val
		}
	}
	return metrics, nil
}
