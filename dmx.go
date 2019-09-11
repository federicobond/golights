package main

import (
	"log"
	"image/color"

	"github.com/Hundemeier/go-sacn/sacn"
)

func dmx() (chan<- [512]byte, func()) {
	//instead of "" you could provide an ip-address that the socket should bind to
	trans, err := sacn.NewTransmitter("", [16]byte{1, 2, 3}, "test")

	//activates the first universe
	ch, err := trans.Activate(1)
	if err != nil {
		log.Fatal(err)
	}
	//deactivate the channel
	cleanup := func () {
		close(ch)
	}

	//send some random data for 10 seconds
	// for i := 0; i < 20; i++ {
	// 	ch <- [512]byte{byte(rand.Int()), byte(i & 0xFF)}
	// 	time.Sleep(500 * time.Millisecond)
	// }

	return ch, cleanup
}

func toDMX(pixels []color.Color) [512]byte {
	var data [512]byte
	for i, p := range pixels {
		r, g, b, _ := p.RGBA()
		data[i * 3 + 0] = byte(r)
		data[i * 3 + 1] = byte(g)
		data[i * 3 + 2] = byte(b)
	}
	return data
}