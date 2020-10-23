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
	"encoding/hex"
	"testing"

	"github.com/HorizenOfficial/rosetta-zen/zen"
	"github.com/HorizenOfficial/rosetta-zen/configuration"
	mocks "github.com/HorizenOfficial/rosetta-zen/mocks/services"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func forceHexDecode(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("could not decode hex %s", s)
	}

	return b
}

func forceMarshalMap(t *testing.T, i interface{}) map[string]interface{} {
	m, err := types.MarshalMap(i)
	if err != nil {
		t.Fatalf("could not marshal map %s", types.PrintStruct(i))
	}

	return m
}

func TestConstructionService(t *testing.T) {
	networkIdentifier = &types.NetworkIdentifier{
		Network:    zen.TestnetNetwork,
		Blockchain: zen.Blockchain,
	}

	cfg := &configuration.Configuration{
		Mode:     configuration.Online,
		Network:  networkIdentifier,
		Params:   zen.TestnetParams,
		Currency: zen.TestnetCurrency,
	}

	mockIndexer := &mocks.Indexer{}
	mockClient := &mocks.Client{}
	servicer := NewConstructionAPIService(cfg, mockClient, mockIndexer)
	ctx := context.Background()

	// Test Derive
	publicKey := &types.PublicKey{
		Bytes: forceHexDecode(
			t,
			"03f892ec106c94bdead9f088797ec2bb6d0f46cc7f7e6a931a0fd76c52aee5d016",
		),
		CurveType: types.Secp256k1,
	}
	deriveResponse, err := servicer.ConstructionDerive(ctx, &types.ConstructionDeriveRequest{
		NetworkIdentifier: networkIdentifier,
		PublicKey:         publicKey,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionDeriveResponse{
		AccountIdentifier: &types.AccountIdentifier{
			Address: "ztmfGwLDqR9bApbXi9Nzb4JuGbhS9Biwn4M",
		},
	}, deriveResponse)

	// Test Preprocess
	ops := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 0,
			},
			Type: zen.InputOpType,
			Account: &types.AccountIdentifier{
				Address: "ztmfGwLDqR9bApbXi9Nzb4JuGbhS9Biwn4M",
			},
			Amount: &types.Amount{
				Value:    "-1143750000",
				Currency: zen.TestnetCurrency,
			},
			CoinChange: &types.CoinChange{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: "507ff99b92b77885a3b308906202e49c0efdc1a75a69fa559f54efe7f6fafce2:0",
				},
				CoinAction: types.CoinSpent,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index: 1,
			},
			Type: zen.OutputOpType,
			Account: &types.AccountIdentifier{
				Address: "ztUWqnHtPBV1kuDS5gcmWi5yGthzCpzWS2G",
			},
			Amount: &types.Amount{
				Value:    "1143740000",
				Currency: zen.TestnetCurrency,
			},
		},
	}
	feeMultiplier := float64(0.75)
	preprocessResponse, err := servicer.ConstructionPreprocess(
		ctx,
		&types.ConstructionPreprocessRequest{
			NetworkIdentifier:      networkIdentifier,
			Operations:             ops,
			SuggestedFeeMultiplier: &feeMultiplier,
		},
	)
	assert.Nil(t, err)
	options := &preprocessOptions{
		Coins: []*types.Coin{
			{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: "507ff99b92b77885a3b308906202e49c0efdc1a75a69fa559f54efe7f6fafce2:0",
				},
				Amount: &types.Amount{
					Value:    "-1143750000",
					Currency: zen.TestnetCurrency,
				},
			},
		},
		EstimatedSize: 114,
		FeeMultiplier: &feeMultiplier,
	}
	assert.Equal(t, &types.ConstructionPreprocessResponse{
		Options: forceMarshalMap(t, options),
	}, preprocessResponse)

	// Test Metadata
	metadata := &constructionMetadata{
		ScriptPubKeys: []*zen.ScriptPubKey{
			{
				ASM:          "OP_DUP OP_HASH160 cafd80252588892a4c340bc26cff7ddf7b8a4170 OP_EQUALVERIFY OP_CHECKSIG",
				Hex:          "76a914cafd80252588892a4c340bc26cff7ddf7b8a417088ac",
				RequiredSigs: 1,
				Type:         "pubkeyhash",
				Addresses: []string{
					"ztmfGwLDqR9bApbXi9Nzb4JuGbhS9Biwn4M",
				},
			},
		},
		ReplayBlockHeight: 0,
		ReplayBlockHash: "0da5ee723b7923feb580518541c6f098206330dbc711a6678922c11f2ccf1abb",
	}

	// Normal Fee
	mockIndexer.On(
		"GetScriptPubKeys",
		ctx,
		options.Coins,
	).Return(
		metadata.ScriptPubKeys,
		nil,
	).Once()
	mockClient.On(
		"SuggestedFeeRate",
		ctx,
		defaultConfirmationTarget,
	).Return(
		zen.MinFeeRate*10,
		nil,
	).Once()
	mockClient.On(
		"GetBestBlock",
		ctx).Return(
		int64(100), nil).Twice()
	mockClient.On(
		"GetHashFromIndex",
		ctx,
		int64(0)).Return(
		"0da5ee723b7923feb580518541c6f098206330dbc711a6678922c11f2ccf1abb", nil).Twice()

	metadataResponse, err := servicer.ConstructionMetadata(ctx, &types.ConstructionMetadataRequest{
		NetworkIdentifier: networkIdentifier,
		Options:           forceMarshalMap(t, options),
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionMetadataResponse{
		Metadata: forceMarshalMap(t, metadata),
		SuggestedFee: []*types.Amount{
			{
				Value:    "855", // 1,420 * 0.75
				Currency: zen.TestnetCurrency,
			},
		},
	}, metadataResponse)

	// Low Fee
	mockIndexer.On(
		"GetScriptPubKeys",
		ctx,
		options.Coins,
	).Return(
		metadata.ScriptPubKeys,
		nil,
	).Once()
	mockClient.On(
		"SuggestedFeeRate",
		ctx,
		defaultConfirmationTarget,
	).Return(
		zen.MinFeeRate,
		nil,
	).Once()
	metadataResponse, err = servicer.ConstructionMetadata(ctx, &types.ConstructionMetadataRequest{
		NetworkIdentifier: networkIdentifier,
		Options:           forceMarshalMap(t, options),
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionMetadataResponse{
		Metadata: forceMarshalMap(t, metadata),
		SuggestedFee: []*types.Amount{
			{
				Value:    "114", // we don't go below minimum fee rate
				Currency: zen.TestnetCurrency,
			},
		},
	}, metadataResponse)

	// Test Payloads
	payloadsResponse, err := servicer.ConstructionPayloads(ctx, &types.ConstructionPayloadsRequest{
		NetworkIdentifier: networkIdentifier,
		Operations:        ops,
		Metadata:          forceMarshalMap(t, metadata),
	})
	val0 := int64(0)
	parseOps := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        0,
				NetworkIndex: &val0,
			},
			Type: zen.InputOpType,
			Account: &types.AccountIdentifier{
				Address: "ztmfGwLDqR9bApbXi9Nzb4JuGbhS9Biwn4M",
			},
			Amount: &types.Amount{
				Value:    "-1143750000",
				Currency: zen.TestnetCurrency,
			},
			CoinChange: &types.CoinChange{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: "507ff99b92b77885a3b308906202e49c0efdc1a75a69fa559f54efe7f6fafce2:0",
				},
				CoinAction: types.CoinSpent,
			},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{
				Index:        1,
				NetworkIndex: &val0,
			},
			Type: zen.OutputOpType,
			Account: &types.AccountIdentifier{
				Address: "ztUWqnHtPBV1kuDS5gcmWi5yGthzCpzWS2G",
			},
			Amount: &types.Amount{
				Value:    "1143740000",
				Currency: zen.TestnetCurrency,
			},
		},
	}

	assert.Nil(t, err)

	signingPayload := &types.SigningPayload{
		Bytes: forceHexDecode(
			t,
			"a0bd933c51ac04de04cbba23630e992da27086f1129610897242ca532a1d082d",
		),
		AccountIdentifier: &types.AccountIdentifier{
			Address: "ztmfGwLDqR9bApbXi9Nzb4JuGbhS9Biwn4M",
		},
		SignatureType: types.Ecdsa,
	}

	unsignedRaw := "7b227472616e73616374696f6e223a22303130303030303030316532666366616636653765663534396635356661363935616137633166643065396365343032363239303038623361333835373862373932396266393766353030303030303030303030666666666666666630313630313632633434303030303030303033633736613931343065656230393135633330653564303362323762313961366133613638313465663632643463303438386163323062623161636632633166633132323839363761363131633764623330363332303938663063363431383535313830623566653233373933623732656561353064303062343030303030303030222c227363726970745075624b657973223a5b7b2261736d223a224f505f445550204f505f484153483136302063616664383032353235383838393261346333343062633236636666376464663762386134313730204f505f455155414c564552494659204f505f434845434b534947222c22686578223a223736613931346361666438303235323538383839326134633334306263323663666637646466376238613431373038386163222c2272657153696773223a312c2274797065223a227075626b657968617368222c22616464726573736573223a5b227a746d6647774c44715239624170625869394e7a62344a7547626853394269776e344d225d7d5d2c22696e7075745f616d6f756e7473223a5b222d31313433373530303030225d2c22696e7075745f616464726573736573223a5b227a746d6647774c44715239624170625869394e7a62344a7547626853394269776e344d225d7d"

	assert.Equal(t, &types.ConstructionPayloadsResponse{
		UnsignedTransaction: unsignedRaw,
		Payloads:            []*types.SigningPayload{signingPayload},
	}, payloadsResponse)

	// Test Parse Unsigned
	parseUnsignedResponse, err := servicer.ConstructionParse(ctx, &types.ConstructionParseRequest{
		NetworkIdentifier: networkIdentifier,
		Signed:            false,
		Transaction:       unsignedRaw,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionParseResponse{
		Operations:               parseOps,
		AccountIdentifierSigners: []*types.AccountIdentifier{},
	}, parseUnsignedResponse)

	// Test Combine
	signedRaw := "7b227472616e73616374696f6e223a2230313030303030303031653266636661663665376566353439663535666136393561613763316664306539636534303236323930303862336133383537386237393239626639376635303030303030303030366234383330343530323231303066623834363634663030656133626636653535366238653165393433303736303463336533386262636639396338346531306662333133316232323562356531303232303262656165333366646362646435643531636630623566663635306361393861396438646435393363353665376461383861323565326336303762363562346330313231303366383932656331303663393462646561643966303838373937656332626236643066343663633766376536613933316130666437366335326165653564303136666666666666666630313630313632633434303030303030303033633736613931343065656230393135633330653564303362323762313961366133613638313465663632643463303438386163323062623161636632633166633132323839363761363131633764623330363332303938663063363431383535313830623566653233373933623732656561353064303062343030303030303030222c22696e7075745f616d6f756e7473223a5b222d31313433373530303030225d7d" // nolint
	combineResponse, err := servicer.ConstructionCombine(ctx, &types.ConstructionCombineRequest{
		NetworkIdentifier:   networkIdentifier,
		UnsignedTransaction: unsignedRaw,
		Signatures: []*types.Signature{
			{
				Bytes: forceHexDecode(
					t,
					"fb84664f00ea3bf6e556b8e1e94307604c3e38bbcf99c84e10fb3131b225b5e12beae33fdcbdd5d51cf0b5ff650ca98a9d8dd593c56e7da88a25e2c607b65b4c", // nolint
				),
				SigningPayload: signingPayload,
				PublicKey:      publicKey,
				SignatureType:  types.Ecdsa,
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionCombineResponse{
		SignedTransaction: signedRaw,
	}, combineResponse)

	// Test Parse Signed
	parseSignedResponse, err := servicer.ConstructionParse(ctx, &types.ConstructionParseRequest{
		NetworkIdentifier: networkIdentifier,
		Signed:            true,
		Transaction:       signedRaw,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionParseResponse{
		Operations: parseOps,
		AccountIdentifierSigners: []*types.AccountIdentifier{
			{Address: "ztmfGwLDqR9bApbXi9Nzb4JuGbhS9Biwn4M"},
		},
	}, parseSignedResponse)

	// Test Hash
	transactionIdentifier := &types.TransactionIdentifier{
		Hash: "ff6dc9082f947e6c2d76cecbab93dbc17591bb62e91194bb780a7dcb734e8184",
	}
	hashResponse, err := servicer.ConstructionHash(ctx, &types.ConstructionHashRequest{
		NetworkIdentifier: networkIdentifier,
		SignedTransaction: signedRaw,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.TransactionIdentifierResponse{
		TransactionIdentifier: transactionIdentifier,
	}, hashResponse)

	// Test Submit
	bitcoinTransaction := "0100000001e2fcfaf6e7ef549f55fa695aa7c1fd0e9ce402629008b3a38578b7929bf97f50000000006b483045022100fb84664f00ea3bf6e556b8e1e94307604c3e38bbcf99c84e10fb3131b225b5e102202beae33fdcbdd5d51cf0b5ff650ca98a9d8dd593c56e7da88a25e2c607b65b4c012103f892ec106c94bdead9f088797ec2bb6d0f46cc7f7e6a931a0fd76c52aee5d016ffffffff0160162c44000000003c76a9140eeb0915c30e5d03b27b19a6a3a6814ef62d4c0488ac20bb1acf2c1fc1228967a611c7db30632098f0c641855180b5fe23793b72eea50d00b400000000" // nolint
	mockClient.On(
		"SendRawTransaction",
		ctx,
		bitcoinTransaction,
	).Return(
		transactionIdentifier.Hash,
		nil,
	)
	submitResponse, err := servicer.ConstructionSubmit(ctx, &types.ConstructionSubmitRequest{
		NetworkIdentifier: networkIdentifier,
		SignedTransaction: signedRaw,
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.TransactionIdentifierResponse{
		TransactionIdentifier: transactionIdentifier,
	}, submitResponse)

	mockClient.AssertExpectations(t)
	mockIndexer.AssertExpectations(t)
}
