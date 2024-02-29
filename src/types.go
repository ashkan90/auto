package src

type NodeId string

type ConnectionId string

type NodeBase struct {
	ID NodeId
}

type ConnectionBase struct {
	ID     ConnectionId
	Source NodeId
	Target NodeId
}

type BaseScheme[NodeData NodeBase, ConnectionData ConnectionBase] struct {
	Node       NodeData
	Connection ConnectionData
}

type NodeExecutor interface {
	Execute(input string, forward func(output string))
}

type NodeData interface {
	Data(inputs func() map[string]any) map[string]any
}
