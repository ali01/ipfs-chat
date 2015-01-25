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
	ChatId   string
	Messages []Message
}

type Message struct {
	Message   string
	Timestamp time.Time
}

func initNode() (*IpfsNode, error) {
	ipfsPath, err := fsrepo.BestKnownPath()
	log.Println(ipfsPath)
	if err != nil {
		return nil, err
	}

	repo := fsrepo.At(ipfsPath)
	if err := repo.Open(); err != nil {
		return nil, err
	}

	node, err := core.NewIPFSNode(context.Background(), core.Online(repo))
	if err != nil {
		return nil, err
	}

	return node, nil
}

func initSession(node *IpfsNode) (*Session, error) {
	hasher := sha256.New()
	hasher.Write([]byte("/ipfs-chat/"))
	hasher.Write([]byte(node.Identity))

	publishId := hasher.Sum(nil)

	mbuf, err := multihash.Encode(publishId, multihash.SHA2_256)
	if err != nil {
		return err
	}

	publishIdString := base58.Encode(mbuf)

	session := &Session{
		PublishId: publishIdString,
		Messages:  make([]Message),
	}

	log.Println("Connecting to peers...")

	peers := node.PeerHost.Network().Peers()
	for len(peers) < 1 {
		peers = node.PeerHost.Network().Peers()
		time.Sleep(1 * time.Second)
	}

	log.Println("Connected.")

	return session
}

func main() {
	node, err := initNode()
	session, err := initSession(node)

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

	if err != nil {
		log.Fatal(err)
	}
}
