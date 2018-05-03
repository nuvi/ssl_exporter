### SSL Exporter

Will get expiration of SSL certificates in Epoch time to be scraped by Prometheus, configured with a yaml file, see configs directory for example.

hit :9197/probe to fire off metric gathering, control of this can be done with scrape_interval on Prometheus server to prevent unnecessarily hitting monitored websites


:9197/metrics example output

```# HELP ssl_expiration SSL certificate name and Expiration
# TYPE ssl_expiration gauge
ssl_expiration{domain="https://google.com"} 1.53064278e+09
# HELP ssl_scrape_up Can we scrape the site
# TYPE ssl_scrape_up gauge
ssl_scrape_up{domain="https://asdfasdfasdfassdfasdf.com"} 0
ssl_scrape_up{domain="https://google.com"} 1```

####Environment Variables
```
	Debug         default false
	ListenAddress default :9197
	MetricsPath   default /metrics
	ProbePath     default /probe
	ConfigPath    default /etc/prometheus/exporters/ssl_exporter/
```