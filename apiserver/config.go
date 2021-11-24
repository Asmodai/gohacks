/*
 * config.go --- Dispatcher configuration.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package apiserver

type Config struct {
	Addr    string `json:"address"`
	Cert    string `json:"cert_file"`
	Key     string `json:"key_file"`
	UseTLS  bool   `json:"use_tls"`
	LogFile string `json:"log_file"`
}

func NewDefaultConfig() *Config {
	return &Config{}
}

func NewConfig(addr, log, cert, key string, tls bool) *Config {
	return &Config{
		Addr:    addr,
		Cert:    cert,
		Key:     key,
		UseTLS:  tls,
		LogFile: log,
	}
}

/* config.go ends here. */
