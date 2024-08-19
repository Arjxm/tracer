package decoder

type Event []Node

type Node struct {
	OnEnter  map[string]interface{} `json:"OnEnter"`
	Children []Node                 `json:"Children"`
	OnExit   map[string]interface{} `json:"OnExit"`
	OpCodes  []interface{}          `json:"OpCodes"`
}

//func decode() {
//	abiFragments, err := ParseContractABI(ABI)
//	fun := GetFuncFragment(abiFragments)
//
//	if err != nil {
//		fmt.Println("Error parsing contract ABI:", err)
//		return
//	}
//
//}

//func GetDecodeFunction(abiFragments []ABIFragment, e []interface{}) {
//	abiFragments, err := ParseContractABI(ABI)
//	fun := GetFuncFragment(abiFragments)
//
//	if err != nil {
//		fmt.Println("Error parsing contract ABI:", err)
//		return
//	}
//}
