// Copyright (c) The PicFight coin developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package picfightcoin

import (
	"github.com/jfixby/difficulty"
	"time"
)

// NetworkPoWLimit is proof-of-work limit parameter.
func NetworkPoWLimit() *difficulty.Difficulty {
	return difficulty.NewDifficultyFromHashString( //
		"00 00 ff ff ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
}

func GenesisBlockPowBits() uint32 {
	return NetworkPoWLimit().ToCompact()
}

func GenesisBlockTimestamp() time.Time {
	return time.Unix(1569336596, 0)
}

func DNSSeeds() []string {
	return []string{
		"eu-01.seed.picfight.org",
		"eu-02.seed.picfight.org",
	}
}

const atomsPerCoin = 1e8

func Premine() map[string]int64 {
	return map[string]int64{
		"JsCVh5SVDQovpW1dswaZNan2mfNWy6uRpPx": 4000000 * atomsPerCoin,
	}
}

const NetworkAddressPrefix = "J"

var (
	// Address encoding magics
	PubKeyAddrID     = [2]byte{0x1b, 0x2d} // starts with Jk
	PubKeyHashAddrID = [2]byte{0x0a, 0x0f} // starts with Js
	PKHEdwardsAddrID = [2]byte{0x09, 0xef} // starts with Je
	PKHSchnorrAddrID = [2]byte{0x09, 0xd1} // starts with JS
	ScriptHashAddrID = [2]byte{0x09, 0xea} // starts with Jc
	PrivateKeyID     = [2]byte{0x22, 0xce} // starts with Pj
)
