// miniredis project server main.go
package main

import (
	"crypto/tls"
	"flag"
	"log"

	"github.com/panupakm/miniredis/server"
)

func main() {
	port := flag.Uint("port", 9988, "port to listen on")
	addr := flag.String("addr", "localhost", "address to listen on")
	certPath := flag.String("cert", "", "certificate PEM file")
	keyPath := flag.String("key", "", "key PEM file")
	restorePath := flag.String("restore", "", "path to restore from")
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

	var config = server.Config{}
	if cert != nil {
		config.Config = tls.Config{
			Certificates: []tls.Certificate{*cert},
			MinVersion:   tls.VersionTLS13,
		}
	}
	config.PersistentPath = *restorePath

	s := InitializeServer()
	s.ListenAndServe(*addr, *port, &config)
}
