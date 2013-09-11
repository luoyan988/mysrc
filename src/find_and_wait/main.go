// Runs a node on UDP port 11221 that attempts to collect 100 peers for an
// infohash, then keeps running as a passive DHT node.
//
// IMPORTANT: if the UDP port is not reachable from the public internet, you
// may see very few results.
//
// To collect 100 peers, it usually has to contact some 10k nodes. This process
// is not instant and should take a minute or two, depending on your network
// connection.
//
//
// There is a builtin web server that can be used to collect debugging stats
// from http://localhost:8711/debug/vars.
package main

import (
	l4g "code.google.com/p/log4go"
	"crypto/rand"
	"encoding/hex"

	"fmt"
	"github.com/nictuku/mydht"

	"os"
	"time"
)

const (
	httpPortTCP = 8711
	//  = 11221
)

var flagc chan bool

func main() {
	/*flag.Parse()
	// Change to l4g.DEBUG to see *lots* of debugging information.
	l4g.AddFilter("stdout", l4g.WARNING, l4g.NewConsoleLogWriter())
	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %v <infohash>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example infohash: d1c5676ae7ac98e8b19f63565905105e3c4c37a2\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	/*ih, err := dht.DecodeInfoHash(flag.Args()[0])
	if err != nil {
		l4g.Critical("DecodeInfoHash error: %v\n", err)
		os.Exit(1)
	}*/

	// This is a hint to the DHT of the minimum number of peers it will try to
	// find for the given node. This is not a reliable limit. In the future this
	// might be moved to "PeersRequest()", so the controlling client can have
	for i := 0; i < 50; i++ {
		dhtPortUDP := 11221 + i
		go runDHT(dhtPortUDP)

	}

	<-flagc

}

// drainresults loops, printing the address of nodes it has found.
func drainresults(n *dht.DHT) {
	fmt.Println("=========================== DHT")
	l4g.Warn("Note that there are many bad nodes that reply to anything you ask.")
	l4g.Warn("Peers found:")
	for r := range n.PeersRequestResults {
		for _, peers := range r {
			for _, x := range peers {
				l4g.Warn("%v", dht.DecodePeerAddress(x))
			}
		}
	}
}
func randHash() string {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		l4g.Exit("nodeId rand:", err)
	}
	hash := hex.EncodeToString(b)
	return hash
}
func runDHT(port int) {
	targetNumPeers := 10
	d, err := dht.NewDHTNode(port, targetNumPeers, false)
	if err != nil {
		l4g.Critical("NewDHTNode error: %v", err)
		os.Exit(1)

	}
	// For debugging.
	//go http.ListenAndServe(fmt.Sprintf(":%d", httpPortTCP), nil)

	go d.DoDHT()
	go drainresults(d)

	for {
		// Give the DHT some time to "warm-up" its routing table.
		newhash := randHash()
		for i := 0; i < 200; i++ {
			time.Sleep(5 * time.Second)
			d.PeersRequest(newhash, false)
		}
	}
	flagc <- true

}
