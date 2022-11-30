package collectors

import (
	"github.com/batmac/nextdns_exporter/nextdns"
)

func NewProtocols(client *nextdns.Client, profileId string) *Collector {
	return NewCollector(
		client,
		profileId,
		&CollectorDefine{
			Name:     "protocols",
			Label:    "protocol",
			Endpoint: "analytics/protocols",
			Parse:    ParseStd,
		},
	)
}
