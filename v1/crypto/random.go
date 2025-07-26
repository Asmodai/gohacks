// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// random.go --- Random string generation.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
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

//
//
//

// * Package:

package crypto

// * Imports:

import (
	"gitlab.com/tozd/go/errors"

	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// * Code:

// ** Functions:

// Check that the underlying OS has a cryptographic randomiser by attempting
// to exercise it.  If the read operation fails, then panic.
func assertAvailablePRNG() {
	buf := make([]byte, 1)

	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: %#v", err))
	}
}

// Generate n number of random bytes from a cryptographic randomiser.
func GenerateRandomBytes(n int) ([]byte, error) {
	buf := make([]byte, n)

	if _, err := rand.Read(buf); err != nil {
		return nil, errors.WithStack(err)
	}

	return buf, nil
}

// Generate a random string of the given length using bytes from the
// cryptographic randomiser.
//
//nolint:mnd
func GenerateRandomString(count int) (string, error) {
	const (
		alnum     = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
		alnumLen  = byte(len(alnum))
		byteLimit = 255
		maxByte   = byteLimit - (byteLimit % alnumLen) // Discard bias.
	)

	ret := make([]byte, 0, count)
	buf := make([]byte, count*2) // Overshoot to reduce re-fills.

	defer func() {
		for i := range buf {
			buf[i] = 0
		}
	}()

	for len(ret) < count {
		if _, err := rand.Read(buf); err != nil {
			return "", errors.WithStack(err)
		}

		for _, b := range buf {
			if b > maxByte {
				continue // Avoid modulo bias.
			}

			ret = append(ret, alnum[b%alnumLen])
			if len(ret) == count {
				break
			}
		}
	}

	return string(ret), nil
}

// Operates the same way as `GenerateRandomString` but encodes the result
// using base64 encoding.
//
//nolint:mnd
func GenerateRandomSafeString(count int) (string, error) {
	b, err := GenerateRandomBytes((count*6 + 7) / 8)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return base64.URLEncoding.EncodeToString(b)[:count], nil
}

// Operates the same way as `GenerateRandomBytes` but encodes the result using
// base64 encoding.
func GenerateRandomSafeBytes(count int) ([]byte, error) {
	b, err := GenerateRandomSafeString(count)
	if err != nil {
		err = errors.WithStack(err)
	}

	return []byte(b), err
}

// Operates the same way as `GenerateRandomBytes` but encodes the result
// as a string of hexadecimal characters.
func GenerateRandomHexString(count int) (string, error) {
	b, err := GenerateRandomBytes(count)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return hex.EncodeToString(b), nil
}

// ** Init magic function:

// Initialise the randomiser.
//
//nolint:gochecknoinits
func init() {
	assertAvailablePRNG()
}

// * random.go ends here.
