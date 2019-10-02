package picfightcoin

type BlockByBlockSupplyEstimator struct {
	cache             map[int64]*int64
	SubsidyCalculator SubsidyCalculator
}

func (c *BlockByBlockSupplyEstimator) Estimate(height int64) int64 {
	if height < 1 {
		return 0
	}
	if c.cache == nil {
		c.cache = make(map[int64]*int64)
		zero := int64(0)
		c.cache[0] = &zero
	}
	cached := c.cache[height]
	if cached != nil {
		return *cached
	}

	highestCachedIndex := int64(0)
	for i := height - 1; ; i-- {
		cached := c.cache[i]
		if cached != nil {
			highestCachedIndex = i
			break
		}
	}

	for i := highestCachedIndex + 1; i <= height; i++ {
		prev := c.Estimate(i - 1)
		curr := prev + c.SubsidyCalculator.CalcBlockSubsidy(i)
		c.cache[i] = &curr
	}

	cached = c.cache[height]
	if cached == nil {
		panic("invalid state")
	}
	return *cached
}

func EstimateSupplyWithCache(cache map[int64]*int64, height int64, CalcBlockSubsidy func(height int64) int64) int64 {
	cached := cache[height]
	if cached != nil {
		return *cached
	}

	if height < 1 {
		return 0
	}

	if height < 10 {
		prevEst := EstimateSupplyWithCache(cache, height-1, CalcBlockSubsidy)
		est := prevEst + CalcBlockSubsidy(height)
		cache[height] = &est
		return est
	}

	highrstCachedIndex := int64(0)
	for i := height - 1; ; i-- {
		cached := cache[i]
		if cached != nil {
			highrstCachedIndex = i
			break
		}
	}
	for i := highrstCachedIndex; i <= height; i++ {
		EstimateSupplyWithCache(cache, i, CalcBlockSubsidy)
	}

	cached = cache[height]
	if cached == nil {
		panic("implement me")
	}
	return *cached
}

func EstimateDecredSupply(params *DecredSubsidyParams, height int64, BlockOneSubsidy int64) int64 {
	if height <= 0 {
		return 0
	}
	// Estimate the supply by calculating the full block subsidy for each
	// reduction interval and multiplying it the number of blocks in the
	// interval then adding the subsidy produced by number of blocks in the
	// current interval.
	supply := BlockOneSubsidy
	reductions := height / params.SubsidyReductionInterval
	subsidy := params.BaseSubsidy
	for i := int64(0); i < reductions; i++ {
		supply += params.SubsidyReductionInterval * subsidy

		subsidy *= params.MulSubsidy
		subsidy /= params.DivSubsidy
	}
	supply += (1 + height%params.SubsidyReductionInterval) * subsidy

	// Blocks 0 and 1 have special subsidy amounts that have already been
	// added above, so remove what their subsidies would have normally been
	// which were also added above.
	supply -= params.BaseSubsidy * 2

	return supply
}
