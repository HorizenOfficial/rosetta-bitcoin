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

package bitcoin

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
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

	blockIdentifier100000 = &types.BlockIdentifier{
		Hash:  "000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506",
		Index: 100000,
	}

	block100000 = &Block{
		Hash:              "000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506",
		Height:            100000,
		PreviousBlockHash: "000000000002d01c1fccc21636b607dfd930d31d01c3a62104612a1719011250",
		Time:              1293623863,
		Size:              957,
		Version:           1,
		MerkleRoot:        "f3e94742aca4b5ef85488dc37c06c3282295ffec960994b2c0d5ac2a25a95766",
		Nonce:             "274148111",
		Bits:              "1b04864c",
		Difficulty:        14484.1623612254,
		Txs: []*Transaction{
			{
				Hex:      "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff08044c86041b020602ffffffff0100f2052a010000004341041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84ac00000000", // nolint
				Hash:     "8c14f0db3df150123e6f3dbbf30f8b955a8249b62ac1d1ff16284aefa3d06d87",
				Size:     135,
				Vsize:    135,
				Version:  1,
				Locktime: 0,
				Inputs: []*Input{
					{
						Coinbase: "044c86041b020602",
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 15.89351625,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_HASH160 228f554bbf766d6f9cc828de1126e3d35d15e5fe OP_EQUAL",
							Hex:          "a914228f554bbf766d6f9cc828de1126e3d35d15e5fe87",
							RequiredSigs: 1,
							Type:         "scripthash",
							Addresses: []string{
								"34qkc2iac6RsyxZVfyE2S5U5WcRsbg2dpK",
							},
						},
					},
					{
						Value: 0,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:  "OP_RETURN aa21a9ed10109f4b82aa3ed7ec9d02a2a90246478b3308c8b85daf62fe501d58d05727a4",
							Hex:  "6a24aa21a9ed10109f4b82aa3ed7ec9d02a2a90246478b3308c8b85daf62fe501d58d05727a4",
							Type: "nulldata",
						},
					},
				},
			},
			{
				Hex:      "0100000001032e38e9c0a84c6046d687d10556dcacc41d275ec55fc00779ac88fdf357a187000000008c493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3ffffffff0200e32321000000001976a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac000fe208010000001976a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac00000000", // nolint
				Hash:     "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
				Size:     259,
				Vsize:    259,
				Version:  1,
				Locktime: 0,
				Inputs: []*Input{
					{
						TxHash: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
							Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 5.56,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
							},
						},
					},
					{
						Value: 44.44,
						Index: 1,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
							},
						},
					},
				},
			},
			{
				Hash:     "fake",
				Hex:      "fake hex",
				Version:  2,
				Size:     421,
				Vsize:    612,
				Locktime: 10,
				Inputs: []*Input{
					{
						TxHash: "503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "00142b2296c588ec413cebd19c3cbc04ea830ead6e78",
							Hex: "1600142b2296c588ec413cebd19c3cbc04ea830ead6e78",
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5",
						Vout:   1,
						ScriptSig: &ScriptSig{
							ASM: "00142b2296c588ec413cebd19c3cbc04ea830ead6e78",
							Hex: "1600142b2296c588ec413cebd19c3cbc04ea830ead6e78",
						},
						Sequence: 4294967295,
					},
					{
						TxHash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
						Vout:   0,
						ScriptSig: &ScriptSig{
							ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
							Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
						},
						Sequence: 4294967295,
					},
				},
				Outputs: []*Output{
					{
						Value: 200.56,
						Index: 0,
						ScriptPubKey: &ScriptPubKey{
							ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
							Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
							RequiredSigs: 1,
							Type:         "pubkeyhash",
							Addresses: []string{
								"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
							},
						},
					},
				},
			},
		},
	}
)

/*
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
				CurrentBlockIdentifier: blockIdentifier1000,
				CurrentBlockTimestamp:  block1000.Time * 1000,
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
				Hash: &blockIdentifier1000.Hash,
			},
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response.json"),
					url:    url,
				},
			},
			expectedBlock: block1000,
			expectedCoins: []string{},
		},
		"lookup by hash 2": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Hash: &blockIdentifier100000.Hash,
			},
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("get_block_response_2.json"),
					url:    url,
				},
			},
			expectedBlock: block100000,
			expectedCoins: []string{
				"87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
				"503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:0",
				"503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:1",
			},
		},
		"lookup by hash (get block api error)": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Hash: &blockIdentifier1000.Hash,
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
				Hash: &blockIdentifier1000.Hash,
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
				Index: &blockIdentifier1000.Index,
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
			expectedBlock: block1000,
			expectedCoins: []string{},
		},
		"lookup by index (out of range)": {
			blockIdentifier: &types.PartialBlockIdentifier{
				Index: &blockIdentifier1000.Index,
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
			expectedBlock: block1000,
			expectedCoins: []string{},
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
*/
func int64Pointer(v int64) *int64 {
	return &v
}

func mustMarshalMap(v interface{}) map[string]interface{} {
	m, _ := types.MarshalMap(v)
	return m
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
		/*
		"block 100000": {
			block: block100000,
			coins: map[string]*storage.AccountCoin{
				"87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0": {
					Account: &types.AccountIdentifier{
						Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
						},
						Amount: &types.Amount{
							Value:    "5000000000",
							Currency: MainnetCurrency,
						},
					},
				},
				"503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:0": {
					Account: &types.AccountIdentifier{
						Address: "3FfQGY7jqsADC7uTVqF3vKQzeNPiBPTqt4",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:0",
						},
						Amount: &types.Amount{
							Value:    "3467607",
							Currency: MainnetCurrency,
						},
					},
				},
				"503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:1": {
					Account: &types.AccountIdentifier{
						Address: "1NdvAyRJLdK5EXs7DV3ebYb5wffdCZk1pD",
					},
					Coin: &types.Coin{
						CoinIdentifier: &types.CoinIdentifier{
							Identifier: "503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:1",
						},
						Amount: &types.Amount{
							Value:    "0",
							Currency: MainnetCurrency,
						},
					},
				},
			},
			expectedBlock: &types.Block{
				BlockIdentifier: blockIdentifier100000,
				ParentBlockIdentifier: &types.BlockIdentifier{
					Hash:  "000000000002d01c1fccc21636b607dfd930d31d01c3a62104612a1719011250",
					Index: 99999,
				},
				Timestamp: 1293623863000,
				Transactions: []*types.Transaction{
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "8c14f0db3df150123e6f3dbbf30f8b955a8249b62ac1d1ff16284aefa3d06d87",
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
									Coinbase: "044c86041b020602",
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
									Address: "34qkc2iac6RsyxZVfyE2S5U5WcRsbg2dpK",
								},
								Amount: &types.Amount{
									Value:    "1589351625",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "8c14f0db3df150123e6f3dbbf30f8b955a8249b62ac1d1ff16284aefa3d06d87:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_HASH160 228f554bbf766d6f9cc828de1126e3d35d15e5fe OP_EQUAL",
										Hex:          "a914228f554bbf766d6f9cc828de1126e3d35d15e5fe87",
										RequiredSigs: 1,
										Type:         "scripthash",
										Addresses: []string{
											"34qkc2iac6RsyxZVfyE2S5U5WcRsbg2dpK",
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
									Address: "6a24aa21a9ed10109f4b82aa3ed7ec9d02a2a90246478b3308c8b85daf62fe501d58d05727a4",
								},
								Amount: &types.Amount{
									Value:    "0",
									Currency: MainnetCurrency,
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:  "OP_RETURN aa21a9ed10109f4b82aa3ed7ec9d02a2a90246478b3308c8b85daf62fe501d58d05727a4",
										Hex:  "6a24aa21a9ed10109f4b82aa3ed7ec9d02a2a90246478b3308c8b85daf62fe501d58d05727a4",
										Type: "nulldata",
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    135,
							Version: 1,
							Vsize:   135,
							Weight:  540,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
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
									Value:    "-5000000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1BNwxHGaFbeUBitpjy2AsKpJ29Ybxntqvb",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
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
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								Amount: &types.Amount{
									Value:    "556000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
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
									Address: "1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
								},
								Amount: &types.Amount{
									Value:    "4444000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:    259,
							Version: 1,
							Vsize:   259,
							Weight:  1036,
						}),
					},
					{
						TransactionIdentifier: &types.TransactionIdentifier{
							Hash: "fake",
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
									Value:    "-3467607",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "3FfQGY7jqsADC7uTVqF3vKQzeNPiBPTqt4",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "00142b2296c588ec413cebd19c3cbc04ea830ead6e78",
										Hex: "1600142b2296c588ec413cebd19c3cbc04ea830ead6e78",
									},
									TxInWitness: []string{
										"304402205f39ccbab38b644acea0776d18cb63ce3e37428cbac06dc23b59c61607aef69102206b8610827e9cb853ea0ba38983662034bd3575cc1ab118fb66d6a98066fa0bed01", // nolint
										"0304c01563d46e38264283b99bb352b46e69bf132431f102d4bd9a9d8dab075e7f",
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
									Value:    "0",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1NdvAyRJLdK5EXs7DV3ebYb5wffdCZk1pD",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "503e4e9824282eb06f1a328484e2b367b5f4f93a405d6e7b97261bafabfb53d5:1",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "00142b2296c588ec413cebd19c3cbc04ea830ead6e78",
										Hex: "1600142b2296c588ec413cebd19c3cbc04ea830ead6e78",
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
									Value:    "-556000000",
									Currency: MainnetCurrency,
								},
								Account: &types.AccountIdentifier{
									Address: "1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinSpent,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptSig: &ScriptSig{
										ASM: "3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
										Hex: "493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3", // nolint
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
									Address: "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
								},
								Amount: &types.Amount{
									Value:    "20056000000",
									Currency: MainnetCurrency,
								},
								CoinChange: &types.CoinChange{
									CoinAction: types.CoinCreated,
									CoinIdentifier: &types.CoinIdentifier{
										Identifier: "fake:0",
									},
								},
								Metadata: mustMarshalMap(&OperationMetadata{
									ScriptPubKey: &ScriptPubKey{
										ASM:          "OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG",
										Hex:          "76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac",
										RequiredSigs: 1,
										Type:         "pubkeyhash",
										Addresses: []string{
											"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn",
											"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx",
										},
									},
								}),
							},
						},
						Metadata: mustMarshalMap(&TransactionMetadata{
							Size:     421,
							Version:  2,
							Vsize:    612,
							Weight:   129992,
							Locktime: 10,
						}),
					},
				},
				Metadata: mustMarshalMap(&BlockMetadata{
					Size:       957,
					Weight:     3828,
					Version:    1,
					MerkleRoot: "f3e94742aca4b5ef85488dc37c06c3282295ffec960994b2c0d5ac2a25a95766",
					MedianTime: 1293622620,
					Nonce:      274148111,
					Bits:       "1b04864c",
					Difficulty: 14484.1623612254,
				}),
			},
		},
		"missing transactions": {
			block:         block100000,
			coins:         map[string]*storage.AccountCoin{},
			expectedError: errors.New("error finding previous tx"),
		},*/
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
}

/*
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
		"invalid range error": {
			responses: []responseFixture{
				{
					status: http.StatusOK,
					body:   loadFixture("invalid_fee_rate.json"),
					url:    url,
				},
			},
			expectedError: errors.New("error getting fee estimate"),
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
*/

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
