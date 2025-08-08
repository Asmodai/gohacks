// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// metadata_test.go --- Metadata tests.
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

package metadata

// * Imports:

import (
	"errors"
	"reflect"
	"testing"
)

// * Code:

// ** Helpers:

type hasClone interface {
	Clone() Metadata
}
type hasTags interface {
	Tags() []string
}
type hasMerge interface {
	Merge(Metadata, bool) Metadata
}

// ** Tests:

func TestReservedKeyRejected(t *testing.T) {
	m := NewMetadata()

	// "version" is reserved; Set should reject it.
	if err := m.Set("version", "1.2.3"); err == nil {
		t.Fatalf("expected Set to reject reserved key 'version'")
	} else if !errors.Is(err, ErrKeyIsReserved) {
		t.Fatalf("expected ErrKeyIsReserved, got %v", err)
	}
}

func TestAmbiguousKeyRejected(t *testing.T) {
	m := NewMetadata()

	// "Version" is a case-insensitive collision with "version".
	if err := m.Set("Version", "1.2.3"); err == nil {
		t.Fatalf("expected Set to reject ambiguous key 'Version'")
	} else if !errors.Is(err, ErrKeyIsAmbiguous) {
		t.Fatalf("expected ErrKeyIsAmbiguous, got %v", err)
	}
}

func TestVisibilityAcceptedAndIgnored(t *testing.T) {
	m := NewMetadata()

	// Valid (case-insensitive) should be stored lowercased.
	m.SetVisibility("PUBLIC")
	if got := m.GetVisibility(); got != "public" {
		t.Fatalf("expected 'public', got %q", got)
	}

	// Invalid should be ignored (no change).
	m.SetVisibility("definitely-not-a-level")
	if got := m.GetVisibility(); got != "public" {
		t.Fatalf("invalid visibility should not change value; got %q", got)
	}
}

func TestTagsSplitAndTrim(t *testing.T) {
	m := NewMetadata()

	m.SetTags("alpha, beta ,Gamma")
	tags := m.(hasTags).Tags()

	want := []string{"alpha", "beta", "Gamma"}
	if !reflect.DeepEqual(tags, want) {
		t.Fatalf("tags mismatch: want %v, got %v", want, tags)
	}
}

func TestCloneCopyOnWriteIsolation(t *testing.T) {
	m := NewMetadata()
	m.Set("foo", "bar") // custom key (allowed)
	m.SetAuthor("Ada Lovelace")

	clone := m.(hasClone).Clone()

	// Mutate original
	_ = m.Set("foo", "baz")
	m.SetAuthor("Grace Hopper")

	// Clone should be unaffected
	val, _ := clone.Get("foo")
	if val != "bar" {
		t.Fatalf("clone should retain old custom key value; got %q", val)
	}

	if got := clone.GetAuthor(); got != "Ada Lovelace" {
		t.Fatalf("clone should retain old author; got %q", got)
	}
}

func TestListReturnsCopy(t *testing.T) {
	m := NewMetadata()
	_ = m.Set("foo", "bar")

	// Grab a copy and mutate it
	cp := m.List()
	cp["foo"] = "whoops"

	// Original should remain unchanged
	val, _ := m.Get("foo")
	if val != "bar" {
		t.Fatalf("List must return a copy; expected 'bar', got %q", val)
	}
}

func TestGetUnknownKey(t *testing.T) {
	m := NewMetadata()

	if v, ok := m.Get("nope"); ok || v != "" {
		t.Fatalf("Get of unknown key should be empty/false; got %q, %v", v, ok)
	}
}

func TestSettersAndGettersRoundTrip(t *testing.T) {
	m := NewMetadata()

	m.SetDoc("docstring")
	m.SetSince("2025-08-07")
	m.SetVersion("1.0.0")
	m.SetDeprecated("use other")
	m.SetProtocol("fs.readable")
	m.SetVisibility("private")
	m.SetExample("foo(bar)")
	m.SetTags("x,y")
	m.SetAuthor("Paul")

	if m.GetDoc() != "docstring" ||
		m.GetSince() != "2025-08-07" ||
		m.GetVersion() != "1.0.0" ||
		m.GetDeprecated() != "use other" ||
		m.GetProtocol() != "fs.readable" ||
		m.GetVisibility() != "private" ||
		m.GetExample() != "foo(bar)" ||
		m.GetTags() != "x"+TagDelimiter+"y" ||
		m.GetAuthor() != "Paul" {
		t.Fatalf("round-trip setter/getter mismatch: %+v", m.List())
	}
}

func TestMerge_NoOverwrite(t *testing.T) {
	dst := NewMetadata()
	dst.SetDoc("old-doc")
	_ = dst.Set("foo", "x") // custom key

	src := NewMetadata()
	src.SetDoc("new-doc")
	_ = src.Set("foo", "y")
	_ = src.Set("bar", "z")

	// Merge without overwrite: existing keys must be preserved, new added.
	dst = dst.(hasMerge).Merge(src, false)

	if dst.GetDoc() != "old-doc" {
		t.Fatalf("doc should remain 'old-doc' without overwrite, got %q", dst.GetDoc())
	}

	if v, _ := dst.Get("foo"); v != "x" {
		t.Fatalf("custom 'foo' should remain 'x' without overwrite, got %q", v)
	}

	if v, _ := dst.Get("bar"); v != "z" {
		t.Fatalf("new key 'bar' should be added with value 'z', got %q", v)
	}
}

func TestMerge_WithOverwrite(t *testing.T) {
	dst := NewMetadata()
	dst.SetDoc("old-doc")
	_ = dst.Set("foo", "x")

	src := NewMetadata()
	src.SetDoc("new-doc")
	_ = src.Set("foo", "y")

	// Overwrite: only custom should update.
	dst = dst.(hasMerge).Merge(src, true)

	if dst.GetDoc() != "old-doc" {
		t.Fatalf("doc should not be overwritten, got %q", dst.GetDoc())
	}

	if v, _ := dst.Get("foo"); v != "y" {
		t.Fatalf("'foo' should be overwritten to 'y', got %q", v)
	}
}

// * metadata_test.go ends here.
