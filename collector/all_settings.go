package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	// "github.com/imdario/mergo"
	"github.com/prometheus/client_golang/prometheus"
)

// AllSettings information struct
type AllSettings struct {
	logger log.Logger
	client *http.Client
	url    *url.URL

	up                              prometheus.Gauge
	readOnlyIndices                 prometheus.Gauge
	totalScrapes, jsonParseFailures prometheus.Counter
}

// NewAllSettings defines Cluster Settings Prometheus metrics
func NewAllSettings(logger log.Logger, client *http.Client, url *url.URL) *AllSettings {
	return &AllSettings{
		logger: logger,
		client: client,
		url:    url,

		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, "allsettings_stats", "up"),
			Help: "Was the last scrape of the ElasticSearch all settings endpoint successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, "allsettings_stats", "total_scrapes"),
			Help: "Current total ElasticSearch all settings scrapes.",
		}),
		readOnlyIndices: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(namespace, "allsettings_stats", "read_only_indices"),
			Help: "Current number of read only indices within cluster",
		}),
		jsonParseFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, "allsettings_stats", "json_parse_failures"),
			Help: "Number of errors while parsing JSON.",
		}),
	}
}

// Describe add Snapshots metrics descriptions
func (cs *AllSettings) Describe(ch chan<- *prometheus.Desc) {
	ch <- cs.up.Desc()
	ch <- cs.totalScrapes.Desc()
	ch <- cs.readOnlyIndices.Desc()
	ch <- cs.jsonParseFailures.Desc()
}

func (cs *AllSettings) getAndParseURL(u *url.URL, data interface{}) error {
	res, err := cs.client.Get(u.String())
	if err != nil {
		return fmt.Errorf("failed to get from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}

	defer func() {
		err = res.Body.Close()
		if err != nil {
			_ = level.Warn(cs.logger).Log(
				"msg", "failed to close http.Client",
				"err", err,
			)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Request failed with code %d", res.StatusCode)
	}

	// mp := map[string]interfcs{}type
	//var mp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(data); err != nil {
		cs.jsonParseFailures.Inc()
		return err
	}
	//fmt.Println(mp)
	fmt.Println(data)
	return nil
}

func (cs *AllSettings) fetchAndDecodeAllSettingsStats() (AllSettingsResponse, error) {

	u := *cs.url
	u.Path = path.Join(u.Path, "/_all/_settings")
	// var asfr AllSettingsFullResponse
	var asr AllSettingsResponse
	// err := cs.getAndParseURL(&u, &asfr)
	err := cs.getAndParseURL(&u, &asr)
	if err != nil {
		return asr, err
	}

	// err = mergo.Merge(&asr, asfr, mergo.WithOverride)
	// if err != nil {
	// 	return asr, err
	// }
	return asr, err
}

// Collect gets cluster settings  metric values
func (cs *AllSettings) Collect(ch chan<- prometheus.Metric) {

	cs.totalScrapes.Inc()
	defer func() {
		ch <- cs.up
		ch <- cs.totalScrapes
		ch <- cs.jsonParseFailures
		ch <- cs.readOnlyIndices
	}()

	asr, err := cs.fetchAndDecodeAllSettingsStats()
	if err != nil {
		cs.readOnlyIndices.Set(0)
		cs.up.Set(0)
		_ = level.Warn(cs.logger).Log(
			"msg", "failed to fetch and decode cluster settings stats",
			"err", err,
		)
		return
	}
	cs.up.Set(1)

	// shardAllocationMap := map[string]int{
	// 	"all":           0,
	// 	"primaries":     1,
	// 	"new_primaries": 2,
	// 	"none":          3,
	// }

	// cs.readOnlyIndices.Set(float64(shardAllocationMap[asr.Cluster.Routing.Allocation.Enabled]))
	var c int = 0
	res2B, _ := json.Marshal(asr)
	fmt.Println(string(res2B))
	// _ = level.Info(cs.logger).Log(asr.Indexes, len(asr.Indexes))

	for key, value := range asr {
		// The read_only_allow_delete is string and not boolean
		_ = level.Info(cs.logger).Log(
			"msg", "Readonly",
			"k", key,
			"value", value.Settings.Index.Blocks.ReadOnly,
		)
		if value.Settings.Index.Blocks.ReadOnly == "true" {
			c++
		}
	}
	cs.readOnlyIndices.Set(float64(c))
}
