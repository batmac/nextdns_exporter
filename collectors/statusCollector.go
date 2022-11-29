package collectors

import (
	"encoding/json"
	"log"

	"github.com/batmac/nextdns_exporter/nextdns"
	"github.com/prometheus/client_golang/prometheus"
)

type StatusCollector struct {
	profile string
	c       *nextdns.Client

	Queries    map[string]int
	statusDesc *prometheus.Desc
}

func (sc StatusCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(sc, ch)
}

func (sc StatusCollector) Collect(ch chan<- prometheus.Metric) {
	err := sc.Fetch()
	if err != nil {
		return
	}
	for tpe, count := range sc.Queries {
		ch <- prometheus.MustNewConstMetric(
			sc.statusDesc,
			prometheus.CounterValue,
			float64(count),
			tpe,
		)
	}
}

func NewStatusCollector(client *nextdns.Client, profileId string) *StatusCollector {
	sc := &StatusCollector{
		profile: profileId,
		c:       client,
		statusDesc: prometheus.NewDesc(
			"nextdns_queries_count",
			"Number of queries.",
			[]string{"type"}, prometheus.Labels{"profile": profileId},
		),
	}
	prometheus.MustRegister(sc)
	return sc
}

func (sc *StatusCollector) Fetch() error {
	type Data struct {
		Status  string `json:"status"`
		Queries int    `json:"queries"`
	}
	type Response struct {
		Data []Data `json:"data"`
	}

	queries := make(map[string]int)
	content := Response{}

	err := json.Unmarshal(sc.c.MustGet("profiles/"+sc.profile+"/analytics/status"), &content)
	if err != nil {
		log.Printf("Error unmarshalling: %v", err)
		return err
	}
	// log.Printf("content: %v", content)

	for _, d := range content.Data {
		queries[d.Status] = d.Queries
	}
	sc.Queries = queries
	return nil
}
