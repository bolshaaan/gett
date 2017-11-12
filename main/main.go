package main

import (
	"github.com/bolshaaan/gett"
	"flag"
)

var addr = flag.String("addr", "localhost:8080", "ip:port to listen")
var pgUrl = flag.String("pgUrl", `postgresql://postgres@127.0.0.1/gett?sslmode=disable`, "postgres url scheme")

func main() {
	flag.Parse()

	gett.StartApp(*addr, *pgUrl)
}