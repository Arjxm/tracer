package decoder

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"io"
	"net/http"
)

type ABIFragment struct {
	Signature       string
	Name            string
	Type            string
	Inputs          []abi.ArgumentMarshaling
	Outputs         []abi.ArgumentMarshaling
	Constant        bool
	Payable         bool
	Anonymous       bool
	StateMutability string
}

const ABI = `[{"inputs":[{"internalType":"address","name":"_WETH","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"clientData","type":"bytes"}],"name":"ClientData","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"reason","type":"string"}],"name":"Error","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"pair","type":"address"},{"indexed":false,"internalType":"uint256","name":"amountOut","type":"uint256"},{"indexed":false,"internalType":"address","name":"output","type":"address"}],"name":"Exchange","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"token","type":"address"},{"indexed":false,"internalType":"uint256","name":"totalAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"totalFee","type":"uint256"},{"indexed":false,"internalType":"address[]","name":"recipients","type":"address[]"},{"indexed":false,"internalType":"uint256[]","name":"amounts","type":"uint256[]"},{"indexed":false,"internalType":"bool","name":"isBps","type":"bool"}],"name":"Fee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"contract IERC20","name":"srcToken","type":"address"},{"indexed":false,"internalType":"contract IERC20","name":"dstToken","type":"address"},{"indexed":false,"internalType":"address","name":"dstReceiver","type":"address"},{"indexed":false,"internalType":"uint256","name":"spentAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"returnAmount","type":"uint256"}],"name":"Swapped","type":"event"},{"inputs":[],"name":"WETH","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"isWhitelist","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"token","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"rescueFunds","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"callTarget","type":"address"},{"internalType":"address","name":"approveTarget","type":"address"},{"internalType":"bytes","name":"targetData","type":"bytes"},{"components":[{"internalType":"contract IERC20","name":"srcToken","type":"address"},{"internalType":"contract IERC20","name":"dstToken","type":"address"},{"internalType":"address[]","name":"srcReceivers","type":"address[]"},{"internalType":"uint256[]","name":"srcAmounts","type":"uint256[]"},{"internalType":"address[]","name":"feeReceivers","type":"address[]"},{"internalType":"uint256[]","name":"feeAmounts","type":"uint256[]"},{"internalType":"address","name":"dstReceiver","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"minReturnAmount","type":"uint256"},{"internalType":"uint256","name":"flags","type":"uint256"},{"internalType":"bytes","name":"permit","type":"bytes"}],"internalType":"struct MetaAggregationRouterV2.SwapDescriptionV2","name":"desc","type":"tuple"},{"internalType":"bytes","name":"clientData","type":"bytes"}],"internalType":"struct MetaAggregationRouterV2.SwapExecutionParams","name":"execution","type":"tuple"}],"name":"swap","outputs":[{"internalType":"uint256","name":"returnAmount","type":"uint256"},{"internalType":"uint256","name":"gasUsed","type":"uint256"}],"stateMutability":"payable","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"callTarget","type":"address"},{"internalType":"address","name":"approveTarget","type":"address"},{"internalType":"bytes","name":"targetData","type":"bytes"},{"components":[{"internalType":"contract IERC20","name":"srcToken","type":"address"},{"internalType":"contract IERC20","name":"dstToken","type":"address"},{"internalType":"address[]","name":"srcReceivers","type":"address[]"},{"internalType":"uint256[]","name":"srcAmounts","type":"uint256[]"},{"internalType":"address[]","name":"feeReceivers","type":"address[]"},{"internalType":"uint256[]","name":"feeAmounts","type":"uint256[]"},{"internalType":"address","name":"dstReceiver","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"minReturnAmount","type":"uint256"},{"internalType":"uint256","name":"flags","type":"uint256"},{"internalType":"bytes","name":"permit","type":"bytes"}],"internalType":"struct MetaAggregationRouterV2.SwapDescriptionV2","name":"desc","type":"tuple"},{"internalType":"bytes","name":"clientData","type":"bytes"}],"internalType":"struct MetaAggregationRouterV2.SwapExecutionParams","name":"execution","type":"tuple"}],"name":"swapGeneric","outputs":[{"internalType":"uint256","name":"returnAmount","type":"uint256"},{"internalType":"uint256","name":"gasUsed","type":"uint256"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"contract IAggregationExecutor","name":"caller","type":"address"},{"components":[{"internalType":"contract IERC20","name":"srcToken","type":"address"},{"internalType":"contract IERC20","name":"dstToken","type":"address"},{"internalType":"address[]","name":"srcReceivers","type":"address[]"},{"internalType":"uint256[]","name":"srcAmounts","type":"uint256[]"},{"internalType":"address[]","name":"feeReceivers","type":"address[]"},{"internalType":"uint256[]","name":"feeAmounts","type":"uint256[]"},{"internalType":"address","name":"dstReceiver","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"minReturnAmount","type":"uint256"},{"internalType":"uint256","name":"flags","type":"uint256"},{"internalType":"bytes","name":"permit","type":"bytes"}],"internalType":"struct MetaAggregationRouterV2.SwapDescriptionV2","name":"desc","type":"tuple"},{"internalType":"bytes","name":"executorData","type":"bytes"},{"internalType":"bytes","name":"clientData","type":"bytes"}],"name":"swapSimpleMode","outputs":[{"internalType":"uint256","name":"returnAmount","type":"uint256"},{"internalType":"uint256","name":"gasUsed","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"addr","type":"address[]"},{"internalType":"bool[]","name":"value","type":"bool[]"}],"name":"updateWhitelist","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]`

//func ABI()  {
//	var abi interface{}
//	abi =[{"inputs":[{"internalType":"address","name":"_WETH","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"clientData","type":"bytes"}],"name":"ClientData","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"reason","type":"string"}],"name":"Error","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"pair","type":"address"},{"indexed":false,"internalType":"uint256","name":"amountOut","type":"uint256"},{"indexed":false,"internalType":"address","name":"output","type":"address"}],"name":"Exchange","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"token","type":"address"},{"indexed":false,"internalType":"uint256","name":"totalAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"totalFee","type":"uint256"},{"indexed":false,"internalType":"address[]","name":"recipients","type":"address[]"},{"indexed":false,"internalType":"uint256[]","name":"amounts","type":"uint256[]"},{"indexed":false,"internalType":"bool","name":"isBps","type":"bool"}],"name":"Fee","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"contract IERC20","name":"srcToken","type":"address"},{"indexed":false,"internalType":"contract IERC20","name":"dstToken","type":"address"},{"indexed":false,"internalType":"address","name":"dstReceiver","type":"address"},{"indexed":false,"internalType":"uint256","name":"spentAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"returnAmount","type":"uint256"}],"name":"Swapped","type":"event"},{"inputs":[],"name":"WETH","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"isWhitelist","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"token","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"rescueFunds","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"callTarget","type":"address"},{"internalType":"address","name":"approveTarget","type":"address"},{"internalType":"bytes","name":"targetData","type":"bytes"},{"components":[{"internalType":"contract IERC20","name":"srcToken","type":"address"},{"internalType":"contract IERC20","name":"dstToken","type":"address"},{"internalType":"address[]","name":"srcReceivers","type":"address[]"},{"internalType":"uint256[]","name":"srcAmounts","type":"uint256[]"},{"internalType":"address[]","name":"feeReceivers","type":"address[]"},{"internalType":"uint256[]","name":"feeAmounts","type":"uint256[]"},{"internalType":"address","name":"dstReceiver","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"minReturnAmount","type":"uint256"},{"internalType":"uint256","name":"flags","type":"uint256"},{"internalType":"bytes","name":"permit","type":"bytes"}],"internalType":"struct MetaAggregationRouterV2.SwapDescriptionV2","name":"desc","type":"tuple"},{"internalType":"bytes","name":"clientData","type":"bytes"}],"internalType":"struct MetaAggregationRouterV2.SwapExecutionParams","name":"execution","type":"tuple"}],"name":"swap","outputs":[{"internalType":"uint256","name":"returnAmount","type":"uint256"},{"internalType":"uint256","name":"gasUsed","type":"uint256"}],"stateMutability":"payable","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"callTarget","type":"address"},{"internalType":"address","name":"approveTarget","type":"address"},{"internalType":"bytes","name":"targetData","type":"bytes"},{"components":[{"internalType":"contract IERC20","name":"srcToken","type":"address"},{"internalType":"contract IERC20","name":"dstToken","type":"address"},{"internalType":"address[]","name":"srcReceivers","type":"address[]"},{"internalType":"uint256[]","name":"srcAmounts","type":"uint256[]"},{"internalType":"address[]","name":"feeReceivers","type":"address[]"},{"internalType":"uint256[]","name":"feeAmounts","type":"uint256[]"},{"internalType":"address","name":"dstReceiver","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"minReturnAmount","type":"uint256"},{"internalType":"uint256","name":"flags","type":"uint256"},{"internalType":"bytes","name":"permit","type":"bytes"}],"internalType":"struct MetaAggregationRouterV2.SwapDescriptionV2","name":"desc","type":"tuple"},{"internalType":"bytes","name":"clientData","type":"bytes"}],"internalType":"struct MetaAggregationRouterV2.SwapExecutionParams","name":"execution","type":"tuple"}],"name":"swapGeneric","outputs":[{"internalType":"uint256","name":"returnAmount","type":"uint256"},{"internalType":"uint256","name":"gasUsed","type":"uint256"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"contract IAggregationExecutor","name":"caller","type":"address"},{"components":[{"internalType":"contract IERC20","name":"srcToken","type":"address"},{"internalType":"contract IERC20","name":"dstToken","type":"address"},{"internalType":"address[]","name":"srcReceivers","type":"address[]"},{"internalType":"uint256[]","name":"srcAmounts","type":"uint256[]"},{"internalType":"address[]","name":"feeReceivers","type":"address[]"},{"internalType":"uint256[]","name":"feeAmounts","type":"uint256[]"},{"internalType":"address","name":"dstReceiver","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"minReturnAmount","type":"uint256"},{"internalType":"uint256","name":"flags","type":"uint256"},{"internalType":"bytes","name":"permit","type":"bytes"}],"internalType":"struct MetaAggregationRouterV2.SwapDescriptionV2","name":"desc","type":"tuple"},{"internalType":"bytes","name":"executorData","type":"bytes"},{"internalType":"bytes","name":"clientData","type":"bytes"}],"name":"swapSimpleMode","outputs":[{"internalType":"uint256","name":"returnAmount","type":"uint256"},{"internalType":"uint256","name":"gasUsed","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"addr","type":"address[]"},{"internalType":"bool[]","name":"value","type":"bool[]"}],"name":"updateWhitelist","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]
//}

func ParseContractABI(abiString string) ([]ABIFragment, error) {
	var abiFragments []ABIFragment
	err := json.Unmarshal([]byte(abiString), &abiFragments)
	if err != nil {
		return nil, err
	}
	return abiFragments, nil
}

func getFunctionSignature(selector string) (map[string]interface{}, error) {
	// Create HTTP client with timeout
	client := &http.Client{}

	// Create channels for results and errors
	resChan := make(chan map[string]interface{}, 2)
	errChan := make(chan error, 2)

	// Goroutine for 4byte.directory API
	go func() {
		url := fmt.Sprintf("https://www.4byte.directory/api/v1/signatures/?hex_signature=%s", selector)
		resp, err := client.Get(url)
		if err == nil {
			defer resp.Body.Close()
			var body map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			if err == nil {
				results := body["results"].([]interface{})
				if len(results) > 0 {
					resChan <- body
					return
				}
			}
		}
		errChan <- err
	}()

	// Goroutine for Openchain API
	go func() {
		url := fmt.Sprintf("https://api.openchain.xyz/signature-database/v1/lookup?function=%s&filter=true", selector)
		resp, err := client.Get(url)
		if err == nil {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(resp.Body)
			var body map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&body)
			if err == nil {
				resChan <- body
				return
			}
		}
		errChan <- err
	}()

	// Standardize the response format
	standardizeResponse := func(resp map[string]interface{}) map[string]interface{} {
		if results, ok := resp["results"]; ok {
			return map[string]interface{}{"results": results}
		}
		if result, ok := resp["result"]; ok {
			if function, ok := result.(map[string]interface{})["function"]; ok {
				if signatures, ok := function.(map[string]interface{})[selector]; ok {
					return map[string]interface{}{"results": signatures}
				}
			}
		}
		return nil
	}

	// Wait for results or errors from both goroutines
	for i := 0; i < 2; i++ {
		select {
		case res := <-resChan:
			if standardRes := standardizeResponse(res); standardRes != nil {
				return standardRes, nil
			}
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
		}
	}

	// Return an error if no results are found
	return nil, fmt.Errorf("no function signature found for selector: %s", selector)
}
