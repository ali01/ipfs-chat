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

	"github.com/jbenet/go-ipfs/p2p/peer"

	"github.com/jbenet/go-ipfs/repo/fsrepo"
)

type ChatId string

type Session struct {
	PublishId      ChatId
	SubscribeIds   []ChatId
	MessageHistory []Message
}

type Message struct {
	Message   string
	Timestamp time.Time
}

func InitNode() (*core.IpfsNode, error) {
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

func InitSession(node *core.IpfsNode, subscribePeers []peer.ID) (*Session, error) {
	publishIdString, err := deriveChatId(node.Identity)
	if err != nil {
		return nil, err
	}

	var subscribeIds []ChatId
	for _, peerId := range subscribePeers {
		chatId, err := deriveChatId(peerId)
		if err != nil {
			return nil, err
		}

		subscribeIds = append(subscribeIds, chatId)
	}

	session := &Session{
		PublishId:    publishIdString,
		SubscribeIds: subscribeIds,
	}

	log.Println("Connecting to peers...")

	peers := node.PeerHost.Network().Peers()
	for len(peers) < 1 {
		peers = node.PeerHost.Network().Peers()
		time.Sleep(1 * time.Second)
	}

	log.Println("Connected.")

	return session, nil
}

func deriveChatId(peerId peer.ID) (ChatId, error) {
	hasher := sha256.New()
	hasher.Write([]byte("/ipfs-chat/"))
	hasher.Write([]byte(peerId))

	chatId := hasher.Sum(nil)

	mbuf, err := multihash.Encode(chatId, multihash.SHA2_256)
	if err != nil {
		return ChatId(0), err
	}

	chatIdString := ChatId(base58.Encode(mbuf))

	return chatIdString, nil
}

func main() {
	// initialize IPFS Node
	node, err := InitNode()
	if err != nil {
		log.Fatal(err)
	}

	// initialize session
	peers := []peer.ID{"QmWWH49ZaWHc8wG9cPUGsnzRbUeNgJpus2aQT4Kou2oz7b"}
	_, err = InitSession(node, peers)
	if err != nil {
		log.Fatal(err)
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
}
