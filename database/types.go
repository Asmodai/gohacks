// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// types.go --- SQL datatype hacks.
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

package database

import (
	"gitlab.com/tozd/go/errors"

	"database/sql"
	"encoding/json"
	"time"
)

//nolint:gochecknoglobals
var (
	zeroTime = time.Time{}
)

// ==[ NullBool ]=====================================================

type NullBool struct {
	sql.NullBool
}

func (x NullBool) MarshalJSON() ([]byte, error) {
	if x.Valid {
		rval, err := json.Marshal(x.Bool)
		if err != nil {
			return []byte{}, errors.WithStack(err)
		}

		return rval, nil
	}

	return []byte("false"), nil
}

// ==[ NullByte ]=====================================================

type NullByte struct {
	sql.NullByte
}

func (x NullByte) MarshalJSON() ([]byte, error) {
	if x.Valid {
		rval, err := json.Marshal(x.Byte)
		if err != nil {
			return []byte{}, errors.WithStack(err)
		}

		return rval, nil
	}

	return []byte(""), nil
}

// ==[ NullFloat64 ]==================================================

type NullFloat64 struct {
	sql.NullFloat64
}

func (x NullFloat64) MarshalJSON() ([]byte, error) {
	if x.Valid {
		rval, err := json.Marshal(x.Float64)
		if err != nil {
			return []byte{}, errors.WithStack(err)
		}

		return rval, nil
	}

	return []byte("0.0"), nil
}

// ==[ NullInt16 ]====================================================

type NullInt16 struct {
	sql.NullInt16
}

func (x NullInt16) MarshalJSON() ([]byte, error) {
	if x.Valid {
		rval, err := json.Marshal(x.Int16)
		if err != nil {
			return []byte{}, errors.WithStack(err)
		}

		return rval, nil
	}

	return []byte("0"), nil
}

// ==[ NullInt32 ]====================================================

type NullInt32 struct {
	sql.NullInt32
}

func (x NullInt32) MarshalJSON() ([]byte, error) {
	if x.Valid {
		rval, err := json.Marshal(x.Int32)
		if err != nil {
			return []byte{}, errors.WithStack(err)
		}

		return rval, nil
	}

	return []byte("0"), nil
}

// ==[ NullInt64 ]====================================================

type NullInt64 struct {
	sql.NullInt64
}

func (x NullInt64) MarshalJSON() ([]byte, error) {
	if x.Valid {
		rval, err := json.Marshal(x.Int64)
		if err != nil {
			return []byte{}, errors.WithStack(err)
		}

		return rval, nil
	}

	return []byte("0"), nil
}

// ==[ NullString ]===================================================

type NullString struct {
	sql.NullString
}

func (x NullString) MarshalJSON() ([]byte, error) {
	if x.Valid {
		rval, err := json.Marshal(x.String)
		if err != nil {
			return []byte{}, errors.WithStack(err)
		}

		return rval, nil
	}

	return []byte("null"), nil
}

// ==[ NullTime ]=====================================================

type NullTime struct {
	sql.NullTime
}

func (x NullTime) MarshalJSON() ([]byte, error) {
	if x.Valid {
		rval, err := json.Marshal(x.Time)
		if err != nil {
			return []byte{}, errors.WithStack(err)
		}

		return rval, nil
	}

	rval, err := json.Marshal(zeroTime.UTC())
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}

	return rval, nil
}

// types.go ends here.
