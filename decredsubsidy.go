package picfightcoin

import (
	"github.com/jfixby/bignum"
	"github.com/jfixby/coin"
	"math"
)

type DecredSubsidyParams struct {
	// Subsidy parameters.
	//
	// Subsidy calculation for exponential reductions:
	// 0 for i in range (0, height / SubsidyReductionInterval):
	// 1     subsidy *= MulSubsidy
	// 2     subsidy /= DivSubsidy
	//
	// Caveat: Don't overflow the int64 register!!

	// BaseSubsidy is the starting subsidy amount for mined blocks.
	BaseSubsidy int64

	// Subsidy reduction multiplier.
	MulSubsidy int64

	// Subsidy reduction divisor.
	DivSubsidy int64

	// SubsidyReductionInterval is the reduction interval in blocks.
	SubsidyReductionInterval int64
}

var decredSubsidy = &decredMainNetSubsidyCalculator{
	subsidyParams: DecredSubsidyParams{
		BaseSubsidy:              3119582664,
		MulSubsidy:               100,
		DivSubsidy:               101,
		SubsidyReductionInterval: 6144,
	},
}

func DecredMainNetSubsidy() SubsidyCalculator {
	return decredSubsidy
}

type decredMainNetSubsidyCalculator struct {
	subsidyCache  map[uint64]int64
	subsidyParams DecredSubsidyParams
}

func (c *decredMainNetSubsidyCalculator) SetEngine(engine bignum.BigNumEngine) {
	panic("implement me")
}

func (c *decredMainNetSubsidyCalculator) ExpectedTotalNetworkSubsidy() coin.Amount {
	return coin.Amount{2103834590794301}
	// value received by block-by-block testing
}

func (decredMainNetSubsidyCalculator) NumberOfGeneratingBlocks() int64 {
	return math.MaxInt64
}

func (c *decredMainNetSubsidyCalculator) PreminedCoins() coin.Amount {
	return coin.Amount{c.BlockOneSubsidy()}
}

func (c *decredMainNetSubsidyCalculator) CalcBlockWorkSubsidy(height int64, voters uint16) int64 {
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

func (c *decredMainNetSubsidyCalculator) CalcStakeVoteSubsidy(height int64) int64 {
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

func (decredMainNetSubsidyCalculator) FirstGeneratingBlockIndex() int64 {
	// 0 - genesis block
	// 1 - premine block
	// and the
	return 2 // - is the first generating block
}

func (c *decredMainNetSubsidyCalculator) CalcBlockTaxSubsidy(height int64, voters uint16) int64 {
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

func (c *decredMainNetSubsidyCalculator) CalcBlockSubsidy(height int64) int64 {
	// Block height 1 subsidy is 'special' and used to
	// distribute initial tokens, if any.
	if height == 1 {
		return c.BlockOneSubsidy()
	}

	iteration := uint64(height / c.subsidyParams.SubsidyReductionInterval)

	if iteration == 0 {
		return c.subsidyParams.BaseSubsidy
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
		cachedValue *= c.subsidyParams.MulSubsidy
		cachedValue /= c.subsidyParams.DivSubsidy

		c.subsidyCache[iteration] = cachedValue

		return cachedValue
	}

	// Calculate the subsidy from scratch and store in the
	// blockSubsidyCache. TODO If there's an older item in the blockSubsidyCache,
	// calculate it from that to save time.
	subsidy := c.subsidyParams.BaseSubsidy
	for i := uint64(0); i < iteration; i++ {
		subsidy *= c.subsidyParams.MulSubsidy
		subsidy /= c.subsidyParams.DivSubsidy
	}

	c.subsidyCache[iteration] = subsidy

	return subsidy
}

func (c *decredMainNetSubsidyCalculator) WorkRewardProportion() uint16 {
	return 6
}

func (c *decredMainNetSubsidyCalculator) TotalSubsidyProportions() uint16 {
	return c.WorkRewardProportion() + c.StakeRewardProportion() + c.BlockTaxProportion()
}

func (c *decredMainNetSubsidyCalculator) TicketsPerBlock() uint16 {
	return 5
}

func (c *decredMainNetSubsidyCalculator) BlockOneSubsidy() int64 {
	return 168000000000000
}

func (c *decredMainNetSubsidyCalculator) StakeRewardProportion() uint16 {
	return 3
}

func (c *decredMainNetSubsidyCalculator) BlockTaxProportion() uint16 {
	return 1
}

func (c *decredMainNetSubsidyCalculator) StakeValidationHeight() int64 {
	return 4096 // ~14 days
}

func (c *decredMainNetSubsidyCalculator) EstimateSupply(height int64) int64 {
	return EstimateDecredSupply(&c.subsidyParams, height, c.BlockOneSubsidy())
}
