// Copyright (c) 2013-2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base58_test

import (
	"testing"

	"github.com/HorizenOfficial/rosetta-zen/zenutil/base58"
)

var checkEncodingStringTests = []struct {
	version uint16
	in      string
	out     string
}{
	{0x2089, "", "HCo8nxRe"},
	{0x2089, " ", "2EWfSfEqw2"},
	{0x2089, "-", "2EWgq6Y1Qs"},
	{0x2089, "0", "2EWhCfRN3G"},
	{0x2089, "1", "2EWhKra8WW"},
	{0x2089, "-1", "6Se2nLNujPM"},
	{0x2089, "11", "6Se4mwP8RKS"},
	{0x2089, "abc", "R2CRaXTNrhq1"},
	{0x2089, "1234598760", "522T7ibo3ettqS2LrcRYuv"},
	{0x2089, "abcdefghijklmnopqrstuvwxyz", "3C1YgcLd1zZnGrDm53kG6zvtdC6kUQe9mRFXpm1L7A3E"},
	{0x2089, "00000000000000000000000000000000000000000000000000000000000000", "5FtGq8oioEKbtgWgcC7DyW2RyqL2p8QkEkAb67B77pteiH531Lnqha2ypkaYZ9LZQS9cHLNGnMHWLKFU4ii9d2mF3QRBM"},
}

func TestBase58Check(t *testing.T) {
	for x, test := range checkEncodingStringTests {
		// test encoding
		if res := base58.CheckEncode([]byte(test.in), test.version); res != test.out {
			t.Errorf("CheckEncode test #%d failed: got %s, want: %s", x, res, test.out)
		}

		// test decoding
		res, version, err := base58.CheckDecode(test.out)
		if err != nil {
			t.Errorf("CheckDecode test #%d failed with err: %v", x, err)
		} else if version != test.version {
			t.Errorf("CheckDecode test #%d failed: got version: %d want: %d", x, version, test.version)
		} else if string(res) != test.in {
			t.Errorf("CheckDecode test #%d failed: got: %s want: %s", x, res, test.in)
		}
	}

	// test the two decoding failure cases
	// case 1: checksum error
	_, _, err := base58.CheckDecode("3MNQE1Y")
	if err != base58.ErrChecksum {
		t.Error("Checkdecode test failed, expected ErrChecksum")
	}
	// case 2: invalid formats (string lengths below 5 mean the version byte and/or the checksum
	// bytes are missing).
	testString := "x"
	for len := 0; len < 3; len++ {
		testString = testString + "x"
		_, _, err = base58.CheckDecode(testString)
		if err != base58.ErrInvalidFormat {
			t.Error("Checkdecode test failed, expected ErrInvalidFormat")
		}
	}

}
