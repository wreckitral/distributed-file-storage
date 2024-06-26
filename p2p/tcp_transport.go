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

type TCPTransportOpts struct {
    ListenAddr string
    HandshakeFunc HandshakeFunc
    Decoder Decoder
}

type TCPTransport struct {
    TCPTransportOpts
    listener net.Listener

    mu sync.RWMutex
    peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
    return &TCPTransport{
        TCPTransportOpts: opts,
    }
}

// this method sets up TCP listener and accept incoming connections
func (t *TCPTransport) ListenAndAccept() error {
    var err error

    t.listener, err = net.Listen("tcp", t.ListenAddr)
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

        fmt.Printf("new incoming connection %+v\n", conn)

        go t.handleConn(conn)
    }
}

type Temp struct {}

// helper func handling individual connection
func (t *TCPTransport) handleConn(conn net.Conn) {
    peer := NewTCPPeer(conn, true)
    
    if err := t.HandshakeFunc(peer); err != nil {
        conn.Close() 
        fmt.Printf("TCP handshake error: %s\n", err) 
        return 
    }
    
    // Read loop
    rpc := &RPC{}
    for {
        if err := t.Decoder.Decode(conn, rpc); err != nil {
            fmt.Printf("tcp error: %s\n", err)
            continue
        }

        rpc.From = conn.RemoteAddr()

        fmt.Printf("message:%+v\n", rpc)
    }
}
