// Copyright (c) 2014-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"time"

	"github.com/HorizenOfficial/rosetta-zen/zend/chaincfg/chainhash"
	"github.com/HorizenOfficial/rosetta-zen/zend/wire"
)

// genesisCoinbaseTx is the coinbase transaction for the genesis blocks for
// the main network, regression test network, and test network (version 3).
var genesisCoinbaseTx = wire.MsgTx{
	Version: 1,
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  chainhash.Hash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x04, 0xff, 0xff, 0x00, 0x1d, 0x01, 0x04, 0x45, /* |.......E| */
				0x54, 0x68, 0x65, 0x20, 0x54, 0x69, 0x6d, 0x65, /* |The Time| */
				0x73, 0x20, 0x30, 0x33, 0x2f, 0x4a, 0x61, 0x6e, /* |s 03/Jan| */
				0x2f, 0x32, 0x30, 0x30, 0x39, 0x20, 0x43, 0x68, /* |/2009 Ch| */
				0x61, 0x6e, 0x63, 0x65, 0x6c, 0x6c, 0x6f, 0x72, /* |ancellor| */
				0x20, 0x6f, 0x6e, 0x20, 0x62, 0x72, 0x69, 0x6e, /* | on brin| */
				0x6b, 0x20, 0x6f, 0x66, 0x20, 0x73, 0x65, 0x63, /* |k of sec|*/
				0x6f, 0x6e, 0x64, 0x20, 0x62, 0x61, 0x69, 0x6c, /* |ond bail| */
				0x6f, 0x75, 0x74, 0x20, 0x66, 0x6f, 0x72, 0x20, /* |out for |*/
				0x62, 0x61, 0x6e, 0x6b, 0x73, /* |banks| */
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: 0x00,
			PkScript: []byte{
				0x41, 0x04, 0x67, 0x8a, 0xfd, 0xb0, 0xfe, 0x55, /* |A.g....U| */
				0x48, 0x27, 0x19, 0x67, 0xf1, 0xa6, 0x71, 0x30, /* |H'.g..q0| */
				0xb7, 0x10, 0x5c, 0xd6, 0xa8, 0x28, 0xe0, 0x39, /* |..\..(.9| */
				0x09, 0xa6, 0x79, 0x62, 0xe0, 0xea, 0x1f, 0x61, /* |..yb...a| */
				0xde, 0xb6, 0x49, 0xf6, 0xbc, 0x3f, 0x4c, 0xef, /* |..I..?L.| */
				0x38, 0xc4, 0xf3, 0x55, 0x04, 0xe5, 0x1e, 0xc1, /* |8..U....| */
				0x12, 0xde, 0x5c, 0x38, 0x4d, 0xf7, 0xba, 0x0b, /* |..\8M...| */
				0x8d, 0x57, 0x8a, 0x4c, 0x70, 0x2b, 0x6b, 0xf1, /* |.W.Lp+k.| */
				0x1d, 0x5f, 0xac, /* |._.| */
			},
		},
	},
	LockTime: 0,
}

// genesisHash is the hash of the first block in the block chain for the main
// network (genesis block).
var genesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x00, 0x07, 0x10, 0x4c, 0xcd, 0xa2, 0x89, 0x42,
	0x79, 0x19, 0xef, 0xc3, 0x9d, 0xc9, 0xe4, 0xd4,
	0x99, 0x80, 0x4b, 0x7b, 0xeb, 0xc2, 0x2d, 0xf5,
	0x5f, 0x8b, 0x83, 0x43, 0x01, 0x26, 0x06, 0x02,
})

// genesisMerkleRoot is the hash of the first transaction in the genesis block
// for the main network.
var genesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x19, 0x61, 0x2b, 0xcf, 0x00, 0xea, 0x76, 0x11,
	0xd3, 0x15, 0xd7, 0xf4, 0x35, 0x54, 0xfa, 0x98,
	0x3c, 0x6e, 0x8c, 0x30, 0xcb, 0xa1, 0x7e, 0x52,
	0xc6, 0x79, 0xe0, 0xe8, 0x0a, 0xbf, 0x7d, 0x42,
})

// genesisBlock defines the genesis block of the block chain which serves as the
// public transaction ledger for the main network.
var genesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},                                                  // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: genesisMerkleRoot,                                                  // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(0x581EA6F5,0),                                            // 2009-01-03 18:15:05 +0000 UTC
		Bits:       0x1f07ffff,                                                         // 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce:      0x000000000000000000000000000000000000000000000000000000000000021d, // 2083236893
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx}, //TODO
}

// regTestGenesisHash is the hash of the first block in the block chain for the
// regression test network (genesis block).
var testnetGenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x03, 0xe1, 0xc4, 0xbb, 0x70, 0x5c, 0x87, 0x1b,
	0xf9, 0xbf, 0xda, 0x3e, 0x74, 0xb7, 0xf8, 0xf8,
	0x6b, 0xff, 0x26, 0x79, 0x93, 0xc2, 0x15, 0xa8,
	0x9d, 0x57, 0x95, 0xe3, 0x70, 0x8e, 0xe5, 0x1f,
})

// regTestGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the regression test network.  It is the same as the merkle root for
// the main network.
var regTestGenesisMerkleRoot = genesisMerkleRoot

// regTestGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the regression test network.
var regTestGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},        // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: regTestGenesisMerkleRoot, // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1494548150, 0), // 2011-02-02 23:16:42 +0000 UTC
		Bits:       0x200f0f0f,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      0x000000000000000000000000000000000000000000000000000000000000003d,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx}, //TODO
}

// testNet3GenesisHash is the hash of the first block in the block chain for the
// test network (version 3).
var regtestGenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x0d, 0xa5, 0xee, 0x72, 0x3b, 0x79, 0x23, 0xfe,
	0xb5, 0x80, 0x51, 0x85, 0x41, 0xc6, 0xf0, 0x98,
	0x20, 0x63, 0x30, 0xdb, 0xc7, 0x11, 0xa6, 0x67,
	0x89, 0x22, 0xc1, 0x1f, 0x2c, 0xcf, 0x1a, 0xbb,
})

// testNet3GenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the test network (version 3).  It is the same as the merkle root
// for the main network.
var testNetGenesisMerkleRoot = genesisMerkleRoot

// testNet3GenesisBlock defines the genesis block of the block chain which
// serves as the public transaction ledger for the test network (version 3).
var testNetGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},                                                  // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: testNetGenesisMerkleRoot,                                           // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1479443947, 0),                                           // 2011-02-02 23:16:42 +0000 UTC
		Bits:       0x2007ffff,                                                         // 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce:      0x0000000000000000000000000000000000000000000000000000000000000013, // 414098458
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx}, //TODO
}
