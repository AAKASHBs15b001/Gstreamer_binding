Audio Mixer and Silence Detection using GStreamer

Overview

This project is a real-time audio processing application built with Golang and GStreamer. It takes two audio streams:

User Microphone Input

Test Tone (440Hz Sine Wave)

The system continuously monitors the microphone input and detects silence. If silence is detected, the test tone is played; otherwise, the tone is muted. This can be useful for applications such as hearing tests, voice-activated systems, or dynamic background audio generation.

Features

🎤 Real-time microphone input processing

🎵 Test tone generation using audiotestsrc

🔊 Dynamic volume control to mute/unmute test tone based on user speech

🎛 Uses audiomixer to combine streams

🖥 Runs on Mac/Linux with GStreamer bindings for Golang

Project Structure

/audio-mixer
│── main.go               # Entry point of the application
│── go.mod                # Go module dependencies
│── go.sum                # Package checksums
│── README.md             # Project documentation

Installation

1. Install GStreamer

Ensure you have GStreamer 1.0 installed with the necessary plugins.

For Mac (Homebrew)

brew install gst-plugins-base gst-plugins-good gst-plugins-bad gst-plugins-ugly gst-libav

For Ubuntu/Debian

sudo apt update && sudo apt install -y \
    gstreamer1.0-plugins-base \
    gstreamer1.0-plugins-good \
    gstreamer1.0-plugins-bad \
    gstreamer1.0-plugins-ugly \
    gstreamer1.0-libav \
    gstreamer1.0-tools

2. Install Golang

Make sure Go is installed (version 1.18+ recommended):

go version

If Go is not installed, follow the official installation guide.


3. Install Dependencies

go mod tidy

5. Run the Application

go run main.go

How It Works

1. Audio Pipeline Setup

The GStreamer pipeline is created as follows:

pipeline, err := gst.NewPipelineFromString(
		"autoaudiosrc name=microphone ! audioconvert ! audioresample ! level interval=100000000 ! queue max-size-time=2000000000 ! volume name=vol_ctrl ! autoaudiosink",
	)

2. Silence Detection

The level plugin monitors the peak audio levels from the microphone.

If the peak level falls below -50 dB, silence is assumed.

The vol_ctrl element is adjusted:

Silence detected → Unmute test tone (volume = 1.0)

User speaking → Mute test tone (volume = 0.0)

3. Volume Control Logic (from silence_detector.go)

if len(peak) > 0 && peak[0] < -50.0 {
    volCtrl.SetProperty("volume", 1.0) // User is silent, play test tone
} else {
    volCtrl.SetProperty("volume", 0.0) // User is speaking, mute test tone
}


