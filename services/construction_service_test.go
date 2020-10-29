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
			"03164f76360ef79e7513eff3095e8b60a5cf98223bed0d3109aaabe5f061be4140",
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
			Address: "ztcHp2reR5d4AhZLLp5bYELzfZXHQERQogi",
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
				Address: "ztcHp2reR5d4AhZLLp5bYELzfZXHQERQogi",
			},
			Amount: &types.Amount{
				Value:    "-1000000000",
				Currency: zen.TestnetCurrency,
			},
			CoinChange: &types.CoinChange{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: "a2b082a14210712ea7d1edd317273154e102a3517c144240da2b8ed696305b08:1",
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
				Address: "ztfPiJyJL3UavuYw5Fiv1V1okdbsmY1b5qX",
			},
			Amount: &types.Amount{
				Value:    "1000000000",
				Currency: zen.TestnetCurrency,
			},
		},
	}
	preprocessResponse, err := servicer.ConstructionPreprocess(
		ctx,
		&types.ConstructionPreprocessRequest{
			NetworkIdentifier:      networkIdentifier,
			Operations:             ops,
		},
	)
	assert.Nil(t, err)
	options := &preprocessOptions{
		Coins: []*types.Coin{
			{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: "a2b082a14210712ea7d1edd317273154e102a3517c144240da2b8ed696305b08:1",
				},
				Amount: &types.Amount{
					Value:    "-1000000000",
					Currency: zen.TestnetCurrency,
				},
			},
		},
		EstimatedSize: 227,
	}
	assert.Equal(t, &types.ConstructionPreprocessResponse{
		Options: forceMarshalMap(t, options),
	}, preprocessResponse)

	// Test Metadata
	metadata := &constructionMetadata{
		ScriptPubKeys: []*zen.ScriptPubKey{
			{
				ASM:          "OP_DUP OP_HASH160 64352ca2f736dc4e7464a65f8b07ef313d7ab53d OP_EQUALVERIFY OP_CHECKSIG b6ce3a2fb53f49ce31bcf2d404cf3bfa88caf71bd5ec0b4b1dc7eef8ad89470d 11 OP_CHECKBLOCKATHEIGHT",
				Hex:          "76a91464352ca2f736dc4e7464a65f8b07ef313d7ab53d88ac20b6ce3a2fb53f49ce31bcf2d404cf3bfa88caf71bd5ec0b4b1dc7eef8ad89470d5bb4",
				RequiredSigs: 1,
				Type:         "pubkeyhashreplay",
				Addresses: []string{
					"ztcHp2reR5d4AhZLLp5bYELzfZXHQERQogi",
				},
			},
		},
		ReplayBlockHeight: 212,
		ReplayBlockHash: "0786aeb320d7eb3c98486eba566b2c2ff893e39abd62e647560b15240f8216f8",
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
		int64(312), nil).Twice()
	mockClient.On(
		"GetHashFromIndex",
		ctx,
		int64(212)).Return(
		"0786aeb320d7eb3c98486eba566b2c2ff893e39abd62e647560b15240f8216f8", nil).Twice()

	metadataResponse, err := servicer.ConstructionMetadata(ctx, &types.ConstructionMetadataRequest{
		NetworkIdentifier: networkIdentifier,
		Options:           forceMarshalMap(t, options),
	})
	assert.Nil(t, err)
	assert.Equal(t, &types.ConstructionMetadataResponse{
		Metadata: forceMarshalMap(t, metadata),
		SuggestedFee: []*types.Amount{
			{
				Value:    "2270", // 1,420 * 0.75
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
				Value:    "227", // we don't go below minimum fee rate
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
				Address: "ztcHp2reR5d4AhZLLp5bYELzfZXHQERQogi",
			},
			Amount: &types.Amount{
				Value:    "-1000000000",
				Currency: zen.TestnetCurrency,
			},
			CoinChange: &types.CoinChange{
				CoinIdentifier: &types.CoinIdentifier{
					Identifier: "a2b082a14210712ea7d1edd317273154e102a3517c144240da2b8ed696305b08:1",
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
				Address: "ztfPiJyJL3UavuYw5Fiv1V1okdbsmY1b5qX",
			},
			Amount: &types.Amount{
				Value:    "1000000000",
				Currency: zen.TestnetCurrency,
			},
		},
	}

	assert.Nil(t, err)

	signingPayload := &types.SigningPayload{
		Bytes: forceHexDecode(
			t,
			"7068aa7955b0469aa51cefcf0ae0a45448fe945a4f606ca373df1f4c3ca8002f",
		),
		AccountIdentifier: &types.AccountIdentifier{
			Address: "ztcHp2reR5d4AhZLLp5bYELzfZXHQERQogi",
		},
		SignatureType: types.Ecdsa,
	}

	unsignedRaw := "7b227472616e73616374696f6e223a2230313030303030303031303835623330393664363865326264613430343231343763353161333032653135343331323731376433656464316137326537313130343261313832623061323031303030303030303066666666666666663031303063613961336230303030303030303365373661393134383633623435353736613133306463396338343838326436366663656165393235363463656230663838616332306638313638323066323431353062353634376536363262643961653339336638326632633662353662613665343839383363656264373230623361653836303730326434303062343030303030303030222c227363726970745075624b657973223a5b7b2261736d223a224f505f445550204f505f484153483136302036343335326361326637333664633465373436346136356638623037656633313364376162353364204f505f455155414c564552494659204f505f434845434b5349472062366365336132666235336634396365333162636632643430346366336266613838636166373162643565633062346231646337656566386164383934373064203131204f505f434845434b424c4f434b4154484549474854222c22686578223a22373661393134363433353263613266373336646334653734363461363566386230376566333133643761623533643838616332306236636533613266623533663439636533316263663264343034636633626661383863616637316264356563306234623164633765656638616438393437306435626234222c2272657153696773223a312c2274797065223a227075626b6579686173687265706c6179222c22616464726573736573223a5b227a746348703272655235643441685a4c4c70356259454c7a665a5848514552516f6769225d7d5d2c22696e7075745f616d6f756e7473223a5b222d31303030303030303030225d2c22696e7075745f616464726573736573223a5b227a746348703272655235643441685a4c4c70356259454c7a665a5848514552516f6769225d7d"

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
	signedRaw := "7b227472616e73616374696f6e223a22303130303030303030313038356233303936643638653262646134303432313437633531613330326531353433313237313764336564643161373265373131303432613138326230613230313030303030303661343733303434303232303632343234663837363563386361303936303134316362306566663132386339633639363762636661313632656239626539323637356533396633346639613730323230373434373830323730393239613130303633303762633538663837643366373132376163383635306264393636313036373837633930626339373663366332633031323130333136346637363336306566373965373531336566663330393565386236306135636639383232336265643064333130396161616265356630363162653431343066666666666666663031303063613961336230303030303030303365373661393134383633623435353736613133306463396338343838326436366663656165393235363463656230663838616332306638313638323066323431353062353634376536363262643961653339336638326632633662353662613665343839383363656264373230623361653836303730326434303062343030303030303030222c22696e7075745f616d6f756e7473223a5b222d31303030303030303030225d7d" // nolint
	combineResponse, err := servicer.ConstructionCombine(ctx, &types.ConstructionCombineRequest{
		NetworkIdentifier:   networkIdentifier,
		UnsignedTransaction: unsignedRaw,
		Signatures: []*types.Signature{
			{
				Bytes: forceHexDecode(
					t,
					"62424f8765c8ca0960141cb0eff128c9c6967bcfa162eb9be92675e39f34f9a7744780270929a1006307bc58f87d3f7127ac8650bd966106787c90bc976c6c2c", // nolint
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
			{Address: "ztcHp2reR5d4AhZLLp5bYELzfZXHQERQogi"},
		},
	}, parseSignedResponse)

	// Test Hash
	transactionIdentifier := &types.TransactionIdentifier{
		Hash: "a589c9941da87b15ffbb419569f38a1d44c805aeaafa167088d270de0fd8eed2",
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
	bitcoinTransaction := "0100000001085b3096d68e2bda4042147c51a302e154312717d3edd1a72e711042a182b0a2010000006a473044022062424f8765c8ca0960141cb0eff128c9c6967bcfa162eb9be92675e39f34f9a70220744780270929a1006307bc58f87d3f7127ac8650bd966106787c90bc976c6c2c012103164f76360ef79e7513eff3095e8b60a5cf98223bed0d3109aaabe5f061be4140ffffffff0100ca9a3b000000003e76a914863b45576a130dc9c84882d66fceae92564ceb0f88ac20f816820f24150b5647e662bd9ae393f82f2c6b56ba6e48983cebd720b3ae860702d400b400000000" // nolint
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
