/*
 * secrets.go --- `Secrets' file support.
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

package secrets

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	SecretsPath string = "/run/secrets"
)

type Secret struct {
	path  string
	value string
}

func New() *Secret {
	return Make("")
}

func Make(file string) *Secret {
	return &Secret{
		path:  SecretsPath + "/" + filepath.Base(file),
		value: "",
	}
}

func (s *Secret) Path() string  { return s.path }
func (s *Secret) Value() string { return s.value }

func (s *Secret) SetPath(val string) error {
	s.path = filepath.Base(val)

	return s.probe()
}

func (s *Secret) Probe() error {
	return s.probe()
}

func (s *Secret) probe() error {
	if s.path == "" {
		return errors.New("No secret path set!")
	}

	fp, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer fp.Close()

	bytes, _ := io.ReadAll(fp)
	if len(bytes) == 0 {
		return errors.New("Zero length secret!")
	}

	s.value = strings.TrimSuffix(string(bytes), "\n")
	return nil
}

/* secrets.go ends here. */
