package main

import (
	"bufio"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// ==================
	f, err := os.Open("pop.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	// ==================

	reader := bufio.NewReader(os.Stdin)

	print("Beats Per Minute: ")
	bpmInput, _ := reader.ReadString('\n')
	bpm, _ := strconv.ParseFloat(strings.TrimRight(bpmInput, "\n"), 64)

	print("Beats Per Bar: ")
	bpbInput, _ := reader.ReadString('\n')
	bpb, _ := strconv.Atoi(strings.TrimRight(bpbInput, "\n"))

	// ==================

	d := time.Duration(float64(time.Minute) / bpm)
	fmt.Println("Delay:", d)
	t := time.NewTicker(d)
	i := 1
	for _ = range t.C {
		i--
		if i == 0 {
			i = bpb
			fmt.Printf("\nTICK ")

			shot := buffer.Streamer(0, buffer.Len())
			speaker.Play(shot)

		} else {
			fmt.Printf("tick ")

			shot := buffer.Streamer(0, buffer.Len())
			speaker.Play(shot)
		}
	}
}
