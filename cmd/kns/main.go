package main

import (
	"flag"
	"log"
	"net/http"

	"kns/internal/knsapi"
	"kns/internal/nameset"
	"kns/internal/protocol"
	"kns/internal/resolver"
	"kns/internal/server"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	api := flag.String("api", knsapi.DefaultBase, "KNS indexer base URL")
	flag.Parse()

	res := &resolver.Resolver{API: knsapi.New(*api, protocol.Mainnet)}
	srv, err := server.New(*addr, res, nameset.New())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("KNS Web4 http://localhost%s — indexer + agent card + silverc artifacts", *addr)
	log.Fatal(http.ListenAndServe(*addr, srv.Handler()))
}
