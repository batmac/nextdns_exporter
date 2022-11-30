# configuration

```
export NEXTDNS_API_KEY=...
```

# help

```
Usage of ./nextdns_exporter:
  -port string
        Port to listen on (default "8080")
  -profiles string
        Comma-separated list of profile IDs to monitor, empty = all
```

# metrics example

```
# HELP nextdns_protocols_count Number of queries, per protocols
# TYPE nextdns_protocols_count counter
nextdns_protocols_count{profile="374fe1",protocol="DNS-over-HTTPS"} 818182
nextdns_protocols_count{profile="374fe1",protocol="TCP"} 4
nextdns_protocols_count{profile="374fe1",protocol="UDP"} 250161
nextdns_protocols_count{profile="aaf63d",protocol="DNS-over-HTTPS"} 63589
# HELP nextdns_status_count Number of queries, per type
# TYPE nextdns_status_count counter
nextdns_status_count{profile="374fe1",type="allowed"} 1416
nextdns_status_count{profile="374fe1",type="blocked"} 124259
nextdns_status_count{profile="374fe1",type="default"} 940866
nextdns_status_count{profile="374fe1",type="relayed"} 1806
nextdns_status_count{profile="aaf63d",type="blocked"} 18
nextdns_status_count{profile="aaf63d",type="default"} 63571
# HELP nextdns_dnssec_count Number of queries, per validated
# TYPE nextdns_dnssec_count counter
nextdns_dnssec_count{profile="374fe1",validated="false"} 854862
nextdns_dnssec_count{profile="374fe1",validated="true"} 87420
nextdns_dnssec_count{profile="aaf63d",validated="false"} 62384
nextdns_dnssec_count{profile="aaf63d",validated="true"} 1187

```
