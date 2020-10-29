// Copyright (c) 2013-2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base58

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

// ErrChecksum indicates that the checksum of a check-encoded string does not verify against
// the checksum.
var ErrChecksum = errors.New("checksum error")

// ErrInvalidFormat indicates that the check-encoded string has an invalid format.
var ErrInvalidFormat = errors.New("invalid format: version and/or checksum bytes missing")

// checksum: first four bytes of sha256^2
func checksum(input []byte) (cksum [4]byte) {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	copy(cksum[:], h2[:4])
	return
}

// CheckEncode prepends a version byte and appends a four byte checksum.
func CheckEncode(input []byte, version uint16) string {
	b := make([]byte, 0, 2+len(input)+4)
	var versionByte = make([]byte, 2)
	binary.BigEndian.PutUint16(versionByte,version)
	b = append(b, versionByte[0])
	b = append(b, versionByte[1])
	b = append(b, input[:]...)
	cksum := checksum(b)
	b = append(b, cksum[:]...)
	return Encode(b)
}

// CheckDecode decodes a string that was encoded with CheckEncode and verifies the checksum.
func CheckDecode(input string) (result []byte, version uint16, err error) {
	decoded := Decode(input)
	b := make([]byte, 0, 2)
	b = append(b, decoded[0])
	b = append(b, decoded[1])
	if len(decoded) < 5 {
		return nil, 0, ErrInvalidFormat
	}
	version = binary.BigEndian.Uint16(b)
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	if checksum(decoded[:len(decoded)-4]) != cksum {
		return nil, 0, ErrChecksum
	}
	payload := decoded[2 : len(decoded)-4]
	result = append(result, payload...)
	return
}
