// main.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-gst/go-gst/gst"
)

func main() {
	gst.Init(nil)

	pipeline, err := CreatePipeline()
	if err != nil {
		fmt.Println("Failed to create pipeline:", err)
		return
	}

	pipeline.SetState(gst.StatePlaying)

	// Handle termination signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("Shutting down pipeline...")
		pipeline.SetState(gst.StateNull)
		os.Exit(0)
	}()

	// Monitor silence detection
	go MonitorSilence(pipeline)

	// Keep running
	for {
		time.Sleep(time.Second)
	}
}



func CreatePipeline() (*gst.Pipeline, error) {
	pipeline, err := gst.NewPipelineFromString(
		"autoaudiosrc name=microphone ! audioconvert ! audioresample ! level interval=100000000 ! queue max-size-time=2000000000 ! volume name=vol_ctrl ! autoaudiosink",
	)
	


	if err != nil {
		return nil, err
	}
	return pipeline, nil
}
func MonitorSilence(pipeline *gst.Pipeline) {
	volCtrl, _ := pipeline.GetElementByName("vol_ctrl")
	if volCtrl == nil {
		fmt.Println("Failed to find volume control element")
		return
	}

	bus := pipeline.GetPipelineBus()
	for {
		msg := bus.TimedPopFiltered(gst.ClockTimeNone, gst.MessageElement)
		if msg != nil && msg.Type() == gst.MessageElement {
			structMsg := msg.GetStructure() // Fix: Only one return value
			if structMsg != nil && structMsg.Name() == "level" {
				value, err := structMsg.GetValue("peak")
				if err != nil {
					fmt.Println("Failed to get 'peak' value:", err)
					continue
				}

				peak, ok := value.([]float64)
				if !ok || len(peak) == 0 {
					fmt.Println("Unexpected or empty 'peak' value")
					continue
				}

				if peak[0] < -50.0 {
					volCtrl.SetProperty("volume", 1.0) // User is silent, play test tone
				} else {
					volCtrl.SetProperty("volume", 0.0) // User is speaking, mute test tone
				}
			}
		}
	}
}
