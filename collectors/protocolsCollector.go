package collectors

import (
	"encoding/json"
	"log"

	"nextdns_exporter/nextdns"

	"github.com/prometheus/client_golang/prometheus"
)

type ProtocolsCollector struct {
	profile string
	c       *nextdns.Client

	Protocols     map[string]int
	protocolsDesc *prometheus.Desc
}

func (pc ProtocolsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(pc, ch)
}

func (pc ProtocolsCollector) Collect(ch chan<- prometheus.Metric) {
	err := pc.Fetch()
	if err != nil {
		return
	}
	for tpe, count := range pc.Protocols {
		ch <- prometheus.MustNewConstMetric(
			pc.protocolsDesc,
			prometheus.CounterValue,
			float64(count),
			tpe,
		)
	}
}

func NewProtocolsCollector(client *nextdns.Client, profileId string) *ProtocolsCollector {
	pc := &ProtocolsCollector{
		profile: profileId,
		c:       client,
		protocolsDesc: prometheus.NewDesc(
			"nextdns_protocols_count",
			"Number of queries per protocol",
			[]string{"type"}, prometheus.Labels{"profile": profileId},
		),
	}
	prometheus.MustRegister(pc)
	return pc
}

func (pc *ProtocolsCollector) Fetch() error {
	type Data struct {
		Protocol string `json:"protocol"`
		Queries  int    `json:"queries"`
	}
	type Response struct {
		Data []Data `json:"data"`
	}

	protocols := make(map[string]int)
	content := Response{}

	r := pc.c.MustGet("profiles/" + pc.profile + "/analytics/protocols")
	// log.Printf("content: %s", r)
	err := json.Unmarshal(r, &content)
	if err != nil {
		log.Printf("Error unmarshalling: %v", err)
		return err
	}
	// log.Printf("content: %v", content)

	for _, d := range content.Data {
		protocols[d.Protocol] = d.Queries
	}
	pc.Protocols = protocols
	return nil
}
