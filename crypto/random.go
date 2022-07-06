/*
 * random.go --- Random string generation.
 *
 * Copyright (c) 2021-2022 Paul Ward <asmodai@gmail.com>
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
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(n int) (string, error) {
	const alnum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	ret := make([]byte, n)

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alnum))))
		if err != nil {
			return "", err
		}

		ret[i] = alnum[num.Int64()]
	}

	return string(ret), nil
}

func GenerateRandomSafeString(n int) (string, error) {
	b, err := GenerateRandomBytes(n)

	return base64.URLEncoding.EncodeToString(b), err
}

func GenerateRandomSafeBytes(n int) ([]byte, error) {
	b, err := GenerateRandomSafeString(n)

	return []byte(b), err
}

func init() {
	assertAvailablePRNG()
}

/* random.go ends here. */
