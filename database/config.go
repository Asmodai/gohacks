/*
 * config.go --- SQL configuration.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package database

import (
	"fmt"
)

/*

SQL configuration structure.

*/
type Config struct {
	Driver        string `json:"driver"`
	Username      string `json:"username"`
	Password      string `json:"password" config_obscure:"true"`
	Hostname      string `json:"hostname"`
	Port          int    `json:"port"`
	Database      string `json:"database"`
	BatchSize     int    `json:"batch_size"`
	SetPoolLimits bool   `json:"set_pool_limits"`
	MaxIdleConns  int    `json:"max_idle_conns"`
	MaxOpenConns  int    `json:"max_open_conns"`

	cachedDsn string `json:"-" config_hide:"true"`
}

// Create a new configuration object.
func NewConfig() *Config {
	return &Config{
		Driver:        "mysql",
		Username:      "",
		Password:      "",
		Hostname:      "",
		Port:          0,
		Database:      "",
		cachedDsn:     "",
		BatchSize:     0,
		SetPoolLimits: true,
		MaxIdleConns:  20,
		MaxOpenConns:  20,
	}
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
