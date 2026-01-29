// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// metadata.go --- Metadata.
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

// TODO: Other possible metadata keys include:
//
//  * `readonly` -- should not mutate.
//  * `experimental` -- Can break things.
//  * `async` -- uses Goroutines.
//  * `transactional` -- can make use of rollback/commit mechanics.

// * Package:

package metadata

// * Imports:

import (
	"slices"
	"strings"
	"sync"

	"github.com/Asmodai/gohacks/generics"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	TagDelimiter = ","

	KeyDoc        = "doc"
	KeySince      = "since"
	KeyVersion    = "version"
	KeyDeprecated = "deprecated"
	KeyProtocol   = "protocol"
	KeyVisibility = "visibility"
	KeyExample    = "example"
	KeyTags       = "tags"
	KeyAuthor     = "author"
)

// * Variables:

var (
	ErrKeyIsInvalid   = errors.Base("metadata key is invalid")
	ErrKeyIsAmbiguous = errors.Base("metadata key is ambiguous")
	ErrKeyIsReserved  = errors.Base("reserved metadata key")

	// Reserved selector metadata keys and their intended semantics.
	//
	// These keys are reserved for internal use and documentation
	// purposes.
	//
	// While not enforced, they are expected to follow consistent
	// formatting and be used by tooling for introspection, generation,
	// and validation.
	//
	// - "doc": Short description of what the selector does. (Markdown
	//          allowed.)
	// - "since": First version or date this selector was introduced.
	//            Format: "v1.2.3" or ISO date (e.g. "2025-08-07").
	// - "version": Current version of the selector logic.
	//              Useful if selector behavior has changed over time.
	// - "deprecated": Optional deprecation notice or replacement advice.
	// - "protocol": Protocol(s) this selector belongs to
	//               (comma-separated).
	// - "visibility": One of: "public", "internal", "private".
	//                 Useful for generating user-facing docs or
	//                 restricting UI tools.
	// - "example": A usage example (inline or structured Markdown).
	// - "tags": Comma-separated labels for filtering or grouping.
	//           E.g. "filesystem,experimental,fastpath"
	// - "author": Who wrote or maintains the selector logic.
	//             Useful for blame or kudos.
	//
	// Tools can recognize and use these for generating CLI docs, debug
	// dumps, live inspector UIs, etc.
	//
	//nolint:gochecknoglobals
	ReservedMetadataKeys = map[string]struct{}{
		KeyDoc:        {},
		KeySince:      {},
		KeyVersion:    {},
		KeyDeprecated: {},
		KeyProtocol:   {},
		KeyVisibility: {},
		KeyExample:    {},
		KeyTags:       {},
		KeyAuthor:     {},
	}

	//nolint:gochecknoglobals
	VisibilityLevels = map[string]struct{}{
		"public":   {},
		"internal": {},
		"private":  {},
	}
)

// * Code:

// ** Types:

type metadata struct {
	mu   sync.RWMutex
	data map[string]string
}

// ** Methods:

func (m *metadata) List() map[string]string {
	copyMap := make(map[string]string, len(m.data))

	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.data {
		copyMap[k] = v
	}

	return copyMap
}

// Internal-only, full-fat set without locking.
func (m *metadata) setWithoutLock(key string, value string) error {
	if isReservedKey(key) {
		return errors.WithMessagef(
			ErrKeyIsReserved,
			"%q",
			key)
	}

	if err := validateKey(key); err != nil {
		return errors.WithStack(err)
	}

	m.data[key] = value

	return nil
}

// Internal-only.  Bypasses sanity checks.
func (m *metadata) set(key string, value string) Metadata {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value

	return m
}

// Internal-only, doesn't check if field exists.
func (m *metadata) get(key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.data[key]
}

func (m *metadata) Set(key string, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.setWithoutLock(key, value)
}

func (m *metadata) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, found := m.data[key]

	return val, found
}

func (m *metadata) SetDoc(val string) Metadata        { return m.set(KeyDoc, val) }
func (m *metadata) SetSince(val string) Metadata      { return m.set(KeySince, val) }
func (m *metadata) SetVersion(val string) Metadata    { return m.set(KeyVersion, val) }
func (m *metadata) SetDeprecated(val string) Metadata { return m.set(KeyDeprecated, val) }
func (m *metadata) SetProtocol(val string) Metadata   { return m.set(KeyProtocol, val) }
func (m *metadata) SetExample(val string) Metadata    { return m.set(KeyExample, val) }
func (m *metadata) SetTags(val string) Metadata       { return m.set(KeyTags, val) }
func (m *metadata) SetAuthor(val string) Metadata     { return m.set(KeyAuthor, val) }

func (m *metadata) SetVisibility(val string) Metadata {
	level := strings.ToLower(strings.TrimSpace(val))

	if _, found := VisibilityLevels[level]; !found {
		return m
	}

	return m.set(KeyVisibility, level)
}

func (m *metadata) SetTagsFromSlice(tags []string) Metadata {
	clean := generics.Map(tags, strings.TrimSpace)

	return m.set("tags", strings.Join(clean, TagDelimiter))
}

func (m *metadata) GetDoc() string        { return m.get(KeyDoc) }
func (m *metadata) GetSince() string      { return m.get(KeySince) }
func (m *metadata) GetVersion() string    { return m.get(KeyVersion) }
func (m *metadata) GetDeprecated() string { return m.get(KeyDeprecated) }
func (m *metadata) GetProtocol() string   { return m.get(KeyProtocol) }
func (m *metadata) GetVisibility() string { return m.get(KeyVisibility) }
func (m *metadata) GetExample() string    { return m.get(KeyExample) }
func (m *metadata) GetTags() string       { return m.get(KeyTags) }
func (m *metadata) GetAuthor() string     { return m.get(KeyAuthor) }

func (m *metadata) Clone() Metadata {
	return &metadata{data: m.List()}
}

func (m *metadata) Tags() []string {
	result := []string{}
	tags := m.get(KeyTags)

	if len(tags) > 0 {
		result = generics.Map(
			strings.Split(tags, TagDelimiter),
			strings.TrimSpace)
	}

	return result
}

func (m *metadata) TagsNormalised() []string {
	seen := map[string]struct{}{}

	for _, t := range m.Tags() {
		t = strings.ToLower(t)

		if t != "" {
			seen[t] = struct{}{}
		}
	}

	out := make([]string, 0, len(seen))

	for t := range seen {
		out = append(out, t)
	}

	slices.Sort(out)

	return out
}

func (m *metadata) Merge(src Metadata, overwrite bool) Metadata {
	if src == nil {
		return m
	}

	// Snapshot the source safely.
	payload := src.List()

	// Fast path: nothing to do
	if len(payload) == 0 {
		return m
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.data == nil {
		m.data = make(map[string]string, len(payload))
	}

	for key, val := range payload {
		if !overwrite {
			if _, exists := m.data[key]; exists {
				continue
			}
		}

		m.setWithoutLock(key, val) //nolint:errcheck
	}

	return m
}

// ** Functions:

func isReservedKey(key string) bool {
	_, reserved := ReservedMetadataKeys[key]

	return reserved
}

func validateKey(key string) error {
	if len(strings.TrimSpace(key)) == 0 {
		return errors.WithStack(ErrKeyIsInvalid)
	}

	for reserved := range ReservedMetadataKeys {
		if strings.EqualFold(key, reserved) {
			return errors.WithMessagef(
				ErrKeyIsAmbiguous,
				"%q is too close to %q",
				key,
				reserved)
		}
	}

	return nil
}

func NewMetadata() Metadata {
	return &metadata{data: make(map[string]string)}
}

// * metadata.go ends here.
