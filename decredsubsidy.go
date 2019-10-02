package picfightcoin

import (
	"github.com/jfixby/bignum"
	"github.com/jfixby/coin"
	"math"
)

type DecredSubsidyCalculator interface {
	BlockOneSubsidy() int64
	BaseSubsidy() int64
	SubsidyReductionInterval() int64
	MulSubsidy() int64
	DivSubsidy() int64
}

var decredSubsidy = &DecredMainNetSubsidyCalculator{}

func DecredMainNetSubsidy() SubsidyCalculator {
	return decredSubsidy
}

type DecredMainNetSubsidyCalculator struct {
	subsidyCache map[uint64]int64
}

func (c *DecredMainNetSubsidyCalculator) SetEngine(engine bignum.BigNumEngine) {
	panic("implement me")
}

func (c *DecredMainNetSubsidyCalculator) ExpectedTotalNetworkSubsidy() coin.Amount {
	return coin.Amount{2103834590794301}
	// value received by block-by-block testing
}

func (DecredMainNetSubsidyCalculator) NumberOfGeneratingBlocks() int64 {
	return math.MaxInt64
}

func (c *DecredMainNetSubsidyCalculator) PreminedCoins() coin.Amount {
	return coin.Amount{c.BlockOneSubsidy()}
}

func (c *DecredMainNetSubsidyCalculator) CalcBlockWorkSubsidy(height int64, voters uint16) int64 {
	subsidy := c.CalcBlockSubsidy(height)

	proportionWork := int64(c.WorkRewardProportion())
	proportions := int64(c.TotalSubsidyProportions())
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

func (c *DecredMainNetSubsidyCalculator) CalcStakeVoteSubsidy(height int64) int64 {
	// Calculate the actual reward for this block, then further reduce reward
	// proportional to StakeRewardProportion.
	// Note that voters/potential voters is 1, so that vote reward is calculated
	// irrespective of block reward.
	subsidy := c.CalcBlockSubsidy(height)

	proportionStake := int64(c.StakeRewardProportion())
	proportions := int64(c.TotalSubsidyProportions())
	subsidy *= proportionStake
	subsidy /= (proportions * int64(c.TicketsPerBlock()))

	return subsidy
}

func (DecredMainNetSubsidyCalculator) FirstGeneratingBlockIndex() int64 {
	// 0 - genesis block
	// 1 - premine block
	// and the
	return 2 // - is the first generating block
}

func (c *DecredMainNetSubsidyCalculator) CalcBlockTaxSubsidy(height int64, voters uint16) int64 {
	if c.BlockTaxProportion() == 0 {
		return 0
	}

	subsidy := c.CalcBlockSubsidy(height)

	proportionTax := int64(c.BlockTaxProportion())
	proportions := int64(c.TotalSubsidyProportions())
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

func (c *DecredMainNetSubsidyCalculator) CalcBlockSubsidy(height int64) int64 {
	// Block height 1 subsidy is 'special' and used to
	// distribute initial tokens, if any.
	if height == 1 {
		return c.BlockOneSubsidy()
	}

	iteration := uint64(height / c.SubsidyReductionInterval())

	if iteration == 0 {
		return c.BaseSubsidy()
	}

	if c.subsidyCache == nil {
		c.subsidyCache = make(map[uint64]int64)
	}

	// First, check the blockSubsidyCache.
	cachedValue, existsInCache := c.subsidyCache[iteration]
	if existsInCache {
		return cachedValue
	}

	// Is the previous one in the blockSubsidyCache? If so, calculate
	// the subsidy from the previous known value and store
	// it in the database and the blockSubsidyCache.
	cachedValue, existsInCache = c.subsidyCache[iteration-1]
	if existsInCache {
		cachedValue *= c.MulSubsidy()
		cachedValue /= c.DivSubsidy()

		c.subsidyCache[iteration] = cachedValue

		return cachedValue
	}

	// Calculate the subsidy from scratch and store in the
	// blockSubsidyCache. TODO If there's an older item in the blockSubsidyCache,
	// calculate it from that to save time.
	subsidy := c.BaseSubsidy()
	for i := uint64(0); i < iteration; i++ {
		subsidy *= c.MulSubsidy()
		subsidy /= c.DivSubsidy()
	}

	c.subsidyCache[iteration] = subsidy

	return subsidy
}

func (c *DecredMainNetSubsidyCalculator) WorkRewardProportion() uint16 {
	return 6
}

func (c *DecredMainNetSubsidyCalculator) SubsidyReductionInterval() int64 {
	return 6144
}

func (c *DecredMainNetSubsidyCalculator) TotalSubsidyProportions() uint16 {
	return c.WorkRewardProportion() + c.StakeRewardProportion() + c.BlockTaxProportion()
}

func (c *DecredMainNetSubsidyCalculator) TicketsPerBlock() uint16 {
	return 5
}

func (c *DecredMainNetSubsidyCalculator) BlockOneSubsidy() int64 {
	return 168000000000000
}

func (c *DecredMainNetSubsidyCalculator) StakeRewardProportion() uint16 {
	return 3
}

func (c *DecredMainNetSubsidyCalculator) BlockTaxProportion() uint16 {
	return 1
}

func (c *DecredMainNetSubsidyCalculator) BaseSubsidy() int64 {
	return 3119582664
}

func (c *DecredMainNetSubsidyCalculator) MulSubsidy() int64 {
	return 100
}

func (c *DecredMainNetSubsidyCalculator) DivSubsidy() int64 {
	return 101
}

func (c *DecredMainNetSubsidyCalculator) StakeValidationHeight() int64 {
	return 4096 // ~14 days
}

func (c *DecredMainNetSubsidyCalculator) EstimateSupply(height int64) int64 {
	return EstimateDecredSupply(c, height)
}
