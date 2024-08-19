package config

import "fmt"

func GetRPCUrl(chainId uint64) (string, error) {
	rpcMap := map[uint64]string{
		1: "https://rpc.ankr.com/eth",
	}

	if url, ok := rpcMap[chainId]; ok {
		return url, nil
	}
	return "no url found", fmt.Errorf("chain id %d not supported", chainId)
}
