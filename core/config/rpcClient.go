package config

import "fmt"

func GetRPCUrl(chainId uint64) (string, error) {
	rpcMap := map[uint64]string{
		1: "http://127.0.0.1:8545",
	}

	if url, ok := rpcMap[chainId]; ok {
		return url, nil
	}
	return "no url found", fmt.Errorf("chain id %d not supported", chainId)
}
