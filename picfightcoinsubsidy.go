package picfightcoin

import (
	"github.com/jfixby/bignum"
	"github.com/jfixby/coin"
	picfightcoin "github.com/jfixby/lineardown"
	"github.com/jfixby/pin"
)

type PicfightCoinSubsidyCalculator struct {
	engine bignum.BigNumEngine
}

func (c *PicfightCoinSubsidyCalculator) SetEngine(engine bignum.BigNumEngine) {
	pin.AssertNotNil("engine", engine)
	c.engine = engine
}

func (c *PicfightCoinSubsidyCalculator) ExpectedTotalNetworkSubsidy() coin.Amount {
	return coin.FromFloat(expectedTotalNetworkSubsidyCoins)
}

func (c *PicfightCoinSubsidyCalculator) NumberOfGeneratingBlocks() int64 {
	return numberOfGeneratingBlocks
}

func (c *PicfightCoinSubsidyCalculator) PreminedCoins() coin.Amount {
	return PremineTotal.Copy()
}

func (c *PicfightCoinSubsidyCalculator) CalcBlockWorkSubsidy(height int64, voters uint16) int64 {

	subsidy := c.CalcBlockSubsidy(height)

	proportionWork := WorkRewardProportion()
	proportions := TotalSubsidyProportions()
	subsidy *= proportionWork
	subsidy /= proportions

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
	subsidy *= StakeRewardProportion()
	subsidy /= (TotalSubsidyProportions() * int64(c.TicketsPerBlock()))
	return subsidy
}

func TotalSubsidyProportions() int64 {
	expected := WorkRewardProportion() + StakeRewardProportion() + BlockTaxProportion() // should be 100%
	pin.AssertTrue("TotalSubsidyProportions==100", expected == 100)
	return expected
}

func WorkRewardProportion() int64 {
	return int64(60) //60%
}

func StakeRewardProportion() int64 {
	return int64(40) //40%
}

func BlockTaxProportion() int64 {
	return int64(0) //0% - no taxation, because we already did the taxation by premining
}

func (c *PicfightCoinSubsidyCalculator) FirstGeneratingBlockIndex() int64 {
	// 0 - genesis block
	// 1 - premine block
	// and the
	return 2 // - is the first generating block
}

func (c *PicfightCoinSubsidyCalculator) CalcBlockTaxSubsidy(height int64, voters uint16) int64 {
	if BlockTaxProportion() == 0 {
		return 0
	}

	subsidy := c.CalcBlockSubsidy(height)

	proportionTax := BlockTaxProportion()
	proportions := TotalSubsidyProportions()
	subsidy *= proportionTax
	subsidy /= proportions

	// Assume all voters 'present' before stake voting is turned on.
	if height < c.StakeValidationHeight() {
		voters = 5
	}

	// If there are no voters, subsidy is 0. The block will fail later anyway.
	if voters == 0 && height >= c.StakeValidationHeight() {
		return 0
	}

	// Adjust for the number of voters. This shouldn't ever overflow if you start
	// with 50 * 10^8 Atoms and voters and potentialVoters are uint16.
	potentialVoters := c.TicketsPerBlock()
	adjusted := (int64(voters) * subsidy) / int64(potentialVoters)

	return adjusted
}

func (c *PicfightCoinSubsidyCalculator) CalcBlockSubsidy(height int64) int64 {

	if height == 1 {
		return PremineTotal.AtomsValue
	}
	if height < 1 {
		return 0
	}
	if height > c.FirstGeneratingBlockIndex()+c.NumberOfGeneratingBlocks() {
		return 0
	}
	engine := c.engine
	pin.AssertTrue("height>1", height > 1)
	index := height - c.FirstGeneratingBlockIndex()
	generateTotalBlocks := c.NumberOfGeneratingBlocks()
	generateTotalCoins := c.ExpectedTotalNetworkSubsidy().AtomsValue - c.PreminedCoins().AtomsValue
	gen := picfightcoin.LinearDownGenerate(engine, generateTotalBlocks, coin.Amount{generateTotalCoins}, index)
	genCoins := coin.FromFloat(gen.ToFloat64())
	return genCoins.AtomsValue
}

func (c *PicfightCoinSubsidyCalculator) StakeValidationHeight() int64 {
	return 4096 // ~14 days
}

func (c *PicfightCoinSubsidyCalculator) TicketsPerBlock() uint16 {
	return 5
}
