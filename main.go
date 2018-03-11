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
	tchan := make(chan dirt.TelemetryAccessor)
	t, _ := dirt.GatherDefaultTelemetry(tchan)
	defer t.Close()

	timeLabel := createTelemetryBox("Total Time")
	lapTimeLabel := createTelemetryBox("Lap Time")
	lapDistance := createTelemetryBox("Lap Distance")

	speed := createTelemetryBox("Speed")
	throttle := createTelemetryGauge("Throttle")
	brake := createTelemetryGauge("Brake")

	rpm := createTelemetryBox("RPM")

	posx := createTelemetryBox("Position X")
	posy := createTelemetryBox("Position Y")
	posz := createTelemetryBox("Position Z")

	velx := createTelemetryBox("Velocity X")
	vely := createTelemetryBox("Velocity Y")
	velz := createTelemetryBox("Velocity Z")

	suslf := createTelemetryBox("Suspension LF")
	susrf := createTelemetryBox("Suspension RF")
	suslr := createTelemetryBox("Suspension LR")
	susrr := createTelemetryBox("Suspension RR")

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(3, 0, timeLabel),
			ui.NewCol(3, 1, lapTimeLabel),
			ui.NewCol(3, 1, lapDistance),
		),
		ui.NewRow(
			ui.NewCol(3, 0, speed),
			ui.NewCol(3, 1, throttle),
			ui.NewCol(3, 1, brake),
		),
		ui.NewRow(
			ui.NewCol(3, 0, rpm),
		),
		ui.NewRow(
			ui.NewCol(3, 0, velx),
			ui.NewCol(3, 1, vely),
			ui.NewCol(3, 1, velz),
		),
		ui.NewRow(
			ui.NewCol(3, 0, posx),
			ui.NewCol(3, 1, posy),
			ui.NewCol(3, 1, posz),
		),
	)

	ui.Body.Align()

	ui.Render(ui.Body)

	for {
		select {
		case <-quit:
			return

		case td := <-tchan:
			timeLabel.Text = fmt.Sprintf("%v s", getFieldValueIgnoreError(td, dirt.TotalTime))
			lapTimeLabel.Text = fmt.Sprintf("%v s", getFieldValueIgnoreError(td, dirt.LapTime))
			lapDistance.Text = fmt.Sprintf("%v m", getFieldValueIgnoreError(td, dirt.LapDistance))
			speed.Text = fmt.Sprintf("%v m/s", getFieldValueIgnoreError(td, dirt.Speed))

			posx.Text = fmt.Sprintf("%v m", getFieldValueIgnoreError(td, dirt.PositionX))
			posy.Text = fmt.Sprintf("%v m", getFieldValueIgnoreError(td, dirt.PositionY))
			posz.Text = fmt.Sprintf("%v m", getFieldValueIgnoreError(td, dirt.PositionZ))

			velx.Text = fmt.Sprintf("%v m", getFieldValueIgnoreError(td, dirt.VelocityX))
			vely.Text = fmt.Sprintf("%v m", getFieldValueIgnoreError(td, dirt.VelocityY))
			velz.Text = fmt.Sprintf("%v m", getFieldValueIgnoreError(td, dirt.VelocityZ))

			suslf.Text = fmt.Sprintf("%v m", getFieldValueIgnoreError(td, dirt.SuspensionPositionFrontLeft))
			susrf.Text = fmt.Sprintf("%v m", getFieldValueIgnoreError(td, dirt.SuspensionPositionFrontRight))
			suslr.Text = fmt.Sprintf("%v m", getFieldValueIgnoreError(td, dirt.SuspensionPositionBackLeft))
			susrr.Text = fmt.Sprintf("%v m", getFieldValueIgnoreError(td, dirt.SuspensionPositionBackRight))

			throttle.Percent = int(getFieldValueIgnoreError(td, dirt.ThrottleInput) * 100)
			brake.Percent = int(getFieldValueIgnoreError(td, dirt.BrakeInput) * 100)

			rpmValue := getFieldValueIgnoreError(td, dirt.EngineRate)
			rpm.Text = fmt.Sprintf("%v", rpmValue*9.5493)

			ui.Render(ui.Body)
			break
		}
	}
}

func getFieldValueIgnoreError(a dirt.TelemetryAccessor, field dirt.TelemetryField) float32 {
	v, _ := a.GetFieldValue(field)

	return v
}

func createTelemetryBox(param string) *ui.Par {
	p := ui.NewPar("?")

	p.Height = 3
	p.Width = 25
	p.TextFgColor = ui.ColorWhite
	p.BorderFg = ui.ColorCyan
	p.BorderLabel = param

	return p
}

func createTelemetryGauge(param string) *ui.Gauge {
	p := ui.NewGauge()

	p.Height = 3
	p.Width = 25
	p.BorderFg = ui.ColorCyan
	p.BorderLabel = param
	p.PercentColor = ui.ColorGreen

	return p
}
