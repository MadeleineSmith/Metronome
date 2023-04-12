package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	//go:embed xylophone.mp3
	res embed.FS
)

type SubdivisionsFlag []int

func (flag *SubdivisionsFlag) String() string { return fmt.Sprint(*flag) }
func (flag *SubdivisionsFlag) Set(v string) error {
	subdivisionsStringSlice := strings.Split(v, ",")

	var subdivisionsIntSlice []int
	for _, subdivision := range subdivisionsStringSlice {
		integerSubdivision, _ := strconv.Atoi(subdivision)

		subdivisionsIntSlice = append(subdivisionsIntSlice, integerSubdivision)
	}

	*flag = append(*flag, subdivisionsIntSlice...)
	return nil
}

func main() {
	buffer := initializeBuffer()

	bpm, bpb, subdivisionsSlice, err := retrieveBeatsInput()
	if err != nil {
		log.Fatal(err)
	}

	initializeMetronome(bpm, bpb, subdivisionsSlice, buffer)
}

func initializeBuffer() *beep.Buffer {
	audioFile, _ := res.Open("xylophone.mp3")

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

func retrieveBeatsInput() (int, int, []int, error) {
	bpm := flag.Int("beats-per-minute", 15, "Number of Beats Per Minute")

	bpb := flag.Int("beats-per-bar", 4, "Number of Beats Per Bar")

	// the index of an item in this slice indicates the beat number:
	var subdivisionsSlice SubdivisionsFlag
	flag.Var(&subdivisionsSlice, "subdivisions", "The Subdivisions for each beat")

	flag.Parse()

	// note: the `beats-per-bar` flag is now technically unnecessary
	// 	(as the `subdivisions` flag tells you how many beats are in a bar),
	//	 but I quite like the explicitness
	if len(subdivisionsSlice) != *bpb {
		return 0, 0, nil, errors.New("num subdivisions must equal num beats per bar 😬")
	}

	return *bpm, *bpb, subdivisionsSlice, nil
}

func initializeMetronome(numBeatsPerMinute int, numBeatsPerBar int, subdivisionsSlice []int, buffer *beep.Buffer) {
	beatsInterval := time.Duration(float64(time.Minute) / float64(numBeatsPerMinute))
	beatsTicker := time.NewTicker(beatsInterval)

	fmt.Printf("Metronome is starting in: %s\n\n", beatsInterval)

	beatsIndex := 0
	for _ = range beatsTicker.C {
		// beatNum starts at 0
		// so for 4 beats in a bar, the value of this on each iteration will be 0, 1, 2, 3, 0 ...
		beatNum := beatsIndex % numBeatsPerBar

		if beatNum == 0 {
			// outputting at beginning of *bar*
			println("========")
		}

		humanBeatNum := beatNum + 1
		fmt.Printf("\u001B[1;33m%d\u001B[0m - \n", humanBeatNum)

		pop := buffer.Streamer(0, buffer.Len())
		speaker.Play(pop)

		// for each beat, I'm creating a new 'subdivisionsTicker' to keep track of the subdivisions
		numSubdivisions := subdivisionsSlice[beatNum]
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
