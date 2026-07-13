//go:build tinygo && esp32s3

// Package statusled drives a single onboard WS2812/NeoPixel LED to show
// pipeline status at a glance, without needing a serial monitor attached.
//
// Vendored instead of using tinygo.org/x/drivers/ws2812: that driver's
// xtensa bit-bang implementation only has hand-cycle-counted assembly for
// 80MHz and 160MHz CPU clocks (verified against both the pinned v0.35.0 and
// the current upstream dev branch -- neither has a 240MHz case). This
// project's esp32s3 runtime always boots at 240MHz (see
// runtime_esp32s3.go's main(), which is fixed, not configurable), so the
// upstream driver returns errUnknownClockSpeed unconditionally on this
// board. This file reuses the exact same instruction sequence as upstream's
// 80MHz case (proven correct there), just with every NOP-padding count
// tripled -- 240MHz is exactly 3x 80MHz, and each non-NOP instruction here
// takes a fixed number of cycles regardless of clock speed, so scaling
// the filler NOPs by the same ratio as the clock speed reproduces the same
// real-world nanosecond timings at 240MHz.
package statusled

import (
	"device"
	"machine"
	"runtime/interrupt"
	"time"
	"unsafe"
)

// SetColor drives a single WS2812/NeoPixel LED on pin with the given RGB
// values (0-255 each). Blocks for the duration of the transfer plus a reset
// latch. Not safe to call concurrently with other WS2812 writes on the same
// pin without external synchronization.
func SetColor(pin machine.Pin, r, g, b byte) {
	// WS2812 wire order is GRB, not RGB.
	writeByte(pin, g)
	writeByte(pin, r)
	writeByte(pin, b)
	time.Sleep(60 * time.Microsecond) // latch/reset (spec requires >50us low)
}

// writeByte sends a single byte MSB-first using the WS2812 protocol.
// Timings target: T0H ~262.5ns, T0L ~837.5ns, T1H ~587.5ns, T1L ~487.5ns
// (scaled 3x from the proven 80MHz case in tinygo.org/x/drivers/ws2812).
func writeByte(pin machine.Pin, c byte) {
	portSet, maskSet := pin.PortMaskSet()
	portClear, maskClear := pin.PortMaskClear()
	mask := interrupt.Disable()

	device.AsmFull(`
		1: // send_bit
			s32i  {maskSet}, {portSet}, 0     // T0H and T1H start here
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			slli  {value}, {value}, 1         // shift value left by 1
			bbsi  {value}, 8, 2f              // branch to skip_store if bit 8 is set
			s32i  {maskClear}, {portClear}, 0 // T0H -> T0L transition
		2: // skip_store
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			s32i  {maskClear}, {portClear}, 0 // T1H -> T1L transition
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			nop
			addi  {i}, {i}, -1
			bnez {i}, 1b                      // send_bit, T1H and T1L end here

			// Restore original values after modifying them in the inline
			// assembly. Not doing that would result in undefined behavior as
			// the compiler doesn't know we're modifying these values.
			movi.n {i}, 8
			slli  {value}, {value}, 8
		`, map[string]interface{}{
		// Note: casting pointers to uintptr here because of what might be
		// an Xtensa backend bug with inline assembly.
		"value":     uint32(c),
		"i":         8,
		"maskSet":   maskSet,
		"portSet":   uintptr(unsafe.Pointer(portSet)),
		"maskClear": maskClear,
		"portClear": uintptr(unsafe.Pointer(portClear)),
	})
	interrupt.Restore(mask)
}
