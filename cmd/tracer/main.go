package main

import (
	"encoding/json"
	"fmt"
	evm_simulator "github.com/Arjxm/tracer/core/evm-simulator"
	"github.com/Arjxm/tracer/core/rpc"
	"github.com/Arjxm/tracer/core/tui"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
)

func main() {

	rpcClt := rpc.NewClient(1)
	sim, err := evm_simulator.NewSimulator(rpcClt)
	if err != nil {
		log.Fatal(err)
	}

	simulation := evm_simulator.TxSimulationReq{
		ChainId: 1,
		TxHash:  "0x0ca14589e6f2512282bfb1b0f49aed1b033e24be3a1c9a8df4327ebbc94aee65",
	}
	stateDB, err := state.New(types.EmptyRootHash, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)

	result, err := sim.Simulate(simulation, stateDB, nil)
	if err != nil {
		fmt.Println(err)
	}

	var trace tui.Event
	err = json.Unmarshal(result.Trace, &trace)
	if err != nil {

	}

	var content string
	for _, node := range trace {
		content += tui.DisplayTree(node, 0)
		content += "\n"
	}

	tui.Display(content)
}
