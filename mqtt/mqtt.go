package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func Connect(clientId string, endpoint string, rootCAPath string, certPath string, keyPath string) (mqtt.Client, error) {
	tlsConfig, err := newTLSConfig(rootCAPath, certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to construct tls config: %v", err)
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%d", endpoint, 443))
	opts.SetTLSConfig(tlsConfig)
	opts.SetClientID(clientId)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect broker: %v", token.Error())
	}

	return client, nil
}

//-----------------------------------------------------------------------------
func Disonnect(client mqtt.Client, quiesce uint) {
	client.Disconnect(quiesce)
}

//-----------------------------------------------------------------------------
func Publish(client mqtt.Client, topic string, qos byte, retained bool, payload interface{}) error {
	token := client.Publish(topic, qos, retained, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		return fmt.Errorf("failed to publish %s: %v", topic, err)
	}
	return nil
}

//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func newTLSConfig(rootCAPath string, certPath string, keyPath string) (*tls.Config, error) {
	rootCA, err := ioutil.ReadFile(rootCAPath)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(rootCA)
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		RootCAs:            pool,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
		NextProtos:         []string{"x-amzn-mqtt-ca"}, // Port 443 ALPN
	}, nil
}
