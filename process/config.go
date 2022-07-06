/*
 * config.go --- Configuration.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
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

package process

import (
	"github.com/Asmodai/gohacks/logger"
)

// Process configuration structure.
type Config struct {
	Name     string         // Pretty name.
	Interval int            // `RunEvery` time interval.
	Function CallbackFn     // `Action` callback.
	OnStart  CallbackFn     // `Start` callback.
	OnStop   CallbackFn     // `Stop` callback.
	OnQuery  QueryFn        // `Query` callback.
	Logger   logger.ILogger // Logger.
}

// Create a default process configuration.
func NewDefaultConfig() *Config {
	return &Config{
		Logger: logger.NewDefaultLogger(),
	}
}

/* config.go ends here. */
