// Package timesync sets the device clock from an NTP server. Baremetal
// TinyGo has no RTC, so time.Now() starts counting from zero at boot (see
// runtime/baremetal.go) -- without this, every Reading's Unix timestamp is
// actually just seconds-since-boot, not a real point in time.
package timesync

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"time"
)

// packetSize is the fixed size of an NTP client/server packet (RFC 5905).
const packetSize = 48

// toUnixEpochOffset converts NTP's epoch (1900-01-01) to Unix's (1970-01-01): 70 years, accounting for leap days.
const toUnixEpochOffset = 2208988800

func Sync(ntpHost string) error {
	conn, err := net.Dial("udp", ntpHost)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := [packetSize]byte{0xe3}
	if _, err := conn.Write(request[:]); err != nil {
		return err
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	response := make([]byte, packetSize)
	n, err := conn.Read(response)
	if err != nil && err != io.EOF {
		return err
	}
	if n != packetSize {
		return fmt.Errorf("unexpected NTP packet size: %d", n)
	}

	secs := uint32(response[40])<<24 | uint32(response[41])<<16 | uint32(response[42])<<8 | uint32(response[43])
	ntpTime := time.Unix(int64(secs)-toUnixEpochOffset, 0)
	runtime.AdjustTimeOffset(-1 * int64(time.Since(ntpTime)))
	return nil
}
