/*
Package adminapi - remote API for FileCoin Secondary Retrieval - Provider Admin library
*/
package adminapi

/*
 * Copyright 2020 ConsenSys Software Inc.
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
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/mocks"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
)

func Test_RequestForceRefresh_Calls_RequestForceRefresh(t *testing.T) {
	// arrange
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHttpCommunicator := mocks.NewMockHttpCommunications(mockCtrl)
	a := NewAdminApiWithDep(mockHttpCommunicator)

	keyPair, _ := fcrcrypto.GenerateRetrievalV1KeyPair()
	pubKeyStr, _ := keyPair.EncodePublicKey()
	fakeProviderInfo  := &register.ProviderRegister{
		NetworkInfoAdmin: "fakeNetworkInfoAdmin",
		SigningKey:       pubKeyStr,
	}

	// assert
	mockHttpCommunicator.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(&fcrmessages.FCRMessage{}, nil).Times(1)

	// act
	_ = a.RequestForceRefresh(fakeProviderInfo, keyPair, fcrcrypto.DecodeKeyVersion(1))

}
