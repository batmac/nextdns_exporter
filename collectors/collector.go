package collectors

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/batmac/nextdns_exporter/nextdns"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	labelValue string
	count      int
	ParseFunc  func([]byte) (map[labelValue]count, error)
)

type CollectorDefine struct {
	Name, Label string // Name of the metric (prefixed and suffixed), discriminig label of the metric
	Endpoint    string // endpoint to fetch data from, without the profile part, starting after the /
	Parse       ParseFunc
}

type Collector struct {
	client   *nextdns.Client
	endpoint string // endpoint to fetch data from, with the profile part, starting after the initial /
	desc     *prometheus.Desc
	fn       ParseFunc

	Data map[labelValue]count // last fetched data
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	err := c.parseWrapper()
	if err != nil {
		return
	}
	for labelValue, count := range c.Data {
		ch <- prometheus.MustNewConstMetric(
			c.desc,
			prometheus.CounterValue,
			float64(count),
			string(labelValue),
		)
	}
}

func NewCollector(client *nextdns.Client, profileId string, cd *CollectorDefine) *Collector {
	if cd.Parse == nil {
		cd.Parse = ParseStd
	}
	c := &Collector{
		client:   client,
		endpoint: "profiles/" + profileId + "/" + cd.Endpoint,
		desc: prometheus.NewDesc(
			prometheus.BuildFQName("nextdns", cd.Name, "count"),
			"Number of queries, per "+cd.Label,
			[]string{cd.Label},
			prometheus.Labels{"profile": profileId},
		),
		fn:   cd.Parse,
		Data: map[labelValue]count{},
	}
	prometheus.MustRegister(c)
	return c
}

func (c *Collector) parseWrapper() error {
	r := c.client.MustGet(c.endpoint)
	// log.Printf("content: %s", r)
	data, err := c.fn(r)
	if err != nil {
		return err
	}
	// log.Printf("content: %v", content)
	c.Data = data
	return nil
}

func ParseStd(content []byte) (map[labelValue]count, error) {
	response := struct {
		Data []map[string]any `json:"data"`
	}{}
	parsedData := make(map[labelValue]count)
	err := json.Unmarshal(content, &response)
	if err != nil {
		log.Printf("Error unmarshalling: %v", err)
		return nil, err
	}
	// log.Printf("response: %v", response)

	for _, elem := range response.Data {
		var queriesCount int
		switch elem["queries"].(type) {
		case float64:
			queriesCount = int(elem["queries"].(float64))
		case int:
			queriesCount = elem["queries"].(int)
		}

		for _, value := range elem {
			switch value.(type) {
			case string:
				parsedData[labelValue(value.(string))] = count(queriesCount)
				break
			case bool:
				parsedData[labelValue(strconv.FormatBool(value.(bool)))] = count(queriesCount)
			default:
				continue
			}
		}
	}
	// log.Printf("data: %v", parsedData)
	return parsedData, nil
}
