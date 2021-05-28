package gocqlastra

import (
	"archive/zip"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/gocql/gocql"
)

type clusterConfig struct {
	config              *gocql.ClusterConfig
	secureConnectBundle string
}

// NewCluster generates a new config for the default cluster implementation
// and implemented required configuration for connecting to Astra Cassandra.
func NewCluster(secureConnectBundle string) (*gocql.ClusterConfig, error) {
	cluster := &clusterConfig{
		config:              gocql.NewCluster(),
		secureConnectBundle: secureConnectBundle,
	}
	if err := parseZipFile(cluster); err != nil {
		return nil, err
	}
	return cluster.config, nil
}

func (c *clusterConfig) setSSLOptions(entries map[string][]byte) error {
	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(entries["ca.crt"])

	cert, err := tls.X509KeyPair(entries["cert"], entries["key"])
	if err != nil {
		return fmt.Errorf("cannot parse a public/private key pair from a pair of PEM encoded data: %v", err)
	}

	tlsConfig := &tls.Config{
		ServerName:   "*.db.astra.datastax.com",
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{cert},
	}
	c.config.SslOpts = &gocql.SslOptions{
		Config:                 tlsConfig,
		EnableHostVerification: true,
	}

	return nil
}

type bundleConfig struct {
	Host    *string `json:"host"`
	CQLPort *int64  `json:"cql_port"`
}

func parseZipFile(cluster *clusterConfig) error {
	r, err := zip.OpenReader(cluster.secureConnectBundle)
	if err != nil {
		return err
	}
	defer func() {
		_ = r.Close()
	}()

	zipEntries := make(map[string][]byte)
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			continue
		}

		buf, err := ioutil.ReadAll(rc)
		if err != nil {
			continue
		}
		_ = rc.Close()

		zipEntries[f.Name] = buf
	}

	if _, ok := zipEntries["config.json"]; !ok {
		return errors.New("config file must be contained in secure bundle")
	}

	var bundleCfg *bundleConfig
	if err := json.Unmarshal(zipEntries["config.json"], &bundleCfg); err != nil {
		return err
	}
	if bundleCfg.Host == nil || bundleCfg.CQLPort == nil {
		return errors.New("config file must include host and cql_port information")
	}

	cluster.config.Hosts = []string{fmt.Sprintf("%s:%d", *bundleCfg.Host, *bundleCfg.CQLPort)}
	return cluster.setSSLOptions(zipEntries)
}
