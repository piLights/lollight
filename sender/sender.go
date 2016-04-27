package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/piLights/dioder-rpc/src/proto"
)

type rgbColor struct {
	R, G, B uint8
}

var (
	input       chan rgbColor
	hostName    = flag.String("host", "", "<host[:port]> of pidioder")
	colorSender LighterGRPC.LighterClient
)

func send(c rgbColor) {
	colorMessage := &LighterGRPC.ColorMessage{true, int32(c.R), int32(c.G), int32(c.B), 100, "dioderLollight", ""}

	success, error := colorSender.SetColor(context.Background(), colorMessage)
	if error != nil {
		log.Fatal(error)
	} else {
		log.Println(success)
	}
}

func stepChannel(cur, target uint8) uint8 {
	switch {
	case cur < target:
		return cur + 1
	case cur > target:
		return cur - 1
	default:
		return cur
	}
}

func fadeStep(cur, target rgbColor) rgbColor {
	return rgbColor{
		stepChannel(cur.R, target.R),
		stepChannel(cur.G, target.G),
		stepChannel(cur.B, target.B),
	}
}

func fader() {
	var current rgbColor
	var last rgbColor

	for {
		select {
		case c := <-input:
			current = c
		default:
			last = fadeStep(last, current)
			send(last)
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func main() {
	flag.Parse()

	if *hostName == "" {
		log.Fatal("No PiDioder address given.")
	}

	//Make the colorSender
	connection, error := grpc.Dial(*hostName, grpc.WithInsecure())
	if error != nil {
		log.Fatal(error)
	}
	defer connection.Close()

	colorSender = LighterGRPC.NewLighterClient(connection)

	input = make(chan rgbColor)

	go fader()

	for {
		var r, g, b int

		_, err := fmt.Scanf("%d %d %d", &r, &g, &b)

		if err != nil {
			log.Fatal(err)
		}

		input <- rgbColor{uint8(r), uint8(g), uint8(b)}
	}
}
