package agent

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Jumpscale/agent2/agent/lib/settings"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type ControllerClient struct {
	URL    string
	Client *http.Client
}

func getKeys(m map[string]*ControllerClient) []string {
	keys := make([]string, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}

	return keys
}

func getHttpClient(security *settings.Security) *http.Client {
	var tlsConfig tls.Config

	if security.CertificateAuthority != "" {
		pem, err := ioutil.ReadFile(security.CertificateAuthority)
		if err != nil {
			log.Fatal(err)
		}

		tlsConfig.RootCAs = x509.NewCertPool()
		tlsConfig.RootCAs.AppendCertsFromPEM(pem)
	}

	if security.ClientCertificate != "" {
		if security.ClientCertificateKey == "" {
			log.Fatal("Missing certificate key file")
		}
		// pem, err := ioutil.ReadFile(security.ClientCertificate)
		// if err != nil {
		//     log.Fatal(err)
		// }

		cert, err := tls.LoadX509KeyPair(security.ClientCertificate,
			security.ClientCertificateKey)
		if err != nil {
			log.Fatal(err)
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			Proxy:               http.ProxyFromEnvironment,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     &tlsConfig,
		},
	}
}

func NewControllerClient(cfg settings.Controller) *ControllerClient {
	client := &ControllerClient{
		URL:    strings.TrimRight(cfg.URL, "/"),
		Client: getHttpClient(&cfg.Security),
	}

	return client
}

func (client *ControllerClient) BuildUrl(gid int, nid int, endpoint string) string {
	return fmt.Sprintf("%s/%d/%d/%s", client.URL,
		gid,
		nid,
		endpoint)
}
