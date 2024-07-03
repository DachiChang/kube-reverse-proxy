package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	//"github.com/sanity-io/litter"
)

var (
	k8s_cert  string
	user_cert string
	user_key  string
	target    string
)

func init() {
	flag.StringVar(&k8s_cert, "k8s_cert", "", "path to authority cert")
	flag.StringVar(&user_cert, "user_cert", "", "path to user cert")
	flag.StringVar(&user_key, "user_key", "", "path to user key")
	flag.StringVar(&target, "target", "https://aks.privatelink.centralus.azmk8s.io:443", "aks control plane URI")
}

func main() {
	flag.Parse()

	if k8s_cert == "" || user_cert == "" || user_key == "" {
		fmt.Println("check -h to give more arguments")
		return
	}

	caCert, _ := os.ReadFile(k8s_cert)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		targetURI := target
		direct := func(req *http.Request) {
			*req = *r
			parsedURI, _ := url.Parse(targetURI)
			//req.Host = parsedURI.Host
			req.URL.Host = parsedURI.Host
			req.URL.Scheme = parsedURI.Scheme
			fmt.Println("request: ", req.URL)
			//req.URL.Path = parsedURI.Path
			//req.URL.RawQuery = parsedURI.RawQuery
			//req.Header.Set("Authorization", "Bearer a4hcviwiappn8ikca9msol8iy6jri3yddqx1xd38lnt8d3pn3w815cxwd14i17ic4p45eryazur58yo0nttu54yrnqzdbugqa9qie6772bk00s7n0yn4qy3xt1aom04q")
			//litter.Dump(req)
		}
		proxy := &httputil.ReverseProxy{Director: direct}
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    caCertPool,
				MinVersion: tls.VersionTLS13,
				MaxVersion: tls.VersionTLS13,
				GetClientCertificate: func(chi *tls.CertificateRequestInfo) (*tls.Certificate, error) {
					cert, _ := tls.LoadX509KeyPair(user_cert, user_key)
					return &cert, nil
				},
			},
		}
		proxy.ServeHTTP(w, r)
	})
	fmt.Println("start serve in: 127.0.0.1:3000")
	http.ListenAndServe(":3000", nil)
}
