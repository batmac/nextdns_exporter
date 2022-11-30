package collectors

import (
	"github.com/batmac/nextdns_exporter/nextdns"
)

func NewDnssec(client *nextdns.Client, profileId string) *Collector {
	return NewCollector(
		client,
		profileId,
		&CollectorDefine{
			Name:     "dnssec",
			Label:    "validated",
			Endpoint: "analytics/dnssec",
			Parse:    ParseStd,
		},
	)
}
