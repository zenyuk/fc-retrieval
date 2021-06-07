package register

import (
	"fmt"

	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

const (
	PrivKey = "015ed053eab6fdf18c03954373ff7f89089992017d56beb8b05305b19800d6afe0"
	PubKey  = "01047799f37b014564e23578447d718e5c70a786b0e4e58ca25cb2a086b822434594d910b9b8c0fcbfe9f4c2db321e874819e0614be5b57fbb5080accd69adb2eaad"
)

var gr GatewayRegister

func init() {
	gr = GatewayRegister{
		NodeID:              "AA",
		Address:             "AA",
		RegionCode:          "AA",
		NetworkInfoGateway:  "AA",
		NetworkInfoProvider: "AA",
		NetworkInfoClient:   "AA",
		NetworkInfoAdmin:    "AA",
		RootSigningKey:      PrivKey,
		SigningKey:          PrivKey,
	}
}

func ExampleGatewayGetNodeID() {
	fmt.Println(gr.GetNodeID())
	// Output: AA
}

func ExampleGatewayGetAddress() {
	fmt.Println(gr.GetAddress())
	// Output: AA
}

func ExampleGatewayGetRegionCode() {
	fmt.Println(gr.GetRegionCode())
	// Output: AA
}

func ExampleGatewayGetNetworkInfoGateway() {
	fmt.Println(gr.GetNetworkInfoGateway())
	// Output: AA
}

func ExampleGatewayGetNetworkInfoProvider() {
	fmt.Println(gr.GetNetworkInfoProvider())
	// Output: AA
}

func ExampleGatewayGetNetworkInfoClient() {
	fmt.Println(gr.GetNetworkInfoClient())
	// Output: AA
}

func ExampleGatewayGetNetworkInfoAdmin() {
	fmt.Println(gr.GetNetworkInfoAdmin())
	// Output: AA
}

func ExampleGatewayGetRootSigningKey() {
	v, err := gr.GetRootSigningKey()
	fmt.Println(v, err)
	// Output: <nil> Incorrect secp256k1 public key length: 32
}

func ExampleGatewayGetSigningKey() {
	v, err := gr.GetSigningKey()
	fmt.Println(v, err)
	// Output: <nil> Incorrect secp256k1 public key length: 32
}

func ExampleRegisterGateway() {
	err := gr.RegisterGateway("")
	fmt.Println(err)
	// Output:
	// Post "/registers/gateway": unsupported protocol scheme ""
}

func ExampleGetRegisteredGateways() {
	v, err := GetRegisteredGateways("")
	fmt.Println(v, err)
	// Output:
	// [] Get "/registers/gateway": unsupported protocol scheme ""
}

func to_be_discussed_ExampleGetGatewayByID_nil() {
	v, err := GetGatewayByID("", nil)
	fmt.Println(v, err)
	// Output:
	// {        } Get "/registers/gateway/00": unsupported protocol scheme ""
}

func ExampleGetGatewayByID_new() {
	nd := &nodeid.NodeID{}
	v, err := GetGatewayByID("", nd)
	fmt.Println(v, err)
	// Output:
	// {        } Get "/registers/gateway/00": unsupported protocol scheme ""
}

