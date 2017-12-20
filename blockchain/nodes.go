package blockchain

import (
	"github.com/satori/go.uuid"
	"net/url"
	"log"
)

type Node struct {
	nodeId uuid.UUID
	nodeUrl url.URL
}

type NodeRegistry struct {
	nodes map[uuid.UUID]Node
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