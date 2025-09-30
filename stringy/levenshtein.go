// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// levenshtein.go --- Levenshtein distance.
//
// Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
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

// Levenshtein returns the edit distance between a and b,
// counting insertions, deletions, and substitutions (all cost = 1).
// It operates on runes (Unicode code points), not bytes.
func Levenshtein(a, b string) int {
	src := []rune(a)
	dst := []rune(b)

	srcLen := len(src)
	dstLen := len(dst)

	if srcLen == 0 {
		return dstLen
	}

	if dstLen == 0 {
		return srcLen
	}

	// Ensure dst is the shorter dimension to keep memory small.
	if dstLen > srcLen {
		src, dst = dst, src
		srcLen, dstLen = dstLen, srcLen
	}

	rowPrev := make([]int, dstLen+1)
	rowCurr := make([]int, dstLen+1)

	for jIdx := 0; jIdx <= dstLen; jIdx++ {
		rowPrev[jIdx] = jIdx
	}

	for iIdx := 1; iIdx <= srcLen; iIdx++ {
		rowCurr[0] = iIdx
		sRune := src[iIdx-1]

		for jIdx := 1; jIdx <= dstLen; jIdx++ {
			dRune := dst[jIdx-1]

			cost := 0
			if sRune != dRune {
				cost = 1
			}

			insCost := rowCurr[jIdx-1] + 1    // insertion
			delCost := rowPrev[jIdx] + 1      // deletion
			subCost := rowPrev[jIdx-1] + cost // substitution

			minCost := insCost
			if delCost < minCost {
				minCost = delCost
			}

			if subCost < minCost {
				minCost = subCost
			}

			rowCurr[jIdx] = minCost
		}

		// swap rows
		rowPrev, rowCurr = rowCurr, rowPrev
	}

	return rowPrev[dstLen]
}

// * levenshtein.go ends here.
