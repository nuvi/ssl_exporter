### SSL Exporter

:9197/metrics example...


```# HELP ssl_expiration SSL certificate name and Expiration
# TYPE ssl_expiration gauge
ssl_expiration{domain="https://google.com"} 1.53064278e+09
# HELP ssl_scrape_up Can we scrape the site
# TYPE ssl_scrape_up gauge
ssl_scrape_up{domain="https://asdfasdfasdfassdfasdf.com"} 0
ssl_scrape_up{domain="https://google.com"} 1```