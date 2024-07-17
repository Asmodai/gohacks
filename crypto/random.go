/*
 * random.go --- Random string generation.
 *
 * Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package crypto

import (
	"gitlab.com/tozd/go/errors"

	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
)

func assertAvailablePRNG() {
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(
			fmt.Sprintf(
				"crypto/rand is unavailable: read() failed with %#v",
				err,
			),
		)
	}
}

func GenerateRandomBytes(n int) ([]byte, error) {
	buf := make([]byte, n)

	_, err := rand.Read(buf)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return buf, nil
}

func GenerateRandomString(count int) (string, error) {
	const alnum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	ret := make([]byte, count)

	for idx := 0; idx < count; idx++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alnum))))
		if err != nil {
			return "", errors.WithStack(err)
		}

		ret[idx] = alnum[num.Int64()]
	}

	return string(ret), nil
}

func GenerateRandomSafeString(count int) (string, error) {
	b, err := GenerateRandomBytes(count)
	if err != nil {
		err = errors.WithStack(err)
	}

	return base64.URLEncoding.EncodeToString(b), err
}

func GenerateRandomSafeBytes(count int) ([]byte, error) {
	b, err := GenerateRandomSafeString(count)
	if err != nil {
		err = errors.WithStack(err)
	}

	return []byte(b), err
}

//nolint:gochecknoinits
func init() {
	assertAvailablePRNG()
}

/* random.go ends here. */
