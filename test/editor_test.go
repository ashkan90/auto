package test

import (
	"encoding/json"
	"github.com/ashkan90/auto-core/src"
	"github.com/wI2L/jsondiff"
	"log"
	"os"
	"testing"
)

func TestEditorDeserialize(t *testing.T) {
	var bus = src.NewEventBus()
	var editor = src.NewNodeEditor(bus)

	n0 := src.NewNode()
	n0.AddControl("valueCtrl", src.NewInputControl(src.InputControlText, &src.InputControlOptions{
		Readonly: src.ToPtr(false),
		Initial:  src.ToPtr("hello"),
		Change: func(value any) {
			log.Println("value input ctrl data has been set", value)
		},
	}))
	n1, _ := editor.AddNode(src.NewNode())
	editor.AddNode(src.NewNode())
	editor.AddNode(src.NewNode())
	editor.AddNode(src.NewNode())

	n0.AddInput("input1", src.NewInput[src.Socket](src.NewSocket("Socket Name 1"), "Input 1 Label", true))
	n0.AddInput("input2", src.NewInput[src.Socket](src.NewSocket("Socket Name 2"), "Input 2 Label", true))

	n0.AddOutput("output", src.NewOutput[src.Socket](src.NewSocket("output"), "Output", true))
	n0.AddOutput("exec", src.NewOutput[src.Socket](src.NewSocket("exec"), "exec", true))

	n1.AddInput("exec", src.NewInput[src.Socket](src.NewSocket("exec"), "exec", true))
	n1.AddOutput("exec", src.NewOutput[src.Socket](src.NewSocket("exec"), "exec", true))

	editor.AddNode(n0)
	editor.AddNode(n1)

	editor.AddConnection(src.NewConnection(n0, "exec", n1, "exec"))

	res := editor.Deserialize()
	if len(res) == 0 {
		t.Error("deserialize didn't succeed")
	}

	log.Println(res)
}

// TestEditorSerialize input variable can be extracted by calling TestEditorDeserialize
func TestEditorSerialize(t *testing.T) {
	var input = []byte(`{"connections":{"0494bcf7073d2072":{"base":{"id":"0494bcf7073d2072","source":"","target":""},"source":"c9db91bfedfe1693","sourceOutput":"exec","target":"0ed475094cc9b7d2","targetInput":"exec"}},"nodes":{"0ed475094cc9b7d2":{"base":{"id":"0ed475094cc9b7d2"},"inputs":{"exec":{"port":{"id":"6b37f68146f6587c","label":"exec","index":0,"multipleConnections":true,"socket":{"name":"exec"}},"control":{"id":"45daa42ebe70810b","index":0},"showControl":true,"label":"exec"}},"outputs":{"exec":{"port":{"id":"1cf48a063b31320c","label":"exec","index":0,"multipleConnections":true,"socket":{"name":"exec"}},"control":{"id":"abc1d37fb47ebbfc","index":0},"showControl":true,"label":"exec"}},"controls":{},"selected":null},"773888f0b4e18562":{"base":{"id":"773888f0b4e18562"},"inputs":{},"outputs":{},"controls":{},"selected":null},"c0369f691e59a00d":{"base":{"id":"c0369f691e59a00d"},"inputs":{},"outputs":{},"controls":{},"selected":null},"c9db91bfedfe1693":{"base":{"id":"c9db91bfedfe1693"},"inputs":{"input1":{"port":{"id":"3bc9dea1d9ad2c27","label":"Input 1 Label","index":0,"multipleConnections":true,"socket":{"name":"Socket Name 1"}},"control":{"id":"bb1b6bbfb661df31","index":0},"showControl":true,"label":"Input 1 Label"},"input2":{"port":{"id":"3bac569b6585d626","label":"Input 2 Label","index":0,"multipleConnections":true,"socket":{"name":"Socket Name 2"}},"control":{"id":"bc401c31a2f04743","index":0},"showControl":true,"label":"Input 2 Label"}},"outputs":{"exec":{"port":{"id":"ec1b539c0d3b0b47","label":"exec","index":0,"multipleConnections":true,"socket":{"name":"exec"}},"control":{"id":"ff52a7712d420f57","index":0},"showControl":true,"label":"exec"},"output":{"port":{"id":"f31510a0a5336a2f","label":"Output","index":0,"multipleConnections":true,"socket":{"name":"output"}},"control":{"id":"c2f227691b1f2b05","index":0},"showControl":true,"label":"Output"}},"controls":{"valueCtrl":{"control":{"id":"1eb0f2ee8f361576","index":0},"type":"text","options":{"readonly":false,"initial":"hello"},"readonly":false,"value":"hello"}},"selected":null},"e5c780721534302e":{"base":{"id":"e5c780721534302e"},"inputs":{},"outputs":{},"controls":{},"selected":null}}}`)
	var bus = src.NewEventBus()

	var jsonEditorData, errData = src.NewJSONEditorData(input)
	if errData != nil {
		t.Error(errData)
	}

	var jsonEditor, errEditor = src.NewJSONEditor(jsonEditorData)
	if errEditor != nil {
		t.Error(errEditor)
	}

	var editor = src.NewNodeEditorFromJSON(bus, jsonEditor)

	if deserialize := editor.Deserialize(); deserialize != string(input) {
		patch, _ := jsondiff.CompareJSON([]byte(deserialize), input)
		b, _ := json.MarshalIndent(patch, "", "    ")

		os.Stdout.Write(b)
		t.Fail()
	}
}
