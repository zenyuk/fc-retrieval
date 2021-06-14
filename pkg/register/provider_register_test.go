package register

import (
	"fmt"

	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

var pr ProviderRegister

func init() {
	pr = ProviderRegister{
		NodeID:             "AA",
		Address:            "AA",
		RegionCode:         "AA",
		NetworkInfoGateway: "AA",
		NetworkInfoClient:  "AA",
		NetworkInfoAdmin:   "AA",
		RootSigningKey:     PrivKey,
		SigningKey:         PrivKey,
	}
}

func ExampleProviderGetNodeID() {
	fmt.Println(pr.GetNodeID())
	// Output: AA
}

func ExampleProviderGetAddress() {
	fmt.Println(pr.GetAddress())
	// Output: AA
}

func ExampleProviderGetRegionCode() {
	fmt.Println(pr.GetRegionCode())
	// Output: AA
}

func ExampleProviderGetNetworkInfoGateway() {
	fmt.Println(pr.GetNetworkInfoGateway())
	// Output: AA
}

func ExampleProviderGetNetworkInfoProvider() {
	fmt.Println(pr.GetNetworkInfoProvider())
	// Output:
}

func ExampleProviderGetNetworkInfoClient() {
	fmt.Println(pr.GetNetworkInfoClient())
	// Output: AA
}

func ExampleProviderGetNetworkInfoAdmin() {
	fmt.Println(pr.GetNetworkInfoAdmin())
	// Output: AA
}

func ExampleProviderGetRootSigningKey() {
	v, err := pr.GetRootSigningKey()
	fmt.Println(v, err)
	// Output: <nil> incorrect secp256k1 public key length: 32
}

func ExampleProviderGetSigningKey() {
	v, err := pr.GetSigningKey()
	fmt.Println(v, err)
	// Output: <nil> incorrect secp256k1 public key length: 32
}

func ExampleRegisterProvider() {
	err := pr.RegisterProvider("")
	fmt.Println(err)
	// Output:
	// Post "/registers/provider": unsupported protocol scheme ""
}

func ExampleGetRegisteredProviders() {
	v, err := GetRegisteredProviders("")
	fmt.Println(v, err)
	// Output:
	// [] Get "/registers/provider": unsupported protocol scheme ""
}

func to_be_discussed_ExampleGetProviderByID_nil() {
	v, err := GetProviderByID("", nil)
	fmt.Println(v, err)
	// Output:
	// {       } Get "/registers/provider/00": unsupported protocol scheme ""

}

func ExampleGetProviderByID_new() {
	nd := &nodeid.NodeID{}
	v, err := GetProviderByID("", nd)
	fmt.Println(v, err)
	// Output:
	// {       } Get "/registers/provider/00": unsupported protocol scheme ""

}

