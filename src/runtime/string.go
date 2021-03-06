package runtime

// This file implements functions related to Go strings.

import (
	"unsafe"
)

// The underlying struct for the Go string type.
type _string struct {
	ptr    *byte
	length lenType
}

// The iterator state for a range over a string.
type stringIterator struct {
	byteindex  lenType
	rangeindex lenType
}

// Return true iff the strings match.
//go:nobounds
func stringEqual(x, y string) bool {
	if len(x) != len(y) {
		return false
	}
	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

// Return true iff x < y.
//go:nobounds
func stringLess(x, y string) bool {
	l := len(x)
	if m := len(y); m < l {
		l = m
	}
	for i := 0; i < l; i++ {
		if x[i] < y[i] {
			return true
		}
		if x[i] > y[i] {
			return false
		}
	}
	return len(x) < len(y)
}

// Add two strings together.
func stringConcat(x, y _string) _string {
	if x.length == 0 {
		return y
	} else if y.length == 0 {
		return x
	} else {
		length := uintptr(x.length + y.length)
		buf := alloc(length)
		memcpy(buf, unsafe.Pointer(x.ptr), uintptr(x.length))
		memcpy(unsafe.Pointer(uintptr(buf)+uintptr(x.length)), unsafe.Pointer(y.ptr), uintptr(y.length))
		return _string{ptr: (*byte)(buf), length: lenType(length)}
	}
}

// Create a string from a []byte slice.
func stringFromBytes(x struct {
	ptr *byte
	len lenType
	cap lenType
}) _string {
	buf := alloc(uintptr(x.len))
	memcpy(buf, unsafe.Pointer(x.ptr), uintptr(x.len))
	return _string{ptr: (*byte)(buf), length: lenType(x.len)}
}

// Convert a string to a []byte slice.
func stringToBytes(x _string) (slice struct {
	ptr *byte
	len lenType
	cap lenType
}) {
	buf := alloc(uintptr(x.length))
	memcpy(buf, unsafe.Pointer(x.ptr), uintptr(x.length))
	slice.ptr = (*byte)(buf)
	slice.len = x.length
	slice.cap = x.length
	return
}

// Create a string from a Unicode code point.
func stringFromUnicode(x rune) _string {
	array, length := encodeUTF8(x)
	// Array will be heap allocated.
	// The heap most likely doesn't work with blocks below 4 bytes, so there's
	// no point in allocating a smaller buffer for the string here.
	return _string{ptr: (*byte)(unsafe.Pointer(&array)), length: length}
}

// Iterate over a string.
// Returns (ok, key, value).
func stringNext(s string, it *stringIterator) (bool, int, rune) {
	if len(s) <= int(it.byteindex) {
		return false, 0, 0
	}
	r, length := decodeUTF8(s, it.byteindex)
	it.byteindex += length
	it.rangeindex += 1
	return true, int(it.rangeindex), r
}

// Convert a Unicode code point into an array of bytes and its length.
func encodeUTF8(x rune) ([4]byte, lenType) {
	// https://stackoverflow.com/questions/6240055/manually-converting-unicode-codepoints-into-utf-8-and-utf-16
	// Note: this code can probably be optimized (in size and speed).
	switch {
	case x <= 0x7f:
		return [4]byte{byte(x), 0, 0, 0}, 1
	case x <= 0x7ff:
		b1 := 0xc0 | byte(x>>6)
		b2 := 0x80 | byte(x&0x3f)
		return [4]byte{b1, b2, 0, 0}, 2
	case x <= 0xffff:
		b1 := 0xe0 | byte(x>>12)
		b2 := 0x80 | byte((x>>6)&0x3f)
		b3 := 0x80 | byte((x>>0)&0x3f)
		return [4]byte{b1, b2, b3, 0}, 3
	case x <= 0x10ffff:
		b1 := 0xf0 | byte(x>>18)
		b2 := 0x80 | byte((x>>12)&0x3f)
		b3 := 0x80 | byte((x>>6)&0x3f)
		b4 := 0x80 | byte((x>>0)&0x3f)
		return [4]byte{b1, b2, b3, b4}, 4
	default:
		// Invalid Unicode code point.
		return [4]byte{0xef, 0xbf, 0xbd, 0}, 3
	}
}

// Decode a single UTF-8 character from a string.
//go:nobounds
func decodeUTF8(s string, index lenType) (rune, lenType) {
	remaining := lenType(len(s)) - index // must be >= 1 before calling this function
	x := s[index]
	switch {
	case x&0x80 == 0x00: // 0xxxxxxx
		return rune(x), 1
	case x&0xe0 == 0xc0: // 110xxxxx
		if remaining < 2 {
			return 0xfffd, 1
		}
		return (rune(x&0x1f) << 6) | (rune(s[index+1]) & 0x3f), 2
	case x&0xf0 == 0xe0: // 1110xxxx
		if remaining < 3 {
			return 0xfffd, 1
		}
		return (rune(x&0x0f) << 12) | ((rune(s[index+1]) & 0x3f) << 6) | (rune(s[index+2]) & 0x3f), 3
	case x&0xf8 == 0xf0: // 11110xxx
		if remaining < 4 {
			return 0xfffd, 1
		}
		return (rune(x&0x07) << 18) | ((rune(s[index+1]) & 0x3f) << 12) | ((rune(s[index+2]) & 0x3f) << 6) | (rune(s[index+3]) & 0x3f), 4
	default:
		return 0xfffd, 1
	}
}
