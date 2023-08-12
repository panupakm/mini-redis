// client to connect to mini redis server
package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/panupakm/miniredis/client"
)

func main() {
	port := flag.Uint("port", 9988, "port to listen on")
	addr := flag.String("addr", "localhost", "address to listen on")
	certPath := flag.String("cert", "", "certificate PEM file")
	flag.Parse()

	fmt.Println(*certPath)
	config := func() *tls.Config {
		if certPath != nil {
			cert, err := os.ReadFile(*certPath)
			if err != nil {
				log.Fatal(err)
			}
			certPool := x509.NewCertPool()
			if ok := certPool.AppendCertsFromPEM(cert); !ok {
				log.Fatalf("unable to parse cert from %s", *certPath)
			}
			fmt.Println("Create certificate")
			return &tls.Config{RootCAs: certPool, MinVersion: tls.VersionTLS13}
		}
		return nil
	}()

	c := client.NewClient()
	if err := c.Connect(fmt.Sprintf("%s:%d", *addr, *port), config); err != nil {
		log.Fatalf("unable to connect to %s:%s", *addr, err)
	}
}
