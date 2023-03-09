// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"context"

	"github.com/HorizenOfficial/rosetta-zen/zen"
	"github.com/coinbase/rosetta-sdk-go/types"
)

const (
	// NodeVersion is the version of
	// zend core we are using.
	NodeVersion = "4.0.0-rc1"

	// HistoricalBalanceLookup indicates
	// that historical balance lookup is supported.
	HistoricalBalanceLookup = true

	// inlineFetchLimit is the maximum number
	// of transactions to fetch inline.
	inlineFetchLimit = 100
)

var (
	// MiddlewareVersion is the version
	// of rosetta-zen. We set this as a
	// variable instead of a constant because
	// we typically need the pointer of this
	// value.
	MiddlewareVersion = "0.0.6"
)

// Client is used by the servicers to get Peer information
// and to submit transactions.
type Client interface {
	GetPeers(context.Context) ([]*types.Peer, error)
	SendRawTransaction(context.Context, string) (string, error)
	SuggestedFeeRate(context.Context, int64) (float64, error)
	RawMempool(context.Context) ([]string, error)
	GetBestBlock (context.Context) (int64, error)
	GetHashFromIndex(context.Context, int64) (string, error)
}

// Indexer is used by the servicers to get block and account data.
type Indexer interface {
	GetBlockLazy(
		context.Context,
		*types.PartialBlockIdentifier,
	) (*types.BlockResponse, error)
	GetBlockTransaction(
		context.Context,
		*types.BlockIdentifier,
		*types.TransactionIdentifier,
	) (*types.Transaction, error)
	GetCoins(
		context.Context,
		*types.AccountIdentifier,
	) ([]*types.Coin, *types.BlockIdentifier, error)
	GetScriptPubKeys(
		context.Context,
		[]*types.Coin,
	) ([]*zen.ScriptPubKey, error)
	GetBalance(
		context.Context,
		*types.AccountIdentifier,
		*types.Currency,
		*types.PartialBlockIdentifier,
	) (*types.Amount, *types.BlockIdentifier, error)
}

type unsignedTransaction struct {
	Transaction    string              `json:"transaction"`
	ScriptPubKeys  []*zen.ScriptPubKey `json:"scriptPubKeys"`
	InputAmounts   []string            `json:"input_amounts"`
	InputAddresses []string            `json:"input_addresses"`
}

type preprocessOptions struct {
	Coins         []*types.Coin `json:"coins"`
	EstimatedSize float64       `json:"estimated_size"`
	FeeMultiplier *float64      `json:"fee_multiplier,omitempty"`
}

type constructionMetadata struct {
	ScriptPubKeys []*zen.ScriptPubKey `json:"script_pub_keys"`
	ReplayBlockHeight int64           `json:"replay_block_height"`
	ReplayBlockHash string            `json:"replay_block_hash"`

}

type signedTransaction struct {
	Transaction  string   `json:"transaction"`
	InputAmounts []string `json:"input_amounts"`
}

// ParseOperationMetadata is returned from
// ConstructionParse.
type ParseOperationMetadata struct {
	ScriptPubKey *zen.ScriptPubKey `json:"scriptPubKey"`
}
