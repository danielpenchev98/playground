package main

import (
	"crypto/tls"
	"crypto/x509"
	"grpc-example/chat"
	"io/ioutil"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func CreateTLSConfig() *tls.Config {
	caCert, err := ioutil.ReadFile("../certs/ca-cert.pem")
	if err != nil {
		log.Fatalf("Couldn't load the ca certificate, reason [%v]", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair("../certs/client-cert.pem", "../certs/client-key.pem")
	if err != nil {
		log.Fatalf("Couldn't load the server certificate, reason [%v]", err)
	}

	return &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}
}

func main() {
	var conn *grpc.ClientConn

	dialOption := grpc.WithTransportCredentials(credentials.NewTLS(CreateTLSConfig()))
	conn, err := grpc.Dial(":9000", dialOption)
	if err != nil {
		log.Fatalf("could not connect :%s", err)
	}

	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	message := chat.Message{
		Body: "Hello from the client!",
	}
	response, err := c.SayHello(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response from Server: %s", response.Body)
}
