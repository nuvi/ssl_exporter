package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/common/log"

	"github.com/ghodss/yaml"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	sslMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ssl_expiration",
		Help: "SSL certificate name and Expiration",
	}, []string{"domain"})

	up = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ssl_scrape_up",
		Help: "Can we scrape the site",
	}, []string{"domain"})
)

func init() {
	prometheus.MustRegister(sslMetric)
	prometheus.MustRegister(up)
}

// Config File Settings
type Config struct {
	Targets []string
}

// Specification CLI Arguments for Changing exporter behavior
type Specification struct {
	Debug         bool   `default:"false"`
	ListenAddress string `default:":9197"`
	MetricsPath   string `default:"/metrics"`
	//ProbePath     string `default:"/probe"`
	ConfigPath string `default:"/etc/prometheus/exporters/ssl_exporter/"`
}

// LoadConfig - Read in Config file.
func (s *Specification) LoadConfig() Config {
	fmt.Println(s.ConfigPath)
	files, err := ioutil.ReadDir(s.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	for _, file := range files {
		filename := fmt.Sprintf(s.ConfigPath+"%s", file.Name())
		yamlFile, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Errorf("Unable to Read Config file", err)
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Errorf("Unable to Unmarshal Config", err)
		}

	}
	return config

}

func uniq(certs []*x509.Certificate) []*x509.Certificate {
	r := []*x509.Certificate{}

	for _, c := range certs {
		if !contains(r, c) {
			r = append(r, c)
		}
	}

	return r
}

func contains(certs []*x509.Certificate, cert *x509.Certificate) bool {
	for _, c := range certs {
		if (c.SerialNumber.String() == cert.SerialNumber.String()) && (c.Issuer.CommonName == cert.Issuer.CommonName) {
			return true
		}
	}
	return false
}

// func probeHandler(w http.ResponseWriter, r *http.Request, cfg Config) {
// 	// time.Sleep(15 * time.Second)
// 	domains := cfg.Targets
// 	for _, domain := range domains {
// 		x := sslStats(domain)
// 		if x == 0.0 {
// 			up.With(prometheus.Labels{"domain": domain}).Set(0)
// 		} else {
// 			sslMetric.With(prometheus.Labels{"domain": domain}).Set(x)
// 			up.With(prometheus.Labels{"domain": domain}).Set(1)
// 		}
// 	}
// }

func sslStats(target string) (expires float64) {

	// Create the HTTP client and make a get request of the target
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: tr,
	}

	resp, err := client.Get(target)
	if err != nil {
		expires := 0.0
		return expires
	}

	peerCertificates := uniq(resp.TLS.PeerCertificates)

	// Loop through returned certificates and create metrics, more stats if we want them in future
	for _, cert := range peerCertificates {
		// subject_cn := cert.Subject.CommonName
		// issuer_cn := cert.Issuer.CommonName
		// subject_dnsn := cert.DNSNames
		// subject_emails := cert.EmailAddresses
		// subject_ips := cert.IPAddresses
		// serial_no := cert.SerialNumber.String()
		expires := float64(cert.NotAfter.UnixNano() / 1e9)
		return expires
	}
	return
}

func main() {
	flag.Parse()
	var s Specification

	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	var cfg = s.LoadConfig()

	log.Info("Starting Server: %s\n", s.ListenAddress)
	log.Info("Metrics Path: %s\n", s.MetricsPath)
	handler := promhttp.Handler()

	go func() {
		domains := cfg.Targets
		for _, domain := range domains {
			x := sslStats(domain)
			if x == 0.0 {
				up.With(prometheus.Labels{"domain": domain}).Set(0)
			} else {
				sslMetric.With(prometheus.Labels{"domain": domain}).Set(x)
				up.With(prometheus.Labels{"domain": domain}).Set(1)
			}
			//only need to gather stats twice per day
		}
		time.Sleep(43200 * time.Second)
	}()

	if s.MetricsPath == "" || s.MetricsPath == "/" {
		http.Handle(s.MetricsPath, handler)
	} else {
		//only gather stats when /probe is hit
		http.Handle(s.MetricsPath, prometheus.Handler())
		// http.HandleFunc(s.ProbePath, func(w http.ResponseWriter, r *http.Request) {
		// 	probeHandler(w, r, cfg)
		// })
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
				<head><title>Prometheus SSL Exporter</title></head>
				<body>
				<h1>Prometheus SSL Exporter</h1>
				<p><a href="` + s.MetricsPath + `">Metrics</a></p>
				</body>
				</html>`))
		})
	}
	err = http.ListenAndServe(s.ListenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

}
