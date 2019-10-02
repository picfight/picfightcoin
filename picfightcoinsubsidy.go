package picfightcoin

import (
	"github.com/jfixby/bignum"
	"github.com/jfixby/coin"
	picfightcoin "github.com/jfixby/lineardown"
	"github.com/jfixby/pin"
	"time"
)

type PicfightCoinSubsidyCalculator struct {
	engine                      bignum.BigNumEngine
	blockSubsidyCache           map[int64]bignum.BigNum
	blockByBlockSupplyEstimator *BlockByBlockSupplyEstimator
}

func (c *PicfightCoinSubsidyCalculator) WorkRewardProportion() uint16 {
	return 6 // 60%
}

func (c *PicfightCoinSubsidyCalculator) StakeRewardProportion() uint16 {
	return 4 // 40%
}

func (c *PicfightCoinSubsidyCalculator) BlockTaxProportion() uint16 {
	return 0 // 0%
}

func (c *PicfightCoinSubsidyCalculator) SubsidyReductionInterval() int64 {
	return 1 // each block
}

func (c *PicfightCoinSubsidyCalculator) EstimateSupply(height int64) int64 {
	if c.blockByBlockSupplyEstimator == nil {
		c.blockByBlockSupplyEstimator = &BlockByBlockSupplyEstimator{
			SubsidyCalculator: c,
		}
	}
	return c.blockByBlockSupplyEstimator.Estimate(height)
}

func (c *PicfightCoinSubsidyCalculator) SetEngine(engine bignum.BigNumEngine) {
	pin.AssertNotNil("engine", engine)
	c.engine = engine
}

func (c *PicfightCoinSubsidyCalculator) ExpectedTotalNetworkSubsidy() coin.Amount {
	return coin.FromFloat(8000000) // 8M
}

func (c *PicfightCoinSubsidyCalculator) NumberOfGeneratingBlocks() int64 {
	targetTimePerBlock := time.Minute * 5
	DAY := time.Hour * 24
	YEAR := DAY * 365
	SubsidyGeneratingPeriod := YEAR * 44
	numberOfGeneratingBlocks := int64(SubsidyGeneratingPeriod / targetTimePerBlock)
	return numberOfGeneratingBlocks
	//return 13
}

func (c *PicfightCoinSubsidyCalculator) PreminedCoins() coin.Amount {
	return PremineTotal.Copy()
}

func (c *PicfightCoinSubsidyCalculator) CalcBlockWorkSubsidy(height int64, voters uint16) int64 {
	blockSubsidy := c.CalcBlockSubsidy(height)
	stakeSubsidy := c.CalcStakeVoteSubsidy(height) * int64(c.TicketsPerBlock())

	subsidy := blockSubsidy - stakeSubsidy

	// Ignore the voters field of the header before we're at a point
	// where there are any voters.
	if height < c.StakeValidationHeight() {
		return subsidy
	}

	// If there are no voters, subsidy is 0. The block will fail later anyway.
	if voters == 0 {
		return 0
	}

	// Adjust for the number of voters. This shouldn't ever overflow if you start
	// with 50 * 10^8 Atoms and voters and potentialVoters are uint16.
	potentialVoters := c.TicketsPerBlock()
	actual := (int64(voters) * subsidy) / int64(potentialVoters)

	return actual
}

func (c *PicfightCoinSubsidyCalculator) CalcStakeVoteSubsidy(height int64) int64 {
	// Calculate the actual reward for this block, then further reduce reward
	// proportional to StakeRewardProportion.
	// Note that voters/potential voters is 1, so that vote reward is calculated
	// irrespective of block reward.
	subsidy := c.CalcBlockSubsidy(height)
	subsidy *= int64(c.StakeRewardProportion())
	subsidy /= 10 * int64(c.TicketsPerBlock())
	return subsidy
}

func (c *PicfightCoinSubsidyCalculator) FirstGeneratingBlockIndex() int64 {
	// 0 - genesis block
	// 1 - premine block
	// and the
	return 2 // - is the first generating block
}

func (c *PicfightCoinSubsidyCalculator) CalcBlockTaxSubsidy(height int64, voters uint16) int64 {
	//0% - no taxation, because we already did the taxation by premining
	return 0
}

func (c *PicfightCoinSubsidyCalculator) CalcBlockSubsidy(height int64) int64 {

	if height < 1 {
		return 0
	}
	if height == 1 {
		return PremineTotal.AtomsValue
	}
	if height > c.FirstGeneratingBlockIndex()+c.NumberOfGeneratingBlocks() {
		return 0
	}

	if c.blockSubsidyCache == nil {
		c.blockSubsidyCache = map[int64]bignum.BigNum{}
	}
	cached := c.blockSubsidyCache[height]
	if cached != nil {
		genCoins := coin.FromFloat(cached.ToFloat64())
		return genCoins.AtomsValue
	}
	engine := c.engine
	index := height - c.FirstGeneratingBlockIndex()
	generateTotalBlocks := c.NumberOfGeneratingBlocks()
	generateTotalCoins := c.ExpectedTotalNetworkSubsidy().AtomsValue - c.PreminedCoins().AtomsValue
	gen := picfightcoin.LinearDownGenerate(engine, generateTotalBlocks, coin.Amount{generateTotalCoins}, index)
	c.blockSubsidyCache[height] = gen
	genCoins := coin.FromFloat(gen.ToFloat64())
	return genCoins.AtomsValue
}

func (c *PicfightCoinSubsidyCalculator) StakeValidationHeight() int64 {
	return 4096 // ~14 days
}

func (c *PicfightCoinSubsidyCalculator) TicketsPerBlock() uint16 {
	return 5
}
