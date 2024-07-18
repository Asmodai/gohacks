/*
 * filesystem.go --- Filesystem utilities.
 *
 * Copyright (c) 2024 Paul Ward <asmodai@gmail.com>
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

package utils

import (
	"os"
)

type Filesystem struct {
	path   string
	info   os.FileInfo
	exists bool
}

func NewFilesystem(path string) *Filesystem {
	obj := &Filesystem{
		path:   path,
		exists: false,
	}

	obj.probe()

	return obj
}

func (fs *Filesystem) Exists() bool {
	return fs.exists
}

func (fs *Filesystem) IsDirectory() bool {
	if fs.info == nil {
		return false
	}

	return fs.info.IsDir()
}

func (fs *Filesystem) IsFile() bool {
	if fs.info == nil {
		return false
	}

	return fs.info.Mode().IsRegular()
}

func (fs *Filesystem) IsOwnerExecutable() bool {
	if fs.info == nil {
		return false
	}

	return fs.info.Mode()&0100 != 0
}

func (fs *Filesystem) IsGroupExecutable() bool {
	if fs.info == nil {
		return false
	}

	return fs.info.Mode()&0010 != 0
}

func (fs *Filesystem) IsOtherExecutable() bool {
	if fs.info == nil {
		return false
	}

	return fs.info.Mode()&0001 != 0
}

func (fs *Filesystem) IsExecutable() bool {
	return fs.IsOwnerExecutable() ||
		fs.IsGroupExecutable() ||
		fs.IsOtherExecutable()
}

func (fs *Filesystem) Name() string {
	if fs.info == nil {
		return ""
	}

	return fs.info.Name()
}

func (fs *Filesystem) Size() int64 {
	if fs.info == nil {
		return 0
	}

	return fs.info.Size()
}

func (fs *Filesystem) Mode() os.FileMode {
	if fs.info == nil {
		return 0
	}

	return fs.info.Mode()
}

func (fs *Filesystem) probe() {
	var (
		err  error
		file *os.File
	)

	if fs.path == "" {
		return
	}

	// We don't care about handling errors.
	if file, err = os.Open(fs.path); err == nil {
		fs.exists = true
		fs.info, _ = file.Stat()

		file.Close()
	}
}

/* filesystem.go ends here. */
