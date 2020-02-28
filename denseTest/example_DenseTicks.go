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

	err = p.Save(800, 600, "denseTicks.png")
	if err != nil {
		panic(err)
	}
}

func testData() (plot.Plotter, plot.Plotter) {
	l, s, _ := plotter.NewLinePoints(plotter.XYs{
		{X: 0, Y: 0},
		{X: 3, Y: 200},
	})
	return l, s
}
