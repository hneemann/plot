package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"log"
)

func main() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.X.Tick.Marker = &plot.DenseTicks{}
	p.Y.Tick.Marker = &plot.DenseTicks{}
	p.Add(plotter.NewGrid())
	p.Add(testData())

	_ = p.Save(100, 75, "denseTicks0.png")
	_ = p.Save(200, 150, "denseTicks1.png")
	_ = p.Save(400, 300, "denseTicks2.png")
	_ = p.Save(800, 600, "denseTicks3.png")
}

func testData() (plot.Plotter, plot.Plotter) {
	l, s, _ := plotter.NewLinePoints(plotter.XYs{
		{X: 0, Y: 0},
		{X: 3, Y: 200},
	})
	return l, s
}
