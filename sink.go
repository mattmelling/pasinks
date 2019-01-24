package main

import (
	"fmt"
	"math"
	"strings"
	"github.com/sqp/pulseaudio"
	"github.com/godbus/dbus"
)

type Sink struct {
	client *pulseaudio.Client
	device *pulseaudio.Object
	pulse *dbus.Conn
}

func (this *Sink) Name () (string) {
	name, _ := this.device.String("Name")
	return name
}

func (this *Sink) State () (int) {
	state, _ := this.device.Uint32("State")
	return int(state)
}

func (this *Sink) Muted () (bool) {
	muted, _ := this.device.Bool("Mute")
	return muted
}

func (this *Sink) StateString () (string) {
	switch (1 + this.State()) {
	case 1:
		return "Running"
	case 2:
		return "Idle"
	case 3:
		return "Suspended"
	default:
		return "Invalid"
	}
}

func (this *Sink) SampleRate () (int) {
	rate, _ := this.device.Uint32("SampleRate")
	return int(rate)
}

func (this *Sink) Driver () (string) {
	driver, _ := this.device.String("Driver")
	return driver
}

func (this *Sink) VolumeRaw () (string) {
	vols, _ := this.device.ListUint32("Volume")
	var total float64
	for _, v := range vols {
		total += float64(v)
	}
	volume := int(math.Round((total / float64(len(vols) * 65536)) * 100))
	return fmt.Sprintf("%d", volume)
}

func (this *Sink) VolumePercent () (string) {
	return fmt.Sprintf("%s%%", this.VolumeRaw())
}

func (this *Sink) ActivePortName () (string) {
	activePort, _ := this.device.ObjectPath("ActivePort")
	object := this.pulse.Object("org.PulseAudio.Core1.DevicePort", activePort)
	name, _ := object.GetProperty("org.PulseAudio.Core1.DevicePort.Name")
	return name.Value().(string)
}

func (this *Sink) ActivePortDescription () (string) {
	activePort, _ := this.device.ObjectPath("ActivePort")
	object := this.pulse.Object("org.PulseAudio.Core1.DevicePort", activePort)
	desc, _ := object.GetProperty("org.PulseAudio.Core1.DevicePort.Description")
	return desc.Value().(string)
}

var formatStrings = map[string]func(*Sink)(string){
	"ActivePortDescription": func (sink *Sink) (string) { return sink.ActivePortDescription() },
	"ActivePortName": func (sink *Sink) (string) { return sink.ActivePortName() },
	"Driver": func (sink *Sink) (string) { return sink.Driver() },
	"Name": func (sink *Sink) (string) { return sink.Name() },
	"VolumeRaw": func (sink *Sink) (string) { return sink.VolumeRaw() },
	"VolumePercent": func (sink *Sink) (string) { return sink.VolumePercent() },
	"StateRaw": func (sink *Sink) (string) { return fmt.Sprintf("%d", sink.State()) },
	"StateString": func (sink *Sink) (string) { return sink.StateString() },
	"SampleRate": func (sink *Sink) (string) { return fmt.Sprintf("%d", sink.SampleRate()) },
	"Muted": func (sink *Sink) (string) { return fmt.Sprintf("%t", sink.Muted()) },
}

func (this *Sink) Format (format string) (string) {
	for k, v := range formatStrings {
		key := fmt.Sprintf("%%%s", k)
		if strings.Contains(format, key) {
			format = strings.Replace(format, key, v(this), -1)
		}
	}
	return format
}
