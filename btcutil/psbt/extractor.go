// Copyright (c) 2018 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package psbt

// The Extractor requires provision of a single PSBT
// in which all necessary signatures are encoded, and
// uses it to construct a fully valid network serialized
// transaction.

import (
	"github.com/HorizenOfficial/rosetta-zen/btcd/wire"
)

// Extract takes a finalized psbt.Packet and outputs a finalized transaction
// instance. Note that if the PSBT is in-complete, then an error
// ErrIncompletePSBT will be returned. As the extracted transaction has been
// fully finalized, it will be ready for network broadcast once returned.
func Extract(p *Packet) (*wire.MsgTx, error) {
	// If the packet isn't complete, then we'll return an error as it
	// doesn't have all the required witness data.
	if !p.IsComplete() {
		return nil, ErrIncompletePSBT
	}

	// First, we'll make a copy of the underlying unsigned transaction (the
	// initial template) so we don't mutate it during our activates below.
	finalTx := p.UnsignedTx.Copy()

	// For each input, we'll now populate any relevant witness and
	// sigScript data.
	for i, tin := range finalTx.TxIn {
		// We'll grab the corresponding internal packet input which
		// matches this materialized transaction input and emplace that
		// final sigScript (if present).
		pInput := p.Inputs[i]
		if pInput.FinalScriptSig != nil {
			tin.SignatureScript = pInput.FinalScriptSig
		}

	}

	return finalTx, nil
}
