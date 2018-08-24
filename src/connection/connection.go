package connection

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/globalsign/mgo"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-mongodb/src/arguments"
	"io/ioutil"
	"net"
	"os"
)

type ConnectionInfo struct {
	Username              string
	Password              string
	AuthSource            string
	Host                  string
	Port                  string
	Ssl                   bool
	SslCaCerts            string
	SslInsecureSkipVerify bool
}

func (c *ConnectionInfo) CreateSession() (*mgo.Session, error) {

	// TODO figure out how port fits into here
	dialInfo := mgo.DialInfo{
		Addrs:    []string{c.Host},
		Username: c.Username,
		Password: c.Password,
		Source:   c.AuthSource,
		FailFast: true,
	}

	if c.Ssl {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: c.SslInsecureSkipVerify,
		}

		if c.SslCaCerts != "" {
			roots := x509.NewCertPool()

			ca, err := ioutil.ReadFile(c.SslCaCerts)
			if err != nil {
				log.Error("Failed to open crt file: %v", err)
			}

			roots.AppendCertsFromPEM(ca)

			tlsConfig.RootCAs = roots
		}

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
	}

	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		log.Error("Failed to dial Mongo instance %s: %v", dialInfo.Addrs[0], err)
		os.Exit(1)
	}
	return session, err

}

func DefaultConnectionInfo() *ConnectionInfo {
	connectionInfo := &ConnectionInfo{
		Username:              arguments.GlobalArgs.Username,
		Password:              arguments.GlobalArgs.Password,
		AuthSource:            arguments.GlobalArgs.AuthSource,
		Host:                  arguments.GlobalArgs.Host,
		Port:                  arguments.GlobalArgs.Port,
		Ssl:                   arguments.GlobalArgs.Ssl,
		SslCaCerts:            arguments.GlobalArgs.SslCaCerts,
		SslInsecureSkipVerify: arguments.GlobalArgs.SslInsecureSkipVerify,
	}

	return connectionInfo

}
