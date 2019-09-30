// Copyright (c) The PicFight coin developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package picfightcoin

import (
	"github.com/jfixby/bignum"
	"github.com/jfixby/coin"
	"github.com/jfixby/difficulty"
	"time"
)

// NetworkPoWLimit is proof-of-work limit parameter.
func NetworkPoWLimit() *difficulty.Difficulty {
	return difficulty.NewDifficultyFromHashString( //
		"00 00 00 ff ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
}

func GenesisBlockPowBits() uint32 {
	return NetworkPoWLimit().ToCompact()
}

func GenesisBlockTimestamp() time.Time {
	return time.Unix(1569855611, 0)
}

func DNSSeeds() []string {
	return []string{
		"eu-01.seed.picfight.org",
		"eu-02.seed.picfight.org",
		//
		"us-01.seed.picfight.org",
		"us-02.seed.picfight.org",
		//
		"ap-01.seed.picfight.org",
		"ap-02.seed.picfight.org",
	}
}

const projectPremineTotal = 4000000.0 // 4M

// tickets_per_block(5) * (mature_time(256) + vote(1) + mature_time(256)) * coins_per_ticket(2)
// 5 * (256 + 1 + 256) * 2 = 5130 (fees excluded)
const projectPreminePoS = 6000

func Premine() map[string]coin.Amount {
	return map[string]coin.Amount{
		"JsKEwugutnHaQ71DaPqfN5hP1VYapSLpw92":// PROJECT PREMINE
		coin.FromFloat(projectPremineTotal - projectPreminePoS),
		"JsXVKweatCh63iX6NsRM5ahGbrb7XRKGwE8":// PROJECT PoS-SECURITY LAYER
		coin.FromFloat(projectPreminePoS),
	}
}

var PremineTotal = calcPremineTotal()

func calcPremineTotal() coin.Amount {
	premine := Premine()
	sum := coin.Amount{0}
	for _, amount := range premine {
		sum.AtomsValue = sum.AtomsValue + amount.AtomsValue
	}
	return sum
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

// Organization related parameters
// Organization address is ?
func OrganizationPkScript() []byte {
	return hexDecode("a914f5916158e3e2c4551c1796708db8367207ed13bb87")
}

// PicfightCoinWire represents the picfight coin network.
const PicfightCoinWire uint32 = 0xd9b488ff

var subsidy = &PicfightCoinSubsidyCalculator{
	engine: bignum.Float64Engine{},
}

func PicFightCoinSubsidy() SubsidyCalculator {
	return subsidy
}
