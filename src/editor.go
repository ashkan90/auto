package src

import (
	"errors"
	"slices"
	"sync"
)

type Root[Scheme BaseScheme[NodeBase, ConnectionBase]] struct {
	Type string
	Data Scheme
}

// NodeEditor, düğümleri ve bağlantıları yöneten bir yapı.
type NodeEditor struct {
	nodes       map[NodeId]NodeInterface
	connections map[ConnectionId]*Connection[ConnectionBase]
	lock        sync.RWMutex
	eventBus    *EventBus
}

// NewNodeEditor, yeni bir NodeEditor örneği oluşturur.
func NewNodeEditor(bus *EventBus) *NodeEditor {
	return &NodeEditor{
		nodes:       make(map[NodeId]NodeInterface),
		connections: make(map[ConnectionId]*Connection[ConnectionBase]),
		eventBus:    bus,
	}
}

// AddNode, bir düğüm ekler.
func (e *NodeEditor) AddNode(node NodeInterface) (NodeInterface, error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	n := node.Node()

	if _, exists := e.nodes[n.E.ID]; exists {
		return nil, errors.New("node already exists")
	}

	e.nodes[n.E.ID] = node
	e.eventBus.Publish(Event{Type: "nodeCreated", Data: node})
	return node, nil
}

// RemoveNode, bir düğümü kaldırır.
func (e *NodeEditor) RemoveNode(nodeID NodeId) (NodeId, error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, exists := e.nodes[nodeID]; !exists {
		return "", errors.New("node does not exist")
	}

	delete(e.nodes, nodeID)
	e.eventBus.Publish(Event{Type: "nodeRemoved", Data: nodeID})
	return nodeID, nil
}

// GetNode, bir düğüm döndürür.
func (e *NodeEditor) GetNode(nodeID NodeId) (NodeInterface, error) {
	e.lock.RLock()
	defer e.lock.RUnlock()

	node, exists := e.nodes[nodeID]
	if !exists {
		return &Node[NodeBase]{}, errors.New("node does not exist")
	}

	return node, nil
}

// GetNodes, tüm düğümleri döndürür.
func (e *NodeEditor) GetNodes() []NodeInterface {
	e.lock.RLock()
	defer e.lock.RUnlock()

	nodes := make([]NodeInterface, 0, len(e.nodes))
	for _, node := range e.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// AddConnection, bir bağlantı ekler.
func (e *NodeEditor) AddConnection(conn *Connection[ConnectionBase]) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, exists := e.connections[conn.E.ID]; exists {
		return errors.New("connection already exists")
	}

	e.connections[conn.E.ID] = conn
	e.eventBus.Publish(Event{Type: "connectionAdded", Data: conn})
	return nil
}

// RemoveConnection, bir bağlantıyı kaldırır.
func (e *NodeEditor) RemoveConnection(connID ConnectionId) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, exists := e.connections[connID]; !exists {
		return errors.New("connection does not exist")
	}

	delete(e.connections, connID)
	e.eventBus.Publish(Event{Type: "connectionRemoved", Data: connID})
	return nil
}

// GetConnection, belirtilen ID'ye sahip bir bağlantıyı döndürür.
func (e *NodeEditor) GetConnection(connID ConnectionId) (*Connection[ConnectionBase], error) {
	e.lock.RLock()
	defer e.lock.RUnlock()

	conn, exists := e.connections[connID]
	if !exists {
		return &Connection[ConnectionBase]{}, errors.New("connection does not exist")
	}

	return conn, nil
}

// GetConnections mevcut tüm bağlantıları döndürür.
func (e *NodeEditor) GetConnections() []*Connection[ConnectionBase] {
	e.lock.RLock()
	defer e.lock.RUnlock()

	conns := make([]*Connection[ConnectionBase], 0, len(e.connections))
	for _, conn := range e.connections {
		conns = append(conns, conn)
	}
	return conns
}

// GetConnectionsTo verilen node'un mevcut tüm bağlantıları döndürür.
func (e *NodeEditor) GetConnectionsTo(nodeID NodeId, inputKeys []string) []*Connection[ConnectionBase] {
	e.lock.RLock()
	defer e.lock.RUnlock()

	conns := make([]*Connection[ConnectionBase], 0, len(e.connections))
	for _, conn := range e.connections {
		if conn.Target == nodeID && slices.Contains(inputKeys, string(conn.TargetInput)) {
			conns = append(conns, conn)
		}
	}
	return conns
}

func (e *NodeEditor) GetBus() *EventBus {
	return e.eventBus
}
