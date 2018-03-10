package main

import (
	"fmt"

	ui "github.com/gizak/termui"
	dirt "github.com/kdelorey/cmtelemetry"
)

func main() {
	err := ui.Init()

	if err != nil {
		panic(err)
	}

	defer ui.Close()

	quit := make(chan struct{})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		close(quit)
		ui.StopLoop()
	})

	go showTelemetry(quit)

	ui.Loop()
}

func showTelemetry(quit chan struct{}) {
	t, _ := dirt.StartDefaultTelemetry()
	defer t.Close()

	timeLabel := createTelemetryBox("Total Time", 0, 0)
	lapTimeLabel := createTelemetryBox("Lap Time", 0, 4)
	lapDistance := createTelemetryBox("Lap Distance", 26, 4)
	speed := createTelemetryBox("Speed", 0, 8)

	posx := createTelemetryBox("Position X", 0, 12)
	posy := createTelemetryBox("Position Y", 26, 12)
	posz := createTelemetryBox("Position Z", 52, 12)

	velx := createTelemetryBox("Velocity X", 0, 16)
	vely := createTelemetryBox("Velocity Y", 26, 16)
	velz := createTelemetryBox("Velocity Z", 52, 16)

	suslf := createTelemetryBox("Suspension LF", 0, 21)
	susrf := createTelemetryBox("Suspension RF", 26, 21)
	suslr := createTelemetryBox("Suspension LR", 0, 25)
	susrr := createTelemetryBox("Suspension RR", 26, 25)

	renderFunc := func() {
		ui.Render(
			timeLabel,
			lapTimeLabel,
			lapDistance,
			speed,
			posx,
			posy,
			posz,
			velx,
			vely,
			velz,
			suslf,
			suslr,
			susrf,
			susrr)
	}

	renderFunc()

	for {
		select {
		case <-quit:
			return

		case <-t.OnFrameReceived():
			timeLabel.Text = fmt.Sprintf("%v s", t.GetFieldValue(dirt.TotalTime))
			lapTimeLabel.Text = fmt.Sprintf("%v s", t.GetFieldValue(dirt.LapTime))
			lapDistance.Text = fmt.Sprintf("%v m", t.GetFieldValue(dirt.LapDistance))
			speed.Text = fmt.Sprintf("%v m/s", t.GetFieldValue(dirt.Speed))

			posx.Text = fmt.Sprintf("%v m", t.GetFieldValue(dirt.PositionX))
			posy.Text = fmt.Sprintf("%v m", t.GetFieldValue(dirt.PositionY))
			posz.Text = fmt.Sprintf("%v m", t.GetFieldValue(dirt.PositionZ))

			velx.Text = fmt.Sprintf("%v m", t.GetFieldValue(dirt.VelocityX))
			vely.Text = fmt.Sprintf("%v m", t.GetFieldValue(dirt.VelocityY))
			velz.Text = fmt.Sprintf("%v m", t.GetFieldValue(dirt.VelocityZ))

			suslf.Text = fmt.Sprintf("%v m", t.GetFieldValue(dirt.SuspensionPositionFrontLeft))
			susrf.Text = fmt.Sprintf("%v m", t.GetFieldValue(dirt.SuspensionPositionFrontRight))
			suslr.Text = fmt.Sprintf("%v m", t.GetFieldValue(dirt.SuspensionPositionBackLeft))
			susrr.Text = fmt.Sprintf("%v m", t.GetFieldValue(dirt.SuspensionPositionBackRight))

			renderFunc()
			break
		}
	}
}

func createTelemetryBox(param string, x int, y int) *ui.Par {
	p := ui.NewPar("?")

	p.Height = 3
	p.Width = 25
	p.X = x
	p.Y = y
	p.TextFgColor = ui.ColorWhite
	p.BorderFg = ui.ColorCyan
	p.BorderLabel = param

	return p
}
