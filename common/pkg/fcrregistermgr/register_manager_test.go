package fcrregistermgr

import (
  "encoding/json"
  "reflect"
  "sync"
  "testing"
  "time"

  "github.com/golang/mock/gomock"
  "github.com/stretchr/testify/assert"

  "github.com/ConsenSys/fc-retrieval/common/pkg/mocks"
  "github.com/ConsenSys/fc-retrieval/common/pkg/nodeid"
  "github.com/ConsenSys/fc-retrieval/common/pkg/register"
)
const fakeRegisterAPIURL = "fakeRegisterAPIURL"
var rm = NewFCRRegisterMgr(fakeRegisterAPIURL, false, false, 1*time.Second)

func TestFCRRegisterMgr_GetAllGateways(t *testing.T) {
  nodeID00, _ := nodeid.NewNodeIDFromHexString("00")
  nodeID01, _ := nodeid.NewNodeIDFromHexString("01")
  nodeID02, _ := nodeid.NewNodeIDFromHexString("02")
  nodeID5A, _ := nodeid.NewNodeIDFromHexString("5A")
  nodeIDFFFF, _ := nodeid.NewNodeIDFromHexString("FFFF")

  gr := map[string]register.GatewayRegistrar{}
  gr[nodeID00.ToString()] = register.NewGatewayRegister(nodeID00.ToString(),"", "", "", "", "", "", "", "")
  gr[nodeID01.ToString()] = register.NewGatewayRegister(nodeID01.ToString(),"", "", "", "", "", "", "", "")
  gr[nodeID02.ToString()] = register.NewGatewayRegister(nodeID02.ToString(),"", "", "", "", "", "", "", "")
  gr[nodeID5A.ToString()] = register.NewGatewayRegister(nodeID5A.ToString(),"", "", "", "", "", "", "", "")
  gr[nodeIDFFFF.ToString()] = register.NewGatewayRegister(nodeIDFFFF.ToString(), "", "", "", "", "", "", "", "")

  type fields struct {
    start                     bool
    gatewayDiscv              bool
    registeredGatewaysMap     map[string]register.GatewayRegistrar
    registeredGatewaysMapLock sync.RWMutex
  }
  tests := []struct {
    name   string
    fields fields

    want []register.GatewayRegistrar
  }{
    {name: "getGateway not started returns nil",
      fields: fields{
        start:                     false,
        gatewayDiscv:              false,
        registeredGatewaysMap:     gr,
        registeredGatewaysMapLock: sync.RWMutex{},
      },

      want: nil,
    },
    {name: "getGateway - gatewayDiscv not started returns nil",
      fields: fields{
        start:                     true,
        gatewayDiscv:              false,
        registeredGatewaysMap:     gr,
        registeredGatewaysMapLock: sync.RWMutex{},
      },

      want: nil,
    },
    {name: "getGateway - started return nodeID00",
      fields: fields{
        start:                     true,
        gatewayDiscv:              true,
        registeredGatewaysMap:     gr,
        registeredGatewaysMapLock: sync.RWMutex{},
      },

      want : []register.GatewayRegistrar{
        register.NewGatewayRegister(nodeID00.ToString(),"", "", "", "", "", "", "", ""),
        register.NewGatewayRegister(nodeID01.ToString(),"", "", "", "", "", "", "", ""),
        register.NewGatewayRegister(nodeID02.ToString(),"", "", "", "", "", "", "", ""),
        register.NewGatewayRegister(nodeID5A.ToString(),"", "", "", "", "", "", "", ""),
        register.NewGatewayRegister(nodeIDFFFF.ToString(), "", "", "", "", "", "", "", ""),
      },
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mgr := &FCRRegisterMgr{
        start:                     tt.fields.start,
        gatewayDiscv:              tt.fields.gatewayDiscv,
        registeredGatewaysMap:     tt.fields.registeredGatewaysMap,
        registeredGatewaysMapLock: tt.fields.registeredGatewaysMapLock,
      }
      got := mgr.GetAllGateways()
      assert.ElementsMatch(t, got, tt.want)
    })
  }
}

func TestFCRRegisterMgr_GetGateway(t *testing.T) {

	nodeID00, _ := nodeid.NewNodeIDFromHexString("00")
	nodeID01, _ := nodeid.NewNodeIDFromHexString("01")
	nodeID02, _ := nodeid.NewNodeIDFromHexString("02")
	nodeID5A, _ := nodeid.NewNodeIDFromHexString("5A")
	nodeIDFFFF, _ := nodeid.NewNodeIDFromHexString("FFFF")

	gr := map[string]register.GatewayRegistrar{}
  gr[nodeID00.ToString()] = register.NewGatewayRegister(
    nodeID00.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )
  gr[nodeID01.ToString()] = register.NewGatewayRegister(
    nodeID01.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )
  gr[nodeID02.ToString()] = register.NewGatewayRegister(
    nodeID02.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )
  gr[nodeID5A.ToString()] = register.NewGatewayRegister(
    nodeID5A.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )
  gr[nodeIDFFFF.ToString()] = register.NewGatewayRegister(
    nodeIDFFFF.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )

	type fields struct {
		start                     bool
		gatewayDiscv              bool
		registeredGatewaysMap     map[string]register.GatewayRegistrar
		registeredGatewaysMapLock sync.RWMutex
	}
	type args struct {
		id *nodeid.NodeID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   register.GatewayRegistrar
	}{
		{name: "getGateway not started returns nil",
			fields: fields{
				start:                     false,
				gatewayDiscv:              false,
				registeredGatewaysMap:     gr,
				registeredGatewaysMapLock: sync.RWMutex{},
			},
			args: struct{ id *nodeid.NodeID }{id: nodeID00},
			want: nil,
		},
		{name: "getGateway - gatewayDiscv not started returns nil",
			fields: fields{
				start:                     true,
				gatewayDiscv:              false,
				registeredGatewaysMap:     gr,
				registeredGatewaysMapLock: sync.RWMutex{},
			},
			args: struct{ id *nodeid.NodeID }{id: nodeID00},
			want: nil,
		},
		{name: "getGateway - started return nodeID00",
			fields: fields{
				start:                     true,
				gatewayDiscv:              true,
				registeredGatewaysMap:     gr,
				registeredGatewaysMapLock: sync.RWMutex{},
			},
			args: struct{ id *nodeid.NodeID }{id: nodeID00},
			want: register.NewGatewayRegister(
        nodeID00.ToString(),
        "address",
        "rootSigningKey",
        "signingKey",
        "regionCode",
        "networkInfoGateway",
        "networkInfoProvider",
        "networkInfoClient",
        "networkInfoAdmin",
      ),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := &FCRRegisterMgr{
				start:                     tt.fields.start,
				gatewayDiscv:              tt.fields.gatewayDiscv,
				registeredGatewaysMap:     tt.fields.registeredGatewaysMap,
				registeredGatewaysMapLock: tt.fields.registeredGatewaysMapLock,
			}
			if got := mgr.GetGateway(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPullGatewaysFromRegisterSrv_ReturnsNil(t *testing.T) {
	nodeID00, _ := nodeid.NewNodeIDFromHexString("00")
	nodeID01, _ := nodeid.NewNodeIDFromHexString("01")
	nodeID02, _ := nodeid.NewNodeIDFromHexString("02")
	nodeID5A, _ := nodeid.NewNodeIDFromHexString("5A")
	nodeIDFFFF, _ := nodeid.NewNodeIDFromHexString("FFFF")

	gr := map[string]register.GatewayRegistrar{}
	gr[nodeID00.ToString()] = register.NewGatewayRegister(
    nodeID00.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )
  gr[nodeID01.ToString()] = register.NewGatewayRegister(
    nodeID01.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )
  gr[nodeID02.ToString()] = register.NewGatewayRegister(
    nodeID02.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )
  gr[nodeID5A.ToString()] = register.NewGatewayRegister(
    nodeID5A.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )
  gr[nodeIDFFFF.ToString()] = register.NewGatewayRegister(
    nodeIDFFFF.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )

	type fields struct {
		start                     bool
		gatewayDiscv              bool
		registeredGatewaysMap     map[string]register.GatewayRegistrar
		registeredGatewaysMapLock sync.RWMutex
	}

	tests := []struct {
		name   string
		fields fields

		want []register.GatewayRegistrar
	}{
		{name: "getGateway not started returns nil",
			fields: fields{
				start:                     false,
				gatewayDiscv:              false,
				registeredGatewaysMap:     gr,
				registeredGatewaysMapLock: sync.RWMutex{},
			},

			want: nil,
		},
		{name: "getGateway - gatewayDiscv not started returns nil",
			fields: fields{
				start:                     true,
				gatewayDiscv:              false,
				registeredGatewaysMap:     gr,
				registeredGatewaysMapLock: sync.RWMutex{},
			},

			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := &FCRRegisterMgr{
				start:                     tt.fields.start,
				gatewayDiscv:              tt.fields.gatewayDiscv,
				registeredGatewaysMap:     tt.fields.registeredGatewaysMap,
				registeredGatewaysMapLock: tt.fields.registeredGatewaysMapLock,
			}
			got := mgr.GetAllGateways()
			assert.ElementsMatch(t, got, tt.want)
		})
	}
}

func TestPullGatewaysFromRegisterSrv_ReturnsNode(t *testing.T) {
  nodeID00, _ := nodeid.NewNodeIDFromHexString("00")
  nodeID01, _ := nodeid.NewNodeIDFromHexString("01")
  nodeID02, _ := nodeid.NewNodeIDFromHexString("02")
  nodeID5A, _ := nodeid.NewNodeIDFromHexString("5A")
  nodeIDFFFF, _ := nodeid.NewNodeIDFromHexString("FFFF")

  gr := map[string]register.GatewayRegistrar{}
  gr[nodeID00.ToString()] = register.NewGatewayRegister(
    nodeID00.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )
  gr[nodeID01.ToString()] = register.NewGatewayRegister(
    nodeID01.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )
  gr[nodeID02.ToString()] = register.NewGatewayRegister(
    nodeID02.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )
  gr[nodeID5A.ToString()] = register.NewGatewayRegister(
    nodeID5A.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )
  gr[nodeIDFFFF.ToString()] = register.NewGatewayRegister(
    nodeIDFFFF.ToString(),
    "address",
    "rootSigningKey",
    "signingKey",
    "regionCode",
    "networkInfoGateway",
    "networkInfoProvider",
    "networkInfoClient",
    "networkInfoAdmin",
  )

  type fields struct {
    start                     bool
    gatewayDiscv              bool
    registeredGatewaysMap     map[string]register.GatewayRegistrar
    registeredGatewaysMapLock sync.RWMutex
  }

  // arrange mock
  mockCtrl := gomock.NewController(t)
  defer mockCtrl.Finish()
  mockHttpCommunicator := mocks.NewMockHttpCommunications(mockCtrl)
  fakeResponse := []*register.GatewayRegister{
    {
      nodeID02.ToString(),
      "address",
      "rootSigningKey",
      "signingKey",
      "regionCode",
      "networkInfoGateway",
      "networkInfoProvider",
      "networkInfoClient",
      "networkInfoAdmin",
    },
    {
      nodeID5A.ToString(),
      "address",
      "rootSigningKey",
      "signingKey",
      "regionCode",
      "networkInfoGateway",
      "networkInfoProvider",
      "networkInfoClient",
      "networkInfoAdmin",
    },
    {
      nodeIDFFFF.ToString(),
      "address",
      "rootSigningKey",
      "signingKey",
      "regionCode",
      "networkInfoGateway",
      "networkInfoProvider",
      "networkInfoClient",
      "networkInfoAdmin",
    },
    {
      nodeID00.ToString(),
      "address",
      "rootSigningKey",
      "signingKey",
      "regionCode",
      "networkInfoGateway",
      "networkInfoProvider",
      "networkInfoClient",
      "networkInfoAdmin",
    },
    {
      nodeID01.ToString(),
      "address",
      "rootSigningKey",
      "signingKey",
      "regionCode",
      "networkInfoGateway",
      "networkInfoProvider",
      "networkInfoClient",
      "networkInfoAdmin",
    },
  }
  fakeResponseBytes, err := json.Marshal(fakeResponse)
  if err != nil{
    t.Fail()
  }
  mockHttpCommunicator.EXPECT().GetJSON(gomock.Any()).Return( fakeResponseBytes, nil).Times(1)

  f := fields{
    start:                     true,
    gatewayDiscv:              true,
    registeredGatewaysMap:     gr,
    registeredGatewaysMapLock: sync.RWMutex{},
  }

  mgr := &FCRRegisterMgr{
    start:                     f.start,
    gatewayDiscv:              f.gatewayDiscv,
    registeredGatewaysMap:     f.registeredGatewaysMap,
    registeredGatewaysMapLock: f.registeredGatewaysMapLock,
    httpCommunicator: mockHttpCommunicator,
  }

  var expectedResult []register.GatewayRegistrar
  for _, g := range fakeResponse {
    expectedResult = append(expectedResult, register.NewGatewayRegister(g.NodeID, g.Address, g.RootSigningKey, g.SigningKey, g.RegionCode, g.NetworkInfoGateway, g.NetworkInfoProvider, g.NetworkInfoClient, g.NetworkInfoAdmin))
  }

  got := mgr.pullGatewaysFromRegisterSrv()
  assert.ElementsMatch(t, got, fakeResponse)
}
