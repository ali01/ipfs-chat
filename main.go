package main

import (
	_ "errors"
	"log"
	"time"

	"code.google.com/p/go.net/context"

	base58 "github.com/jbenet/go-base58"
	"github.com/jbenet/go-ipfs/core"

	"crypto/sha256"
	multihash "github.com/jbenet/go-multihash"

	// "github.com/jbenet/go-ipfs/core/coreunix"

	// "github.com/jbenet/go-ipfs/routing/dht"

	// "github.com/jbenet/go-ipfs/p2p/host"

	"github.com/jbenet/go-ipfs/repo/fsrepo"
)

type Session struct {
	messages []string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	ipfsPath, err := fsrepo.BestKnownPath()
	log.Println(ipfsPath)
	if err != nil {
		return err
	}

	repo := fsrepo.At(ipfsPath)
	if err := repo.Open(); err != nil {
		return err
	}

	node, err := core.NewIPFSNode(context.Background(), core.Online(repo))
	if err != nil {
		return err
	}

	hash := sha256.New()
	hash.Write([]byte("/ipfs-chat/"))
	hash.Write([]byte(node.Identity))

	hashValue := hash.Sum(nil)

	mbuf, err := multihash.Encode(hashValue, multihash.SHA2_256)
	if err != nil {
		return err
	}

	hashString := base58.Encode(mbuf)

	// connect to peers

	peers := node.PeerHost.Network().Peers()
	for len(peers) < 1 {
		peers = node.PeerHost.Network().Peers()
		log.Println(peers)
		time.Sleep(1 * time.Second)
	}

	// TODO: take in list of Peer.ID keys to listen on

	// TODO: derive DHT key to publish on from Chat-Prefix + Peer.ID

	// Publish Data, get Data, and interleave

	// dht := node.Routing
	// err = dht.PutValue(context.Background(), "test key", []byte("test value"))
	// if err != nil {
	// 	return err
	// }

	// value, err := dht.GetValue(context.Background(), "test key")
	// if err != nil {
	// 	return err
	// }

	// log.Println(string(value))

	return nil
}
