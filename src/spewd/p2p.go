package main

import (
	"blockchain"
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"strconv"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"

	ma "github.com/multiformats/go-multiaddr"
)

// P2pMsgType represents the message type
type P2pMsgType uint64

const (
	// P2pMsgIdentity identifies the peer
	P2pMsgIdentity P2pMsgType = iota
	// P2pMsgTree sends the base (merkle?) tree
	P2pMsgTree = iota
	// P2pMsgChain sends a chain fragment for syncing
	P2pMsgChain = iota
)

// P2pIdentity represents the peers identity
type P2pIdentity struct {
	ID string
}

// P2pMsg represents the message
type P2pMsg struct {
	Type P2pMsgType
	Data interface{}
}

// UnmarshalJSON unmarshales a P2pMsg handling the Data
func (msg *P2pMsg) UnmarshalJSON(buf []byte) error {
	fmt.Println("GET TYPE")
	var raw struct{
		Type P2pMsgType
		Data *json.RawMessage
	}

	if err := json.Unmarshal(buf, &raw); err != nil {
		return nil
	}

	switch msg.Type {
	case P2pMsgIdentity:
		fmt.Println("ID PLEASE")
		var id P2pIdentity
		err := json.Unmarshal(*raw.Data, &id)
		if err != nil {
			return err
		}
		msg.Data = id
		return nil
	}
	return fmt.Errorf("unknown data type")
}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, ma.Multiaddr, error) {
	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort),
		),
		libp2p.Identity(priv),
	}

	if !secio {
		opts = append(opts, libp2p.NoEncryption())
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fulladdr := addr.Encapsulate(hostAddr)
	/*
	if secio {
		log.Printf("Now run \"go run main.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
	} else {
		log.Printf("Now run \"go run main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	}
	*/

	return basicHost, fulladdr, nil
}

func getStreamHandler(addr ma.Multiaddr) func(s net.Stream) {
	return func(s net.Stream) {
		//cid := make(chan P2pIdentity)
		//ctree := make(chan blockchain.Tree)
		log.Println("Got a new stream!")
		// Create a buffer stream for non blocking read and write.
		rw := bufio.NewReadWriter(
			bufio.NewReader(s), bufio.NewWriter(s))
		go readData(rw)
		go writeData(addr, rw)
		// stream 's' will stay open until you close it (or the other side closes it).
	}
}

func readData(rw *bufio.ReadWriter) {
	msg := &P2pMsg{}
	decoder := json.NewDecoder(rw)
	if err := decoder.Decode(&msg); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
	fmt.Printf("PEER: %s\n", msg.Data.(P2pIdentity).ID)
	for {
		var fragment blockchain.Fragment
		if err := decoder.Decode(&fragment); err == nil {
			fmt.Printf("replace with chain: %d...\n",
				len(fragment.Blocks))
			if err := Blockchain.Replace(&fragment); err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
			}
		} else {
			fmt.Printf("[E] %s\n", err.Error())
			return
		}
	}
}

func writeChain(encoder *json.Encoder, rw *bufio.ReadWriter) error {
	last := Blockchain.Get(-1)
	if last.Index <= 0 {
		return nil
	}
	fmt.Println("send out blocks...")
	fragment := Blockchain.GetFragment(1, last.Index)
	if err := encoder.Encode(fragment); err != nil {
		return err
	}
	rw.WriteString("\n")
	rw.Flush()
	return nil
}

func writeData(fulladdr ma.Multiaddr, rw *bufio.ReadWriter) {
	t5m := time.NewTicker(15 * time.Second)
	encoder := json.NewEncoder(rw)

	log.Printf("I am %s\n", fulladdr)
	encoder.Encode(&P2pMsg{
		Type: P2pMsgIdentity,
		Data: &P2pIdentity{
			ID: fulladdr.String(),
		},
	})
	rw.WriteString("\n")
	rw.Flush()
	for {
		select {
		case <-t5m.C:
			if err := writeChain(encoder, rw); err != nil {
				fmt.Printf("WRITE-ERROR: %s\n", err.Error())
				return
			}
		}
	}
}

func startPeer() {
	port, _ := strconv.Atoi(os.Getenv("NODE_ADDR"))
	basic, fulladdr, _ := makeBasicHost(port, true, 0)

	basic.SetStreamHandler("/p2p/1.0.0", getStreamHandler(fulladdr))
	if len(os.Args) <= 1 {
		log.Printf("I am %s\n", fulladdr)
		select{}
	} else {
		// The following code extracts target's peer ID from the
		// given multiaddress
		addr, err := ma.NewMultiaddr(os.Args[1])
		if err != nil {
			panic(err)
		}

		pid, err := addr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			panic(err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			panic(err)
		}
		tpaddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		taddr := addr.Decapsulate(tpaddr)

		basic.Peerstore().AddAddr(peerid, taddr, pstore.PermanentAddrTTL)
		s, err := basic.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		// Create a thread to read and write data.
		go writeData(fulladdr, rw)
		go readData(rw)

		select{}
	}
}
