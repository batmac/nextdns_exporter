package collectors

import (
	"github.com/batmac/nextdns_exporter/nextdns"
)

func NewStatus(client *nextdns.Client, profileId string) *Collector {
	return NewCollector(
		client,
		profileId,
		&CollectorDefine{
			Name:     "status",
			Label:    "type",
			Endpoint: "analytics/status",
			Parse:    ParseStd,
		},
	)
}
