/*
 * conversion.go --- Various conversion utilities.
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

package conversion

func BToKiB(b uint64) uint64 {
	return b / 1024
}

func BToMiB(b uint64) uint64 {
	return b / 1.049e6
}

func BToGiB(b uint64) uint64 {
	return b / 1.074e9
}

func BToTiB(b uint64) uint64 {
	return b / 1.1e12
}

func KiBToB(b uint64) uint64 {
	return b * 1024
}

func MiBToB(b uint64) uint64 {
	return b * 1.049e6
}

func GiBToB(b uint64) uint64 {
	return b * 1.074e9
}

func TiBToB(b uint64) uint64 {
	return b * 1.1e12
}

/* conversion.go ends here. */
