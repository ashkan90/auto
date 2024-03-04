package test

import (
	"github.com/ashkan90/auto-core/src"
	"log"
)

type TextNode struct {
	src.NodeInterface
}

func (n *TextNode) Data(inputs func() map[string]any) map[string]any {
	var m = make(map[string]any)
	n.Node().Controls.Range(func(key, value any) bool {
		m[key.(string)] = value.(src.ControlInterface).GetValue()
		return true
	})
	log.Println("[TextNode.Data]", m)
	return m
}

func (n *TextNode) Execute(input string, forward func(output string)) {
	forward("exec")
}

func NewTextNode() src.NodeInterface {
	return &TextNode{
		NodeInterface: src.NewNode(),
	}
}
