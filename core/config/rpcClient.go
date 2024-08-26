package config

import "fmt"

func GetRPCUrl(chainId uint64) (string, error) {
	rpcMap := map[uint64]string{
		1: "https://wider-maximum-vineyard.quiknode.pro/a389516789aeeed43be7b300fe02630d452c0356",
	}

	if url, ok := rpcMap[chainId]; ok {
		return url, nil
	}
	return "no url found", fmt.Errorf("chain id %d not supported", chainId)
}
