/*
 * config.go --- Process configuration.
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

package process

// Process configuration structure.
type Config struct {
	Name     string     // Pretty name.
	Interval int        // `RunEvery` time interval.
	ActionFn CallbackFn // `Action` callback.
	StopFn   OnStopFn   // `Stop` callback.
}

// Create a default process configuration.
func NewDefaultConfig() *Config {
	return &Config{}
}

/* config.go ends here. */
