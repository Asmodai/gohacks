/*
 * config.go --- SQL configuration.
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

package database

import (
	"github.com/Asmodai/gohacks/secrets"

	"fmt"
)

/*
SQL configuration structure.
*/
type Config struct {
	Driver         string `json:"driver"`
	Username       string `json:"username"`
	UsernameSecret string `json:"username_secret"`
	Password       string `json:"password" config_obscure:"true"`
	PasswordSecret string `json:"password_secret"`
	Hostname       string `json:"hostname"`
	Port           int    `json:"port"`
	Database       string `json:"database"`
	BatchSize      int    `json:"batch_size"`
	SetPoolLimits  bool   `json:"set_pool_limits"`
	MaxIdleConns   int    `json:"max_idle_conns"`
	MaxOpenConns   int    `json:"max_open_conns"`

	cachedDsn string `json:"-" config_hide:"true"`
}

// Create a new configuration object.
func NewConfig() *Config {
	return &Config{
		Driver:         "mysql",
		Username:       "",
		UsernameSecret: "",
		Password:       "",
		PasswordSecret: "",
		Hostname:       "",
		Port:           0,
		Database:       "",
		cachedDsn:      "",
		BatchSize:      0,
		SetPoolLimits:  true,
		MaxIdleConns:   20,
		MaxOpenConns:   20,
	}
}

// Check if we should use secrets files.
func (c *Config) checkSecrets() error {
	if c.UsernameSecret != "" {
		secret := secrets.Make(c.UsernameSecret)

		if err := secret.Probe(); err != nil {
			return err
		}

		c.Username = secret.Value()
	}

	if c.PasswordSecret != "" {
		secret := secrets.Make(c.PasswordSecret)

		if err := secret.Probe(); err != nil {
			return err
		}

		c.Password = secret.Value()
	}

	return nil
}

func (c *Config) Validate() error {
	if err := c.checkSecrets(); err != nil {
		return err
	}

	return nil
}

// Return the DSN for this database configuration.
func (c *Config) ToDSN() string {
	if c.cachedDsn != "" {
		return c.cachedDsn
	}

	var s string = ""

	if c.Username != "" {
		s += c.Username

		if c.Password != "" {
			s += ":" + c.Password
		}

		s += "@"
	}

	s += "tcp(" + c.Hostname
	if c.Port > 0 {
		s += fmt.Sprintf(":%d", c.Port)
	}
	s += ")/" + c.Database

	// Set session timezone to UTC.
	s += "?parseTime=True&loc=UTC&time_zone='-00:00'"

	c.cachedDsn = s

	return s
}

/* config.go ends here. */
