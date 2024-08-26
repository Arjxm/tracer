package evm_simulator

import (
	"github.com/Arjxm/tracer/core/rpc"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"testing"
)

func TestSimulate(t *testing.T) {

	rpcClt := rpc.NewClient(1)
	sim, err := NewSimulator(rpcClt)
	if err != nil {
		log.Fatal(err)
	}

	simulation := TxSimulationReq{
		ChainId: 1,
		TxHash:  "0x0ca14589e6f2512282bfb1b0f49aed1b033e24be3a1c9a8df4327ebbc94aee65",
	}
	stateDB, err := state.New(types.EmptyRootHash, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := sim.Simulate(simulation, stateDB, nil)
	if err != nil {
		t.Fatal(err)
	}
	//
	//log.Println("-----------------------------------------------------------")
	//// just log the returned value for now
	//log.Println(hexutil.Encode(result.ReturnedData))
	log.Println(result.GasUsed, result.GasLimit, result.Trace)

	//for _, l := range result.Record.AccessList {
	//	log.Println("ADDRESS: ", l.Address.Hex())
	//	for _, st := range l.StorageKeys {
	//		log.Println(st.Hex())
	//	}
	//}

	//codeLen := stateDB.GetCodeSize(contractAddr)
	//if codeLen == 0 {
	//	t.Fatal("code of contract is zero")
	//}

	// check state value

}
