package verifiers

import "testing"

var codes = [][]int{
	{1, 2, 3},
	{3, 3, 3},
	{3, 2, 2},
	{1, 5, 1},
	{5, 4, 1},
}

func TestEqualsNumberFirstPosition(t *testing.T) {
	blueEqualsThree := EqualsNumber(0, 3)
	blueEqualsThreeExpected := []bool{false, true, true, false, false}

	for i, expected := range blueEqualsThreeExpected {
		if blueEqualsThree.Verify(codes[i]...) != expected {
			t.Fatalf("%v blue equals three expected: %v", codes[i], expected)
		}
	}
}

func TestEqualsNumberSecondPosition(t *testing.T) {
	yellowEqualsFive := EqualsNumber(1, 5)
	yellowEqualsFiveExpected := []bool{false, false, false, true, false}

	for i, expected := range yellowEqualsFiveExpected {
		if yellowEqualsFive.Verify(codes[i]...) != expected {
			t.Fatalf("%v yellow equals 5 expected: %v", codes[i], expected)
		}
	}
}

func TestGreaterThan(t *testing.T) {
	secondGreaterThanThree := GreaterThanNumber(1, 3)
	expectations := []bool{false, false, false, true, true}

	for i, expected := range expectations {
		if secondGreaterThanThree.Verify(codes[i]...) != expected {
			t.Fatalf("%v second > 3 expected: %v", codes[i], expected)
		}
	}
}

func TestLessThan(t *testing.T) {
	thirdLessThanThree := LessThanNumber(2, 3)
	expectations := []bool{false, false, true, true, true}

	for i, expected := range expectations {
		if thirdLessThanThree.Verify(codes[i]...) != expected {
			t.Fatalf("%v third < 3 expected: %v", codes[i], expected)
		}
	}
}

func TestEven(t *testing.T) {
	secondPositionEven := Even(1)
	expectations := []bool{true, false, true, false, true}

	for i, expected := range expectations {
		if secondPositionEven.Verify(codes[i]...) != expected {
			t.Fatalf("%v second is even expected; %v", codes[i], expected)
		}
	}
}

func TestOdd(t *testing.T) {
	secondPositionOdd := Odd(1)
	expectations := []bool{false, true, false, true, false}

	for i, expected := range expectations {
		if secondPositionOdd.Verify(codes[i]...) != expected {
			t.Fatalf("%v second is odd expected: %v", codes[i], expected)
		}
	}
}

func TestNumberAppearsTimes(t *testing.T) {
	oneAppearsTwice := NumberAppearsTimes(1, 2)
	expectations := []bool{false, false, false, true, false}
	for i, expected := range expectations {
		if oneAppearsTwice.Verify(codes[i]...) != expected {
			t.Fatalf("%v expected one twice", codes[i])
		}
	}
}

func TestPositionLessThanPosition(t *testing.T) {
	firstLessThanThird := PositionLessThanPosition(0, 2)
	expectations := []bool{true, false, false, false, false}

	for i, expected := range expectations {
		if firstLessThanThird.Verify(codes[i]...) != expected {
			t.Fatalf("%v expected first position to be less than third", codes[i])
		}
	}
}

func TestPositionEqualsPosition(t *testing.T) {
	firstEqualsThird := PositionEqualsPosition(0, 2)
	expectations := []bool{false, true, false, true, false}

	for i, expected := range expectations {
		if firstEqualsThird.Verify(codes[i]...) != expected {
			t.Fatalf("%v expected first to equal thrid", codes[i])
		}
	}
}

func TestPositionGreaterThanPosition(t *testing.T) {
	secondGreaterThanThird := PositionGreaterThanPosition(1, 2)
	expectations := []bool{false, false, false, true, true}

	for i, expected := range expectations {
		if secondGreaterThanThird.Verify(codes[i]...) != expected {
			t.Fatalf("%v expected second to be greater than third", codes[i])
		}
	}
}

func TestPositionIsSmallest(t *testing.T) {
	secondIsSmallest := PositionIsSmallest(1)
	expectations := []bool{false, false, false, false, false}

	for i, expected := range expectations {
		if secondIsSmallest.Verify(codes[i]...) != expected {
			t.Fatalf("%v expected second to be smallest", codes[i])
		}
	}
}

func TestPositionIsSmallestOrEqual(t *testing.T) {
	secondIsSmallestOrEqual := PositionIsSmallestOrEqual(1)
	expectations := []bool{false, true, true, false, false}

	for i, expected := range expectations {
		if secondIsSmallestOrEqual.Verify(codes[i]...) != expected {
			t.Fatalf("%v expected second to be smallest or equal", codes[i])
		}
	}
}

func TestPositionIsLargest(t *testing.T) {
	secondIsLargest := PositionIsLargest(1)
	expectations := []bool{false, false, false, true, false}

	for i, expected := range expectations {
		if secondIsLargest.Verify(codes[i]...) != expected {
			t.Fatalf("%v expected second to be largest", codes[i])
		}
	}
}

func TestPositionIsLargestOrEqual(t *testing.T) {
	secondIsLargestOrEqual := PositionIsLargestOrEqual(1)
	expectations := []bool{false, true, false, true, false}

	for i, expected := range expectations {
		if secondIsLargestOrEqual.Verify(codes[i]...) != expected {
			t.Fatalf("%v expected second to be largest or equal, expected: %v", codes[i], expected)
		}
	}
}

func TestMoreEvens(t *testing.T) {
	moreEvens := MoreEvens()
	expectations := []bool{false, false, true, false, false}

	for i, expected := range expectations {
		if moreEvens.Verify(codes[i]...) != expected {
			t.Fatalf("%v more evens than odds? expected: %v", codes[i], expected)
		}
	}
}

func TestMoreOdds(t *testing.T) {
	moreOdds := MoreOdds()
	expectations := []bool{true, true, false, true, true}

	for i, expected := range expectations {
		if moreOdds.Verify(codes[i]...) != expected {
			t.Fatalf("%v more odds than evens? expected: %v", codes[i], expected)
		}
	}
}

func TestNumberOfEvens(t *testing.T) {
	oneEven := NumberOfEvens(1)
	expectations := []bool{true, false, false, false, true}

	for i, expected := range expectations {
		if oneEven.Verify(codes[i]...) != expected {
			t.Fatalf("%v one even? expected: %v", codes[i], expected)
		}
	}
}

func TestSummationIsEven(t *testing.T) {
	sumIsEven := SummationIsEven()
	expectations := []bool{true, false, false, false, true}

	for i, expected := range expectations {
		if sumIsEven.Verify(codes[i]...) != expected {
			t.Fatalf("%v summation is even? expected: %v", codes[i], expected)
		}
	}
}

func TestSummationIsOdd(t *testing.T) {
	sumIsOdd := SummationIsOdd()
	expectations := []bool{false, true, true, true, false}

	for i, expected := range expectations {
		if sumIsOdd.Verify(codes[i]...) != expected {
			t.Fatalf("%v summation is odd? expected: %v", codes[i], expected)
		}
	}
}

func TestSumOfTwoPositionsLessThanNumber(t *testing.T) {
	sumOfFirstAndSecondLessThanFive := SumOfTwoPositionsLessThanNumber(0, 1, 5)
	expectations := []bool{true, false, false, false, false}

	for i, expected := range expectations {
		if sumOfFirstAndSecondLessThanFive.Verify(codes[i]...) != expected {
			t.Fatalf("%v sum of first and second less than 5? expected: %v", codes[i], expected)
		}
	}
}

func TestSumOfTwoPositionsEqualsNumber(t *testing.T) {
	sumOfFirstAndSecondEqualsFive := SumOfTwoPositionsEqualsNumber(0, 1, 5)
	expectations := []bool{false, false, true, false, false}

	for i, expected := range expectations {
		if sumOfFirstAndSecondEqualsFive.Verify(codes[i]...) != expected {
			t.Fatalf("%v sum of first and second equals 5? expected: %v", codes[i], expected)
		}
	}
}

func TestSumOfTwoPositionsGreaterThanNumber(t *testing.T) {
	sumOfFirstAndSecondGreaterThanFive := SumOfTwoPositionsGreaterThanNumber(0, 1, 5)
	expectations := []bool{false, true, false, true, true}

	for i, expected := range expectations {
		if sumOfFirstAndSecondGreaterThanFive.Verify(codes[i]...) != expected {
			t.Fatalf("%v sum of first and second greater than 5? expected: %v", codes[i], expected)
		}
	}
}

func TestRepeatsTimes(t *testing.T) {
	repeatsTwice := RepeatsTimes(2)
	expectations := []bool{false, false, true, true, false}

	for i, expected := range expectations {
		if repeatsTwice.Verify(codes[i]...) != expected {
			t.Fatalf("%v repeats twice? expected: %v", codes[i], expected)
		}
	}
}

func TestRepeatsThreeTimes(t *testing.T) {
	repeatsThree := RepeatsTimes(3)
	expectations := []bool{false, true, false, false, false}

	for i, expected := range expectations {
		if repeatsThree.Verify(codes[i]...) != expected {
			t.Fatalf("%v repeats three times? expected: %v", codes[i], expected)
		}
	}
}

func TestNoPairs(t *testing.T) {
	noPairs := NoPairs()
	expectations := []bool{true, true, false, false, true}

	for i, expected := range expectations {
		if noPairs.Verify(codes[i]...) != expected {
			t.Fatalf("%v has no pairs? expected: %v", codes[i], expected)
		}
	}
}

func TestAscending(t *testing.T) {
	asc := Ascending()
	expectations := []bool{true, false, false, false, false}

	for i, expected := range expectations {
		if asc.Verify(codes[i]...) != expected {
			t.Fatalf("%v is ascending? expected: %v", codes[i], expected)
		}
	}
}

func TestDecending(t *testing.T) {
	dsc := Descending()
	expectations := []bool{false, false, false, false, true}

	for i, expected := range expectations {
		if dsc.Verify(codes[i]...) != expected {
			t.Fatalf("%v is descending? expected: %v", codes[i], expected)
		}
	}
}

func TestNoOrder(t *testing.T) {
	noOrder := NoOrder()
	expectations := []bool{false, true, true, true, false}

	for i, expected := range expectations {
		if noOrder.Verify(codes[i]...) != expected {
			t.Fatalf("%v has no order? expected: %v", codes[i], expected)
		}
	}
}

func TestSummationLessThanNumber(t *testing.T) {
	summationLessThanNine := SummationLessThanNumber(9)
	expectations := []bool{true, false, true, true, false}

	for i, expected := range expectations {
		if summationLessThanNine.Verify(codes[i]...) != expected {
			t.Fatalf("%v summation less than 9? expected: %v", codes[i], expected)
		}
	}
}

func TestSummationEqualsNumber(t *testing.T) {
	summationEqualsSix := SummationEqualsNumber(6)
	expectations := []bool{true, false, false, false, false}

	for i, expected := range expectations {
		if summationEqualsSix.Verify(codes[i]...) != expected {
			t.Fatalf("%v summation equals 6? expected: %v", codes[i], expected)
		}
	}
}

func TestSummationGreaterThanNumber(t *testing.T) {
	summationGreaterThanNine := SummationGreaterThanNumber(9)
	expectations := []bool{false, false, false, false, true}

	for i, expected := range expectations {
		if summationGreaterThanNine.Verify(codes[i]...) != expected {
			t.Fatalf("%v summation greater than 9? expected: %v", codes[i], expected)
		}
	}
}

func TestAscendingSequenceOfSize(t *testing.T) {
	ascendingSequenceOfSizeThree := AscendingSequenceOfSize(3)
	expectations := []bool{true, false, false, false, false}

	for i, expected := range expectations {
		if ascendingSequenceOfSizeThree.Verify(codes[i]...) != expected {
			t.Fatalf("%v has an ascending sequence of size 3? expected: %v", codes[i], expected)
		}
	}
}

func TestDescendingSequenceOfSize(t *testing.T) {
	descendingSequenceOfSizeThree := DescendingSequenceOfSize(2)
	expectations := []bool{false, false, true, false, true}

	for i, expected := range expectations {
		if descendingSequenceOfSizeThree.Verify(codes[i]...) != expected {
			t.Fatalf("%v has a descending sequence of size 3? expected: %v", codes[i], expected)
		}
	}
}

func TestAscendingOrDescendingSequenceOfSize(t *testing.T) {
	NoSequence := AscendingOrDecendingSequenceOfSize(1)
	expectations := []bool{false, true, false, true, false}

	for i, expected := range expectations {
		if NoSequence.Verify(codes[i]...) != expected {
			t.Fatalf("%v has no sequence? expected: %v", codes[i], expected)
		}
	}
}

func TestSummationIsMultipleOf(t *testing.T) {
	summationIsMultipleOfThree := SummationIsMultipleOf(3)
	expectations := []bool{true, true, false, false, false}

	for i, expected := range expectations {
		if summationIsMultipleOfThree.Verify(codes[i]...) != expected {
			t.Fatalf("%v summation is multiple of 3? expected: %v", codes[i], expected)
		}
	}
}
