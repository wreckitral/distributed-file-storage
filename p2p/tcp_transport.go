package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node 
type TCPPeer struct {
    // conn is the underlying connections of the peer
    conn net.Conn
    
    // if we dial and retrieve a connection, outbound := true
    // if we accept and retrieve a connection, outbound := false
    outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
    return &TCPPeer{
        conn: conn,
        outbound: outbound,
    }
}

type TCPTransport struct {
    listenAddress string
    listener net.Listener

    mu sync.RWMutex
    peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
    return &TCPTransport{
        listenAddress: listenAddr,
    }
}

// this method sets up TCP listener and accept incoming connections
func (t *TCPTransport) ListenAndAccept() error {
    var err error

    t.listener, err = net.Listen("tcp", t.listenAddress)
    if err != nil {
        return err
    }

    go t.startAcceptLoop()

    return nil
}

// helper method for accepting and handling incoming connections
func (t *TCPTransport) startAcceptLoop() {
    for {
        conn, err := t.listener.Accept()
        if err != nil {
            fmt.Printf("TCP accept error: %s\n", err) 
        }

        go t.handleConn(conn)
    }
}

// helper func handling individual connection
func (t *TCPTransport) handleConn(conn net.Conn) {
    peer := NewTCPPeer(conn, true)
    fmt.Printf("new incoming connection %+v\n", peer)
}
