/*
 * callbacks.go --- Default callbacks.
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

// Default signal handler.
func defaultHandler(*Application) {
}

// Default main loop.
func defaultMainLoop(*Application) {
}

// Special case for HUP support.
func defaultOnHUP(app *Application) {
	app.logger.Info(
		"Default SIGHUP handler invoked.",
	)
}

/* callbacks.go ends here. */
