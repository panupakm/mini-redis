// miniredis project server main.go
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"

	"github.com/panupakm/miniredis/internal/db"
	"github.com/panupakm/miniredis/internal/pubsub"
	"github.com/panupakm/miniredis/server"
)

func main() {
	port := flag.Uint("port", 9988, "port to listen on")
	addr := flag.String("addr", "localhost", "address to listen on")
	certPath := flag.String("cert", "", "certificate PEM file")
	keyPath := flag.String("key", "", "key PEM file")
	flag.Parse()

	cert := func() *tls.Certificate {
		if *certPath == "" || *keyPath == "" {
			return nil
		}
		cert, err := tls.LoadX509KeyPair(*certPath, *keyPath)
		if err != nil {
			log.Fatal(err)
		}
		return &cert
	}()

	var config *tls.Config
	if cert != nil {
		config = &tls.Config{Certificates: []tls.Certificate{*cert}}
	} else {
		config = nil
	}

	s := server.NewServer(*addr, *port, db.NewDb(), pubsub.NewPubSub(), config)
	fmt.Println("Server started")
	s.ListenAndServe()
}
