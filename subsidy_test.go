package picfightcoin

import (
	"fmt"
	"github.com/jfixby/bignum"
	"github.com/jfixby/coin"
	"testing"
)

func TestPremine(t *testing.T) {
	block1_subsidy := PicFightCoinSubsidy().CalcBlockSubsidy(1)

	// Block height 1 subsidy is 'special' and used to
	// distribute initial tokens, if any.
	block1_spremined := calcPremineTotal().AtomsValue

	if block1_subsidy != block1_spremined {
		t.Errorf("Premine mismatch: got %v expected %v ",
			block1_subsidy,
			block1_spremined,
		)
		t.Fail()
	}
}

func TestPicfightCoinSubsidy(t *testing.T) {
	calc := PicFightCoinSubsidy()
	calc.SetEngine(bignum.BigDecimalEngine{})
	expected := calc.ExpectedTotalNetworkSubsidy().AtomsValue
	fullSubsidyCheck(t, calc, expected)
}

func TestDecredSubsidy(t *testing.T) {
	calc := DecredSubsidy
	expected := calc.ExpectedTotalNetworkSubsidy().AtomsValue
	fullSubsidyCheck(t, calc, expected)
}

func fullSubsidyCheck(t *testing.T, calc SubsidyCalculator, expected int64) {

	cache := map[int64]int64{}
	for i := int64(0); ; i++ {
		blockIndex := i

		work := calc.CalcBlockWorkSubsidy(blockIndex,
			calc.TicketsPerBlock())
		stake := calc.CalcStakeVoteSubsidy(blockIndex) * int64(calc.TicketsPerBlock())
		tax := calc.CalcBlockTaxSubsidy(blockIndex, calc.TicketsPerBlock())
		if (i%10000 == 0) {
			fmt.Println(fmt.Sprintf("block: %v/%v: %v", i, calc.NumberOfGeneratingBlocks(), work+stake+tax))
		}
		if (work+stake+tax) == 0 && i > 0 {
			break
		}

		cache[i] = (work + stake + tax)

	}

	totalSubsidy := coin.Amount{0}
	for i := int64(0); i < int64(len(cache)); i++ {
		k := int64(len(cache)) - 1 - i
		totalSubsidy.AtomsValue = totalSubsidy.AtomsValue + cache[k]
	}
	fmt.Println(fmt.Sprintf("total: %v", totalSubsidy.AtomsValue))
	expectedTotal := coin.Amount{expected}
	if totalSubsidy.AtomsValue != expectedTotal.AtomsValue {
		t.Errorf("Bad total subsidy; want %v, got %v",
			expectedTotal.AtomsValue,
			totalSubsidy.AtomsValue,
		)
		t.Errorf("Bad total subsidy; want %v, got %v",
			expectedTotal,
			totalSubsidy,
		)
	}
}

// originalTestExpected is value from the original decred/dcrd repo
// most likely is invalid due to incorrect testing
const originalTestExpected int64 = 2099999999800912

func TestDecredSubsidyOriginal(t *testing.T) {
	calc := DecredSubsidy
	expected := calc.ExpectedTotalNetworkSubsidy().AtomsValue
	expected = originalTestExpected
	originalDecredSubsidyCheck(t, calc, expected)
}

func originalDecredSubsidyCheck(t *testing.T, calc *DecredMainNetSubsidyCalculator, expected int64) {
	totalSubsidy := calc.BlockOneSubsidy()
	for i := int64(0); ; i++ {
		// Genesis block or first block.
		if i == 0 || i == 1 {
			continue
		}

		if i%calc.SubsidyReductionInterval() == 0 {
			numBlocks := calc.SubsidyReductionInterval()
			// First reduction internal, which is reduction interval - 2
			// to skip the genesis block and block one.
			if i == calc.SubsidyReductionInterval() {
				numBlocks -= 2
			}
			height := i - numBlocks

			work := calc.CalcBlockWorkSubsidy(height, calc.TicketsPerBlock())
			stake := calc.CalcStakeVoteSubsidy(height) * int64(calc.TicketsPerBlock())
			tax := calc.CalcBlockTaxSubsidy(height, calc.TicketsPerBlock())
			if (work + stake + tax) == 0 {
				break
			}
			totalSubsidy += ((work + stake + tax) * numBlocks)

			// First reduction internal, subtract the stake subsidy for
			// blocks before the staking system is enabled.
			if i == calc.SubsidyReductionInterval() {
				totalSubsidy -= stake * (calc.StakeValidationHeight() - 2)
			}
		}
	}
	if totalSubsidy != expected {
		t.Errorf("Bad total subsidy; want %v, got %v", expected, totalSubsidy)
	}
}
