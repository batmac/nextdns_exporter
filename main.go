package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/batmac/nextdns_exporter/collectors"
	"github.com/batmac/nextdns_exporter/nextdns"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	argPort     = flag.String("port", "8080", "Port to listen on")
	argProfiles = flag.String("profiles", "", "Comma-separated list of profile IDs to monitor, empty = all")
)

func main() {
	flag.Parse()
	apiKey := os.Getenv("NEXTDNS_API_KEY")
	c := nextdns.NewClient(apiKey)

	var profileIds []string
	if (*argProfiles) == "" {
		var err error
		profileIds, err = getProfileIDs(c)
		if err != nil {
			log.Printf("Error getting profiles list: %v", err)
			return
		}
		log.Printf("Discovered Profiles: %v", profileIds)
	} else {
		profileIds = strings.Split(*argProfiles, ",")
	}

	for _, id := range profileIds {
		_ = collectors.NewStatus(c, id)
		_ = collectors.NewProtocols(c, id)
		_ = collectors.NewDnssec(c, id)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Listening on port %v\n", *argPort)
	log.Fatal(http.ListenAndServe(":"+*argPort, nil))

	/*
		 	for _, id := range ids {
				fmt.Printf("getting profile %v\n", id)
				contentStatus := mustGetByProfile(c, id, "analytics/status")
				fmt.Printf("contentStatus: %v\n", contentStatus)

				contentDomains := mustGetByProfile(c, id, "analytics/domains?root=true")
				fmt.Printf("contentDomains: %v\n", contentDomains)

				contentReasons := mustGetByProfile(c, id, "analytics/reasons")
				fmt.Printf("contentReasons: %v\n", contentReasons)

				contentDevices := mustGetByProfile(c, id, "analytics/devices")
				fmt.Printf("contentDevices: %v\n", contentDevices)

				contentProtocols := mustGetByProfile(c, id, "analytics/protocols")
				fmt.Printf("contentProtocols: %v\n", contentProtocols)

				contentQueryTypes := mustGetByProfile(c, id, "analytics/queryTypes")
				fmt.Printf("contentQueryTypes: %v\n", contentQueryTypes)

				contentIpVersions := mustGetByProfile(c, id, "analytics/ipVersions")
				fmt.Printf("contentIpVersions: %v\n", contentIpVersions)

				contentDnssec := mustGetByProfile(c, id, "analytics/dnssec")
				fmt.Printf("contentDnssec: %v\n", contentDnssec)

				contentEncryption := mustGetByProfile(c, id, "analytics/encryption")
				fmt.Printf("contentEncryption: %v\n", contentEncryption)

				contentCountries := mustGetByProfile(c, id, "analytics/destinations?type=countries")
				fmt.Printf("contentCountries: %v\n", contentCountries)

				contentGafam := mustGetByProfile(c, id, "analytics/destinations?type=gafam")
				fmt.Printf("contentGafam: %v\n", contentGafam)

			}
	*/
}
