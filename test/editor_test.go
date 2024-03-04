package test

import (
	"github.com/ashkan90/auto-core/src"
	"log"
	"testing"
)

func TestEditorDeserialize(t *testing.T) {
	var bus = src.NewEventBus()
	var editor = src.NewNodeEditor(bus)

	n0 := src.NewNode()
	n1, _ := editor.AddNode(src.NewNode())
	editor.AddNode(src.NewNode())
	editor.AddNode(src.NewNode())
	editor.AddNode(src.NewNode())

	n0.AddInput("input1", src.NewInput(src.NewSocket("Socket Name 1"), "Input 1 Label", true))
	n0.AddInput("input2", src.NewInput(src.NewSocket("Socket Name 2"), "Input 2 Label", true))

	n0.AddOutput("output", src.NewPort(src.NewSocket("Socket Name 3"), "Output Label", true))
	n1.AddOutput("exec", src.NewPort(src.NewSocket("Socket Name 3"), "Execute Label", true))
	editor.AddNode(n0)

	log.Println(editor.Deserialize())
}
