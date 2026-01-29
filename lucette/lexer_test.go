// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// lexer_test.go --- Lexer tests.
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

package lucette

// * Imports:

import (
	"bufio"
	"hash/fnv"
	"io"
	"math/rand"
	"strings"
	"testing"

	"gitlab.com/tozd/go/errors"
	"golang.org/x/text/unicode/norm"
)

// * Constants:

// * Variables:

// * Code:

// ** Functions:

func makeInvalidUTF8Corpus() [][]byte {
	return [][]byte{
		{0xC0, 0xAF},             // overlong '/'
		{0xE0, 0x80, 0xAF},       // overlong
		{0xF0, 0x80, 0x80, 0xAF}, // overlong
		{0x80},                   // lone continuation
		{0xC2},                   // truncated 2-byte
		{0xE1, 0x80},             // truncated 3-byte
		{0xF1, 0x80, 0x80},       // truncated 4-byte
	}
}

// Add combining marks after random runes.
func addZalgo(base string, rng *rand.Rand) string {
	var b strings.Builder
	for _, r := range base {
		b.WriteRune(r)
		if r == '\n' { // avoid attaching to newlines.
			continue
		}
		if rng.Intn(3) == 0 { // ~33% chance to zalgo this rune.
			n := rng.Intn(8) // up to 7 combining marks.
			for i := 0; i < n; i++ {
				b.WriteRune(randCombining(rng))
			}
		}
	}
	return b.String()
}

func randCombining(rng *rand.Rand) rune {
	ranges := [...]struct{ lo, hi rune }{
		{0x0300, 0x036F}, // Combining Diacritical Marks
		{0x1AB0, 0x1AFF}, // Combining Diacritical Marks Extended
		{0x1DC0, 0x1DFF}, // Combining Diacritical Marks Supplement
		{0x20D0, 0x20FF}, // Combining Diacritical Marks for Symbols
		{0xFE20, 0xFE2F}, // Combining Half Marks
	}
	r := ranges[rng.Intn(len(ranges))] // pick a single range
	span := int(r.hi - r.lo + 1)       // inclusive size
	if span <= 0 {                     // paranoia guard
		return r.lo
	}
	return r.lo + rune(rng.Intn(span))
}

// Randomly inject bidi control chars and zero-widths.
func injectInvisibles(s string, rng *rand.Rand) string {
	invis := []rune{
		0x200B, // ZWSP
		0x200C, // ZWNJ
		0x200D, // ZWJ
		// Bidi controls (a subset; enough to catch bugs)
		0x202A, // LRE
		0x202B, // RLE
		0x202D, // LRO
		0x202E, // RLO
		0x202C, // PDF
		0x2066, // LRI
		0x2067, // RLI
		0x2068, // FSI
		0x2069, // PDI
	}
	var out []rune
	for _, r := range s {
		out = append(out, r)
		if rng.Intn(10) == 0 { // 10% chance to inject after rune
			out = append(out, invis[rng.Intn(len(invis))])
		}
	}
	return string(out)
}

// ** Tests:

// *** Reader:

func TestLexer_Reader(t *testing.T) {
	t.Run("Reader", func(t *testing.T) {
		input := "abcdef!hijklm"
		lx := &lexer{reader: bufio.NewReader(strings.NewReader(input))}

		t.Run("readRune", func(t *testing.T) {
			res, err := lx.readRune()
			if err != nil {
				t.Fatalf("Unexpected error: %#v", err)
			}

			if res != 'a' {
				t.Errorf("Unexpected rune: %#v", res)
			}
		})

		t.Run("unreadRune", func(t *testing.T) {
			t.Run("Works", func(t *testing.T) {
				if err := lx.unreadRune(); err != nil {
					t.Fatalf("Unexpected error: %#v", err)
				}
			})

			t.Run("Errors on double unread", func(t *testing.T) {
				err := lx.unreadRune()

				if err == nil {
					t.Fatal("Expected an error")
				}

				if !errors.Is(err, ErrDoubleUnread) {
					t.Fatalf("Unexpected error: %#v", err)
				}
			})

			t.Run("unread before any read errors", func(t *testing.T) {
				lx := &lexer{reader: bufio.NewReader(strings.NewReader("x"))}
				if err := lx.unreadRune(); err == nil {
					t.Fatal("expected error")
				}
			})

			t.Run("read, read, unread gives back last rune", func(t *testing.T) {
				lx := &lexer{reader: bufio.NewReader(strings.NewReader("ab"))}

				r1, _ := lx.readRune()
				r2, _ := lx.readRune()

				if r1 != 'a' || r2 != 'b' {
					t.Fatal("setup")
				}

				if err := lx.unreadRune(); err != nil {
					t.Fatal(err)
				}

				r3, _ := lx.readRune()
				if r3 != 'b' {
					t.Fatalf("expected 'b', got %q", r3)
				}
			})
		})

		t.Run("peekRune", func(t *testing.T) {
			t.Run("Works", func(t *testing.T) {
				res, err := lx.peekRune()
				if err != nil {
					t.Fatalf("Unexpected error: %#v", err)
				}

				// We unread, so we're back at the start.
				if res != 'a' {
					t.Errorf("Unexpected rune: %c", res)
				}
			})

			t.Run("Does not advance", func(t *testing.T) {
				res, err := lx.peekRune()
				if err != nil {
					t.Fatalf("Unexpected error: %#v", err)
				}

				// Should still be at start.
				if res != 'a' {
					t.Errorf("Unexpected rune: %c", res)
				}
			})

			t.Run("peek then read returns same rune", func(t *testing.T) {
				lx := &lexer{reader: bufio.NewReader(strings.NewReader("xyz"))}

				p, err := lx.peekRune()
				if err != nil {
					t.Fatal(err)
				}

				r, err := lx.readRune()
				if err != nil {
					t.Fatal(err)
				}

				if p != r {
					t.Fatalf("peek %q != read %q", p, r)
				}
			})
		})

		t.Run("Errors at end", func(t *testing.T) {
			var err error

			for err == nil {
				_, err = lx.readRune()

				if err != nil {
					if !errors.Is(err, io.EOF) {
						t.Fatalf("Unexpected error %#v",
							err)
					}
				}
			}
		})

		t.Run("Peek errors at end", func(t *testing.T) {
			_, err := lx.peekRune()

			if err == nil {
				t.Fatal("Expected an error")
			}

			if !errors.Is(err, io.EOF) {
				t.Fatalf("Unexpected error: %#v", err)
			}
		})

		t.Run("readWhile", func(t *testing.T) {
			lx := &lexer{reader: bufio.NewReader(strings.NewReader(input))}

			t.Run("Works", func(t *testing.T) {
				stopAt := '!'
				want := "abcdef"

				res, err := lx.readWhile(func(r rune) bool {
					return r != stopAt
				})

				if err != nil {
					t.Fatalf("Unexpected error: %#v", err)
				}

				if res != want {
					t.Errorf("Mismatch.  %q != %q",
						res,
						want)
				}
			})

			t.Run("readWhile leaves stop rune unread", func(t *testing.T) {
				res, err := lx.readRune()
				if err != nil {
					t.Fatal(err)
				}

				if res != '!' {
					t.Fatalf("expected stop rune '!', got %q",
						res)
				}
			})

			t.Run("Error at EOF", func(t *testing.T) {
				stopAt := '#' // Will never be found.

				_, err := lx.readWhile(func(r rune) bool {
					return r != stopAt
				})

				if err == nil {
					t.Fatal("Expected an error")
				}

				if !errors.Is(err, io.EOF) {
					t.Errorf("Unexpected error: %#v", err)
				}
			})

			t.Run(
				"readWhile immediate stop returns empty and does not consume",
				func(t *testing.T) {
					lx := &lexer{reader: bufio.NewReader(strings.NewReader("!boom"))}

					s, err := lx.readWhile(func(r rune) bool {
						return r != '!'
					})

					if err != nil {
						t.Fatal(err)
					}

					if s != "" {
						t.Fatalf("got %q", s)
					}

					r, _ := lx.readRune()
					if r != '!' {
						t.Fatalf("expected '!', got %q", r)
					}
				})
		})

		t.Run("unicode rune and newline positions", func(t *testing.T) {
			lx := &lexer{reader: bufio.NewReader(strings.NewReader("α\nβ"))}

			r, _ := lx.readRune()
			if r != 'α' {
				t.Fatal("expected alpha")
			}

			if lx.currPos.Column != 1 {
				t.Fatal("Not at column 1")
			}

			r, _ = lx.readRune()
			if r != '\n' {
				t.Fatal("expected newline")
			}

			if err := lx.unreadRune(); err != nil {
				t.Fatal(err)
			}

			r, _ = lx.readRune()
			if r != '\n' {
				t.Fatal("expected newline after unread")
			}
		})

		t.Run("comparators lookahead", func(t *testing.T) {
			lx := &lexer{reader: bufio.NewReader(strings.NewReader("<="))}

			r1, _ := lx.readRune()
			if r1 != '<' {
				t.Fatal()
			}

			p, _ := lx.peekRune()
			if p != '=' {
				t.Fatal("expected to peek '='")
			}

			r2, _ := lx.readRune()
			if r2 != '=' {
				t.Fatal("expected to read '='")
			}
		})
	})
}

// ** Fuzzers:

func fuzzums(t *testing.T, s string, rng *rand.Rand) {
	lx := &lexer{reader: bufio.NewReader(strings.NewReader(s))}
	runes := []rune(s)

	i := 0
	lastReadOK := false
	unreadBuffered := false

	// Limit the chaos so test time stays reasonable.
	steps := len(runes)*3 + 16
	if steps < 32 {
		steps = 32
	}

	for step := 0; step < steps; step++ {
		switch rng.Intn(3) {
		case 0: // peek
			r, err := lx.peekRune()
			if i >= len(runes) {
				if !errors.Is(err, io.EOF) {
					t.Fatalf("peek EOF: want io.EOF, got %v", err)
				}
				lastReadOK = false
				continue
			}
			if err != nil {
				t.Fatalf("peek err: %v", err)
			}
			if r != runes[i] {
				t.Fatalf("peek mismatch: got %q want %q at i=%d", r, runes[i], i)
			}
			// idempotent peek
			r2, err2 := lx.peekRune()
			if err2 != nil || r2 != r {
				t.Fatalf("second peek mismatch: %q / %v", r2, err2)
			}
			// peek does not change model state
			lastReadOK = false

		case 1: // read
			r, err := lx.readRune()
			if i >= len(runes) {
				if !errors.Is(err, io.EOF) {
					t.Fatalf("read EOF: want io.EOF, got %v", err)
				}
				lastReadOK = false
				continue
			}
			if err != nil {
				t.Fatalf("read err: %v", err)
			}
			if r != runes[i] {
				t.Fatalf("read mismatch: got %q want %q at i=%d", r, runes[i], i)
			}
			i++
			lastReadOK = true
			unreadBuffered = false

		case 2: // unread
			err := lx.unreadRune()
			// Legal only if last op was a successful read and we haven't unread already.
			if !lastReadOK || unreadBuffered {
				if err == nil {
					t.Fatalf("unread should error (lastReadOK=%v unreadBuffered=%v)", lastReadOK, unreadBuffered)
				}
			} else {
				if err != nil {
					t.Fatalf("unread err: %v", err)
				}
				i--                   // model rewinds one rune
				lastReadOK = false    // you haven’t re-read it yet
				unreadBuffered = true // guard double-unread
				// immediate read must return same rune
				r, err2 := lx.readRune()
				if err2 != nil || r != runes[i] {
					t.Fatalf("post-unread read mismatch: %q/%v want %q at i=%d", r, err2, runes[i], i)
				}
				i++
				lastReadOK = true
				unreadBuffered = false
			}
		}
	}

	// Consume the rest; it must equal the suffix runes[i:].
	var tail []rune
	for {
		r, err := lx.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatalf("tail read err: %v", err)
		}
		tail = append(tail, r)
	}
	want := string(runes[i:])
	got := string(tail)
	if got != want {
		t.Fatalf("tail mismatch: got %q want %q (i=%d len=%d)", got, want, i, len(runes))
	}
}

func FuzzLexer_ReadUnreadPeek(f *testing.F) {
	seeds := []string{
		"",
		"abc",
		"a\nb",
		"line1\r\nline2",
		"a\tb",
		"αβγ",
		"😀x",
		"˙puǝıɹɟ pןo ʎɯ ƃuıɹʇs ʇxǝʇ oןןǝH",
		"⧦ⅇǁǁ☉ ╬ⅇ⨳╬ 𝕤╬ℾⅈℼ𓉛 ⩕ℽ ☉ǁⅆ ⨎ℾⅈⅇℼⅆ.",
		"🅷🄴🅻🄻🅾 🆃🄴🆇🅃 🅂🆃🅁🅸🄽🅶 🅼🅈 🄾🅻🄳 🄵🆁🄸🅴🄽🅳.",
		"ᕼᗕᖶᖶᗝ ᐪᗕ᙭ᐪ ᔑᐪᖇᓵᐱᘜ ᗑᖿ ᗝᖶᐅ ᖴᖇᓵᗕᐱᐅ.",
		"ꖾꗍꝆꝆꗞ ꖡꗍꘉꖡ ꕷꖡ𐝥ꕯꖦꗱ ꕮꔇ ꗞꝆꕒ ꘘ𐝥ꕯꗍꖦꕒ.",
		"𖦙𖠢ꛚꛚ𖥕 𖢧𖠢𖧦𖢧 𖨚𖢧𖦪𖥣ꛘꛪ 𖢑ꚲ 𖥕ꛚ𖦧 𖨨𖦪𖥣𖠢ꛘ𖦧.",
		strings.Repeat("Z", 128),
		"🇬🇧",      // regional-indicator pair
		"👩‍❤️‍👨",  // ZWJ sequence
		"✈", "✈️", // text vs emoji presentation
		"🏳️‍🌈", // rainbow flag with VS-16 + ZWJ
		"᚛ᚄᚓᚐᚋᚒᚄ ᚑᚄᚂᚑᚏᚅ᚜",
		"Powerلُلُصّبُلُلصّبُررً ॣ ॣh ॣ ॣ冗",
		"🏳0🌈️",
		"జ్ఞ‌ా",
		"ثم نفس سقطت وبالتحديد،, جزيرتي باستخدام أن دنو. إذ هنا؟ الستار وتنصيب كان. أهّل ايطاليا، بريطانيا-فرنسا قد أخذ. سليمان، إتفاقية بين ما, يذكر الحدود أي بعد, معاملة بولندا، الإطلاق عل إيو.",
		/* DO NOT EDIT THIS --> */ "test⁠test‫",
	}

	for _, s := range seeds {
		f.Add(s)
	}

	for _, b := range makeInvalidUTF8Corpus() {
		f.Add(string(b))
	}

	f.Fuzz(func(t *testing.T, s string) {
		// Deterministic RNG by hashing input
		h := fnv.New64a()
		_, _ = h.Write([]byte(s))
		rng := rand.New(rand.NewSource(int64(h.Sum64())))

		// Base variants
		variants := []string{
			s,
			addZalgo(s, rng),
			injectInvisibles(s, rng),
			injectInvisibles(addZalgo(s, rng), rng),
			norm.NFC.String(addZalgo(s, rng)),
			norm.NFD.String(addZalgo(s, rng)),
		}

		for _, v := range variants {
			fuzzums(t, v, rng)
		}
	})
}

// * lexer_test.go ends here.
