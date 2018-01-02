package blockchain

import (
	"github.com/satori/go.uuid"
	"net/url"
	"log"
	"strings"
	"net"
)

type Node struct {
	nodeId uuid.UUID
	nodeUrl url.URL
}

type NodeRegistry struct {
	nodes map[uuid.UUID]Node
}

func newNode() *Node {
	ip := getOutboundIP()
	url := url.URL{Scheme: "http", Host: ip}
	nodeId := uuid.NewV1()

	return &Node{nodeId:nodeId, nodeUrl:url}
}

func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{}
}

func (nr *NodeRegistry) RegisterNode(node Node)  {
	_, ok := nr.nodes[node.nodeId]
	if !ok {
		log.Printf("Registering new node with ID %s", node.nodeId.String())
	}
	nr.nodes[node.nodeId] = node
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