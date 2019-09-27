// Copyright (c) The PicFight coin developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package picfightcoin

import (
	"bytes"
	"github.com/jfixby/difficulty"
	"strings"
	"testing"

	"github.com/decred/base58"
)

// TestPowLimitsAreConsistent ensures NetworkPoWLimit and PowLimitBits are consistent
// with each other
func TestPowLimitsAreConsistent(t *testing.T) {
	powLimitBigInt := NetworkPoWLimit().ToBigInt()
	powLimitCompact := NetworkPoWLimit().ToCompact()

	toBig := difficulty.CompactToBig(powLimitCompact)
	toCompact := difficulty.BigToCompact(powLimitBigInt)

	// Check params.PowLimitBits matches params.NetworkPoWLimit converted
	// into the compact form
	if toCompact != powLimitCompact {
		t.Fatalf("NetworkPoWLimit values mismatch:\n"+
			"params.NetworkPoWLimit    :%064x\n"+
			"                   :%x\n"+
			"params.PowLimitBits:%064x\n"+
			"                   :%x\n"+
			"params.NetworkPoWLimit is not consistent with the params.PowLimitBits",
			powLimitBigInt, toCompact, toBig, powLimitCompact)
	}
}

// TestGenesisBlockRespectsNetworkPowLimit ensures genesis.Header.Bits value
// is within the network PoW limit.
//
// Genesis header bits define starting difficulty of the network.
// Header bits of each block define target difficulty of the subsequent block.
//
// The first few solved blocks of the network will inherit the genesis block
// bits value before the difficulty reajustment takes place.
//
// Solved block shouldn't be rejected due to the PoW limit check.
//
// This test ensures these blocks will respect the network PoW limit.
func TestGenesisBlockRespectsNetworkPowLimit(t *testing.T) {
	bits := GenesisBlockPowBits()

	// Header bits as big.Int
	bitsAsBigInt := difficulty.CompactToBig(bits)

	// network PoW limit
	powLimitBigInt := NetworkPoWLimit().ToBigInt()

	if bitsAsBigInt.Cmp(powLimitBigInt) > 0 {
		t.Fatalf("Genesis block fails the consensus:\n"+
			"genesis.Header.Bits:%x\n"+
			"                   :%064x\n"+
			"params.NetworkPoWLimit    :%064x\n"+
			"genesis.Header.Bits "+
			"should respect network PoW limit",
			bits, bitsAsBigInt, powLimitBigInt)
	}
}

// checkPrefix checks if targetString starts with the given prefix
func checkPrefix(t *testing.T, prefix string, targetString string) {
	if strings.Index(targetString, prefix) != 0 {
		t.Logf("Address prefix mismatch: expected <%s> received <%s>",
			prefix, targetString)
		t.FailNow()
	}
}

// checkInterval creates two corner cases defining interval
// of all key values: [ xxxx000000000...0cccc , xxxx111111111...1cccc ],
// where xxxx - is the encoding magic, and cccc is a checksum.
// The interval is mapped to corresponding interval in base 58.
// Then prefixes are checked for mismatch.
func checkInterval(t *testing.T, desiredPrefix string, keySize int, magic [2]byte) {
	// min and max possible keys
	// all zeroes
	minKey := bytes.Repeat([]byte{0x00}, keySize)
	// all ones
	maxKey := bytes.Repeat([]byte{0xff}, keySize)

	base58interval := [2]string{
		base58.CheckEncode(minKey, magic),
		base58.CheckEncode(maxKey, magic),
	}
	checkPrefix(t, desiredPrefix, base58interval[0])
	checkPrefix(t, desiredPrefix, base58interval[1])
}

// TestAddressPrefixesAreConsistent ensures address encoding magics and
// NetworkAddressPrefix are consistent with each other.
// This test will light red when a new network is started with incorrect values.
func TestAddressPrefixesAreConsistent(t *testing.T) {
	P := NetworkAddressPrefix

	// Desired prefixes
	Pk := P + "k"
	Ps := P + "s"
	Pe := P + "e"
	PS := P + "S"
	Pc := P + "c"
	pk := "Pj"

	checkInterval(t, Pk, 33, PubKeyAddrID)
	checkInterval(t, Ps, 20, PubKeyHashAddrID)
	checkInterval(t, Pe, 20, PKHEdwardsAddrID)
	checkInterval(t, PS, 20, PKHSchnorrAddrID)
	checkInterval(t, Pc, 20, ScriptHashAddrID)
	checkInterval(t, pk, 33, PrivateKeyID)
}
