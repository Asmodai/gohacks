// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// config.go --- Rate-limited HTTP client configuration.
//
// Copyright (c) 2026 Paul Ward <paul@lisphacker.uk>
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

package rlhttp

// * Imports:

import (
	"time"

	"github.com/Asmodai/gohacks/errx"
	"github.com/Asmodai/gohacks/types"
)

// * Constants:

const (
	// A sane default client timeout value.
	//
	// This is in line with Go, as well as with services like Kubernetes,
	// AWS, et al.
	DefaultClientTimeout = types.Duration(30 * time.Second)

	// Default burst value.
	DefaultBurst int = 1
)

// * Variables:

var (
	ErrInvalidSettings = errx.Base("invalid rate limiter settings")
	ErrInvalidLimiter  = errx.Base("invalid rate limit")
)

// * Code:
// ** Type:

/*
Rate limiter configuration.

The limiter spaces requests at an interval of Every/Max (i.e. `Max` requests
per `Every`).

If rate limiting is enabled, then the following conditions hold true:

1) If there is no `Timeout` then a default of 30 seconds is used,
2) If there is no `Burst` then a default of 1 is used.
3) If there is no `Every` then a validation error shall be raised.
4) If there is no `Max` then a validation error shall be raised.
5) If `Burst` is greater than `Max` then a validation error shall be raised.

It is advised to call `Validate` after populating `Config` and checking if
there are any raised errors.

If `Validate` does raise errors, those will need to be addressed first.
*/
type Config struct {
	Enabled bool           `json:"enabled"` // Rate limiting enabled?
	Timeout types.Duration `json:"timeout"` // Request timeout.
	Every   types.Duration `json:"every"`   // Time measure.
	Burst   int            `json:"burst"`   // Number of bursts
	Max     int            `json:"max"`     // Max requests per measure.
}

// ** Methods:

// Validate rate limiter settings.
// TODO: Change `Validate` to `Validate/Normalise`. Don't forget!
func (c *Config) Validate() []error {
	errs := []error{}

	if !c.Enabled {
		// No point validating unless enabled.
		return nil
	}

	if c.Timeout <= 0 {
		c.Timeout = DefaultClientTimeout
	}

	if c.Burst <= 0 {
		c.Burst = DefaultBurst
	}

	if c.Every <= 0 {
		errs = append(
			errs,
			errx.WithMessage(
				ErrInvalidSettings,
				"'every' cannot be 0 or less",
			),
		)
	}

	if c.Max <= 0 {
		errs = append(
			errs,
			errx.WithMessage(
				ErrInvalidSettings,
				"'max' cannot be 0 or less",
			),
		)
	}

	if c.Burst > c.Max {
		errs = append(
			errs,
			errx.WithMessage(
				ErrInvalidSettings,
				"'burst' cannot be higher than 'max'",
			),
		)
	}

	return errs
}

// ** Functions:

// Create a new default rate limiter configuration.
func NewDefaultConfig() *Config {
	return &Config{
		Enabled: false,
	}
}

// * config.go ends here.
