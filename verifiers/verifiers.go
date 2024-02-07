package verifiers

import "fmt"

type Verifier struct {
	Verify      func(n ...int) bool
	Description string
}

func (v *Verifier) String() string {
	return v.Description
}

func EqualsNumber(position Position, equalTo int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return n[position] == equalTo
		},
		Description: fmt.Sprintf("%v = %v", position, equalTo),
	}
}

func GreaterThanNumber(position Position, greaterThan int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return n[position] > greaterThan
		},
		Description: fmt.Sprintf("%v > %v", position, greaterThan),
	}
}

func LessThanNumber(position Position, lessThan int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return n[position] < lessThan
		},
		Description: fmt.Sprintf("%v < %v", position, lessThan),
	}
}

func Even(position Position) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return n[position]%2 == 0
		},
		Description: fmt.Sprintf("%v is even", position),
	}
}

func Odd(position Position) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return n[position]%2 == 1
		},
		Description: fmt.Sprintf("%v is odd", position),
	}
}

func NumberAppearsTimes(numberToCount, appearsTimes int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			count := 0
			for position := 0; position < len(n); position++ {
				if n[position] == numberToCount {
					count += 1
				}
			}

			return count == appearsTimes
		},
		Description: fmt.Sprintf("%v appears %v times", numberToCount, appearsTimes),
	}
}

func PositionLessThanPosition(position Position, isLessThanPosition Position) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return n[position] < n[isLessThanPosition]
		},
		Description: fmt.Sprintf("%v < %v", position, isLessThanPosition),
	}
}

func PositionGreaterThanPosition(position Position, isGreaterThanPosition Position) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return n[position] > n[isGreaterThanPosition]
		},
		Description: fmt.Sprintf("%v > %v", position, isGreaterThanPosition),
	}
}

func PositionEqualsPosition(position Position, isEqualToPosition Position) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return n[position] == n[isEqualToPosition]
		},
		Description: fmt.Sprintf("%v = %v", position, isEqualToPosition),
	}
}

func PositionIsSmallest(position Position) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			positionValue := n[position]
			for i := 0; i < len(n); i++ {
				if position != Position(i) && n[i] <= positionValue {
					return false
				}
			}

			return true
		},
		Description: fmt.Sprintf("%v is smallest", position),
	}
}

func PositionIsSmallestOrEqual(position Position) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			positionValue := n[position]
			for i := 0; i < len(n); i++ {
				if position != Position(i) && n[i] < positionValue {
					return false
				}
			}

			return true
		},
		Description: fmt.Sprintf("%v is smallest or equal", position),
	}
}

func PositionIsLargest(position Position) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			positionValue := n[position]
			for i := 0; i < len(n); i++ {
				if position != Position(i) && n[i] >= positionValue {
					return false
				}
			}

			return true
		},
		Description: fmt.Sprintf("%v is largest", position),
	}
}

func PositionIsLargestOrEqual(position Position) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			positionValue := n[position]
			for i := 0; i < len(n); i++ {
				if position != Position(i) && n[i] > positionValue {
					return false
				}
			}

			return true
		},
		Description: fmt.Sprintf("%v is largest or equal", position),
	}
}

func MoreEvens() *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			if isEven(len(n)) {
				panic("MoreEvens requires an odd number of numbers")
			}

			return numberOfEvens(n...) > len(n)/2
		},
		Description: "more evens",
	}
}

func MoreOdds() *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			if isEven(len(n)) {
				panic("MoreOdds requires an odd number of numbers")
			}

			return numberOfEvens(n...) <= len(n)/2
		},
		Description: "more odds",
	}
}

func NumberOfEvens(nEvensToCheckFor int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return numberOfEvens(n...) == nEvensToCheckFor
		},
		Description: fmt.Sprintf("number of evens = %v", nEvensToCheckFor),
	}
}

func SummationIsEven() *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return isEven(sum(n...))
		},
		Description: "sum is even",
	}
}

func SummationIsOdd() *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return !isEven(sum(n...))
		},
		Description: "sum is odd",
	}
}

func SumOfTwoPositionsLessThanNumber(position1, position2 Position, lessThan int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return n[position1]+n[position2] < lessThan
		},
		Description: fmt.Sprintf("%v + %v < %v", position1, position2, lessThan),
	}
}

func SumOfTwoPositionsEqualsNumber(position1, position2 Position, equals int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return n[position1]+n[position2] == equals
		},
		Description: fmt.Sprintf("%v + %v = %v", position1, position2, equals),
	}
}

func SumOfTwoPositionsGreaterThanNumber(position1, position2 Position, greaterThan int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return n[position1]+n[position2] > greaterThan
		},
		Description: fmt.Sprintf("%v + %v > %v", position1, position2, greaterThan),
	}
}

func RepeatsTimes(repeatsTimes int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			counts := make(map[int]int)
			for _, number := range n {
				counts[number] += 1
			}

			max := 0
			for _, count := range counts {
				if count > max {
					max = count
				}
			}

			return max == repeatsTimes
		},
		Description: fmt.Sprintf("any number repeats %v times", repeatsTimes),
	}
}

func NoPairs() *Verifier {
	pair := RepeatsTimes(2).Verify
	return &Verifier{
		Verify: func(n ...int) bool {
			return !pair(n...)
		},
		Description: "no pairs",
	}
}

func Ascending() *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			for i := 0; i < len(n)-1; i++ {
				if n[i] >= n[i+1] {
					return false
				}
			}

			return true
		},
		Description: "ascending",
	}
}

func Descending() *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			for i := 0; i < len(n)-1; i++ {
				if n[i] <= n[i+1] {
					return false
				}
			}

			return true
		},
		Description: "descending",
	}
}

func NoOrder() *Verifier {
	asc := Ascending().Verify
	dsc := Descending().Verify
	return &Verifier{
		Verify: func(n ...int) bool {
			return !asc(n...) && !dsc(n...)
		},
		Description: "no order",
	}
}

func SummationLessThanNumber(lessThan int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return sum(n...) < lessThan
		},
		Description: fmt.Sprintf("sum < %v", lessThan),
	}
}

func SummationEqualsNumber(equals int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return sum(n...) == equals
		},
		Description: fmt.Sprintf("sum = %v", equals),
	}
}

func SummationGreaterThanNumber(greaterThan int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return sum(n...) > greaterThan
		},
		Description: fmt.Sprintf("sum > %v", greaterThan),
	}
}

// AscendingSequenceOfSize use size of 1 for no sequence
func AscendingSequenceOfSize(size int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return sizeOfAscendingSequence(n...) == size
		},
		Description: fmt.Sprintf("%v numbers in ascending sequence", size),
	}
}

// DescendingSequenceOfSize use size of 1 for no sequence
func DescendingSequenceOfSize(size int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return sizeOfDescendingSequence(n...) == size
		},
		Description: fmt.Sprintf("%v numbers in descending sequence", size),
	}
}

// AscendingOrDecendingSequenceOfSize use size of 1 for no sequence
func AscendingOrDecendingSequenceOfSize(size int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			sizeAscending := sizeOfAscendingSequence(n...)
			sizeDescending := sizeOfDescendingSequence(n...)
			largerSize := sizeAscending
			if sizeDescending > sizeAscending {
				largerSize = sizeDescending
			}

			return largerSize == size
		},
		Description: fmt.Sprintf("%v numbers in ascending or descending sequence", size),
	}
}

func SummationIsMultipleOf(multipleOf int) *Verifier {
	return &Verifier{
		Verify: func(n ...int) bool {
			return sum(n...)%multipleOf == 0
		},
		Description: fmt.Sprintf("sum is multiple of %v", multipleOf),
	}
}

func sizeOfAscendingSequence(n ...int) int {
	sequenceSize := 1
	for i := 0; i < len(n)-1; i++ {
		if n[i]+1 == n[i+1] {
			sequenceSize += 1
		}
	}

	return sequenceSize
}

func sizeOfDescendingSequence(n ...int) int {
	sequenceSize := 1
	for i := 0; i < len(n)-1; i++ {
		if n[i]-1 == n[i+1] {
			sequenceSize += 1
		}
	}

	return sequenceSize
}

func sum(numbers ...int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}

	return sum
}

func isEven(n int) bool {
	return n%2 == 0
}

func numberOfEvens(numbers ...int) int {
	nEvens := 0
	for _, n := range numbers {
		if isEven(n) {
			nEvens += 1
		}
	}

	return nEvens
}
