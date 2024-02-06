package verifiers

import "fmt"

type Position int

const blueColor = "\033[34m"
const yellowColor = "\033[33m"
const purpleColor = "\033[35m"
const resetColor = "\033[0m"

const (
	Blue   = Position(0)
	Yellow = Position(1)
	Purple = Position(2)
)

func (p Position) String() string {
	switch p {
	case 0:
		return blueColor + "▲" + resetColor
	case 1:
		return yellowColor + "■" + resetColor
	case 2:
		return purpleColor + "●" + resetColor
	default:
		return fmt.Sprintf("pos %v", int(p))
	}
}
