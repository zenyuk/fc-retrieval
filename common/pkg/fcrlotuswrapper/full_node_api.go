package fcrlotuswrapper

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
	"context"
	"net/http"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

type Config struct {
	Address   string
	AuthToken string
}

// NewFullNodeAPI json-rpc api to lotus.
func (c *Config) NewFullNodeAPI(ctx context.Context) (*LotusNodeApi, error) {
	headers := http.Header{"Authorization": []string{"Bearer " + c.AuthToken}}
	var nodeAPI LotusNodeApi
	closer, err := jsonrpc.NewMergeClient(ctx, c.Address, "Filecoin", []interface{}{
		&nodeAPI.Methods,
	}, headers)
	if err != nil {

		return nil, err
	}

	nodeAPI.closer = closer

	return &nodeAPI, nil
}

func (l *LotusNodeApi) Close() {
	l.closer()
}

func (l *LotusNodeApi) ChainHead(ctx context.Context) (*types.TipSet, error) {
	return l.Methods.ChainHead(ctx)
}

func (l *LotusNodeApi) WalletNew(ctx context.Context, kt types.KeyType) (address.Address, error) {
	return l.Methods.WalletNew(ctx, kt)
}

func (l *LotusNodeApi) MpoolPushMessage(ctx context.Context, msg *types.Message, spec *api.MessageSendSpec) (*types.SignedMessage, error) {
	return l.Methods.MpoolPushMessage(ctx, msg, spec)
}
