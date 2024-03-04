package src

import (
	"encoding/json"
	"errors"
	"github.com/ashkan90/auto-core/utils"
	"sync"
)

type JSONEditor struct {
	Nodes       map[NodeId]NodeInterface                     `json:"nodes"`
	Connections map[ConnectionId]*Connection[ConnectionBase] `json:"connections"`
}

type JSONEditorData struct {
	Nodes       map[NodeId]*JSONEditorNode             `json:"nodes"`
	Connections map[ConnectionId]*JSONEditorConnection `json:"connections"`
}

type JSONEditorNode struct {
	Base     JSONEditorNodeBase `json:"base"`
	Inputs   map[string]any     `json:"inputs"`
	Outputs  map[string]any     `json:"outputs"`
	Controls map[string]any     `json:"controls"`
	Selected *bool              `json:"selected"`
}

type JSONEditorNodeBase struct {
	Id NodeId `json:"id"`
}

type JSONEditorConnection struct {
	Base         JSONEditorConnectionBase `json:"base"`
	Source       string                   `json:"source"`
	SourceOutput string                   `json:"sourceOutput"`
	Target       string                   `json:"target"`
	TargetInput  string                   `json:"targetInput"`
}

type JSONEditorConnectionBase struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

func NewJSONEditor(jsonData *JSONEditorData) (*JSONEditor, error) {
	if jsonData == nil {
		return nil, errors.New("empty input given")
	}

	var editor = &JSONEditor{
		Nodes:       make(map[NodeId]NodeInterface),
		Connections: make(map[ConnectionId]*Connection[ConnectionBase]),
	}

	for id, node := range jsonData.Nodes {
		var inlineNode = &Node[NodeBase]{
			E:        NodeBase{ID: id},
			Inputs:   utils.NewSyncMap(),
			Outputs:  utils.NewSyncMap(),
			Controls: utils.NewSyncMap(),
			Selected: node.Selected,
			mu:       &sync.Mutex{},
		}

		for inputId, inputValue := range node.Inputs {
			inputValueCpy := inputValue.(map[string]any)
			if inputValueCpy == nil {
				continue
			}
			inputValuePort := inputValueCpy["port"].(map[string]any)
			inputValueControl := inputValueCpy["control"].(map[string]any)

			inlineNode.Inputs.Add(inputId, &Input[Socket]{
				Port: &Port[Socket]{
					Id:                  PortId(inputValuePort["id"].(string)),
					Label:               inputValuePort["label"].(string),
					Index:               int(inputValuePort["index"].(float64)),
					MultipleConnections: inputValuePort["multipleConnections"].(bool),
					Socket: Socket{
						Name: inputValuePort["socket"].(map[string]any)["name"].(string),
					},
				},
				Control: &Control{
					Id:    inputValueControl["id"].(string),
					Index: int(inputValueControl["index"].(float64)),
				},
				ShowControl: inputValueCpy["showControl"].(bool),
				Label:       inputValueCpy["label"].(string),
			})
		}

		for outputId, outputValue := range node.Outputs {
			outputValueCpy := outputValue.(map[string]any)
			if outputValue == nil {
				continue
			}

			outputValuePort := outputValueCpy["port"].(map[string]any)
			outputValueControl := outputValueCpy["control"].(map[string]any)

			inlineNode.Outputs.Add(outputId, &Output[Socket]{
				Port: &Port[Socket]{
					Id:                  PortId(outputValuePort["id"].(string)),
					Label:               outputValuePort["label"].(string),
					Index:               int(outputValuePort["index"].(float64)),
					MultipleConnections: outputValuePort["multipleConnections"].(bool),
					Socket: Socket{
						Name: outputValuePort["socket"].(map[string]any)["name"].(string),
					},
				},
				Control: &Control{
					Id:    outputValueControl["id"].(string),
					Index: int(outputValueControl["index"].(float64)),
				},
				ShowControl: outputValueCpy["showControl"].(bool),
				Label:       outputValueCpy["label"].(string),
			})
		}

		for controlId, controlValue := range node.Controls {
			controlValueCpy := controlValue.(map[string]any)
			if controlValueCpy == nil {
				continue
			}

			controlValueCtrl := controlValueCpy["control"].(map[string]any)
			controlValueOpts := controlValueCpy["options"].(map[string]any)

			inlineNode.Controls.Add(controlId, &InputControl{
				Control: &Control{
					Id:    controlValueCtrl["id"].(string),
					Index: int(controlValueCtrl["index"].(float64)),
				},
				Type: InputControlType(controlValueCpy["type"].(string)),
				Options: &InputControlOptions{
					Readonly: ToPtr(controlValueOpts["readonly"].(bool)),
					Initial:  controlValueOpts["initial"],
				},
				Readonly: ToPtr(controlValueCpy["readonly"].(bool)),
				Value:    ToPtr(controlValueCpy["value"]),
			})
		}

		editor.Nodes[id] = inlineNode
	}

	for id, connection := range jsonData.Connections {
		var inlineConnection = &Connection[ConnectionBase]{
			E: ConnectionBase{
				ID:     id,
				Source: NodeId(connection.Base.Source),
				Target: NodeId(connection.Base.Target),
			},
			Source:       NodeId(connection.Source),
			SourceOutput: NodeId(connection.SourceOutput),
			Target:       NodeId(connection.Target),
			TargetInput:  NodeId(connection.TargetInput),
		}

		editor.Connections[id] = inlineConnection
	}

	return editor, nil
}

func NewJSONEditorData(input []byte) (*JSONEditorData, error) {
	if len(input) == 0 {
		return nil, errors.New("empty input given")
	}

	var editorData JSONEditorData

	err := json.Unmarshal(input, &editorData)
	if err != nil {
		return nil, err
	}

	return &editorData, nil
}
