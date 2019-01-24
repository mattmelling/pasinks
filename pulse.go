package main

import (
	"github.com/godbus/dbus"
)

func GetPulseaudioBus () (*dbus.Conn) {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	l, err := conn.Object("org.PulseAudio1", "/org/pulseaudio/server_lookup1").GetProperty("org.PulseAudio.ServerLookup1.Address")
	pulse, err := dbus.Dial((l.Value().(string)))

	e := pulse.Auth(nil)
	if e != nil {
		pulse.Close()
		panic(e)
	}
	
	if err != nil {
		panic(err)
	}
	return pulse
}
