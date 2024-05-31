package main

import (
	"log"

	"github.com/wreckitral/distributed-file-storage/p2p"
)

func main() {
    tr := p2p.NewTCPTransport(":7000")

    if err := tr.ListenAndAccept(); err != nil {
        log.Fatal(err)
    }

    select {}
}
