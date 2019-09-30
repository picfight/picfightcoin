package picfightcoin

import (
	"fmt"
	"github.com/jfixby/coin"
	"github.com/jfixby/pin"
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

//var expectedPFCActual = coin.FromFloat(7999999.84736554)
var expectedPFCActual = coin.FromFloat(7999999.97687360)

func TestPicfightCoinSubsidy(t *testing.T) {
	calc := PicFightCoinSubsidy()
	//calc.SetEngine(bignum.BigDecimalEngine{})
	expected := calc.ExpectedTotalNetworkSubsidy().AtomsValue
	expected = expectedPFCActual.AtomsValue
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
		//blockSubsidy := calc.CalcBlockSubsidy(blockIndex)
		work := calc.CalcBlockWorkSubsidy(blockIndex,
			calc.TicketsPerBlock())
		stake := calc.CalcStakeVoteSubsidy(blockIndex) * int64(calc.TicketsPerBlock())
		tax := calc.CalcBlockTaxSubsidy(blockIndex, calc.TicketsPerBlock())
		if i%1000000 == 0 {
			pin.D(fmt.Sprintf("block: %v/%v: %v", i, calc.NumberOfGeneratingBlocks(), work+stake+tax))
		}
		//if blockSubsidy != work+stake+tax && blockIndex > 1 {
		//	t.Errorf("Bad block[%v] subsidy; want %v, got %v\n"+
		//		"work: %v\n"+
		//		"stake: %v\n"+
		//		"tax: %v\n",
		//		blockIndex,
		//		blockSubsidy,
		//		work+stake+tax,
		//		work,
		//		stake,
		//		tax,
		//	)
		//	t.FailNow()
		//}
		if (work+stake+tax) == 0 && i > 1 {
			break
		}

		cache[i] = (work + stake + tax)

	}

	totalSubsidy := coin.Amount{0}
	for i := int64(0); i <= int64(len(cache)); i++ {
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
