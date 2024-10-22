package evm_simulator

import (
	"errors"
	"fmt"
	evm "github.com/Arjxm/tracer/core/evm"
	"github.com/Arjxm/tracer/core/evm/runtime"
	"github.com/Arjxm/tracer/core/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"math/big"
	"strconv"
	"sync"
)

type TxSimulationReq struct {
	ChainId uint64
	TxHash  string
}

type TxSimulation struct {
	From        common.Address
	To          common.Address
	BlockNumber *big.Int
	GasLimit    uint64
	GasPrice    *big.Int
	Value       *big.Int
	Input       []byte
	Code        []byte
	ChainId     uint64
}

type TxSimulationResult struct {
	GasUsed      uint64
	ReturnedData []byte
	GasLimit     uint64
	Trace        []byte
}

type Simulator struct {
	RpcClient *rpc.Client
}

func NewSimulator(RpcClient *rpc.Client) (*Simulator, error) {
	return &Simulator{RpcClient: RpcClient}, nil
}

func (s *Simulator) Simulate(simulationReq TxSimulationReq, stateDB *state.StateDB, recordInitializer *runtime.RecordToInitiateState) (*TxSimulationResult, error) {
	traceRecoder := NewCustomTracer()

	tx, _ := s.RpcClient.GetTxByHash(simulationReq.TxHash)
	if tx == nil {
		fmt.Printf("Error fetching transaction")
	}
	gasLimit, _ := strconv.ParseInt(tx["gas"].(string), 0, 64)
	gasPrice, _ := strconv.ParseInt(tx["gasPrice"].(string), 0, 64) // 0 allows automatic base detection (hex or decimal)
	value, _ := strconv.ParseInt(tx["value"].(string), 0, 64)       // Similar conversion for value

	simulation := TxSimulation{
		From: common.HexToAddress(tx["from"].(string)), // Convert string to common.Address
		To:   common.HexToAddress(tx["to"].(string)),
		BlockNumber: func() *big.Int {
			v, _ := new(big.Int).SetString(tx["blockNumber"].(string)[2:], 16)
			return v.Sub(v, big.NewInt(128))
		}(),
		GasLimit: uint64(gasLimit),
		GasPrice: big.NewInt(gasPrice),
		Value:    big.NewInt(value),
		Input:    hexutil.MustDecode(tx["input"].(string)),
		ChainId:  simulationReq.ChainId,
	}

	cfg := TxSimulationConfig(simulation, traceRecoder)

	var (
		blk     = ""
		err     error
		code    = simulation.Code
		balance = big.NewInt(0)
		wg      sync.WaitGroup
	)

	if simulation.BlockNumber.Cmp(big.NewInt(0)) > 0 {
		blk = "0x" + simulation.BlockNumber.Text(16)
	} else {
		// fetch latest block number
	}

	if len(code) == 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			codeSize := stateDB.GetCodeSize(simulation.To)
			if codeSize == 0 {
				// fetch code of address
				code, err = s.RpcClient.GetCode(simulation.To.Hex(), blk)
				if err != nil {
					return
				}
			} else {
				code = stateDB.GetCode(simulation.To)
			}
		}()
	}

	if simulation.Value.Cmp(big.NewInt(0)) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fromBalance := stateDB.GetBalance(simulation.From)
			if fromBalance.Cmp(common.U2560) <= 0 {
				balance, err = s.RpcClient.GetBalance(simulation.From.Hex(), blk)
				if err != nil {
					return
				}

				if balance.Cmp(simulation.Value) <= 0 {
					err = errors.New("insuficient balance to proceed with simulation")
					return
				}
			}
		}()
	}

	wg.Wait()
	if err != nil {
		return nil, err
	}

	var recordToInit *evm.RecordToInitiateState
	if recordInitializer != nil {
		recordToInit = &evm.RecordToInitiateState{
			AddressCodeSet:    recordInitializer.AddressCodeSet,
			AddressBalanceSet: recordInitializer.AddressBalanceSet,
			AddressStorageSet: recordInitializer.AddressStorageSet,
			AccessList:        recordInitializer.AccessList,
		}
	}

	gas, err := s.RpcClient.EstimateGas(simulation.From, simulation.To, simulation.Value, simulation.Input)
	if err != nil {
	} else {
		cfg.GasLimit = gas
	}

	result, err := runtime.Execute(simulation.To, balance, code, simulation.Input, cfg, stateDB, recordToInit)
	if err != nil {
		err = traceRecoder.SaveResultToJSON()
		if err != nil {
			return nil, err
		}
		return nil, err

	}
	err = traceRecoder.SaveResultToJSON()
	if err != nil {
		return nil, err
	}

	return &TxSimulationResult{
		ReturnedData: result.Ret,
		GasUsed:      result.GasUsed,
		Trace:        traceRecoder.GetResultFromJSON(),
	}, nil
}
