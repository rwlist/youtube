package lists

import "strings"

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const baseLength = 5

func last(s string) uint8 {
	return s[len(s)-1]
}

func validXord(xord string) bool {
	if len(xord) == 0 {
		return false
	}
	if last(xord) == alphabet[0] {
		return false
	}
	return true
}

func splitXord(a, b string) string {
	if a == "" && b == "" {
		return strings.Repeat(string(alphabet[len(alphabet)/2]), baseLength)
	}

	if a == "" {
		return string(subXord(b))
	}

	if b == "" {
		var vals []uint8
		growx8(&vals, baseLength, len(a)+1, last(alphabet))
		b = string(vals)
	}

	if a >= b {
		panic("a >= b")
	}

	vals := []uint8(a)

	critical := 0
	for i := 0; i < len(vals); i++ {
		if vals[i] == last(alphabet) {
			critical++
		} else {
			break
		}
	}

	var c []uint8
	if critical != len(vals) {
		c = up1(vals)
		if string(c) >= b {
			c = []uint8(a)
		}
	} else {
		c = []uint8(a)
	}
	newMinLength := 2 * critical
	if newMinLength < baseLength {
		newMinLength = baseLength
	}
	needGrow := 0
	if string(c) == a {
		needGrow = 1
	}
	growx8(&c, needGrow, newMinLength, alphabet[0])

	if !(string(c) > a) {
		panic("up1 or growx8 is broken")
	}

	common := len(c)
	if common > len(b) {
		common = len(b)
	}

	if string(c[:common]) > b[:common] {
		panic("c is bigger, shouldn't be")
	}

	// short path if appending is not needed
	if critical == 0 && string(c) < b && validXord(string(c)) {
		return string(c)
	}

	growx8(&c, 0, len(b), alphabet[0])
	c = append(c, alphabet[1])

	res := string(c)
	if res <= a || res >= b {
		panic("res is not in (a, b)")
	}
	if !validXord(res) {
		panic("res is not valid")
	}

	return res
}

func up1(vals []uint8) []uint8 {
	for i := len(vals) - 1; i >= 0; i-- {
		if vals[i] != last(alphabet) {
			newVals := make([]uint8, i+1)
			copy(newVals, vals)
			newVals[i]++
			return newVals
		}
	}

	panic("cannot up1")
}

func subXord(xord string) []uint8 {
	if !validXord(xord) {
		panic("invalid xord: " + xord)
	}

	vals := []uint8(xord)
	vals[len(vals)-1]--

	critical := 0
	for i := 0; i < len(vals); i++ {
		if vals[i] == alphabet[0] {
			critical++
		} else {
			break
		}
	}

	if vals[len(vals)-1] == alphabet[0] {
		if critical == 0 {
			for vals[len(vals)-1] == alphabet[0] {
				vals = vals[:len(vals)-1]
			}
			vals[len(vals)-1]--
			growx8(&vals, 1, baseLength, last(alphabet))
		} else {
			newLength := 2 * critical
			if newLength < baseLength {
				newLength = baseLength
			}
			growx8(&vals, 1, baseLength, last(alphabet))
		}
	}
	return vals
}

func growx8(x8 *[]uint8, min int, target int, char uint8) {
	for len(*x8) < target || min > 0 {
		*x8 = append(*x8, char)
		min--
	}
}
