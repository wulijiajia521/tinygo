// +build wasm,!arm,!avr

package runtime

type timeUnit int64

const tickMicros = 1

var timestamp timeUnit

var line []byte

//go:export _start
func start() {
}

func putchar(c byte) {
	switch c {
	case '\r':
		// ignore
	case '\n':
		// write line
		line = line[:0]
	default:
		line = append(line, c)
	}
}

func sleepTicks(d timeUnit) {
	// TODO: actually sleep here for the given time.
	timestamp += d
}

func ticks() timeUnit {
	return timestamp
}

// Align on word boundary.
func align(ptr uintptr) uintptr {
	return (ptr + 3) &^ 3
}

func abort() {
	// TODO
}
