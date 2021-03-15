package integration

/*
 * Copyright 2021 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/cidoffer"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/pkg/fcrgatewayadmin"
	"github.com/stretchr/testify/assert"
)

// GenerateDummyMessage to generate a dummy message
func GenerateDummyMessage() *fcrmessages.FCRMessage {
	providerID, _ := nodeid.NewRandomNodeID()
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()
	contentID, _ := cid.NewRandomContentID()
	pieceCIDs := []cid.ContentID{*contentID}

	cidOffer := cidoffer.CidGroupOffer{
		NodeID:     providerID,
		Cids:       pieceCIDs,
		Price:      42,
		Expiry:     expiryDate,
		QoS:        42,
		MerkleTrie: nil,
		Signature:  "",
	}

	dummyMessage, err := fcrmessages.EncodeProviderPublishGroupCIDRequest(
		rand.Int63n(100000),
		&cidOffer,
	)
	if err != nil {
		logging.Error("Error when encoding message: %v", err)
	}
	return dummyMessage
}

func TestSignMessage(t *testing.T) {
	gatewayRootKey, err := fcrgatewayadmin.CreateKey()
	if err != nil {
		panic(err)
	}
	_, err1 := gatewayRootKey.EncodePublicKey()
	if err1 != nil {
		panic(err1)
	}
	gatewayRetrievalPrivateKey, err := fcrgatewayadmin.CreateKey()
	if err != nil {
		panic(err)
	}
	gatewayRetrievalSigningKey, err2 := gatewayRetrievalPrivateKey.EncodePublicKey()
	if err2 != nil {
		panic(err2)
	}
	pubKey, err3 := fcrcrypto.DecodePublicKey(gatewayRetrievalSigningKey)
	if err3 != nil {
		panic(err3)
	}
	privatekeyversion := fcrcrypto.DecodeKeyVersion(1)

	message := GenerateDummyMessage()
	message.SignMessage(func(msg interface{}) (string, error) {
		return fcrcrypto.SignMessage(gatewayRetrievalPrivateKey, privatekeyversion, msg)
	})

	// Verify the response
	ok, err := message.VerifySignature(func(sig string, msg interface{}) (bool, error) {
		return fcrcrypto.VerifyMessage(pubKey, sig, msg)
	})

	assert.Equal(t, true, ok)
}
