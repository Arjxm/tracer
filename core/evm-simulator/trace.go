package evm_simulator

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Arjxm/tracer/core/evm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"math/big"
)

type TracerEvent struct {
	OnEnter  *OnEnterEvent
	OpCodes  []*OpCodeEvent
	Children []*TracerEvent
	OnExit   *OnExitEvent
}

type OnEnterEvent struct {
	Depth int
	Type  string
	From  common.Address
	To    common.Address
	Input string
	Gas   uint64
	Value *big.Int
}

type OpCodeEvent struct {
	PC        uint64
	OpCode    string
	Gas       uint64
	Cost      uint64
	RData     []byte
	Depth     int
	Err       error
	Stack     []string
	Memory    []string
	ScopeData ScopeData
}

type OnExitEvent struct {
	Depth    int
	Output   string
	GasUsed  uint64
	Err      string
	Reverted bool
}

type ScopeData struct {
	MemoryData string
	StackData  []string
	Caller     string
	Address    string
	CallValue  string
	CallInput  string
}

type CustomTracer struct {
	Events         []*TracerEvent
	CurrentEvents  []*TracerEvent
	TraceCompleted bool
	JSONData       []byte
}

func NewCustomTracer() *CustomTracer {
	return &CustomTracer{
		Events:        make([]*TracerEvent, 0),
		CurrentEvents: make([]*TracerEvent, 0),
	}
}

func (t *CustomTracer) OnEnter(depth int, typ byte, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	callType := evm.OpToString(typ)
	event := &TracerEvent{
		OnEnter: &OnEnterEvent{
			Depth: depth,
			Type:  callType,
			From:  from,
			To:    to,
			Input: fmt.Sprintf("%x", input),
			Gas:   gas,
			Value: value,
		},
		OpCodes:  make([]*OpCodeEvent, 0),
		Children: make([]*TracerEvent, 0),
	}

	if len(t.CurrentEvents) > 0 {
		parent := t.CurrentEvents[len(t.CurrentEvents)-1]
		parent.Children = append(parent.Children, event)
	} else {
		t.Events = append(t.Events, event)
	}

	t.CurrentEvents = append(t.CurrentEvents, event)

	// fmt.Printf("OnEnter: Depth: %d, Type: %s, From: %s, To: %s, Input: %x, Gas: %d, Value: %s\n", depth, callType, from.String(), to.String(), input, gas, value.String())
}

func (t *CustomTracer) OnExit(depth int, output []byte, gasUsed uint64, err error, reverted bool) {
	if len(t.CurrentEvents) > 0 {
		event := t.CurrentEvents[len(t.CurrentEvents)-1]
		event.OnExit = &OnExitEvent{
			Depth:    depth,
			Output:   fmt.Sprintf("%x", output),
			GasUsed:  gasUsed,
			Err:      fmt.Sprintf("%v", err),
			Reverted: reverted,
		}
		t.CurrentEvents = t.CurrentEvents[:len(t.CurrentEvents)-1]
	}

	// fmt.Printf("OnExit: Depth: %d, Output: %x, GasUsed: %d, Err: %v, Reverted: %v\n", depth, output, gasUsed, err, reverted)
}

func (t *CustomTracer) OnOpcode(pc uint64, op byte, gas, cost uint64, scope tracing.OpContext, rData []byte, depth int, err error) {
	opCode := evm.OpToString(op)

	if len(t.CurrentEvents) > 0 {
		event := t.CurrentEvents[len(t.CurrentEvents)-1]

		scopeData := &ScopeData{
			Caller:    scope.Caller().String(),
			Address:   scope.Address().String(),
			CallValue: fmt.Sprintf("%x", scope.CallValue()),
			CallInput: fmt.Sprintf("%x", scope.CallInput()),
		}

		stack := scope.StackData()
		stackData := make([]string, len(stack))
		for i := len(stack) - 1; i >= 0; i-- {
			stackData[len(stack)-i-1] = stack[i].Hex()
		}

		memory := scope.MemoryData()
		memoryData := make([]string, 0)
		for i := 0; i < len(memory); i += 32 {
			end := i + 32
			if end > len(memory) {
				end = len(memory)
			}
			memoryData = append(memoryData, fmt.Sprintf("0x%04x: %s", i, hex.EncodeToString(memory[i:end])))
		}

		event.OpCodes = append(event.OpCodes, &OpCodeEvent{
			PC:        pc,
			OpCode:    opCode,
			Gas:       gas,
			Cost:      cost,
			RData:     rData,
			Depth:     depth,
			Err:       err,
			ScopeData: *scopeData,
			Stack:     stackData,
			Memory:    memoryData,
		})
	}

}

func (t *CustomTracer) OnFault(pc uint64, op byte, gas, cost uint64, scope tracing.OpContext, depth int, err error) {
	// Add your implementation here if needed
}

func (t *CustomTracer) SaveResultToJSON() error {
	data, err := json.MarshalIndent(t.Events, "", "  ")
	if err != nil {
		return err
	}

	t.JSONData = data

	return nil
}

func (t *CustomTracer) GetResultFromJSON() []byte {
	if len(t.JSONData) == 0 {
		return nil
	}

	err := json.Unmarshal(t.JSONData, &t.Events)
	if err != nil {
		return nil
	}

	return t.JSONData
}
