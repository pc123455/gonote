package utils

import (
	crand "crypto/rand"
	"fmt"
	mrand "math/rand"
	"time"
)

type UUID [16]byte

var seeded = false

func Rand() UUID {
	var x UUID
	randBytes(x[:])
	x[6] = (x[6] & 0x0F) | 0x40
	x[8] = (x[8] & 0x3F) | 0x80
	return x
}

func (this *UUID) Hex() string {
	x := [16]byte(*this)
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		x[0], x[1], x[2], x[3], x[4],
		x[5], x[6],
		x[7], x[8],
		x[9], x[10], x[11], x[12], x[13], x[14], x[15])
}

func randBytes(x []byte) {
	length := len(x)
	n, err := crand.Read(x)
	if n != length || err != nil {
		if !seeded {
			mrand.Seed(time.Now().UnixNano())
		}

		for i := 0; i < length; i++ {
			x[i] = byte(mrand.Int31n(256))
		}
	}
}
