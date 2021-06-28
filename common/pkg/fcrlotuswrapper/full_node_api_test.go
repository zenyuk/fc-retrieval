package fcrlotuswrapper

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"
)

type testLotusNode struct{}

func (l *testLotusNode) ChainHead(context.Context) (*types.TipSet, error) {
	blh := testBlockHeader()
	return types.NewTipSet([]*types.BlockHeader{blh})
}

func (l *testLotusNode) WalletNew(_ context.Context, _ types.KeyType) (address.Address, error) {
	return address.NewIDAddress(1234564)
}

func (l *testLotusNode) MpoolPushMessage(_ context.Context, _ *types.Message, _ *api.MessageSendSpec) (*types.SignedMessage, error) {
	smsg := testMpoolPushMessageOut()
	return smsg, nil
}

func TestAPIInvalidSchemaError(t *testing.T) {

	rpcServer := jsonrpc.NewServer()
	handler := &testLotusNode{}
	rpcServer.Register("Filecoin", handler)
	testServ := httptest.NewServer(rpcServer)

	addr := testServ.Listener.Addr()
	listenAddr := "invalidSchema://" + addr.String()

	bgCtx := context.Background()

	ctx, cancel := context.WithCancel(bgCtx)
	defer cancel()

	c := &Config{
		Address:   listenAddr,
		AuthToken: "",
	}

	_, err := c.NewFullNodeAPI(ctx)

	assert.NotNil(t, err)
}

func TestAPI(t *testing.T) {

	rpcServer := jsonrpc.NewServer()
	handler := &testLotusNode{}
	rpcServer.Register("Filecoin", handler)
	testServ := httptest.NewServer(rpcServer)

	addr := testServ.Listener.Addr()
	listenAddr := "http://" + addr.String()

	bgCtx := context.Background()

	ctx, cancel := context.WithCancel(bgCtx)
	defer cancel()

	c := &Config{
		Address:   listenAddr,
		AuthToken: "",
	}

	nodeAPI, err := c.NewFullNodeAPI(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer nodeAPI.Close()

	t.Run("test ChainHead", func(t *testing.T) {
		head, err := nodeAPI.ChainHead(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, abi.ChainEpoch(15095854), head.Height())
	})

	t.Run("test WalletNew", func(t *testing.T) {
		walletNew, err := nodeAPI.WalletNew(ctx, types.KTBLS)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, false, walletNew.Empty())
		assert.Equal(t, "f01234564", walletNew.String())
	})

	t.Run("test MpoolPushMessage", func(t *testing.T) {
		msg, err := nodeAPI.MpoolPushMessage(ctx, testMpoolPushMessageIn(), nil)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 61, msg.ChainLength())
	})
}

func testMpoolPushMessageIn() *types.Message {
	toAdd, _ := address.NewIDAddress(1234565)
	fromAdd, _ := address.NewIDAddress(1234560)
	msg := &types.Message{
		Version:    0,
		To:         toAdd,
		From:       fromAdd,
		Nonce:      0,
		Value:      types.FromFil(30),
		GasLimit:   0,
		GasFeeCap:  types.NewInt(0),
		GasPremium: types.NewInt(1),
		Method:     0,
		Params:     nil,
	}
	return msg
}

func testMpoolPushMessageOut() *types.SignedMessage {
	fromAdd, _ := address.NewIDAddress(1234565)
	toAdd, _ := address.NewIDAddress(1234566)

	smsg := &types.SignedMessage{
		Message: types.Message{
			To:         toAdd,
			From:       fromAdd,
			Params:     []byte("some bytes, idk"),
			Method:     1235126,
			Value:      types.NewInt(123123),
			GasFeeCap:  types.NewInt(1234),
			GasPremium: types.NewInt(132414234),
			GasLimit:   100_000_000,
			Nonce:      123123,
		},
	}
	return smsg
}

func testBlockHeader() *types.BlockHeader {
	addr, err := address.NewIDAddress(12512063)
	if err != nil {
		panic(err)
	}

	c, err := cid.Decode("bafyreicmaj5hhoy5mgqvamfhgexxyergw7hdeshizghodwkjg6qmpoco7i")
	if err != nil {
		panic(err)
	}

	return &types.BlockHeader{
		Miner: addr,
		Ticket: &types.Ticket{
			VRFProof: []byte("vrf proof0000000vrf proof0000000"),
		},
		ElectionProof: &types.ElectionProof{
			VRFProof: []byte("vrf proof0000000vrf proof0000000"),
		},
		Parents:               []cid.Cid{c, c},
		ParentMessageReceipts: c,
		BLSAggregate:          &crypto.Signature{Type: crypto.SigTypeBLS, Data: []byte("bls signature")},
		ParentWeight:          types.NewInt(123125126212),
		Messages:              c,
		Height:                15095854,
		ParentStateRoot:       c,
		BlockSig:              &crypto.Signature{Type: crypto.SigTypeBLS, Data: []byte("block signature")},
		ParentBaseFee:         types.NewInt(7451555677934999),
	}
}
