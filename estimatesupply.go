package picfightcoin

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

func EstimateDecredSupply(c *DecredSubsidyParams, height int64, BlockOneSubsidy int64) int64 {
	if height <= 0 {
		return 0
	}

	// Estimate the supply by calculating the full block subsidy for each
	// reduction interval and multiplying it the number of blocks in the
	// interval then adding the subsidy produced by number of blocks in the
	// current interval.
	supply := BlockOneSubsidy
	reductions := height / c.SubsidyReductionInterval
	subsidy := c.BaseSubsidy
	for i := int64(0); i < reductions; i++ {
		supply += c.SubsidyReductionInterval * subsidy

		subsidy *= c.MulSubsidy
		subsidy /= c.DivSubsidy
	}
	supply += (1 + height%c.SubsidyReductionInterval) * subsidy

	// Blocks 0 and 1 have special subsidy amounts that have already been
	// added above, so remove what their subsidies would have normally been
	// which were also added above.
	supply -= c.BaseSubsidy * 2

	return supply
}
