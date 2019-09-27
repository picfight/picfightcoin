package picfightcoin

import (
	"github.com/jfixby/difficulty"
	"time"
)

// PicfightCoinPowLimit is proof-of-work limit parameter.
func PicfightCoinPowLimit() *difficulty.Difficulty {
	return difficulty.NewDifficultyFromHashString( //
		"00 00 ff ff ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
}

func GenesisBlockTimestamp() time.Time {
	return time.Unix(1569336596, 0)
}
