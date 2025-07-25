// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// filesystem.go --- Filesystem utilities.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

package utils

// * Imports:

import (
	"os"
)

// * Constants:

const (
	maskOwnerExec os.FileMode = 0100
	maskGroupExec os.FileMode = 0010
	maskOtherExec os.FileMode = 0001
)

// * Code:

// ** Types:

/*
Filesystem

Allow fort he checking of filesystem entities without having to open, stat,
close continuously.

This is a wrapper around various methods from the `os` package that is
designed to allow repeated querying without having to re-open and re-stat the
entity.

If used within a struct, this will not provide any sort of live updates and it
does not utilise `epoll` (or any other similar facility) to update if a file
or directory structure changes.

Maybe one day it will facilitate live updates et al to track file changes, but
that day is not today.
*/
type Filesystem struct {
	path   string
	info   os.FileInfo
	exists bool
}

// ** Methods:

// Does the entity exist?
func (fs *Filesystem) Exists() bool {
	return fs.exists
}

// Is the entity a directory?
func (fs *Filesystem) IsDirectory() bool {
	if fs.info == nil {
		return false
	}

	return fs.info.IsDir()
}

// Is the entity a file?
func (fs *Filesystem) IsFile() bool {
	if fs.info == nil {
		return false
	}

	return fs.info.Mode().IsRegular()
}

// Helper for file permission functions.
func (fs *Filesystem) permissionHelper(mask os.FileMode) bool {
	if fs.info == nil {
		return false
	}

	return fs.info.Mode()&mask != 0
}

// Does the entity's permission bits include "owner executable"?
func (fs *Filesystem) IsOwnerExecutable() bool {
	return fs.permissionHelper(maskOwnerExec)
}

// Does the entity's permission bits include "group executable"?
func (fs *Filesystem) IsGroupExecutable() bool {
	return fs.permissionHelper(maskGroupExec)
}

// Does the entity's permission bits include "other executable"?
func (fs *Filesystem) IsOtherExecutable() bool {
	return fs.permissionHelper(maskOtherExec)
}

// Does the entity's permission bits include any type of "executable"?
func (fs *Filesystem) IsExecutable() bool {
	return fs.IsOwnerExecutable() ||
		fs.IsGroupExecutable() ||
		fs.IsOtherExecutable()
}

// Return the entity's name.
func (fs *Filesystem) Name() string {
	if fs.info == nil {
		return ""
	}

	return fs.info.Name()
}

// Return the entity's size.
func (fs *Filesystem) Size() int64 {
	if fs.info == nil {
		return 0
	}

	return fs.info.Size()
}

// Return the entity's mode.
func (fs *Filesystem) Mode() os.FileMode {
	if fs.info == nil {
		return 0
	}

	return fs.info.Mode()
}

// Probe the entity for its information.
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

// ** Functions:

// Create a new filesystem object.
func NewFilesystem(path string) *Filesystem {
	obj := &Filesystem{
		path:   path,
		exists: false,
	}

	obj.probe()

	return obj
}

// * filesystem.go ends here.
