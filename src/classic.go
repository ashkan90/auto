package src

import (
	"encoding/json"
	"github.com/ashkan90/auto-core/utils"
	"log"
	"sync"
)

type PortId string

type Socket struct {
	Name string `json:"name"`
}

func NewSocket(name string) Socket {
	return Socket{
		Name: name,
	}
}

type PortInterface interface {
	GetId() string
}

type Port[S Socket] struct {
	Id                  PortId `json:"id"`
	Label               string `json:"label"`
	Index               int    `json:"index"`
	MultipleConnections bool   `json:"multipleConnections"`
	Socket              S      `json:"socket"`
}

func (p *Port[S]) GetId() string {
	return string(p.Id)
}

func NewPort[S Socket](socket S, label string, multipleConnections bool) *Port[S] {
	return &Port[S]{
		Id:                  PortId(GetUID()),
		Label:               label,
		MultipleConnections: multipleConnections,
		Socket:              socket,
	}
}

type Input[S Socket] struct {
	// Port instance
	Port PortInterface `json:"port"`
	// Control instance
	Control ControlInterface `json:"control"`
	// ShowControl Whether the control is visible. Can be managed dynamically by extensions. Default is `true`
	ShowControl bool `json:"showControl"`
	// Label input label
	Label string `json:"label"`
}

type InputInterface interface {
	GetId() string
}

func (i *Input[S]) GetId() string {
	return i.Port.GetId()
}

func NewInput[S Socket](socket S, label string, multipleConnections bool) *Input[S] {
	return &Input[S]{
		Port:        NewPort[S](socket, label, multipleConnections),
		Control:     NewControl(),
		Label:       label,
		ShowControl: true,
	}
}

func (i *Input[S]) AddControl(control *Control) {
	if i.Control != nil {
		log.Panic("control already added")
	}
	i.Control = control
}

func (i *Input[S]) RemoveControl() {
	i.Control = nil
}

type Output[S Socket] struct {
	// Port instance
	Port PortInterface `json:"port"`
	// Control instance
	Control ControlInterface `json:"control"`
	// ShowControl Whether the control is visible. Can be managed dynamically by extensions. Default is `true`
	ShowControl bool `json:"showControl"`
	// Label input label
	Label string `json:"label"`
}

func (o *Output[S]) GetId() string {
	return o.Port.GetId()
}

func NewOutput[S Socket](socket S, label string, multipleConnections bool) PortInterface {
	return &Output[S]{
		Port:        NewPort[S](socket, label, multipleConnections),
		Control:     NewControl(),
		Label:       label,
		ShowControl: true,
	}
}

type ControlInterface interface {
	GetId() string
	GetValue() any
}

type Control struct {
	// Control Id unique string generated by `getUID` function
	Id string `json:"id"`
	// Control Index used for sorting controls. Default is `0`
	Index int `json:"index"`
}

func NewControl() ControlInterface {
	return &Control{
		Id: GetUID(),
	}
}

func (c *Control) GetId() string {
	return c.Id
}

func (c *Control) GetValue() any {
	return nil
}

type InputControlOptions struct {
	Readonly *bool           `json:"readonly"`
	Initial  any             `json:"initial"`
	Change   func(value any) `json:"-"`
}

type InputControlType string

const (
	InputControlText   InputControlType = "text"
	InputControlNumber InputControlType = "number"
)

type InputControl struct {
	Control  ControlInterface     `json:"control"`
	Type     InputControlType     `json:"type"`
	Options  *InputControlOptions `json:"options"`
	Readonly *bool                `json:"readonly"`
	Value    any                  `json:"value"`
}

func NewInputControl(_type InputControlType, opt *InputControlOptions) *InputControl {
	return &InputControl{
		Control:  NewControl(),
		Type:     _type,
		Options:  opt,
		Value:    opt.Initial,
		Readonly: opt.Readonly,
	}
}

func (ic *InputControl) GetId() string {
	return ic.Control.GetId()
}

func (ic *InputControl) GetValue() any {
	return ic.Value
}

func (ic *InputControl) SetValue(value any) {
	ic.Value = &value
	if ic.Options.Change != nil {
		ic.Options.Change(value)
	}
}

type Node[Base NodeBase] struct {
	E        Base           `json:"base"`
	Inputs   *utils.SyncMap `json:"inputs"`
	Outputs  *utils.SyncMap `json:"outputs"`
	Controls *utils.SyncMap `json:"controls"`
	Selected *bool          `json:"selected"`
	mu       *sync.Mutex
}

type NodeInterface interface {
	NodeExecutor
	NodeData
	Node() *Node[NodeBase]
	FromModule() bool

	HasInput(k string) bool
	AddInput(k string, input InputInterface)
	RemoveInput(k string)

	HasOutput(k string) bool
	AddOutput(k string, output PortInterface)
	RemoveOutput(k string)

	HasControl(k string) bool
	AddControl(k string, control ControlInterface)
	RemoveControl(k string)
}

func (n *Node[Base]) String() string {
	strs, _ := json.Marshal(n)
	return string(strs)
}

func (n *Node[Base]) HasInput(k string) bool {
	_, ok := n.Inputs.Get(k)
	return ok
}

func (n *Node[Base]) AddInput(k string, input InputInterface) {
	n.Inputs.Add(k, input)
}

func (n *Node[Base]) RemoveInput(k string) {
	n.Inputs.Delete(k)
}

func (n *Node[Base]) HasOutput(k string) bool {
	_, ok := n.Outputs.Get(k)
	return ok
}

func (n *Node[Base]) AddOutput(k string, output PortInterface) {
	n.Outputs.Add(k, output)
}

func (n *Node[Base]) RemoveOutput(k string) {
	n.Outputs.Delete(k)
}

func (n *Node[Base]) HasControl(k string) bool {
	_, ok := n.Controls.Get(k)
	return ok
}

func (n *Node[Base]) AddControl(k string, control ControlInterface) {
	n.Controls.Add(k, control)
}

func (n *Node[Base]) RemoveControl(k string) {
	n.Controls.Delete(k)
}

func (n *Node[Base]) FromModule() bool {
	return false
}

func (n *Node[Base]) Node() *Node[Base] {
	return n
}

func NewNode() *Node[NodeBase] {
	return &Node[NodeBase]{
		E:        NodeBase{ID: NodeId(GetUID())},
		Inputs:   utils.NewSyncMap(),
		Outputs:  utils.NewSyncMap(),
		Controls: utils.NewSyncMap(),
		mu:       &sync.Mutex{},
	}
}

func (n *Node[NodeBase]) Execute(_ string, forward func(output string)) {
	forward("exec")
}

func (n *Node[NodeBase]) Data(inputs func() map[string]any) map[string]any {
	return inputs()
}

type Connection[Base ConnectionBase] struct {
	E            Base   `json:"base"`
	Source       NodeId `json:"source"`
	SourceOutput NodeId `json:"sourceOutput"`
	Target       NodeId `json:"target"`
	TargetInput  NodeId `json:"targetInput"`
}

func NewConnection(source NodeInterface, sourceOutput NodeId, target NodeInterface, targetInput NodeId) *Connection[ConnectionBase] {
	var sourceNode, targetNode = source.Node(), target.Node()
	if _, ok := sourceNode.Outputs.Get(string(sourceOutput)); !ok {
		log.Panicf("source node doesn't have output with key %s", sourceOutput)
	}
	if _, ok := targetNode.Inputs.Get(string(targetInput)); !ok {
		log.Panicf("target node doesn't have input with key %s", targetInput)
	}

	return &Connection[ConnectionBase]{
		E:            ConnectionBase{ID: ConnectionId(GetUID())},
		Source:       sourceNode.E.ID,
		SourceOutput: sourceOutput,
		Target:       targetNode.E.ID,
		TargetInput:  targetInput,
	}
}
