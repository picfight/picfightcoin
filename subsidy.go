package picfightcoin

import (
	"github.com/jfixby/bignum"
	"github.com/jfixby/coin"
)

type SubsidyCalculator interface {
	ExpectedTotalNetworkSubsidy() coin.Amount
	NumberOfGeneratingBlocks() int64
	PreminedCoins() coin.Amount
	FirstGeneratingBlockIndex() int64
	CalcBlockWorkSubsidy(height int64, voters uint16) int64
	CalcStakeVoteSubsidy(height int64) int64
	CalcBlockTaxSubsidy(height int64, voters uint16) int64
	CalcBlockSubsidy(height int64) int64
	TicketsPerBlock() uint16
	SetEngine(engine bignum.BigNumEngine)
	EstimateSupply(height int64) int64
	WorkRewardProportion() uint16
	StakeRewardProportion() uint16
	BlockTaxProportion() uint16
	StakeValidationHeight() int64
}
