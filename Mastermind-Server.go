package main

import "github.com/paradoxxl/gomastermind/server"

func main() {
	server.NewServer()
}

/*
import (

	"github.com/paradoxxl/gomastermind/server"
	"flag"
)



var genCert = flag.Bool("c", false, "genCert")
var keyFile = flag.String("kf", "priv.key", "keyFile")
var pemFile = flag.String("pf", "pub.pem", "pemFile")
var port = flag.Int("p", 8758, "port")


func main() {

	flag.Parse()

	if *genCert {
		server.GenCert(server.DefaultCert, *keyFile, *pemFile)
	}
	srv := server.NewMastermindServer(*pemFile, *keyFile, *port)
	defer srv.Close()

}
*/