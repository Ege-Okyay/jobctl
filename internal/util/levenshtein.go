package util

import "unicode/utf8"

// Levenshtein calculates the Levenshtein distance between two strings.
// The distance is the number of single-character edits (insertions, deletions, or substitutions)
// required to change one string into the other.
func Levenshtein(str1, str2 []rune) int {
	str1len := utf8.RuneCountInString(string(str1))
	str2len := utf8.RuneCountInString(string(str2))

	column := make([]int, str1len+1)

	for y := 1; y <= str1len; y++ {
		column[y] = y
	}

	for x := 1; x <= str2len; x++ {
		column[0] = x
		lastKey := x - 1

		for y := 1; y <= str1len; y++ {
			oldKey := column[y]

			incr := 0
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = min(column[y]+1, column[y-1]+1, lastKey+incr)
			lastKey = oldKey
		}
	}
	return column[str1len]
}
