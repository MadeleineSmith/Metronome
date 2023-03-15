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

	bpm, bpb, subdivisionsMap := retrieveBeatsInput()

	initializeMetronome(bpm, bpb, subdivisionsMap, buffer)
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

func retrieveBeatsInput() (int, int, map[int]int) {
	reader := bufio.NewReader(os.Stdin)

	print("Beats Per Minute (default 15): ")
	bpmInput, _ := reader.ReadString('\n')
	bpm := 15
	if bpmInput != "\n" {
		bpm, _ = strconv.Atoi(strings.TrimRight(bpmInput, "\n"))
	}

	print("Beats Per Bar (default 4): ")
	bpbInput, _ := reader.ReadString('\n')
	bpb := 4
	if bpbInput != "\n" {
		bpb, _ = strconv.Atoi(strings.TrimRight(bpbInput, "\n"))
	}

	subdivisionsMap := make(map[int]int)
	for i := 0; i < bpb; i++ {
		fmt.Printf("Num subdivisions for beat number %d (default 4): ", i+1)

		numSubdivisionsInput, _ := reader.ReadString('\n')
		numSubdivisions := 4
		if numSubdivisionsInput != "\n" {
			numSubdivisions, _ = strconv.Atoi(strings.TrimRight(numSubdivisionsInput, "\n"))
		}

		subdivisionsMap[i] = numSubdivisions
	}
	println()

	return bpm, bpb, subdivisionsMap
}

func initializeMetronome(numBeatsPerMinute int, numBeatsPerBar int, subdivisionsMap map[int]int, buffer *beep.Buffer) {
	beatsInterval := time.Duration(float64(time.Minute) / float64(numBeatsPerMinute))
	beatsTicker := time.NewTicker(beatsInterval)

	beatsIndex := 0
	for _ = range beatsTicker.C {
		// beatNum starts at 0
		// so for 4 beats in a bar, the value of this on each iteration will be 0, 1, 2, 3, 0 ...
		beatNum := beatsIndex % numBeatsPerBar

		println("========")
		humanBeatNum := beatNum + 1
		fmt.Printf("%d - \n", humanBeatNum)

		pop := buffer.Streamer(0, buffer.Len())
		speaker.Play(pop)

		// for each beat, I'm creating a new 'subdivisionsTicker' to keep track of the subdivisions
		numSubdivisions := subdivisionsMap[beatNum]
		subdivisionsInterval := beatsInterval / time.Duration(numSubdivisions)

		subdivisionsTicker := time.NewTicker(subdivisionsInterval)
		subdivisionsIndex := 0

		for _ = range subdivisionsTicker.C {
			// again, subdivisionNum starts at 0
			subdivisionNum := subdivisionsIndex % numSubdivisions
			humanSubdivisionNum := subdivisionNum + 2

			fmt.Printf("  - %d \n", humanSubdivisionNum)

			pop := buffer.Streamer(0, buffer.Len())
			speaker.Play(pop)

			if subdivisionNum == (numSubdivisions - 2) {
				subdivisionsTicker.Stop()
				break
			}

			subdivisionsIndex++
		}

		beatsIndex++
	}
}
