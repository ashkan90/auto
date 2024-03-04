package src

type NodeId string

type ConnectionId string

type NodeBase struct {
	ID NodeId `json:"id"`
}

type ConnectionBase struct {
	ID     ConnectionId `json:"id"`
	Source NodeId       `json:"source"`
	Target NodeId       `json:"target"`
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
