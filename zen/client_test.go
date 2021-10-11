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

package zen

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/storage"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

const (
	url = "/"
)

func forceMarshalMap(t *testing.T, i interface{}) map[string]interface{} {
	m, err := types.MarshalMap(i)
	if err != nil {
		t.Fatalf("could not marshal map %s", types.PrintStruct(i))
	}

	return m
}

var (
	blockIdentifier717983 = &types.BlockIdentifier{
		Hash:  "005d3821c522b528f42fa16187d70ccb59170e2dcd72e9242d54d967e63b6ffe",
		Index: 717983,
	}

	block717983 = &Block{
		Hash:              "005d3821c522b528f42fa16187d70ccb59170e2dcd72e9242d54d967e63b6ffe",
		Height:            717983,
		PreviousBlockHash: "0067f80ce10d4255932b7f8c9baf7bd0dcfd408c312d33144be0ea12caf7f7f0",
		Time:              1601465727,
		Size:              2271,
		Version:           3,
		MerkleRoot:        "97c960c90e0b6bc30d2629f06d114f1c49aadb0e3d9bd70eb4f0f9ed1ea69279",
		Nonce:             "00002e570d64b4b3ea1c30dec68b2dff255eb3148656f06f5e018ae739a400eb",
		Bits:              "1f754920",
		Difficulty:        17.46160923,
		Txs: []*Transaction{
			{
				Hex:      "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff06039ff40a0102ffffffff04903eb42c000000001976a914557662a6b307f95aa00311c074f7feebb955d45188ac80b2e60e0000000017a9148d3468b6686ac59caf9ad94e547a737b09fa102787405973070000000017a914fc1d7f04db5e2c05b051e0decc85effe6bc539d587405973070000000017a9148b85fc1e171a4c7994c088b91d5a75dff9e56cad8700000000", // nolint
				Hash:     "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88",
				Size:     187,
				Version:  1,
				Locktime: 0,
				Inputs: []*Input{
					{
						Coinbase: "039ff40a0102",
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 7.5001,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "OP_DUP OP_HASH160 557662a6b307f95aa00311c074f7feebb955d451 OP_EQUALVERIFY OP_CHECKSIG", // nolint
							Hex:  "76a914557662a6b307f95aa00311c074f7feebb955d45188ac",         // nolint
							RequiredSigs: 1,
							Type: "pubkeyhash",
							Addresses: []string{
								"ztawr1vEZ6pZRtLqNy2C9u7EK7JN2gP8W6z",
							},
						},
					},
					{
						Value: 2.5,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "OP_HASH160 8d3468b6686ac59caf9ad94e547a737b09fa1027 OP_EQUAL", // nolint
							Hex:  "a9148d3468b6686ac59caf9ad94e547a737b09fa102787",         // nolint
							RequiredSigs: 1,
							Type: "scripthash",
							Addresses: []string{
								"zrFzxutppvxEdjyu4QNjogBMjtC1py9Hp1S",
							},
						},
					},
					{
						Value: 1.25,
						Index: 2,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "OP_HASH160 fc1d7f04db5e2c05b051e0decc85effe6bc539d5 OP_EQUAL", // nolint
							Hex:  "a914fc1d7f04db5e2c05b051e0decc85effe6bc539d587",         // nolint
							RequiredSigs: 1,
							Type: "scripthash",
							Addresses: []string{
								"zrS7QUB2eDbbKvyP43VJys3t7RpojW8GdxH",
							},
						},
					},
					{
						Value: 1.25,
						Index: 3,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "OP_HASH160 8b85fc1e171a4c7994c088b91d5a75dff9e56cad OP_EQUAL", // nolint
							Hex:  "a9148b85fc1e171a4c7994c088b91d5a75dff9e56cad87",         // nolint
							RequiredSigs: 1,
							Type: "scripthash",
							Addresses: []string{
								"zrFr5HVm7woVq3oFzkMEdJdbfBchfPAPDsP",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000031afda1ec75afe8f9c163059ded874fdfcd8ea8db513f2d36fff310c235f50194000000006a473044022059135f673a4919ab56775064cc82080ead1c74d8f0ebd943062b247c5946cf88022048f26c94a15752fa04d8bfff7388dd65d57485acd2395e539a50b2ca8e278700012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ceffffffff3767ef09fac4ef1f2b9b9fd26b9fa10657d03b9495bfb68c7d234eec02fee814000000006a4730440220527c59b1d2dbb87b71e01c9d1489f110727fc3120e5306539bd4668ed1063d30022079b6ca4ff77de3ab953bb0d896b74bb60c8ceca28248340201e701da0d1fd12b012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ceffffffffbc4619665eed136bc292fdbf963767b6627e2165876fa1482d4fe9a09b2f294c010000006a47304402202d3b75ed231c1fe478c471452a0385c5cdc9fe2e337d5ee62cacd8a26d013e5002207d864a38e013d8c61b1972bd7bf78a53accd9b8d600fbbd7c79c21b2171fd8cb012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ceffffffff020065cd1d000000003f76a914b87cc09d17751ffeab924a82134665ae4202cbfc88ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab4f023e398010000003f76a914fd2831ec8fc1bf3ccdeadbe9fcdb515aac90476188ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab400000000", // nolint
				Hash:     "67c76a34cb6bde6f9628fdc8348c23191d3222e88386ed05c97e3c63384a01af",
				Size:     595,
				Version:  1,
				Locktime: 0,
				Inputs:   []*Input{
					{
						TxHash: "9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a",
						Vout:      0,
						ScriptSig: &ScriptSig{
							ASM: "3044022059135f673a4919ab56775064cc82080ead1c74d8f0ebd943062b247c5946cf88022048f26c94a15752fa04d8bfff7388dd65d57485acd2395e539a50b2ca8e27870001 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce",
							Hex: "473044022059135f673a4919ab56775064cc82080ead1c74d8f0ebd943062b247c5946cf88022048f26c94a15752fa04d8bfff7388dd65d57485acd2395e539a50b2ca8e278700012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce",
						},
						Sequence:  4294967295,
					},
					{
						TxHash: "14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737",
						Vout:      0,
						ScriptSig: &ScriptSig{
							ASM: "30440220527c59b1d2dbb87b71e01c9d1489f110727fc3120e5306539bd4668ed1063d30022079b6ca4ff77de3ab953bb0d896b74bb60c8ceca28248340201e701da0d1fd12b01 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce",
							Hex: "4730440220527c59b1d2dbb87b71e01c9d1489f110727fc3120e5306539bd4668ed1063d30022079b6ca4ff77de3ab953bb0d896b74bb60c8ceca28248340201e701da0d1fd12b012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce",
						},
						Sequence:  4294967295,
					},
					{
						TxHash: "4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc",
						Vout:      1,
						ScriptSig: &ScriptSig{
							ASM: "304402202d3b75ed231c1fe478c471452a0385c5cdc9fe2e337d5ee62cacd8a26d013e5002207d864a38e013d8c61b1972bd7bf78a53accd9b8d600fbbd7c79c21b2171fd8cb01 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce",
							Hex: "47304402202d3b75ed231c1fe478c471452a0385c5cdc9fe2e337d5ee62cacd8a26d013e5002207d864a38e013d8c61b1972bd7bf78a53accd9b8d600fbbd7c79c21b2171fd8cb012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce",
						},
						Sequence:  4294967295,
					},
				}, // all we care about in this test is the outputs
				Outputs: []*Output{
					{
						Value: 5,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 b87cc09d17751ffeab924a82134665ae4202cbfc OP_EQUALVERIFY OP_CHECKSIG bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a00 717682 OP_CHECKBLOCKATHEIGHT",
							Hex:          "76a914b87cc09d17751ffeab924a82134665ae4202cbfc88ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab4",
							RequiredSigs: 1,
							Type:         "pubkeyhashreplay",
							Addresses: []string{
								"ztjySYJL8g9i6wc2YTusbDpPZSpPM5xuTua",
							},
						},
					},
					{
						Value: 68.5999,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "OP_DUP OP_HASH160 fd2831ec8fc1bf3ccdeadbe9fcdb515aac904761 OP_EQUALVERIFY OP_CHECKSIG bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a00 717682 OP_CHECKBLOCKATHEIGHT",
							Hex:  "76a914fd2831ec8fc1bf3ccdeadbe9fcdb515aac90476188ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab4",
							RequiredSigs: 1,
							Type: "pubkeyhashreplay",
							Addresses: []string{
								"ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
							},
						},
					},
				},
			},
		},
	}

	blockIdentifier717984 = &types.BlockIdentifier{
		Hash:  "005d3821c522b528f42fa16187d70ccb59170e2dcd72e9242d54d967e63b6fxe",
		Index: 717984,
	}

	block717984 = &Block{
		Hash:              "005d3821c522b528f42fa16187d70ccb59170e2dcd72e9242d54d967e63b6fxe",
		Height:            717984,
		PreviousBlockHash: "0067f80ce10d4255932b7f8c9baf7bd0dcfd408c312d33144be0ea12caf7f7f0",
		Time:              1601465727,
		Size:              2271,
		Version:           3,
		MerkleRoot:        "97c960c90e0b6bc30d2629f06d114f1c49aadb0e3d9bd70eb4f0f9ed1ea69279",
		Nonce:             "00002e570d64b4b3ea1c30dec68b2dff255eb3148656f06f5e018ae739a400eb",
		Bits:              "1f754920",
		Difficulty:        17.46160923,
		Txs: []*Transaction{
			{
				Hex:      "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff06039ff40a0102ffffffff04903eb42c000000001976a914557662a6b307f95aa00311c074f7feebb955d45188ac80b2e60e0000000017a9148d3468b6686ac59caf9ad94e547a737b09fa102787405973070000000017a914fc1d7f04db5e2c05b051e0decc85effe6bc539d587405973070000000017a9148b85fc1e171a4c7994c088b91d5a75dff9e56cad8700000000", // nolint
				Hash:     "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88",
				Size:     187,
				Version:  1,
				Locktime: 0,
				Inputs: []*Input{
					{
						Coinbase: "039ff40a0102",
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 7.5001,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "OP_DUP OP_HASH160 557662a6b307f95aa00311c074f7feebb955d451 OP_EQUALVERIFY OP_CHECKSIG", // nolint
							Hex:  "76a914557662a6b307f95aa00311c074f7feebb955d45188ac",         // nolint
							RequiredSigs: 1,
							Type: "pubkeyhash",
							Addresses: []string{
								"ztawr1vEZ6pZRtLqNy2C9u7EK7JN2gP8W6z",
							},
						},
					},
					{
						Value: 2.5,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "OP_HASH160 8d3468b6686ac59caf9ad94e547a737b09fa1027 OP_EQUAL", // nolint
							Hex:  "a9148d3468b6686ac59caf9ad94e547a737b09fa102787",         // nolint
							RequiredSigs: 1,
							Type: "scripthash",
							Addresses: []string{
								"zrFzxutppvxEdjyu4QNjogBMjtC1py9Hp1S",
							},
						},
					},
					{
						Value: 1.25,
						Index: 2,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "OP_HASH160 fc1d7f04db5e2c05b051e0decc85effe6bc539d5 OP_EQUAL", // nolint
							Hex:  "a914fc1d7f04db5e2c05b051e0decc85effe6bc539d587",         // nolint
							RequiredSigs: 1,
							Type: "scripthash",
							Addresses: []string{
								"zrS7QUB2eDbbKvyP43VJys3t7RpojW8GdxH",
							},
						},
					},
					{
						Value: 1.25,
						Index: 3,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "OP_HASH160 8b85fc1e171a4c7994c088b91d5a75dff9e56cad OP_EQUAL", // nolint
							Hex:  "a9148b85fc1e171a4c7994c088b91d5a75dff9e56cad87",         // nolint
							RequiredSigs: 1,
							Type: "scripthash",
							Addresses: []string{
								"zrFr5HVm7woVq3oFzkMEdJdbfBchfPAPDsP",
							},
						},
					},
				},
			},
			{
				Hex:      "01000000031afda1ec75afe8f9c163059ded874fdfcd8ea8db513f2d36fff310c235f50194000000006a473044022059135f673a4919ab56775064cc82080ead1c74d8f0ebd943062b247c5946cf88022048f26c94a15752fa04d8bfff7388dd65d57485acd2395e539a50b2ca8e278700012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ceffffffff3767ef09fac4ef1f2b9b9fd26b9fa10657d03b9495bfb68c7d234eec02fee814000000006a4730440220527c59b1d2dbb87b71e01c9d1489f110727fc3120e5306539bd4668ed1063d30022079b6ca4ff77de3ab953bb0d896b74bb60c8ceca28248340201e701da0d1fd12b012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ceffffffffbc4619665eed136bc292fdbf963767b6627e2165876fa1482d4fe9a09b2f294c010000006a47304402202d3b75ed231c1fe478c471452a0385c5cdc9fe2e337d5ee62cacd8a26d013e5002207d864a38e013d8c61b1972bd7bf78a53accd9b8d600fbbd7c79c21b2171fd8cb012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ceffffffff020065cd1d000000003f76a914b87cc09d17751ffeab924a82134665ae4202cbfc88ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab4f023e398010000003f76a914fd2831ec8fc1bf3ccdeadbe9fcdb515aac90476188ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab400000000", // nolint
				Hash:     "67c76a34cb6bde6f9628fdc8348c23191d3222e88386ed05c97e3c63384a01af",
				Size:     595,
				Version:  1,
				Locktime: 0,
				Inputs:   []*Input{
					{
						TxHash: "9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a",
						Vout:      0,
						ScriptSig: &ScriptSig{
							ASM: "3044022059135f673a4919ab56775064cc82080ead1c74d8f0ebd943062b247c5946cf88022048f26c94a15752fa04d8bfff7388dd65d57485acd2395e539a50b2ca8e27870001 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce",
							Hex: "473044022059135f673a4919ab56775064cc82080ead1c74d8f0ebd943062b247c5946cf88022048f26c94a15752fa04d8bfff7388dd65d57485acd2395e539a50b2ca8e278700012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce",
						},
						Sequence:  4294967295,
					},
					{
						TxHash: "14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737",
						Vout:      0,
						ScriptSig: &ScriptSig{
							ASM: "30440220527c59b1d2dbb87b71e01c9d1489f110727fc3120e5306539bd4668ed1063d30022079b6ca4ff77de3ab953bb0d896b74bb60c8ceca28248340201e701da0d1fd12b01 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce",
							Hex: "4730440220527c59b1d2dbb87b71e01c9d1489f110727fc3120e5306539bd4668ed1063d30022079b6ca4ff77de3ab953bb0d896b74bb60c8ceca28248340201e701da0d1fd12b012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce",
						},
						Sequence:  4294967295,
					},
					{
						TxHash: "4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc",
						Vout:      1,
						ScriptSig: &ScriptSig{
							ASM: "304402202d3b75ed231c1fe478c471452a0385c5cdc9fe2e337d5ee62cacd8a26d013e5002207d864a38e013d8c61b1972bd7bf78a53accd9b8d600fbbd7c79c21b2171fd8cb01 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce",
							Hex: "47304402202d3b75ed231c1fe478c471452a0385c5cdc9fe2e337d5ee62cacd8a26d013e5002207d864a38e013d8c61b1972bd7bf78a53accd9b8d600fbbd7c79c21b2171fd8cb012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce",
						},
						Sequence:  4294967295,
					},
				}, // all we care about in this test is the outputs
				Outputs: []*Output{
					{
						Value: 5,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 b87cc09d17751ffeab924a82134665ae4202cbfc OP_EQUALVERIFY OP_CHECKSIG bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a00 717682 OP_CHECKBLOCKATHEIGHT",
							Hex:          "76a914b87cc09d17751ffeab924a82134665ae4202cbfc88ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab4",
							RequiredSigs: 1,
							Type:         "pubkeyhashreplay",
							Addresses: []string{
								"ztjySYJL8g9i6wc2YTusbDpPZSpPM5xuTua",
							},
						},
					},
					{
						Value: 68.5999,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "OP_DUP OP_HASH160 fd2831ec8fc1bf3ccdeadbe9fcdb515aac904761 OP_EQUALVERIFY OP_CHECKSIG bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a00 717682 OP_CHECKBLOCKATHEIGHT",
							Hex:  "76a914fd2831ec8fc1bf3ccdeadbe9fcdb515aac90476188ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab4",
							RequiredSigs: 1,
							Type: "pubkeyhashreplay",
							Addresses: []string{
								"ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
							},
						},
					},
				},
			},
		},
		Certs: []*Certificate{
			{
				Hash:     "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e89",
				Version:  1,
				Inputs:   []*Input{
					{
						TxHash: "62091923e9805a8650d752b3b83e0d56ce70e775ee67c080feade7e5ee677ad9",
						Vout:      0,
						ScriptSig: &ScriptSig{
							ASM: "3044022014d8dee1da3821dce95e48060f8f38394aee00f84d03a8203611ff3e703c10a002205ce62cffdc12dd26742489120d50d071ff08f993b9cca0b31a73e0f20f20cb5d01 0241b92fed18a3ded2b98459b5432982a0712912ad86b929ec6feb19655824b7cc",
							Hex: "473044022014d8dee1da3821dce95e48060f8f38394aee00f84d03a8203611ff3e703c10a002205ce62cffdc12dd26742489120d50d071ff08f993b9cca0b31a73e0f20f20cb5d01210241b92fed18a3ded2b98459b5432982a0712912ad86b929ec6feb19655824b7cc",
						},
						Sequence:  4294967295,
					},
				},
				Cert: &Cert{
					Scid:                        "2f1f1b22ef02396fcb5fcff08915767b57206d3dffca92b211ae4eed3c5f1db7",
					EpochNumber:                 0,
					Quality:                     3,
					EndEpochCumScTxCommTreeRoot: "695becea24179b26ef68a98ef31b4a0550da05c1017c48252520412bb1dd552c",
					ScProof:                     "0202904752173865bebc0bd372227a68edd2b3aa8aa84367e0f213a1022d0ccf102080a709dae6a636502c3c51de17666923e9b4e831cbfe88806cccdf9655d7cf3b2400000285427df898179ffae71a48423f2ae2521a8a7d2700116ec97d26c3321ad6ff2600a43d09da18984c5a977d4f32a80b917aead29b410b3c6ad5bfa835324c6bc70d8000022bde3e3e649f7ce94ecd4f7603015f9daa56ee2e9fe2247d9d77da23d5c0763a00a43d09da18984c5a977d4f32a80b917aead29b410b3c6ad5bfa835324c6bc70d0000026eb5db3716925e48180ce2702c88c2784c59298c57a7accd2b056ef5913803248094c4db317a8f83b3aa08ffef9feab7c9ba0d6d408d507206a8246fc961a21d3e800002f37e2b27410218c6df9114586fe42de4cc03a2ba407cba8418ba4527afc920028045151836e4b72e563b07315e14d8e1536e77b70247b2702a174a59450f4b862780000466b9f31f49bda90b1cc8aed4432b67dde9a86761a703c062c154592bbe4a2a0d807099f56c454c7204991f8ee6051144e006dafaffceca902d3edf20be6e6bde0100c54b2164aba82689a707e87e3c4a7a56fc7d1b16cf42ec4975d1834f34a0751580c20cc7c7f8ca3e10a278a5c720eef78df837323f7bb2c741b901587f28043000800004372c642bc2fb49b9306afd2bfd74bf05e0395000c6b69eaef9bb1aab9a08ef0900213d5783027660cc5387a14c93c9e804414d0fbec998c4bb2ca1a03f4b4b083e802cb58576c6a8a502d0556d3d8c79d5ae7ded0a1bd09b58a4ff6eaafe04260702801019fbbeade638586385acaa63f1f8782dd470c566d035de8aaab2c08c71fb3200000804a21c358f43c072650f4f5145f9e0258c19994f5a8f8746e197d2c0e8b2c716003a721a33ec908b12e180ef205ae94da0ea42e9519ab7e83428d49e8cb3a0ec0e00e58efd81003c37e3798b43eddb7022c8bcfaa5284881e7aefb7645f1d9da7a20006c7ffc231c5b46c273e049b530471e7ec195b69a8f333eb9425da2fd154bc832800b30685b5c7bd1bf8142674056b5a316b1981caf85ae79a58ca4492620bc1d20001813007c7f60b0d495a7686585163df1a899c6c4800219ef7cbba688130ce20f00dc3d2f4c56e447056d11088160f2d7e72a20b76234f4a39d7cde82795bf2e81780b9b62117d09f5b60faafb6c2507918f2f76926baa135e470ff83672454e31f1280002c098dbb08397cf4aa53065a456715945695a7e3fe110d11d2f3a2d1db50af3d241305728391639953707c9a723d8ceb2e06cce0e0de05fb46aedb460fb98723594be628a5a48ad275173967d33cd73261a55ab9da28837b906d3a346287ab0f48b4831037278effa45da4a676e7514c86e4ef4dfd7956ba284aa5e1c625ff03edb0c9e25b28f8a24b6fed8b695f72403f7c6b9140b2b9b01cbc51db7e99131b577b0e9762d3ae58da67a7cb1da9e3fb94ba1105b8a20ba139e1e0843bde863dfd6eb9f4509ab67841f4404b8695afad95a69425a9c2b70d15d0b23c0d65a9255704efbfd56203b43a2e505a67cc9500920bc5f15245be4671f432338087c03701000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007bfa27460ccd311113e78e913838d2b968345b3c1a92c96bc2ca02caef29571af5977747082b8877d9d1a253797c5f6bac9efc6a789c1a9ceb433bbb4620cb254e70c386f7092a380113f8943701532c3799ddd0e404651d8e18625a825a7d33b589048566f56013b40f8971f4776033ce7b3926965db122c7b5e6c9520a943156a2cdceafb743f39a77acf056c301aca51731fc8fb644a1c60f6056fde21a3bdfc0dd715b94e81148af173e521816e27ba64e57cb38024d5db06b1aa8437416c90520d0619a5f43485e4bb7af3452d68162127ffafe946eadc8b095caba3a3572f7e39073618bbfbc14f457e827052f22e9021d05fa3099f36e0081f4f7e2241830e52f5adbcb7727d4fc63f173229b0e6926dddf0944bed52382dad2c0b41beacf1ad08789742992854208362d6868f196d92220f6bb412adc7d252d3f4b2409b334491e898d0aded26ed833f796e761264510e590e416b1ec1c56ab819e530380f0815290b10b569460756008773f3bcc7e7dcbdc82253419d4187a8c8d64950f8025019e89f96fd21220175c76e00eba5bc137cc0040860c41e7adc9a1f40e580c00741e9eddcdb7c2251228e0b942b5c080db2b3ab209e34994b208037820b5972780f6a51909c1a359793905a438a756df510d7a7e58fc69fa3bdae39af268522e3c009a32012723088edb87c5244d0698be9d1d037bcca7bd08b6eada3dff0fb8ac150041db23f4f8c0be03743d1440110eb2e1745df66950b52ab8e2a0dd2d467757350052dd493f29450c8ff346584557597ae1511c4302099730b9daf877857f233230807dec2ccc48cc0520353fddfba517d56d826b3d10e45def9ffa029596d753631b80d6ea74ad0bd894b5dbfdf6b988934f9b4b71c211a2ad92d19949c75d660d983580b02d89efeeb00d068ac551afe287daac217fe3d86a703120357646dd8e00f13f807a9c5d47f45bc559eebb23e4b839a9b858ce118129f8ae43147b4b5d9ff0283d8082c939b2cb977ae8692aaf2539af81c91dd529f33f8508b6e772bec792a47b2800071c054a7e7ac694e3a64a897c14bd67d64ed8d7ce521e4c147885c6d31bbe2980013e97b2d322be251600ac3c696db2c9d47c6358affc7511f788e9ac342bfe3e00ad5316f20f3bd775f1ff8f59cfdc9035e1fef6b2251ddaf126598460930e8e0980cabfd9c4b7d1a13be4dfd9fb696ed63968408f690c7fb7282b460fb63b03b51680ae49cf365ef6c73b4e44535772352d9b6f65eb856735d1eddcea075eb12e6407800abafa92b85d1d01dd9ad2d0a318a033783ad65a46e5742d4b72bbc0ff13ff1a00e2a63170a86d48d772f2fc170635e5d4b9ea43369b39f69a8e7b1fe992a1652a00000832f50c2508b14a836745de1e7f458c846f584102029f654fc8e09e8851f25f3000d9da6339bd2da525d0db2372daecc0f8133004e438b2c2b0187380692745363800c4d082a5f075ab7669157b078b929d38e21a173eafbc1fad2bb19ad53220ba2c0062a9a8daf498afbca88efadcbe3f8416aa94b2b29af1c7ee66a1890b57566c0b00e51cb31b11a0db013ce33fc72c7f787a9ab3640d1e3effb27dc94b7f8e862324807f6126ddde2dd9e32c7fc9bebae73cbdec96bb6e862d681340845d1b8289743e00322f613fe364e99b46732d2aff3ac605c3df7a819caf48cc9021adaa185ad4338030a66acc013c768006367bfc9dab30e4d7630d44e59dbe30a2a39f806906c13d00147364065231024fad4054d477e200df46a7decc9606c83ffde210bbb6b2551f397966b77c39fa2f85559cb218d824ae87814c0d11bb5ddc8af616a673da96b30198f399d1e82d8c866d371c5ab959f26fb7e1b99fe769776d0dc798186fc4ec3283e4510d913ba1a35e7f4007f3d8d37cb1b5d654b09c2e02829b118581f417011f030e2b25c8cd73c03ae6494cbc9d0da62832c61f093f12d2fd895fdc9a622ed425832ef7e8fe1b176160f06e8520d88cf63e736e42c0bbf72c06613876951719068def3a75bd2ba4b8f5141eadcb7bab5552a0b5b765259c9cf5ca0fb5283b953e198ca1efd4012b2fcc5393b25fe6620d2d232547d5ea1429065d2893923001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009f5553d0b388af509300d71c3b322fb5404133d74ab2de7a2f45e50ca7fee83a88179c0ace01549f66456be7621705ae9dc0ba766f140232211cce5306683522fde38862f55f8f45159a63fcfd478d1433596f88a066a310ca6f490343ff883da8e50b9bca33f40d850c063d71e571fff1cdc373b007d52adef46061bd1bf32ee1ccb6a353c1d3df07a0166f5a11cd3c185fe3d714d36ddb6d98c97458d90422fd297e1d1d87e01a8caa87e6f0c7e8ce76c13ed89176f8130f3ebd3ad57de427c438530af4ab738fb7cd75a100dab13af8fcb57aa687b2e2da67610c7d013e393ec7acf5edb8cc11028cc9ca26c7d8c807034a8559784d1d25989ef382fec106",
					VFieldElementCertificateField: []string{
						"ab000100",
						"ccccdddd0000",
						"0100",
					},
					VBitVectorCertificateField: []string{
						"021f8b08000000000002ff017f0080ff44c7e21ba1c7c0a29de006cb8074e2ba39f15abfef2525a4cbb3f235734410bda21cdab6624de769ceec818ac6c2d3a01e382e357dce1f6e9a0ff281f0fedae0efe274351db37599af457984dcf8e3ae4479e0561341adfff4746fbe274d90f6f76b8a2552a6ebb98aee918c7ceac058f4c1ae0131249546ef5e22f4187a07da02ca5b7f000000",
					},
					FtScFee: 10.00000000,
					MbtrScFee: 20.00000000,
				},
				Outputs: []*Output{
					{
						Value: 0.24895145,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 ec54fedd6a312d5c536046323bfabb9d2a475d7a OP_EQUALVERIFY OP_CHECKSIG 4ca064b46515f3f00e846e6c1b45ef36a082ea786783096d2cb6169556756e08 21 OP_CHECKBLOCKATHEIGHT",
							Hex:          "76a914ec54fedd6a312d5c536046323bfabb9d2a475d7a88ac204ca064b46515f3f00e846e6c1b45ef36a082ea786783096d2cb6169556756e080115b4",
							RequiredSigs: 1,
							Type:         "pubkeyhashreplay",
							Addresses: []string{
								"ztpha3vQzv7eTdBvPC1oWnouuManmCEVbTT",
							},
						},
					},
					{
						Value: 1.00000000,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "OP_DUP OP_HASH160 4aeea9b9beec0af6eb8e6e8d6015a8a679590553 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:  "76a9144aeea9b9beec0af6eb8e6e8d6015a8a67959055388ac",
							RequiredSigs: 1,
							Type: "pubkeyhash",
							Addresses: []string{
								"ztZzAfqxzua7EDHUMFq6hpQPhXyC1XPJMUs",
							},
						},
						BackwardTransfer: true,
					},
				},
			},
		},
	}
)


func TestNetworkStatus(t *testing.T) {
	tests := map[string]struct {
		responses []responseFixture

		expectedStatus *types.NetworkStatusResponse
		expectedError  error
	}{
		"successful": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_blockchain_info_response.json"),
					url:    url,
				},
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response.json"),
					url:    url,
				},
				{
					status: http.StatusOK,
					body:   loadFixture("get_peer_info_response.json"),
					url:    url,
				},
			},
			expectedStatus: &types.NetworkStatusResponse{
				CurrentBlockIdentifier: blockIdentifier717983,
				CurrentBlockTimestamp:  block717983.Time * 1000,
				GenesisBlockIdentifier: MainnetGenesisBlockIdentifier,
				Peers: []*types.Peer{
					{
						PeerID: "77.93.223.9:8333",
						Metadata: forceMarshalMap(t, &PeerInfo{
							Addr:           "77.93.223.9:8333",
							Version:        70015,
							SubVer:         "/Satoshi:0.14.2/",
							StartingHeight: 643579,
							RelayTxes:      true,
							LastSend:       1597606676,
							LastRecv:       1597606677,
							BanScore:       0,
							SyncedHeaders:  644046,
							SyncedBlocks:   644046,
						}),
					},
					{
						PeerID: "172.105.93.179:8333",
						Metadata: forceMarshalMap(t, &PeerInfo{
							Addr:           "172.105.93.179:8333",
							RelayTxes:      true,
							LastSend:       1597606678,
							LastRecv:       1597606676,
							Version:        70015,
							SubVer:         "/Satoshi:0.18.1/",
							StartingHeight: 643579,
							BanScore:       0,
							SyncedHeaders:  644046,
							SyncedBlocks:   644046,
						}),
					},
				},
			},
		},
		"blockchain warming up error": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("rpc_in_warmup_response.json"),
					url:    url,
				},
			},
			expectedError: errors.New("rpc in warmup"),
		},
		"blockchain info error": {
			responses: []responseFixture{
				{
					status: http.StatusInternalServerError,
					body:   "{}",
					url:    url,
				},
			},
			expectedError: errors.New("invalid response: 500 Internal Server Error"),
		},
		"peer info not accessible": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_blockchain_info_response.json"),
					url:    url,
				},
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response.json"),
					url:    url,
				},
				{
					status: http.StatusInternalServerError,
					body:   "{}",
					url:    url,
				},
			},
			expectedError: errors.New("invalid response: 500 Internal Server Error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			responses := make(chan responseFixture, len(test.responses))
			for _, response := range test.responses {
				responses <- response
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := <-responses
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("POST", r.Method)
				assert.Equal(response.url, r.URL.RequestURI())

				w.WriteHeader(response.status)
				fmt.Fprintln(w, response.body)
			}))

			client := NewClient(ts.URL, MainnetGenesisBlockIdentifier, MainnetCurrency)
			status, err := client.NetworkStatus(context.Background())
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedStatus, status)
			}
		})
	}
}

func TestGetRawBlock(t *testing.T) {
	tests := map[string]struct {
		blockIdentifier *types.PartialBlockIdentifier
		responses       []responseFixture

		expectedBlock *Block
		expectedCoins []string
		expectedError error
	}{
		"lookup by hash": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Hash: &blockIdentifier717983.Hash,
			},
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response.json"),
					url:    url,
				},
			},
			expectedBlock: block717983,
			expectedCoins: []string{"9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0", "14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0", "4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1"},
		},
		"lookup by hash - block with certificate": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Hash: &blockIdentifier717984.Hash,
			},
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response_with_certificate.json"),
					url:    url,
				},
			},
			expectedBlock: block717984,
			expectedCoins: []string{"9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0", "14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0", "4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1","62091923e9805a8650d752b3b83e0d56ce70e775ee67c080feade7e5ee677ad9:0"},
		},
		"lookup by hash (get block api error)": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Hash: &blockIdentifier717983.Hash,
			},
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_not_found_response.json"),
					url:    url,
				},
			},
			expectedError: ErrBlockNotFound,
		},
		"lookup by hash (get block internal error)": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Hash: &blockIdentifier717983.Hash,
			},
			responses: []responseFixture{
				{
					status: http.StatusInternalServerError,
					body:   "{}",
					url:    url,
				},
			},
			expectedBlock: nil,
			expectedError: errors.New("invalid response: 500 Internal Server Error"),
		},
		"lookup by index": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Index: &blockIdentifier717983.Index,
			},
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_hash_response.json"),
					url:    url,
				},
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response.json"),
					url:    url,
				},
			},
			expectedBlock: block717983,
			expectedCoins: []string{"9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0", "14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0", "4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1"},
		},
		"lookup by index (out of range)": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Index: &blockIdentifier717983.Index,
			},
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_hash_out_of_range_response.json"),
					url:    url,
				},
			},
			expectedError: errors.New("height out of range"),
		},
		"current block lookup": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_blockchain_info_response.json"),
					url:    url,
				},
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response.json"),
					url:    url,
				},
			},
			expectedBlock: block717983,
			expectedCoins:[]string{"9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0", "14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0", "4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1"},
		},
		"current block lookup (can't get current info)": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("rpc_in_warmup_response.json"),
					url:    url,
				},
			},
			expectedError: errors.New("unable to get blockchain info"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			responses := make(chan responseFixture, len(test.responses))
			for _, response := range test.responses {
				responses <- response
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := <-responses
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("POST", r.Method)
				assert.Equal(response.url, r.URL.RequestURI())

				w.WriteHeader(response.status)
				fmt.Fprintln(w, response.body)
			}))

			client := NewClient(ts.URL, MainnetGenesisBlockIdentifier, MainnetCurrency)
			block, coins, err := client.GetRawBlock(context.Background(), test.blockIdentifier)
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedBlock, block)
				assert.Equal(test.expectedCoins, coins)
			}
		})
	}
}

func int64Pointer(v int64) *int64 {
	return &v
}

func mustMarshalMap(v interface{}) map[string]interface{} {
	m, _ := types.MarshalMap(v)
	return m
}

func TestParseBlockRegtest(t *testing.T) {
	tests := map[string]struct {
		block *Block
		coins map[string]*storage.AccountCoin

		expectedBlock *types.Block
		expectedError error
	}{
		"block717983": {
			block: block717983,
			coins: map[string]*storage.AccountCoin{
				"9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0": {
					Account: &types.AccountIdentifier{
						Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0",
						},
						Amount: &types.Amount{
							Value:    "60000000",
							Currency: MainnetCurrency,
						},
					},
				},
				"14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0": {
					Account: &types.AccountIdentifier{
						Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0",
						},
						Amount: &types.Amount{
							Value:    "200000000",
							Currency: MainnetCurrency,
						},
					},
				},
				"4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1": {
					Account: &types.AccountIdentifier{
						Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1",
						},
						Amount: &types.Amount{
							Value:    "7100000000",
							Currency: MainnetCurrency,
						},
					},
				},
			},
			expectedBlock: &types.Block{
				BlockIdentifier: blockIdentifier717983,
				ParentBlockIdentifier: &types.BlockIdentifier{
					Hash:  "0067f80ce10d4255932b7f8c9baf7bd0dcfd408c312d33144be0ea12caf7f7f0",
					Index: 717982,
				},
				Timestamp: 1601465727000,
				Transactions: []*types.Transaction{
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   CoinbaseOpType,
								Status: SuccessStatus,
								Metadata: mustMarshalMap(&OperationMetadata{
									Coinbase: "039ff40a0102",
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "ztawr1vEZ6pZRtLqNy2C9u7EK7JN2gP8W6z", // nolint
								},
								Amount: &types.Amount{
									Value:    "750010000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_DUP OP_HASH160 557662a6b307f95aa00311c074f7feebb955d451 OP_EQUALVERIFY OP_CHECKSIG", // nolint
										Hex:  "76a914557662a6b307f95aa00311c074f7feebb955d45188ac",         // nolint
										RequiredSigs: 1,
										Type: "pubkeyhash",
										Addresses: []string{
											"ztawr1vEZ6pZRtLqNy2C9u7EK7JN2gP8W6z",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "zrFzxutppvxEdjyu4QNjogBMjtC1py9Hp1S", // nolint
								},
								Amount: &types.Amount{
									Value:    "250000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_HASH160 8d3468b6686ac59caf9ad94e547a737b09fa1027 OP_EQUAL", // nolint
										Hex:  "a9148d3468b6686ac59caf9ad94e547a737b09fa102787",         // nolint
										RequiredSigs: 1,
										Type: "scripthash",
										Addresses: []string{
											"zrFzxutppvxEdjyu4QNjogBMjtC1py9Hp1S",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        3,
									NetworkIndex: int64Pointer(2),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "zrS7QUB2eDbbKvyP43VJys3t7RpojW8GdxH", // nolint
								},
								Amount: &types.Amount{
									Value:    "125000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88:2",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_HASH160 fc1d7f04db5e2c05b051e0decc85effe6bc539d5 OP_EQUAL", // nolint
										Hex:  "a914fc1d7f04db5e2c05b051e0decc85effe6bc539d587",         // nolint
										RequiredSigs: 1,
										Type: "scripthash",
										Addresses: []string{
											"zrS7QUB2eDbbKvyP43VJys3t7RpojW8GdxH",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        4,
									NetworkIndex: int64Pointer(3),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "zrFr5HVm7woVq3oFzkMEdJdbfBchfPAPDsP", // nolint
								},
								Amount: &types.Amount{
									Value:    "125000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88:3",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_HASH160 8b85fc1e171a4c7994c088b91d5a75dff9e56cad OP_EQUAL", // nolint
										Hex:  "a9148b85fc1e171a4c7994c088b91d5a75dff9e56cad87",         // nolint
										RequiredSigs: 1,
										Type: "scripthash",
										Addresses: []string{
											"zrFr5HVm7woVq3oFzkMEdJdbfBchfPAPDsP",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    187,
							Version: 1,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "67c76a34cb6bde6f9628fdc8348c23191d3222e88386ed05c97e3c63384a01af",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: SuccessStatus,
								Amount: &types.Amount{
									Value:    "-60000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3044022059135f673a4919ab56775064cc82080ead1c74d8f0ebd943062b247c5946cf88022048f26c94a15752fa04d8bfff7388dd65d57485acd2395e539a50b2ca8e27870001 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
										Hex: "473044022059135f673a4919ab56775064cc82080ead1c74d8f0ebd943062b247c5946cf88022048f26c94a15752fa04d8bfff7388dd65d57485acd2395e539a50b2ca8e278700012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(1),
								},
								Type:   InputOpType,
								Status: SuccessStatus,
								Amount: &types.Amount{
									Value:    "-200000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "30440220527c59b1d2dbb87b71e01c9d1489f110727fc3120e5306539bd4668ed1063d30022079b6ca4ff77de3ab953bb0d896b74bb60c8ceca28248340201e701da0d1fd12b01 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
										Hex: "4730440220527c59b1d2dbb87b71e01c9d1489f110727fc3120e5306539bd4668ed1063d30022079b6ca4ff77de3ab953bb0d896b74bb60c8ceca28248340201e701da0d1fd12b012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(2),
								},
								Type:   InputOpType,
								Status: SuccessStatus,
								Amount: &types.Amount{
									Value:    "-7100000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "304402202d3b75ed231c1fe478c471452a0385c5cdc9fe2e337d5ee62cacd8a26d013e5002207d864a38e013d8c61b1972bd7bf78a53accd9b8d600fbbd7c79c21b2171fd8cb01 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
										Hex: "47304402202d3b75ed231c1fe478c471452a0385c5cdc9fe2e337d5ee62cacd8a26d013e5002207d864a38e013d8c61b1972bd7bf78a53accd9b8d600fbbd7c79c21b2171fd8cb012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        3,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "ztjySYJL8g9i6wc2YTusbDpPZSpPM5xuTua",
								},
								Amount: &types.Amount{
									Value:    "500000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "67c76a34cb6bde6f9628fdc8348c23191d3222e88386ed05c97e3c63384a01af:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_DUP OP_HASH160 b87cc09d17751ffeab924a82134665ae4202cbfc OP_EQUALVERIFY OP_CHECKSIG bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a00 717682 OP_CHECKBLOCKATHEIGHT",
										Hex:  "76a914b87cc09d17751ffeab924a82134665ae4202cbfc88ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab4",
										RequiredSigs: 1,
										Type: "pubkeyhashreplay",
										Addresses: []string{
											"ztjySYJL8g9i6wc2YTusbDpPZSpPM5xuTua",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        4,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
								},
								Amount: &types.Amount{
									Value:    "6859990000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "67c76a34cb6bde6f9628fdc8348c23191d3222e88386ed05c97e3c63384a01af:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_DUP OP_HASH160 fd2831ec8fc1bf3ccdeadbe9fcdb515aac904761 OP_EQUALVERIFY OP_CHECKSIG bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a00 717682 OP_CHECKBLOCKATHEIGHT",
										Hex:  "76a914fd2831ec8fc1bf3ccdeadbe9fcdb515aac90476188ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab4",
										RequiredSigs: 1,
										Type: "pubkeyhashreplay",
										Addresses: []string{
											"ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    595,
							Version: 1,
						}),
					},
				},
				Metadata: mustMarshalMap(&BlockMetadata{
					Size:       2271,
					Version:    3,
					MerkleRoot: "97c960c90e0b6bc30d2629f06d114f1c49aadb0e3d9bd70eb4f0f9ed1ea69279",
					Nonce:      "00002e570d64b4b3ea1c30dec68b2dff255eb3148656f06f5e018ae739a400eb",
					Bits:       "1f754920",
					Difficulty: 17.46160923,
				}),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			client := NewClient("", RegtestGenesisBlockIdentifier, MainnetCurrency)
			block, err := client.ParseBlock(context.Background(), test.block, test.coins)
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedBlock, block)
			}
		})
	}
}

func TestParseBlock(t *testing.T) {
	tests := map[string]struct {
		block *Block
		coins map[string]*storage.AccountCoin

		expectedBlock *types.Block
		expectedError error
	}{
		"block717983": {
			block: block717983,
			coins: map[string]*storage.AccountCoin{
				"9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0": {
					Account: &types.AccountIdentifier{
						Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0",
						},
						Amount: &types.Amount{
							Value:    "60000000",
							Currency: MainnetCurrency,
						},
					},
				},
				"14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0": {
					Account: &types.AccountIdentifier{
						Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0",
						},
						Amount: &types.Amount{
							Value:    "200000000",
							Currency: MainnetCurrency,
						},
					},
				},
				"4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1": {
					Account: &types.AccountIdentifier{
						Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1",
						},
						Amount: &types.Amount{
							Value:    "7100000000",
							Currency: MainnetCurrency,
						},
					},
				},
			},
			expectedBlock: &types.Block{
				BlockIdentifier: blockIdentifier717983,
				ParentBlockIdentifier: &types.BlockIdentifier{
					Hash:  "0067f80ce10d4255932b7f8c9baf7bd0dcfd408c312d33144be0ea12caf7f7f0",
					Index: 717982,
				},
				Timestamp: 1601465727000,
				Transactions: []*types.Transaction{
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   CoinbaseOpType,
								Status: SuccessStatus,
								Metadata: mustMarshalMap(&OperationMetadata{
									Coinbase: "039ff40a0102",
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "ztawr1vEZ6pZRtLqNy2C9u7EK7JN2gP8W6z", // nolint
									SubAccount: &types.SubAccountIdentifier{
										Address:  "coinbase",
									},
								},
								Amount: &types.Amount{
									Value:    "750010000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_DUP OP_HASH160 557662a6b307f95aa00311c074f7feebb955d451 OP_EQUALVERIFY OP_CHECKSIG", // nolint
										Hex:  "76a914557662a6b307f95aa00311c074f7feebb955d45188ac",         // nolint
										RequiredSigs: 1,
										Type: "pubkeyhash",
										Addresses: []string{
											"ztawr1vEZ6pZRtLqNy2C9u7EK7JN2gP8W6z",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "zrFzxutppvxEdjyu4QNjogBMjtC1py9Hp1S", // nolint
									SubAccount: &types.SubAccountIdentifier{
										Address:  "coinbase",
									},
								},
								Amount: &types.Amount{
									Value:    "250000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_HASH160 8d3468b6686ac59caf9ad94e547a737b09fa1027 OP_EQUAL", // nolint
										Hex:  "a9148d3468b6686ac59caf9ad94e547a737b09fa102787",         // nolint
										RequiredSigs: 1,
										Type: "scripthash",
										Addresses: []string{
											"zrFzxutppvxEdjyu4QNjogBMjtC1py9Hp1S",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        3,
									NetworkIndex: int64Pointer(2),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "zrS7QUB2eDbbKvyP43VJys3t7RpojW8GdxH", // nolint
									SubAccount: &types.SubAccountIdentifier{
										Address:  "coinbase",
									},
								},
								Amount: &types.Amount{
									Value:    "125000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88:2",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_HASH160 fc1d7f04db5e2c05b051e0decc85effe6bc539d5 OP_EQUAL", // nolint
										Hex:  "a914fc1d7f04db5e2c05b051e0decc85effe6bc539d587",         // nolint
										RequiredSigs: 1,
										Type: "scripthash",
										Addresses: []string{
											"zrS7QUB2eDbbKvyP43VJys3t7RpojW8GdxH",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        4,
									NetworkIndex: int64Pointer(3),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "zrFr5HVm7woVq3oFzkMEdJdbfBchfPAPDsP", // nolint
									SubAccount: &types.SubAccountIdentifier{
										Address:  "coinbase",
									},
								},
								Amount: &types.Amount{
									Value:    "125000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88:3",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_HASH160 8b85fc1e171a4c7994c088b91d5a75dff9e56cad OP_EQUAL", // nolint
										Hex:  "a9148b85fc1e171a4c7994c088b91d5a75dff9e56cad87",         // nolint
										RequiredSigs: 1,
										Type: "scripthash",
										Addresses: []string{
											"zrFr5HVm7woVq3oFzkMEdJdbfBchfPAPDsP",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    187,
							Version: 1,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "67c76a34cb6bde6f9628fdc8348c23191d3222e88386ed05c97e3c63384a01af",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: SuccessStatus,
								Amount: &types.Amount{
									Value:    "-60000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3044022059135f673a4919ab56775064cc82080ead1c74d8f0ebd943062b247c5946cf88022048f26c94a15752fa04d8bfff7388dd65d57485acd2395e539a50b2ca8e27870001 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
										Hex: "473044022059135f673a4919ab56775064cc82080ead1c74d8f0ebd943062b247c5946cf88022048f26c94a15752fa04d8bfff7388dd65d57485acd2395e539a50b2ca8e278700012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(1),
								},
								Type:   InputOpType,
								Status: SuccessStatus,
								Amount: &types.Amount{
									Value:    "-200000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "30440220527c59b1d2dbb87b71e01c9d1489f110727fc3120e5306539bd4668ed1063d30022079b6ca4ff77de3ab953bb0d896b74bb60c8ceca28248340201e701da0d1fd12b01 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
										Hex: "4730440220527c59b1d2dbb87b71e01c9d1489f110727fc3120e5306539bd4668ed1063d30022079b6ca4ff77de3ab953bb0d896b74bb60c8ceca28248340201e701da0d1fd12b012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(2),
								},
								Type:   InputOpType,
								Status: SuccessStatus,
								Amount: &types.Amount{
									Value:    "-7100000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "304402202d3b75ed231c1fe478c471452a0385c5cdc9fe2e337d5ee62cacd8a26d013e5002207d864a38e013d8c61b1972bd7bf78a53accd9b8d600fbbd7c79c21b2171fd8cb01 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
										Hex: "47304402202d3b75ed231c1fe478c471452a0385c5cdc9fe2e337d5ee62cacd8a26d013e5002207d864a38e013d8c61b1972bd7bf78a53accd9b8d600fbbd7c79c21b2171fd8cb012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        3,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "ztjySYJL8g9i6wc2YTusbDpPZSpPM5xuTua",
								},
								Amount: &types.Amount{
									Value:    "500000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "67c76a34cb6bde6f9628fdc8348c23191d3222e88386ed05c97e3c63384a01af:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_DUP OP_HASH160 b87cc09d17751ffeab924a82134665ae4202cbfc OP_EQUALVERIFY OP_CHECKSIG bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a00 717682 OP_CHECKBLOCKATHEIGHT",
										Hex:  "76a914b87cc09d17751ffeab924a82134665ae4202cbfc88ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab4",
										RequiredSigs: 1,
										Type: "pubkeyhashreplay",
										Addresses: []string{
											"ztjySYJL8g9i6wc2YTusbDpPZSpPM5xuTua",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        4,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
								},
								Amount: &types.Amount{
									Value:    "6859990000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "67c76a34cb6bde6f9628fdc8348c23191d3222e88386ed05c97e3c63384a01af:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_DUP OP_HASH160 fd2831ec8fc1bf3ccdeadbe9fcdb515aac904761 OP_EQUALVERIFY OP_CHECKSIG bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a00 717682 OP_CHECKBLOCKATHEIGHT",
										Hex:  "76a914fd2831ec8fc1bf3ccdeadbe9fcdb515aac90476188ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab4",
										RequiredSigs: 1,
										Type: "pubkeyhashreplay",
										Addresses: []string{
											"ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    595,
							Version: 1,
						}),
					},
				},
				Metadata: mustMarshalMap(&BlockMetadata{
					Size:       2271,
					Version:    3,
					MerkleRoot: "97c960c90e0b6bc30d2629f06d114f1c49aadb0e3d9bd70eb4f0f9ed1ea69279",
					Nonce:      "00002e570d64b4b3ea1c30dec68b2dff255eb3148656f06f5e018ae739a400eb",
					Bits:       "1f754920",
					Difficulty: 17.46160923,
				}),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			client := NewClient("", TestnetGenesisBlockIdentifier, MainnetCurrency)

			block, err := client.ParseBlock(context.Background(), test.block, test.coins)
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedBlock, block)
			}
		})
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			client := NewClient("", MainnetGenesisBlockIdentifier, MainnetCurrency)

			block, err := client.ParseBlock(context.Background(), test.block, test.coins)
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedBlock, block)
			}
		})
	}
}

func TestParseBlockWithCertificate(t *testing.T) {
	tests := map[string]struct {
		block *Block
		coins map[string]*storage.AccountCoin

		expectedBlock *types.Block
		expectedError error
	}{
		"block717984": {
			block: block717984,
			coins: map[string]*storage.AccountCoin{
				"9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0": {
					Account: &types.AccountIdentifier{
						Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0",
						},
						Amount: &types.Amount{
							Value:    "60000000",
							Currency: MainnetCurrency,
						},
					},
				},
				"14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0": {
					Account: &types.AccountIdentifier{
						Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0",
						},
						Amount: &types.Amount{
							Value:    "200000000",
							Currency: MainnetCurrency,
						},
					},
				},
				"4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1": {
					Account: &types.AccountIdentifier{
						Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1",
						},
						Amount: &types.Amount{
							Value:    "7100000000",
							Currency: MainnetCurrency,
						},
					},
				},
				"62091923e9805a8650d752b3b83e0d56ce70e775ee67c080feade7e5ee677ad9:0": {
					Account: &types.AccountIdentifier{
						Address: "ztpha3vQzv7eTdBvPC1oWnouuManmCEVbTT",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "62091923e9805a8650d752b3b83e0d56ce70e775ee67c080feade7e5ee677ad9:0",
						},
						Amount: &types.Amount{
							Value:    "7100000000",
							Currency: MainnetCurrency,
						},
					},
				},
			},
			expectedBlock: &types.Block{
				BlockIdentifier: blockIdentifier717984,
				ParentBlockIdentifier: &types.BlockIdentifier{
					Hash:  "0067f80ce10d4255932b7f8c9baf7bd0dcfd408c312d33144be0ea12caf7f7f0",
					Index: 717983,
				},
				Timestamp: 1601465727000,
				Transactions: []*types.Transaction{
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   CoinbaseOpType,
								Status: SuccessStatus,
								Metadata: mustMarshalMap(&OperationMetadata{
									Coinbase: "039ff40a0102",
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "ztawr1vEZ6pZRtLqNy2C9u7EK7JN2gP8W6z", // nolint
									SubAccount: &types.SubAccountIdentifier{
										Address:  "coinbase",
									},
								},
								Amount: &types.Amount{
									Value:    "750010000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_DUP OP_HASH160 557662a6b307f95aa00311c074f7feebb955d451 OP_EQUALVERIFY OP_CHECKSIG", // nolint
										Hex:  "76a914557662a6b307f95aa00311c074f7feebb955d45188ac",         // nolint
										RequiredSigs: 1,
										Type: "pubkeyhash",
										Addresses: []string{
											"ztawr1vEZ6pZRtLqNy2C9u7EK7JN2gP8W6z",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "zrFzxutppvxEdjyu4QNjogBMjtC1py9Hp1S", // nolint
									SubAccount: &types.SubAccountIdentifier{
										Address:  "coinbase",
									},
								},
								Amount: &types.Amount{
									Value:    "250000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_HASH160 8d3468b6686ac59caf9ad94e547a737b09fa1027 OP_EQUAL", // nolint
										Hex:  "a9148d3468b6686ac59caf9ad94e547a737b09fa102787",         // nolint
										RequiredSigs: 1,
										Type: "scripthash",
										Addresses: []string{
											"zrFzxutppvxEdjyu4QNjogBMjtC1py9Hp1S",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        3,
									NetworkIndex: int64Pointer(2),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "zrS7QUB2eDbbKvyP43VJys3t7RpojW8GdxH", // nolint
									SubAccount: &types.SubAccountIdentifier{
										Address:  "coinbase",
									},
								},
								Amount: &types.Amount{
									Value:    "125000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88:2",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_HASH160 fc1d7f04db5e2c05b051e0decc85effe6bc539d5 OP_EQUAL", // nolint
										Hex:  "a914fc1d7f04db5e2c05b051e0decc85effe6bc539d587",         // nolint
										RequiredSigs: 1,
										Type: "scripthash",
										Addresses: []string{
											"zrS7QUB2eDbbKvyP43VJys3t7RpojW8GdxH",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        4,
									NetworkIndex: int64Pointer(3),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "zrFr5HVm7woVq3oFzkMEdJdbfBchfPAPDsP", // nolint
									SubAccount: &types.SubAccountIdentifier{
										Address:  "coinbase",
									},
								},
								Amount: &types.Amount{
									Value:    "125000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e88:3",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_HASH160 8b85fc1e171a4c7994c088b91d5a75dff9e56cad OP_EQUAL", // nolint
										Hex:  "a9148b85fc1e171a4c7994c088b91d5a75dff9e56cad87",         // nolint
										RequiredSigs: 1,
										Type: "scripthash",
										Addresses: []string{
											"zrFr5HVm7woVq3oFzkMEdJdbfBchfPAPDsP",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    187,
							Version: 1,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "67c76a34cb6bde6f9628fdc8348c23191d3222e88386ed05c97e3c63384a01af",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: SuccessStatus,
								Amount: &types.Amount{
									Value:    "-60000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "9401f535c210f3ff362d3f51dba88ecddf4f87ed9d0563c1f9e8af75eca1fd1a:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3044022059135f673a4919ab56775064cc82080ead1c74d8f0ebd943062b247c5946cf88022048f26c94a15752fa04d8bfff7388dd65d57485acd2395e539a50b2ca8e27870001 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
										Hex: "473044022059135f673a4919ab56775064cc82080ead1c74d8f0ebd943062b247c5946cf88022048f26c94a15752fa04d8bfff7388dd65d57485acd2395e539a50b2ca8e278700012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(1),
								},
								Type:   InputOpType,
								Status: SuccessStatus,
								Amount: &types.Amount{
									Value:    "-200000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "14e8fe02ec4e237d8cb6bf95943bd05706a19f6bd29f9b2b1fefc4fa09ef6737:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "30440220527c59b1d2dbb87b71e01c9d1489f110727fc3120e5306539bd4668ed1063d30022079b6ca4ff77de3ab953bb0d896b74bb60c8ceca28248340201e701da0d1fd12b01 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
										Hex: "4730440220527c59b1d2dbb87b71e01c9d1489f110727fc3120e5306539bd4668ed1063d30022079b6ca4ff77de3ab953bb0d896b74bb60c8ceca28248340201e701da0d1fd12b012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        2,
									NetworkIndex: int64Pointer(2),
								},
								Type:   InputOpType,
								Status: SuccessStatus,
								Amount: &types.Amount{
									Value:    "-7100000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "4c292f9ba0e94f2d48a16f8765217e62b6673796bffd92c26b13ed5e661946bc:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "304402202d3b75ed231c1fe478c471452a0385c5cdc9fe2e337d5ee62cacd8a26d013e5002207d864a38e013d8c61b1972bd7bf78a53accd9b8d600fbbd7c79c21b2171fd8cb01 03ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
										Hex: "47304402202d3b75ed231c1fe478c471452a0385c5cdc9fe2e337d5ee62cacd8a26d013e5002207d864a38e013d8c61b1972bd7bf78a53accd9b8d600fbbd7c79c21b2171fd8cb012103ae26fe63b19c80972b6ffbd47e9f3b3e202740e5e349b0e23fd712927b0792ce", // nolint
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        3,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "ztjySYJL8g9i6wc2YTusbDpPZSpPM5xuTua",
								},
								Amount: &types.Amount{
									Value:    "500000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "67c76a34cb6bde6f9628fdc8348c23191d3222e88386ed05c97e3c63384a01af:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_DUP OP_HASH160 b87cc09d17751ffeab924a82134665ae4202cbfc OP_EQUALVERIFY OP_CHECKSIG bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a00 717682 OP_CHECKBLOCKATHEIGHT",
										Hex:  "76a914b87cc09d17751ffeab924a82134665ae4202cbfc88ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab4",
										RequiredSigs: 1,
										Type: "pubkeyhashreplay",
										Addresses: []string{
											"ztjySYJL8g9i6wc2YTusbDpPZSpPM5xuTua",
										},
									},
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        4,
									NetworkIndex: int64Pointer(1),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
								},
								Amount: &types.Amount{
									Value:    "6859990000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "67c76a34cb6bde6f9628fdc8348c23191d3222e88386ed05c97e3c63384a01af:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_DUP OP_HASH160 fd2831ec8fc1bf3ccdeadbe9fcdb515aac904761 OP_EQUALVERIFY OP_CHECKSIG bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a00 717682 OP_CHECKBLOCKATHEIGHT",
										Hex:  "76a914fd2831ec8fc1bf3ccdeadbe9fcdb515aac90476188ac20bd1d792d97a7da359adbc2fdadd04536f79aad9afc5821c4340043f7fb302a000372f30ab4",
										RequiredSigs: 1,
										Type: "pubkeyhashreplay",
										Addresses: []string{
											"ztrEXsPLywPcxE3Sn9qdWV6tYkBH4HnYwin",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    595,
							Version: 1,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e89",
						},
						Operations: []*types.Operation{
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        0,
									NetworkIndex: int64Pointer(0),
								},
								Type:   InputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "ztpha3vQzv7eTdBvPC1oWnouuManmCEVbTT",
								},
								Amount: &types.Amount{
									Value:    "-7100000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "62091923e9805a8650d752b3b83e0d56ce70e775ee67c080feade7e5ee677ad9:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3044022014d8dee1da3821dce95e48060f8f38394aee00f84d03a8203611ff3e703c10a002205ce62cffdc12dd26742489120d50d071ff08f993b9cca0b31a73e0f20f20cb5d01 0241b92fed18a3ded2b98459b5432982a0712912ad86b929ec6feb19655824b7cc",
										Hex: "473044022014d8dee1da3821dce95e48060f8f38394aee00f84d03a8203611ff3e703c10a002205ce62cffdc12dd26742489120d50d071ff08f993b9cca0b31a73e0f20f20cb5d01210241b92fed18a3ded2b98459b5432982a0712912ad86b929ec6feb19655824b7cc",
									},
									Sequence: 4294967295,
								}),
							},
							{
								OperationIdentifier: &types.OperationIdentifier{
									Index:        1,
									NetworkIndex: int64Pointer(0),
								},
								Type:   OutputOpType,
								Status: SuccessStatus,
								Account: &types.AccountIdentifier{
									Address: "ztpha3vQzv7eTdBvPC1oWnouuManmCEVbTT",
									SubAccount: &types.SubAccountIdentifier{
										Address:  "coinbase",
									},
								},
								Amount: &types.Amount{
									Value:    "24895145",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "afa747bcb78e22e5550e880d0803a5fa4cdbc7e04ff303a4b14da2c36e348e89:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM: "OP_DUP OP_HASH160 ec54fedd6a312d5c536046323bfabb9d2a475d7a OP_EQUALVERIFY OP_CHECKSIG 4ca064b46515f3f00e846e6c1b45ef36a082ea786783096d2cb6169556756e08 21 OP_CHECKBLOCKATHEIGHT",
										Hex: "76a914ec54fedd6a312d5c536046323bfabb9d2a475d7a88ac204ca064b46515f3f00e846e6c1b45ef36a082ea786783096d2cb6169556756e080115b4",
										RequiredSigs: 1,
										Type: "pubkeyhashreplay",
										Addresses: []string{
											"ztpha3vQzv7eTdBvPC1oWnouuManmCEVbTT",
										},
									},
								}),
							},
						},
					},
				},
				Metadata: mustMarshalMap(&BlockMetadata{
					Size:       2271,
					Version:    3,
					MerkleRoot: "97c960c90e0b6bc30d2629f06d114f1c49aadb0e3d9bd70eb4f0f9ed1ea69279",
					Nonce:      "00002e570d64b4b3ea1c30dec68b2dff255eb3148656f06f5e018ae739a400eb",
					Bits:       "1f754920",
					Difficulty: 17.46160923,
				}),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			client := NewClient("", TestnetGenesisBlockIdentifier, MainnetCurrency)

			block, err := client.ParseBlock(context.Background(), test.block, test.coins)
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedBlock, block)
			}
		})
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			client := NewClient("", MainnetGenesisBlockIdentifier, MainnetCurrency)

			block, err := client.ParseBlock(context.Background(), test.block, test.coins)
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedBlock, block)
			}
		})
	}
}

func TestSuggestedFeeRate(t *testing.T) {
	tests := map[string]struct {
		responses []responseFixture

		expectedRate  float64
		expectedError error
	}{
		"successful": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("fee_rate.json"),
					url:    url,
				},
			},
			expectedRate: float64(0.00001),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			responses := make(chan responseFixture, len(test.responses))
			for _, response := range test.responses {
				responses <- response
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := <-responses
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("POST", r.Method)
				assert.Equal(response.url, r.URL.RequestURI())

				w.WriteHeader(response.status)
				fmt.Fprintln(w, response.body)
			}))

			client := NewClient(ts.URL, MainnetGenesisBlockIdentifier, MainnetCurrency)
			rate, err := client.SuggestedFeeRate(context.Background(), 1)
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedRate, rate)
			}
		})
	}
}

func TestRawMempool(t *testing.T) {
	tests := map[string]struct {
		responses []responseFixture

		expectedTransactions []string
		expectedError        error
	}{
		"successful": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("raw_mempool.json"),
					url:    url,
				},
			},
			expectedTransactions: []string{
				"9cec12d170e97e21a876fa2789e6bfc25aa22b8a5e05f3f276650844da0c33ab",
				"37b4fcc8e0b229412faeab8baad45d3eb8e4eec41840d6ac2103987163459e75",
				"7bbb29ae32117597fcdf21b464441abd571dad52d053b9c2f7204f8ea8c4762e",
			},
		},
		"500 error": {
			responses: []responseFixture{
				{
					status: http.StatusInternalServerError,
					body:   "{}",
					url:    url,
				},
			},
			expectedError: errors.New("invalid response: 500 Internal Server Error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				assert = assert.New(t)
			)

			responses := make(chan responseFixture, len(test.responses))
			for _, response := range test.responses {
				responses <- response
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := <-responses
				assert.Equal("application/json", r.Header.Get("Content-Type"))
				assert.Equal("POST", r.Method)
				assert.Equal(response.url, r.URL.RequestURI())

				w.WriteHeader(response.status)
				fmt.Fprintln(w, response.body)
			}))

			client := NewClient(ts.URL, MainnetGenesisBlockIdentifier, MainnetCurrency)
			txs, err := client.RawMempool(context.Background())
			if test.expectedError != nil {
				assert.Contains(err.Error(), test.expectedError.Error())
			} else {
				assert.NoError(err)
				assert.Equal(test.expectedTransactions, txs)
			}
		})
	}
}


// loadFixture takes a file name and returns the response fixture.
func loadFixture(fileName string) string {
	content, err := ioutil.ReadFile(fmt.Sprintf("client_fixtures/%s", fileName))
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

type responseFixture struct {
	status int
	body   string
	url    string
}
