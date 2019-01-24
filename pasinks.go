package main

import (
	"fmt"
	"flag"
	"github.com/sqp/pulseaudio"
)

type Config struct {
	name string
	state int
	format string
}

func getSinks (client *pulseaudio.Client, config Config) ([]Sink) {
	sinkList, err := client.Core().ListPath("Sinks")
	if err != nil {
		panic(err)
	}

	pulse := GetPulseaudioBus()

	var sinks []Sink
	for _, sink := range sinkList {
		s := Sink{
			device: client.Device(sink),
			client: client,
			pulse: pulse,
		}
		if config.name != "" && s.Name() != config.name {
			continue
		}
		if config.state != -1 && s.State() != config.state {
			continue
		}
		sinks = append(sinks, s)
	}

	return sinks
}

func main() {
	format := flag.String("format", "%Name: %VolumePercent", "Output format")
	name := flag.String("name", "", "Limit output to devices named --name")
	running := flag.Bool("running", false, "Limit output to devices in RUNNING state")
	flag.Parse()

	config := Config {
		name: *name,
		state: -1,
	}
	if *running {
		config.state = 0
	}

	pulse, err := pulseaudio.New()
	if err != nil {
		panic(err)
	}

	sinks := getSinks(pulse, config)
	for _, sink := range sinks {
		fmt.Println(sink.Format(*format))
	}
}
