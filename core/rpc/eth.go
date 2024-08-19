package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Arjxm/tracer/core/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
)

type Client struct {
	RpcUrl string
}

type Request struct {
	ID      int           `json:"id"`
	JSONRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type Response struct {
	ID      int             `json:"id"`
	JSONRpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Err     *ErrResponse    `json:"error,omitempty"`
}

type ErrResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (e *ErrResponse) Error() string {
	return fmt.Sprintf(`{"code": "%d", "message": "%s"}`, e.Code, e.Message)
}

func NewClient(chainId uint64) *Client {
	rpcUrl, err := config.GetRPCUrl(chainId)
	if err != nil {
		panic(err)
	}
	return &Client{RpcUrl: rpcUrl}
}

func (c *Client) GetCode(address, blk string) ([]byte, error) {
	// try to convert block into number
	blkNumber, ok := new(big.Int).SetString(strings.TrimLeft(blk, "0x"), 16)
	if !ok || blkNumber.Cmp(big.NewInt(0)) <= 0 {
		blk = "latest"
	}

	params := []interface{}{
		address, blk,
	}

	rpcResp, err := rpcPost(c.RpcUrl, "eth_getCode", params)
	if err != nil {
		return nil, err
	}

	resultB, _ := rpcResp.Result.MarshalJSON()

	var result string
	err = json.Unmarshal(resultB, &result)
	if err != nil {
		return nil, err
	}

	return hexutil.MustDecode(result), nil
}

func (c *Client) GetStorageAt(address, position, blk string) (common.Hash, error) {
	blkNumber, ok := new(big.Int).SetString(strings.TrimLeft(blk, "0x"), 16)
	if !ok || blkNumber.Cmp(big.NewInt(0)) <= 0 {
		blk = "latest"
	}

	params := []interface{}{
		address, position, blk,
	}

	rpcResp, err := rpcPost(c.RpcUrl, "eth_getStorageAt", params)
	if err != nil {
		return common.Hash{}, err
	}

	resultB, _ := rpcResp.Result.MarshalJSON()

	var result string
	err = json.Unmarshal(resultB, &result)
	if err != nil {
		return common.Hash{}, err
	}

	return common.HexToHash(result), nil
}

func (c *Client) GetCodeAndStorageAt(address, position, blk string) ([]byte, common.Hash, error) {
	// fetch code and storage
	code, err := c.GetCode(address, blk)
	if err != nil {
		return nil, common.Hash{}, err
	}

	storage, err := c.GetStorageAt(address, position, blk)
	if err != nil {
		return nil, common.Hash{}, err
	}

	return code, storage, nil
}

func (c *Client) GetBalance(address, blk string) (*big.Int, error) {
	blkNumber, ok := new(big.Int).SetString(strings.TrimLeft(blk, "0x"), 16)
	if !ok || blkNumber.Cmp(big.NewInt(0)) <= 0 {
		blk = "latest"
	}

	params := []interface{}{
		address, blk,
	}

	rpcResp, err := rpcPost(c.RpcUrl, "eth_getBalance", params)
	if err != nil {
		return nil, err
	}

	resultB, _ := rpcResp.Result.MarshalJSON()

	var result string
	err = json.Unmarshal(resultB, &result)
	if err != nil {
		return nil, err
	}

	balance, ok := new(big.Int).SetString(result[2:], 16)
	if !ok {
		return nil, fmt.Errorf("invalid balance received in response: %s", result)
	}

	return balance, nil
}

func (c *Client) EstimateGas(from common.Address, to common.Address, value *big.Int, input []byte) (uint64, error) {
	params := []interface{}{
		map[string]interface{}{
			"from":  from.Hex(),
			"to":    to.Hex(),
			"value": (*hexutil.Big)(value).String(),
			"data":  hexutil.Bytes(input).String(),
		},
	}

	log.Printf("params: %v\n", params)

	rpcResp, err := rpcPost(c.RpcUrl, "eth_estimateGas", params)
	if err != nil {
		return 0, fmt.Errorf("RPC call failed: %w", err)
	}

	if rpcResp.Err != nil {
		return 0, fmt.Errorf("RPC error: %s", rpcResp.Err.Error())
	}

	var hexGas hexutil.Uint64
	err = json.Unmarshal(rpcResp.Result, &hexGas)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal gas estimate: %w", err)
	}

	return uint64(hexGas), nil
}

func (c *Client) GetTxByHash(hash string) (map[string]interface{}, error) {
	params := []interface{}{
		hash,
	}

	log.Printf("params: %v\n", params)

	rpcResp, err := rpcPost(c.RpcUrl, "eth_getTransactionByHash", params)
	if err != nil {
		return nil, fmt.Errorf("RPC call failed: %w", err)
	}

	if rpcResp.Err != nil {
		return nil, fmt.Errorf("RPC error: %s", rpcResp.Err.Error())
	}

	var result map[string]interface{}
	err = json.Unmarshal(rpcResp.Result, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal gas estimate: %w", err)
	}

	fmt.Printf("rpcResp: %v\n", result)
	return result, nil
}

func rpcPost(rpcRpcUrl, method string, params []interface{}) (*Response, error) {
	payload := Request{
		ID:      1,
		JSONRpc: "2.0",
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(data)

	resp, err := http.Post(rpcRpcUrl, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	err = json.Unmarshal(b, &result)

	return &result, err
}
