package fcrregistermgr

import (
	"reflect"
	"sync"
	"testing"

	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/stretchr/testify/assert"
)

// func TestFCRRegisterMgr_GetGatewaysNearCID(t *testing.T) {

// 	contentID, _ := cid.NewContentIDFromHexString("01")
// 	nodeID00, _ := nodeid.NewNodeIDFromHexString("00")
// 	nodeID01, _ := nodeid.NewNodeIDFromHexString("01")
// 	nodeID02, _ := nodeid.NewNodeIDFromHexString("02")
// 	nodeID5A, _ := nodeid.NewNodeIDFromHexString("5A")
// 	nodeIDFFFF, _ := nodeid.NewNodeIDFromHexString("FFFF")

// 	gr := map[string]register.GatewayRegister{}
// 	gr[nodeID00.ToString()] = register.GatewayRegister{NodeID: nodeID00.ToString()}
// 	gr[nodeID01.ToString()] = register.GatewayRegister{NodeID: nodeID01.ToString()}
// 	gr[nodeID02.ToString()] = register.GatewayRegister{NodeID: nodeID02.ToString()}
// 	gr[nodeID5A.ToString()] = register.GatewayRegister{NodeID: nodeID5A.ToString()}
// 	gr[nodeIDFFFF.ToString()] = register.GatewayRegister{NodeID: nodeIDFFFF.ToString()}

// 	type fields struct {
// 		registeredGatewaysMap     map[string]register.GatewayRegister
// 		registeredGatewaysMapLock sync.RWMutex
// 		closestGatewaysIDs        [][]byte
// 		closestGatewaysMapLock    sync.RWMutex
// 	}
// 	type args struct {
// 		cID        *cid.ContentID
// 		numDHT     int
// 		excludeIDs []*nodeid.NodeID
// 	}
// 	var tests = []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    []register.GatewayRegister
// 		wantErr bool
// 	}{
// 		{
// 			name: "test nil closest gateways",
// 			fields: fields{
// 				registeredGatewaysMapLock: sync.RWMutex{},
// 				closestGatewaysIDs:        nil,
// 			}, args: args{
// 				cID:        nil,
// 				numDHT:     10,
// 				excludeIDs: nil,
// 			},
// 			want:    []register.GatewayRegister{},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test empty closest gateways",
// 			fields: fields{
// 				registeredGatewaysMapLock: sync.RWMutex{},
// 				closestGatewaysIDs:        nil,
// 			}, args: args{
// 				cID:        nil,
// 				numDHT:     10,
// 				excludeIDs: nil,
// 			},
// 			want:    []register.GatewayRegister{},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test one closest gateways",
// 			fields: fields{
// 				registeredGatewaysMap: gr,
// 				closestGatewaysIDs:    [][]byte{nodeID00.ToBytes()},
// 			}, args: args{
// 				cID:        contentID,
// 				numDHT:     5,
// 				excludeIDs: nil,
// 			},
// 			want:    []register.GatewayRegister{{NodeID: nodeID00.ToString()}},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test one closest gateways excluded",
// 			fields: fields{
// 				registeredGatewaysMapLock: sync.RWMutex{},
// 				registeredGatewaysMap:     gr,
// 				closestGatewaysIDs:        [][]byte{nodeIDFFFF.ToBytes()},
// 			}, args: args{
// 				cID:        contentID,
// 				numDHT:     12,
// 				excludeIDs: []*nodeid.NodeID{nodeIDFFFF},
// 			},
// 			want:    []register.GatewayRegister{},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test two closestGatewaysIDs, and one closest gateways excluded",
// 			fields: fields{
// 				registeredGatewaysMapLock: sync.RWMutex{},
// 				registeredGatewaysMap:     gr,
// 				closestGatewaysIDs:        [][]byte{nodeIDFFFF.ToBytes(), nodeID01.ToBytes()},
// 			}, args: args{
// 				cID:        contentID,
// 				numDHT:     12,
// 				excludeIDs: []*nodeid.NodeID{nodeIDFFFF},
// 			},
// 			want:    []register.GatewayRegister{{NodeID: nodeID01.ToString()}},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test gateways clockwise order",
// 			fields: fields{
// 				registeredGatewaysMapLock: sync.RWMutex{},
// 				registeredGatewaysMap:     gr,
// 				closestGatewaysIDs:        [][]byte{nodeID00.ToBytes(), nodeID02.ToBytes(), nodeID5A.ToBytes(), nodeIDFFFF.ToBytes()},
// 			}, args: args{
// 				cID:        contentID,
// 				numDHT:     12,
// 				excludeIDs: []*nodeid.NodeID{},
// 			},
// 			want: []register.GatewayRegister{
// 				{NodeID: nodeID02.ToString()},
// 				{NodeID: nodeID5A.ToString()},
// 				{NodeID: nodeIDFFFF.ToString()},
// 				{NodeID: nodeID00.ToString()},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test gateways clockwise order, one excluded",
// 			fields: fields{
// 				registeredGatewaysMapLock: sync.RWMutex{},
// 				registeredGatewaysMap:     gr,
// 				closestGatewaysIDs:        [][]byte{nodeID00.ToBytes(), nodeID02.ToBytes(), nodeID5A.ToBytes(), nodeIDFFFF.ToBytes()},
// 			}, args: args{
// 				cID:        contentID,
// 				numDHT:     12,
// 				excludeIDs: []*nodeid.NodeID{nodeIDFFFF},
// 			},
// 			want: []register.GatewayRegister{
// 				{NodeID: nodeID02.ToString()},
// 				{NodeID: nodeID5A.ToString()},
// 				{NodeID: nodeID00.ToString()},
// 			},
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mgr := &FCRRegisterMgr{
// 				registeredGatewaysMap: tt.fields.registeredGatewaysMap,
// 				closestGatewaysIDs:    tt.fields.closestGatewaysIDs,
// 			}

// 			got, err := mgr.GetGatewaysNearCID(tt.args.cID, tt.args.numDHT, tt.args.excludeIDs)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetGatewaysNearCID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetGatewaysNearCID() got = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestFCRRegisterMgr_GetGateway(t *testing.T) {

	nodeID00, _ := nodeid.NewNodeIDFromHexString("00")
	nodeID01, _ := nodeid.NewNodeIDFromHexString("01")
	nodeID02, _ := nodeid.NewNodeIDFromHexString("02")
	nodeID5A, _ := nodeid.NewNodeIDFromHexString("5A")
	nodeIDFFFF, _ := nodeid.NewNodeIDFromHexString("FFFF")

	gr := map[string]register.GatewayRegister{}
	gr[nodeID00.ToString()] = register.GatewayRegister{NodeID: nodeID00.ToString()}
	gr[nodeID01.ToString()] = register.GatewayRegister{NodeID: nodeID01.ToString()}
	gr[nodeID02.ToString()] = register.GatewayRegister{NodeID: nodeID02.ToString()}
	gr[nodeID5A.ToString()] = register.GatewayRegister{NodeID: nodeID5A.ToString()}
	gr[nodeIDFFFF.ToString()] = register.GatewayRegister{NodeID: nodeIDFFFF.ToString()}

	type fields struct {
		start                     bool
		gatewayDiscv              bool
		registeredGatewaysMap     map[string]register.GatewayRegister
		registeredGatewaysMapLock sync.RWMutex
	}
	type args struct {
		id *nodeid.NodeID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *register.GatewayRegister
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
			want: &register.GatewayRegister{NodeID: nodeID00.ToString()},
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

func TestFCRRegisterMgr_GetAllGateways(t *testing.T) {
	nodeID00, _ := nodeid.NewNodeIDFromHexString("00")
	nodeID01, _ := nodeid.NewNodeIDFromHexString("01")
	nodeID02, _ := nodeid.NewNodeIDFromHexString("02")
	nodeID5A, _ := nodeid.NewNodeIDFromHexString("5A")
	nodeIDFFFF, _ := nodeid.NewNodeIDFromHexString("FFFF")

	gr := map[string]register.GatewayRegister{}
	gr[nodeID00.ToString()] = register.GatewayRegister{NodeID: nodeID00.ToString()}
	gr[nodeID01.ToString()] = register.GatewayRegister{NodeID: nodeID01.ToString()}
	gr[nodeID02.ToString()] = register.GatewayRegister{NodeID: nodeID02.ToString()}
	gr[nodeID5A.ToString()] = register.GatewayRegister{NodeID: nodeID5A.ToString()}
	gr[nodeIDFFFF.ToString()] = register.GatewayRegister{NodeID: nodeIDFFFF.ToString()}

	type fields struct {
		start                     bool
		gatewayDiscv              bool
		registeredGatewaysMap     map[string]register.GatewayRegister
		registeredGatewaysMapLock sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields

		want []register.GatewayRegister
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

			want: []register.GatewayRegister{

				{NodeID: nodeID02.ToString()},
				{NodeID: nodeID5A.ToString()},
				{NodeID: nodeIDFFFF.ToString()},
				{NodeID: nodeID00.ToString()},
				{NodeID: nodeID01.ToString()},
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
