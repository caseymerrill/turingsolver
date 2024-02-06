package verifiers

import "strings"

type VerifierCard struct {
	Verifiers []*Verifier
}

var Cards = []VerifierCard{
	// 1
	{
		Verifiers: []*Verifier{
			EqualsNumber(0, 1),
			GreaterThanNumber(0, 1),
		},
	},
	// 2
	{
		Verifiers: []*Verifier{
			LessThanNumber(0, 3),
			EqualsNumber(0, 3),
			GreaterThanNumber(0, 3),
		},
	},
	// 3
	{
		Verifiers: []*Verifier{
			LessThanNumber(1, 3),
			EqualsNumber(1, 3),
			GreaterThanNumber(1, 3),
		},
	},
	// 4
	{
		Verifiers: []*Verifier{
			LessThanNumber(1, 4),
			EqualsNumber(1, 4),
			GreaterThanNumber(1, 4),
		},
	},
	// 5
	{
		Verifiers: []*Verifier{
			Even(0),
			Odd(0),
		},
	},
	// 6
	{
		Verifiers: []*Verifier{
			Even(1),
			Odd(1),
		},
	},
	// 7
	{
		Verifiers: []*Verifier{
			Even(2),
			Odd(2),
		},
	},
	// 8
	{
		Verifiers: []*Verifier{
			NumberAppearsTimes(1, 0),
			NumberAppearsTimes(1, 1),
			NumberAppearsTimes(1, 2),
			NumberAppearsTimes(1, 3),
		},
	},
	// 9
	{
		Verifiers: []*Verifier{
			NumberAppearsTimes(3, 0),
			NumberAppearsTimes(3, 1),
			NumberAppearsTimes(3, 2),
			NumberAppearsTimes(3, 3),
		},
	},
	// 10
	{
		Verifiers: []*Verifier{
			NumberAppearsTimes(4, 0),
			NumberAppearsTimes(4, 1),
			NumberAppearsTimes(4, 2),
			NumberAppearsTimes(4, 3),
		},
	},
	// 11
	{
		Verifiers: []*Verifier{
			PositionLessThanPosition(0, 1),
			PositionEqualsPosition(0, 1),
			PositionGreaterThanPosition(0, 1),
		},
	},
	// 12
	{
		Verifiers: []*Verifier{
			PositionLessThanPosition(0, 2),
			PositionEqualsPosition(0, 2),
			PositionGreaterThanPosition(0, 2),
		},
	},
	// 13
	{
		Verifiers: []*Verifier{
			PositionLessThanPosition(1, 2),
			PositionEqualsPosition(1, 2),
			PositionGreaterThanPosition(1, 2),
		},
	},
	// 14
	{
		Verifiers: []*Verifier{
			PositionIsSmallest(0),
			PositionIsSmallest(1),
			PositionIsSmallest(2),
		},
	},
	// 15
	{
		Verifiers: []*Verifier{
			PositionIsLargest(0),
			PositionIsLargest(1),
			PositionIsLargest(2),
		},
	},
	// 16
	{
		Verifiers: []*Verifier{
			MoreEvens(),
			MoreOdds(),
		},
	},
	// 17
	{
		Verifiers: []*Verifier{
			NumberOfEvens(0),
			NumberOfEvens(1),
			NumberOfEvens(2),
			NumberOfEvens(3),
		},
	},
	// 18
	{
		Verifiers: []*Verifier{
			SummationIsEven(),
			SummationIsOdd(),
		},
	},
	// 19
	{
		Verifiers: []*Verifier{
			SumOfTwoPositionsLessThanNumber(0, 1, 6),
			SumOfTwoPositionsEqualsNumber(0, 1, 6),
			SumOfTwoPositionsGreaterThanNumber(0, 1, 6),
		},
	},
	// 20
	{
		Verifiers: []*Verifier{
			RepeatsTimes(3),
			RepeatsTimes(2),
			RepeatsTimes(1),
		},
	},
	// 21
	{
		Verifiers: []*Verifier{
			NoPairs(),
			RepeatsTimes(2),
		},
	},
	// 22
	{
		Verifiers: []*Verifier{
			Ascending(),
			Descending(),
			NoOrder(),
		},
	},
	// 23
	{
		Verifiers: []*Verifier{
			SummationLessThanNumber(6),
			SummationEqualsNumber(6),
			SummationGreaterThanNumber(6),
		},
	},
	// 24
	{
		Verifiers: []*Verifier{
			AscendingSequenceOfSize(3),
			AscendingSequenceOfSize(2),
			AscendingSequenceOfSize(1),
		},
	},
	// 25
	{
		Verifiers: []*Verifier{
			AscendingOrDecendingSequenceOfSize(1),
			AscendingOrDecendingSequenceOfSize(2),
			AscendingOrDecendingSequenceOfSize(3),
		},
	},
	// 26
	{
		Verifiers: []*Verifier{
			LessThanNumber(0, 3),
			LessThanNumber(1, 3),
			LessThanNumber(2, 3),
		},
	},
	// 27
	{
		Verifiers: []*Verifier{
			LessThanNumber(0, 4),
			LessThanNumber(1, 4),
			LessThanNumber(2, 4),
		},
	},
	// 28
	{
		Verifiers: []*Verifier{
			EqualsNumber(0, 1),
			EqualsNumber(1, 1),
			EqualsNumber(2, 1),
		},
	},
	// 29
	{
		Verifiers: []*Verifier{
			EqualsNumber(0, 3),
			EqualsNumber(1, 3),
			EqualsNumber(2, 3),
		},
	},
	// 30
	{
		Verifiers: []*Verifier{
			EqualsNumber(0, 4),
			EqualsNumber(1, 4),
			EqualsNumber(2, 4),
		},
	},
	// 31
	{
		Verifiers: []*Verifier{
			GreaterThanNumber(0, 1),
			GreaterThanNumber(1, 1),
			GreaterThanNumber(2, 1),
		},
	},
	// 32
	{
		Verifiers: []*Verifier{
			GreaterThanNumber(0, 3),
			GreaterThanNumber(1, 3),
			GreaterThanNumber(2, 3),
		},
	},
	// 33
	{
		Verifiers: []*Verifier{
			Even(0),
			Odd(0),
			Even(1),
			Odd(1),
			Even(2),
			Odd(2),
		},
	},
	// 34
	{
		Verifiers: []*Verifier{
			PositionIsSmallestOrEqual(0),
			PositionIsSmallestOrEqual(1),
			PositionIsSmallestOrEqual(2),
		},
	},
	// 35
	{
		Verifiers: []*Verifier{
			PositionIsLargestOrEqual(0),
			PositionIsLargestOrEqual(1),
			PositionIsLargestOrEqual(2),
		},
	},
	// 36
	{
		Verifiers: []*Verifier{
			SummationIsMultipleOf(3),
			SummationIsMultipleOf(4),
			SummationIsMultipleOf(5),
		},
	},
	// 37
	{
		Verifiers: []*Verifier{
			SumOfTwoPositionsEqualsNumber(0, 1, 4),
			SumOfTwoPositionsEqualsNumber(0, 2, 4),
			SumOfTwoPositionsEqualsNumber(1, 2, 4),
		},
	},
	// 38
	{
		Verifiers: []*Verifier{
			SumOfTwoPositionsEqualsNumber(0, 1, 6),
			SumOfTwoPositionsEqualsNumber(0, 2, 6),
			SumOfTwoPositionsEqualsNumber(1, 2, 6),
		},
	},
	// 39
	{
		Verifiers: []*Verifier{
			EqualsNumber(0, 1),
			EqualsNumber(1, 1),
			EqualsNumber(2, 1),
			GreaterThanNumber(0, 1),
			GreaterThanNumber(1, 1),
			GreaterThanNumber(2, 1),
		},
	},
	// 40
	{
		Verifiers: []*Verifier{
			LessThanNumber(0, 3),
			LessThanNumber(1, 3),
			LessThanNumber(2, 3),
			EqualsNumber(0, 3),
			EqualsNumber(1, 3),
			EqualsNumber(2, 3),
			GreaterThanNumber(0, 3),
			GreaterThanNumber(1, 3),
			GreaterThanNumber(2, 3),
		},
	},
	// 41
	{
		Verifiers: []*Verifier{
			LessThanNumber(0, 4),
			LessThanNumber(1, 4),
			LessThanNumber(2, 4),
			EqualsNumber(0, 4),
			EqualsNumber(1, 4),
			EqualsNumber(2, 4),
			GreaterThanNumber(0, 4),
			GreaterThanNumber(1, 4),
			GreaterThanNumber(2, 4),
		},
	},
	// 42
	{
		Verifiers: []*Verifier{
			PositionIsSmallest(0),
			PositionIsSmallest(1),
			PositionIsSmallest(2),
			PositionIsLargest(0),
			PositionIsLargest(1),
			PositionIsLargest(2),
		},
	},
	// 43
	{
		Verifiers: []*Verifier{
			PositionLessThanPosition(0, 1),
			PositionLessThanPosition(0, 2),
			PositionEqualsPosition(0, 1),
			PositionEqualsPosition(0, 2),
			PositionGreaterThanPosition(0, 1),
			PositionGreaterThanPosition(0, 2),
		},
	},
	// 44
	{
		Verifiers: []*Verifier{
			PositionLessThanPosition(1, 0),
			PositionLessThanPosition(1, 2),
			PositionEqualsPosition(1, 0),
			PositionEqualsPosition(1, 2),
			PositionGreaterThanPosition(1, 0),
			PositionGreaterThanPosition(1, 2),
		},
	},
	// 45
	{
		Verifiers: []*Verifier{
			NumberAppearsTimes(1, 0),
			NumberAppearsTimes(1, 1),
			NumberAppearsTimes(1, 2),
			NumberAppearsTimes(3, 0),
			NumberAppearsTimes(3, 1),
			NumberAppearsTimes(3, 2),
		},
	},
	// 46
	{
		Verifiers: []*Verifier{
			NumberAppearsTimes(3, 0),
			NumberAppearsTimes(3, 1),
			NumberAppearsTimes(3, 2),
			NumberAppearsTimes(4, 0),
			NumberAppearsTimes(4, 1),
			NumberAppearsTimes(4, 2),
		},
	},
	// 47
	{
		Verifiers: []*Verifier{
			NumberAppearsTimes(1, 0),
			NumberAppearsTimes(1, 1),
			NumberAppearsTimes(1, 2),
			NumberAppearsTimes(4, 0),
			NumberAppearsTimes(4, 1),
			NumberAppearsTimes(4, 2),
		},
	},
	// 48
	{
		Verifiers: []*Verifier{
			PositionLessThanPosition(0, 1),
			PositionLessThanPosition(0, 2),
			PositionLessThanPosition(1, 2),
			PositionEqualsPosition(0, 1),
			PositionEqualsPosition(0, 2),
			PositionEqualsPosition(1, 2),
			PositionGreaterThanPosition(0, 1),
			PositionGreaterThanPosition(0, 2),
			PositionGreaterThanPosition(1, 2),
		},
	},
}

func (vc VerifierCard) String() string {
	verifierDescriptions := make([]string, len(vc.Verifiers))
	for i, verifier := range vc.Verifiers {
		verifierDescriptions[i] = verifier.Description
	}
	return strings.Join(verifierDescriptions, " | ")
}
