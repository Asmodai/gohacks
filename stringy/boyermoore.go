// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// boyermoore.go --- Boyer/Moore string match algorithm.
//
// Copyright (c) 2025-2026 Paul Ward <paul@lisphacker.uk>
//
// Author:     Paul Ward <paul@lisphacker.uk>
// Maintainer: Paul Ward <paul@lisphacker.uk>
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * Comments:

// * Package:

package stringy

// * Code:

// Returns the byte index of the first occurrence of needle in haystack, or -1.
func BMH(needle, haystack []byte) int {
	nlen := len(needle)
	hlen := len(haystack)

	if nlen == 0 {
		return 0
	}

	if hlen < nlen {
		return -1
	}

	// Build bad-char shift table for bytes.
	var shift [256]int

	// default shift
	for idx := range shift {
		shift[idx] = nlen
	}

	// last char keeps default
	for idx := range nlen - 1 {
		shift[needle[idx]] = nlen - 1 - idx
	}

	idx := 0
	for idx <= hlen-nlen {
		// Compare right-to-left
		sub := nlen - 1

		for sub >= 0 && haystack[idx+sub] == needle[sub] {
			sub--
		}

		if sub < 0 {
			return idx
		}

		idx += shift[haystack[idx+nlen-1]]
	}

	return -1
}

func BMHRunes(needle, haystack []rune) int {
	nlen := len(needle)
	hlen := len(haystack)

	if nlen == 0 {
		return 0
	}

	if hlen < nlen {
		return -1
	}

	// Sparse table: last index of rune (excluding last char).
	last := make(map[rune]int, nlen)

	for idx := range nlen - 1 {
		last[needle[idx]] = idx
	}

	idx := 0
	for idx <= hlen-nlen {
		sub := nlen - 1

		for sub >= 0 && haystack[idx+sub] == needle[sub] {
			sub--
		}

		if sub < 0 {
			return idx
		}

		chr := haystack[idx+nlen-1]

		if pos, ok := last[chr]; ok {
			idx += nlen - 1 - pos
		} else {
			idx += nlen
		}
	}

	return -1
}

// Boyer–Moore search (bytes):

// -

// Returns the byte index of the first occurrence of pattern in text, or -1.
func IndexBM(pattern, text []byte) int {
	pLen, tLen := len(pattern), len(text)
	if pLen == 0 {
		return 0
	}

	if tLen < pLen {
		return -1
	}

	// Preprocess
	lastPos := buildBadCharTable(pattern)
	suffix, prefix := buildGoodSuffixTables(pattern)

	// Main search
	var idx int // alignment of pattern against text

	for idx <= tLen-pLen {
		pos := pLen - 1
		for pos >= 0 && pattern[pos] == text[idx+pos] {
			pos--
		}

		if pos < 0 {
			return idx
		}

		// Bad-character shift
		badChar := int(text[idx+pos])
		badShift := pos - lastPos[badChar]

		if badShift < 1 {
			badShift = 1
		}

		// Good-suffix shift
		goodShift := 0

		if pos < pLen-1 {
			goodShift = moveByGoodSuffix(pos, pLen, suffix, prefix)
		}

		if badShift > goodShift {
			idx += badShift
		} else {
			idx += goodShift
		}
	}

	return -1
}

// FindAllBM returns all (overlapping) byte indices where pattern occurs.
func FindAllBM(pattern, text []byte) []int {
	pLen, tLen := len(pattern), len(text)

	if pLen == 0 || tLen < pLen {
		return nil
	}

	lastPos := buildBadCharTable(pattern)
	suffix, prefix := buildGoodSuffixTables(pattern)

	var (
		hits []int
		idx  int
	)

	for idx <= tLen-pLen {
		pos := pLen - 1
		for pos >= 0 && pattern[pos] == text[idx+pos] {
			pos--
		}

		if pos < 0 {
			hits = append(hits, idx)
			// For overlapping matches, shift by 1; for
			// non-overlapping, shift by pLen.
			// Choose one; here we do overlapping:
			idx++

			continue
		}

		badChar := int(text[idx+pos])
		badShift := pos - lastPos[badChar]
		goodShift := 0

		if badShift < 1 {
			badShift = 1
		}

		if pos < pLen-1 {
			goodShift = moveByGoodSuffix(pos, pLen, suffix, prefix)
		}

		if badShift > goodShift {
			idx += badShift
		} else {
			idx += goodShift
		}
	}

	return hits
}

// buildBadCharTable returns the last occurrence table for bytes (size 256).
// lastPos[b] = greatest index i where pattern[i] == b, or -1 if b never appears.
func buildBadCharTable(pattern []byte) [256]int {
	var lastPos [256]int

	for idx := range 256 {
		lastPos[idx] = -1
	}

	for idx := range pattern {
		lastPos[int(pattern[idx])] = idx
	}

	return lastPos
}

// buildGoodSuffixTables computes the "suffix" and "prefix" tables used for
// the good-suffix rule.
//
// suffix[k] = start index of a substring in pattern that matches the suffix
// of length k of pattern, and that substring's rightmost end is before
// pattern's last char. -1 if none.
//
// prefix[k] = true if a prefix of pattern of length k matches a suffix of pattern.
func buildGoodSuffixTables(pattern []byte) ([]int, []bool) {
	pLen := len(pattern)
	suffix := make([]int, pLen)
	prefix := make([]bool, pLen)

	for idx := range pLen {
		suffix[idx] = -1
		prefix[idx] = false
	}

	for idx := range pLen - 1 { // idx: last index of the "front" substring
		pos := idx
		sufflen := 0 // length of matched suffix

		for pos >= 0 && pattern[pos] == pattern[pLen-1-sufflen] {
			pos--
			sufflen++

			// record start index of this matched substring
			suffix[sufflen] = pos + 1
		}

		if pos == -1 {
			// matched all the way to the beginning: also a
			// prefix.
			prefix[sufflen] = true
		}
	}

	return suffix, prefix
}

// moveByGoodSuffix computes how far to shift when mismatch occurs at
// position "mismatchIdx".
//
// mismatchIdx is the index in pattern where comparison failed (0..pLen-1).
func moveByGoodSuffix(mismatchIdx, pLen int, suffix []int, prefix []bool) int {
	// length of good suffix that matched: from mismatchIdx+1 to end
	goodLen := pLen - 1 - mismatchIdx
	if goodLen <= 0 {
		return 0
	}

	// Case 1: there is another occurrence of the good suffix in
	// pattern (not at the end).
	if suffix[goodLen] != -1 {
		return mismatchIdx - suffix[goodLen] + 1
	}

	// Case 2: find the longest suffix of the good suffix that is also a
	// prefix.
	for res := mismatchIdx + 2; res <= pLen-1; res++ { //nolint:mnd
		if prefix[pLen-res] {
			return res
		}
	}

	// Case 3: no match, shift by full pattern length
	return pLen
}

// * boyermoore.go ends here.
