package blockchain

import (
	"github.com/satori/go.uuid"
	"net/url"
	"log"
	"strings"
	"net"
)

type Node struct {
	NodeId  uuid.UUID `json:"id,string"`
	NodeUrl url.URL   `json:"url,string"`
}

type NodeRegistry struct {
	Nodes map[uuid.UUID]Node
}

func newNode() *Node {
	ip := getOutboundIP()
	url := url.URL{Scheme: "http", Host: ip}
	nodeId, _ := uuid.NewV1()

	return &Node{NodeId:nodeId, NodeUrl:url}
}

func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{make(map[uuid.UUID]Node)}
}

func (nr *NodeRegistry) RegisterNode(node Node)  {
	_, ok := nr.Nodes[node.NodeId]
	if !ok {
		log.Printf("Registering new node with ID %s", node.NodeId.String())
	}
	nr.Nodes[node.NodeId] = node
}

func getOutboundIP() (string) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	if err != nil {
		log.Printf("No outbound IP found, using localhost instead")
		return "localhost"
	}

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}