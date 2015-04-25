package main

import (
	"fmt"
	"net"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
)

func main() {
	conn, err := net.Dial("tcp", ":6363")
	if err != nil {
		fmt.Println(err)
		return
	}
	face := ndn.NewFace(conn, nil)
	defer face.Close()

	f := mux.NewFetcher()
	f.Use(mux.SHA256Verifier)
	f.Use(mux.Cacher)
	f.Use(mux.Logger)
	f.Use(mux.Assembler)
	spew.Dump(f.Fetch(face, ndn.NewName("/hello")))
}
