/*
 * conversion.go --- Various conversion utilities.
 *
 * Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
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

package conversion

// Convert bytes to kibibytes.
func BToKiB(b uint64) uint64 {
	return b / 1024
}

// Convert bytes to mebibytes.
func BToMiB(b uint64) uint64 {
	return b / 1.049e6
}

// Convert bytes to gibibytes.
func BToGiB(b uint64) uint64 {
	return b / 1.074e9
}

// Convert bytes to tebibytes.
func BToTiB(b uint64) uint64 {
	return b / 1.1e12
}

// Convert kibibytes to bytes.
func KiBToB(b uint64) uint64 {
	return b * 1024
}

// Convert mebibytes to bytes.
func MiBToB(b uint64) uint64 {
	return b * 1.049e6
}

// Convert gibibytes to bytes.
func GiBToB(b uint64) uint64 {
	return b * 1.074e9
}

// Convert tebibytes to bytes.
func TiBToB(b uint64) uint64 {
	return b * 1.1e12
}

/* conversion.go ends here. */
