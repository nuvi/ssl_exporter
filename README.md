## Prometheus SSL Exporter

Will get expiration of SSL certificates in Epoch time to be scraped by Prometheus, configured with a yaml file, see configs directory for example.

####:9197/metrics example output

```# HELP ssl_expiration SSL certificate name and Expiration
# TYPE ssl_expiration gauge
ssl_expiration{domain="https://google.com"} 1.53064278e+09
# HELP ssl_scrape_up Can we scrape the site
# TYPE ssl_scrape_up gauge
ssl_scrape_up{domain="https://asdfasdfasdfassdfasdf.com"} 0
ssl_scrape_up{domain="https://google.com"} 1
```

####Environment Variables
```
	Debug         default false
	ListenAddress default :9197
	MetricsPath   default /metrics
	ConfigPath    default /etc/prometheus/exporters/ssl_exporter/

```

#### Prometheus Alert rule, if cert is set to expire in 30 days
```
ssl_expiration - time() < 86400 * 30
```