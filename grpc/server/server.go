package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"grpc-example/chat"
)

func CreateTLSConfig() *tls.Config {
	caCert, err := ioutil.ReadFile("../certs/ca-cert.pem")
	if err != nil {
		log.Fatalf("Couldn't load the ca certificate, reason [%v]", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair("../certs/server-cert.pem", "../certs/server-key.pem")
	if err != nil {
		log.Fatalf("Couldn't load the server certificate, reason [%v]", err)
	}

	return &tls.Config{
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
	}
}

func main() {
	config := CreateTLSConfig()
	creds := credentials.NewTLS(config)

	s := chat.Server{}
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	chat.RegisterChatServiceServer(grpcServer, &s)

	go func() {
		lis, err := net.Listen("tcp", ":9000")
		if err != nil {
			log.Fatalf("Failed to listen on port 9000: %v", err)
		}

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to server gRPC server over port 9000: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	grpcServer.GracefulStop()

}
