// +build wasm,!arm,!avr

package runtime

type timeUnit int64

const tickMicros = 1

var timestamp timeUnit

var line []byte

const (
	logLevelError   = 1
	logLevelWarning = 3
	logLevelInfo    = 6
)

// CommonWA: log_write
func _Cfunc_log_write(level int32, ptr *uint8, len int32)

//go:export _start
func start() {
	initAll()
}

//go:export cwa_main
func main() {
	mainWrapper()
}

//go:linkname _writeLog runtime.writeLog
func _writeLog(level int32, ptr *uint8, len, cap lenType) {
	_Cfunc_log_write(level, ptr, int32(len))
}

// hack around slice types
func writeLog(level int32, line []byte)

func putchar(c byte) {
	switch c {
	case '\r':
		// ignore
	case '\n':
		// write line
		writeLog(logLevelInfo, line)
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
