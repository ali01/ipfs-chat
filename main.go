package main

import (
	"bufio"
	_ "errors"
	"fmt"
	"log"
	"os"
	"time"

	"code.google.com/p/go.net/context"
	"github.com/gogo/protobuf/proto"

	base58 "github.com/jbenet/go-base58"
	"github.com/jbenet/go-ipfs/core"

	"crypto/sha256"
	multihash "github.com/jbenet/go-multihash"

	// "github.com/jbenet/go-ipfs/core/coreunix"

	// "github.com/jbenet/go-ipfs/routing/dht"

	"github.com/jbenet/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-logging"
	"github.com/jbenet/go-ipfs/p2p/peer"
	"github.com/jbenet/go-ipfs/repo/fsrepo"
	u "github.com/jbenet/go-ipfs/util"
)

type ChatId string

type Session struct {
	Name         string
	PublishId    ChatId
	SubscribeIds []ChatId

	Stream *Stream
	Node   *core.IpfsNode
}

func InitNode() (*core.IpfsNode, error) {
	ipfsPath, err := fsrepo.BestKnownPath()
	if err != nil {
		return nil, err
	}

	repo := fsrepo.At(ipfsPath)
	if err := repo.Open(); err != nil {
		return nil, err
	}

	u.SetAllLoggers(logging.CRITICAL)
	node, err := core.NewIPFSNode(context.Background(), core.Online(repo))
	if err != nil {
		return nil, err
	}

	return node, nil
}

func InitSession(node *core.IpfsNode, name string,
	subscribePeers []peer.ID) (*Session, error) {

	log.Printf("Initializing Session with Peer ID: %s", node.Identity.Pretty())

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
		Name:         name,
		PublishId:    publishIdString,
		SubscribeIds: subscribeIds,

		Stream: &Stream{},
		Node:   node,
	}

	waitForPeers(session.Node)

	PublishSession(session)

	// value, err := dht.GetValue(context.Background(), "test key")
	// if err != nil {
	// 	return nil, err
	// }

	return session, nil
}

func PublishSession(session *Session) error {
	streamBytes, err := proto.Marshal(session.Stream)
	if err != nil {
		return err
	}

	dht := session.Node.Routing
	err = dht.PutValue(context.Background(), u.Key(session.PublishId), streamBytes)
	if err != nil {
		return err
	}

	return nil
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

func waitForPeers(node *core.IpfsNode) {
	log.Println("Connecting to peers...")

	peers := node.PeerHost.Network().Peers()
	for len(peers) < 1 {
		peers = node.PeerHost.Network().Peers()
		time.Sleep(1 * time.Second)
	}

	log.Println("Connected.")
}

func main() {
	// initialize IPFS Node
	node, err := InitNode()
	if err != nil {
		log.Fatal(err)
	}

	// initialize session
	name := "alive"
	peers := []peer.ID{"QmWWH49ZaWHc8wG9cPUGsnzRbUeNgJpus2aQT4Kou2oz7b"}
	session, err := InitSession(node, name, peers)
	if err != nil {
		log.Fatal(err)
	}

	for _, subscribeId := range session.SubscribeIds {
		go func(subscribeId ChatId) {

		}(subscribeId)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")

		msgString, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		message := &Message{
			Message:   proto.String(msgString),
			Timestamp: proto.Int64(time.Now().UnixNano()),
			PeerId:    proto.String(node.Identity.Pretty()),
			Name:      proto.String(session.Name),
		}

		session.Stream.Message = append(session.Stream.Message, message)

		PublishSession(session)
	}

	// Publish Data, get Data, and interleave

}
