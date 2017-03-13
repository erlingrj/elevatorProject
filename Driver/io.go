package Driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"

func IoInit() int {
	return int(C.io_init())
}

func IoSetBit(channel int) {
	C.io_set_bit(C.int(channel))
}

func IoClearBit(channel int) {
	C.io_clear_bit(C.int(channel))
}

func IoWriteAnalog(channel int, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}

func IoReadBit(channel int) int {
	return int(C.io_read_bit(C.int(channel)))
}

func IoReadAnalog(channel int) int {
	return int(C.io_read_analog(C.int(channel)))
}
