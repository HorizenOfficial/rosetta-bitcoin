// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package zenutil_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/HorizenOfficial/rosetta-zen/zend/chaincfg"
	"github.com/HorizenOfficial/rosetta-zen/zend/wire"
	"github.com/HorizenOfficial/rosetta-zen/zenutil"
	"golang.org/x/crypto/ripemd160"
)

type CustomParamStruct struct {
	Net              wire.BitcoinNet
	PubKeyHashAddrID uint16
	ScriptHashAddrID uint16
}

var CustomParams = CustomParamStruct{
	Net:              0xdbb6c0fb, // litecoin mainnet HD version bytes
	PubKeyHashAddrID: 0x2089,     // starts with L
	ScriptHashAddrID: 0x2096,     // starts with M
}

// We use this function to be able to test functionality in DecodeAddress for
// defaultNet addresses
func applyCustomParams(params chaincfg.Params, customParams CustomParamStruct) chaincfg.Params {
	params.Net = customParams.Net
	params.PubKeyHashAddrID = customParams.PubKeyHashAddrID
	params.ScriptHashAddrID = customParams.ScriptHashAddrID
	return params
}

var customParams = applyCustomParams(chaincfg.MainNetParams, CustomParams)

func TestAddresses(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		encoded string
		valid   bool
		result  zenutil.Address
		f       func() (zenutil.Address, error)
		net     *chaincfg.Params
	}{
		// Positive P2PKH test.
		{
			name:    "mainnet p2pkh",
			addr:    "zncPA6vFaz5wMyHeuQUcVfJk8FS4gGYQXRb",
			encoded: "zncPA6vFaz5wMyHeuQUcVfJk8FS4gGYQXRb",
			valid:   true,
			result: zenutil.TstAddressPubKeyHash(
				[ripemd160.Size]byte{
					0x7b, 0xec, 0x3c, 0x2b, 0x8f, 0x04, 0xdd, 0xbd, 0xa3, 0xc2,
					0x38, 0xe3, 0xea, 0xe6, 0xfb, 0xcf, 0xfa, 0x75, 0x7a, 0x25},
				chaincfg.MainNetParams.PubKeyHashAddrID),
			f: func() (zenutil.Address, error) {
				pkHash := []byte{
					0x7b, 0xec, 0x3c, 0x2b, 0x8f, 0x04, 0xdd, 0xbd, 0xa3, 0xc2,
					0x38, 0xe3, 0xea, 0xe6, 0xfb, 0xcf, 0xfa, 0x75, 0x7a, 0x25}
				return zenutil.NewAddressPubKeyHash(pkHash, &chaincfg.MainNetParams)
			},
			net: &chaincfg.MainNetParams,
		},
		// Positive P2PSH test.
		{

			name:    "mainnet p2psh",
			addr:    "zsnf7iK6WK863vkEEtf6yVrXPGYDGCU9xL5",
			encoded: "zsnf7iK6WK863vkEEtf6yVrXPGYDGCU9xL5",
			valid:   true,
			result: zenutil.TstAddressScriptHash(
				[ripemd160.Size]byte{
					0x59, 0xb6, 0xa1, 0x74, 0xdc, 0x66, 0xe3, 0x73, 0x5f, 0x8a,
					0xc8, 0x1f, 0x4a, 0x7b, 0x02, 0x6e, 0xdb, 0xd8, 0x21, 0x06},
				chaincfg.MainNetParams.ScriptHashAddrID),
			f: func() (zenutil.Address, error) {
				pkHash := []byte{
					0x59, 0xb6, 0xa1, 0x74, 0xdc, 0x66, 0xe3, 0x73, 0x5f, 0x8a,
					0xc8, 0x1f, 0x4a, 0x7b, 0x02, 0x6e, 0xdb, 0xd8, 0x21, 0x06}

				return zenutil.NewAddressScriptHashFromHash(pkHash, &chaincfg.MainNetParams)
			},
			net: &chaincfg.MainNetParams,
		},
	}

	if err := chaincfg.Register(&customParams); err != nil {
		panic(err)
	}

	for _, test := range tests {
		// Decode addr and compare error against valid.
		decoded, err := zenutil.DecodeAddress(test.addr, test.net)
		if (err == nil) != test.valid {
			t.Errorf("%v: decoding test failed: %v", test.name, err)
			return
		}
		if err == nil {
			// Ensure the stringer returns the same address as the
			// original.
			if decodedStringer, ok := decoded.(fmt.Stringer); ok {
				addr := test.addr

				// For Segwit addresses the string representation
				// will always be lower case, so in that case we
				// convert the original to lower case first.
				if strings.Contains(test.name, "segwit") {
					addr = strings.ToLower(addr)
				}

				if addr != decodedStringer.String() {
					t.Errorf("%v: String on decoded value does not match expected value: %v != %v",
						test.name, test.addr, decodedStringer.String())
					return
				}
			}

			// Encode again and compare against the original.
			encoded := decoded.EncodeAddress()
			if test.encoded != encoded {
				t.Errorf("%v: decoding and encoding produced different addressess: %v != %v",
					test.name, test.encoded, encoded)
				return
			}

			// Perform type-specific calculations.
			var saddr []byte
			switch d := decoded.(type) {
			case *zenutil.AddressPubKeyHash:
				saddr = zenutil.TstAddressSAddr(encoded)
			case *zenutil.AddressScriptHash:
				saddr = zenutil.TstAddressSAddr(encoded)

			case *zenutil.AddressPubKey:
				// Ignore the error here since the script
				// address is checked below.
				saddr, _ = hex.DecodeString(d.String())
			}

			// Check script address, as well as the Hash160 method for P2PKH and
			// P2SH addresses.
			if !bytes.Equal(saddr, decoded.ScriptAddress()) {
				t.Errorf("%v: script addresses do not match:\n%x != \n%x",
					test.name, saddr, decoded.ScriptAddress())
				return
			}
			switch a := decoded.(type) {
			case *zenutil.AddressPubKeyHash:
				if h := a.Hash160()[:]; !bytes.Equal(saddr, h) {
					t.Errorf("%v: hashes do not match:\n%x != \n%x",
						test.name, saddr, h)
					return
				}

			case *zenutil.AddressScriptHash:
				if h := a.Hash160()[:]; !bytes.Equal(saddr, h) {
					t.Errorf("%v: hashes do not match:\n%x != \n%x",
						test.name, saddr, h)
					return
				}

			}

			// Ensure the address is for the expected network.
			if !decoded.IsForNet(test.net) {
				t.Errorf("%v: calculated network does not match expected",
					test.name)
				return
			}
		}

		if !test.valid {
			// If address is invalid, but a creation function exists,
			// verify that it returns a nil addr and non-nil error.
			if test.f != nil {
				_, err := test.f()
				if err == nil {
					t.Errorf("%v: address is invalid but creating new address succeeded",
						test.name)
					return
				}
			}
			continue
		}
		// Valid test, compare address created with f against expected result.
		addr, err := test.f()
		if err != nil {
			t.Errorf("%v: address is valid but creating new address failed with error %v",
				test.name, err)
			return
		}
		if !reflect.DeepEqual(addr, test.result) {
			t.Errorf("%v: created address does not match expected result",
				test.name)
			return
		}
	}
}
