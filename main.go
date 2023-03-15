package main

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	//go:embed pop.mp3
	res embed.FS
)

func main() {
	println("Metronome has started!")

	buffer := initializeBuffer()

	bpm, bpb := retrieveBeatsInput()

	initializeMetronome(bpm, bpb, buffer)
}

func initializeBuffer() *beep.Buffer {
	audioFile, _ := res.Open("pop.mp3")

	streamer, format, err := mp3.Decode(audioFile)
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	return buffer
}

func retrieveBeatsInput() (int, int) {
	reader := bufio.NewReader(os.Stdin)

	print("Beats Per Minute (default 60): ")
	bpmInput, _ := reader.ReadString('\n')
	bpm := 60
	if bpmInput != "\n" {
		bpm, _ = strconv.Atoi(strings.TrimRight(bpmInput, "\n"))
	}

	print("Beats Per Bar (default 4): ")
	bpbInput, _ := reader.ReadString('\n')
	bpb := 4
	if bpbInput != "\n" {
		bpb, _ = strconv.Atoi(strings.TrimRight(bpbInput, "\n"))
	}

	return bpm, bpb
}

func initializeMetronome(numBeatsPerMinute int, numBeatsPerBar int, buffer *beep.Buffer) {
	beatsInterval := time.Duration(float64(time.Minute) / float64(numBeatsPerMinute))
	fmt.Println("beatsInterval:", beatsInterval)

	numSubdivisions := 4
	subdivisionsInterval := beatsInterval / time.Duration(numSubdivisions)

	beatsTicker := time.NewTicker(beatsInterval)

	i := 0
	for beatTime := range beatsTicker.C {
		// beatNum starts at 0
		// so for 4 beats in a bar, the value of this on each iteration will be 0, 1, 2, 3, 0 ...
		beatNum := i % numBeatsPerBar

		switch beatNum {
		case 0:
			fmt.Printf("BEAT: %d \n", beatNum)

			pop := buffer.Streamer(0, buffer.Len())

			//louderPop := &effects.Volume{
			//	Streamer: pop,
			//	Base:     1.5,
			//	Volume:   1,
			//	Silent:   false,
			//}
			speaker.Play(pop)

			fmt.Printf("BEAT TIME: %s \n", beatTime)

		default:
			fmt.Printf("BEAT: %d \n", beatNum)

			pop := buffer.Streamer(0, buffer.Len())
			speaker.Play(pop)

			fmt.Printf("BEAT TIME: %s \n", beatTime)
		}

		subdivisionsTicker := time.NewTicker(subdivisionsInterval)
		anotherIndex := 0

		for subTime := range subdivisionsTicker.C {
			subdivisionNum := anotherIndex % numSubdivisions

			fmt.Printf("SUB: %d \n", subdivisionNum)

			pop := buffer.Streamer(0, buffer.Len())
			speaker.Play(pop)
			fmt.Printf("SUBDIVISION TIME: %s \n", subTime)

			if subdivisionNum == (numSubdivisions - 2) {
				subdivisionsTicker.Stop()
				break
			}

			anotherIndex++
		}

		i++
	}
}
