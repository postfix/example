package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-ndn/mux"
	"github.com/go-ndn/ndn"
	"github.com/go-ndn/persist"
)

func main() {
	conn, err := net.Dial("tcp", ":6363")
	if err != nil {
		fmt.Println(err)
		return
	}
	recv := make(chan *ndn.Interest)
	face := ndn.NewFace(conn, recv)
	defer face.Close()

	var key ndn.Key
	pem, _ := ioutil.ReadFile("key/default.pri")
	key.DecodePrivateKey(pem)

	register := func(name string) {
		ndn.SendControl(face, "rib", "register", &ndn.Parameters{
			Name: ndn.NewName(name),
		}, &key)
	}

	register("/hello")
	register("/file")

	m := mux.New()
	m.Use(mux.Logger)
	m.Use(mux.Segmentor(10))
	m.Use(persist.Cacher("test.db"))
	m.HandleFunc("/hello", func(w mux.Sender, i *ndn.Interest) {
		spew.Dump(i)
		w.SendData(&ndn.Data{
			Name:    ndn.NewName("/hello"),
			Content: []byte(time.Now().UTC().String()),
		})
	})
	m.Handle("/file", mux.FileServer("/etc"), mux.PrefixTrimmer("/file"))
	m.Run(face, recv)
}
