/*
 * loop.go --- Application loop code.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
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

package app

import (
	"time"
)

// Main loop.
func (app *Application) loop() {
	app.running = true

	// While we're running...
	for app.running == true {
		// Check for parent context cancellation
		select {
		case <-app.ctx.Done():
			// Stop the happening train!
			app.running = false

		default:
		}

		app.MainLoop(app)
		time.Sleep(EventLoopSleep)
	}

	// No longer running, so shut things down.
	app.OnExit(app)
	app.procmgr.StopAll()
	app.logger.Info(
		"Application is terminating.",
		"type", "stop",
	)
}

/* loop.go ends here. */
