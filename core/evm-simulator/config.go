package evm_simulator

import (
	"fmt"
	"github.com/Arjxm/tracer/core/evm/runtime"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/vm"
	"math/big"
)

func vmConfig(tracer *CustomTracer) vm.Config {
	return vm.Config{
		Tracer: &tracing.Hooks{
			OnEnter: tracer.OnEnter,
			OnExit:  tracer.OnExit,
			OnFault: func(pc uint64, op byte, gas, cost uint64, scope tracing.OpContext, depth int, err error) {
				fmt.Printf("OnFault: PC: %d, OpCode: 0x%02x, Gas: %d, Cost: %d, Depth: %d, Err: %v\n", pc, op, gas, cost, depth, err)
			},
			//OnOpcode: tracer.OnOpcode,
			//OnBalanceChange: func(addr common.Address, prev, new *big.Int, reason tracing.BalanceChangeReason) {
			//	fmt.Printf("OnBalanceChange: Address: %s, Prev: %s, New: %s, Reason: %s\n", addr.Hex(), prev.String(), new.String(), reason)
			//},
			//
			//OnGasChange: func(old, new uint64, reason tracing.GasChangeReason) {
			//	fmt.Printf("OnGasChange: Old: %d, New: %d", old, new)
			//},
			//OnNonceChange: func(addr common.Address, prev, new uint64) {
			//	fmt.Printf("OnNonceChange: Address: %s, Prev: %d, New: %d\n", addr.Hex(), prev, new)
			//},
			//OnCodeChange: func(addr common.Address, prevCodeHash common.Hash, prevCode []byte, codeHash common.Hash, code []byte) {
			//	fmt.Printf("OnCodeChange: Address: %s, PrevCodeHash: %s, CodeHash: %s\n", addr.Hex(), prevCodeHash.Hex(), codeHash.Hex())
			//},
			//OnStorageChange: func(addr common.Address, slot common.Hash, prev, new common.Hash) {
			//	fmt.Printf("OnStorageChange: Address: %s, Slot: %s, Prev: %s, New: %s\n", addr.Hex(), slot.Hex(), prev.Hex(), new.Hex())
			//},
			//OnLog: func(log *types.Log) {
			//	fmt.Printf("OnLog: Address: %s, Topics: %v, Data: %s\n", log.Address.Hex(), log.Topics, hex.EncodeToString(log.Data))
			//},
		},
	}
}

func TxSimulationConfig(simulation TxSimulation, tracerRecord *CustomTracer) *runtime.Config {
	return &runtime.Config{
		Debug:       true,
		Origin:      simulation.From,
		BlockNumber: simulation.BlockNumber,
		GasLimit:    simulation.GasLimit,
		GasPrice:    simulation.GasPrice,
		Value:       simulation.Value,
		ChainId:     simulation.ChainId,
		EVMConfig:   vmConfig(tracerRecord),
		Time:        1723484999,
		BaseFee:     big.NewInt(3310633170),
	}
}
