package plot

import (
	"gonum.org/v1/plot/vg"
	"math"
	"strconv"
)

var _ SizeTicker = &DenseTicks{}

// DenseTicks creates tick marks as dense as possible
type DenseTicks struct {
	vks       int
	delta     float64
	fineStep  int
	stepWidth float64
	log       int

	axisSize    vg.Length
	stringSizer func(string) vg.Length
}

func (mt *DenseTicks) SetAxis(axis Axis, axisLen vg.Length, orientation orientation) {
	mt.axisSize = axisLen
	if orientation == horizontal {
		mt.stringSizer = func(s string) vg.Length {
			return axis.Tick.Label.Font.Width(s + "M")
		}
	} else {
		mt.stringSizer = func(s string) vg.Length {
			return axis.Tick.Label.Font.Size * 2
		}
	}
}

var finer = []float64{1, 0.5, 0.25, 0.2}
var logCorr = []int{0, 1, 2, 1}

func (mt *DenseTicks) Ticks(min, max float64) []Tick {

	mt.delta = max - min
	mt.log = int(math.Log10(mt.delta))
	mt.stepWidth = exp10(mt.log)
	mt.fineStep = 0
	mt.vks = int(math.Floor(math.Max(math.Log10(math.Abs(min)), math.Log10(math.Abs(max)))) + 1)
	if mt.vks < 1 {
		mt.vks = 1
	}
	if min < 0 {
		mt.vks++
	}

	mt.stepWidth *= 10
	mt.log++ // start to small

	for mt.checkTextWidth(mt.getPixels(mt.axisSize), mt.vks, mt.getNks()) {
		mt.inc()
	}
	mt.dec()

	mt.stepWidth *= finer[mt.fineStep]

	startTick := math.Ceil(min/mt.stepWidth) * mt.stepWidth

	eps := mt.stepWidth / 10000
	nks := mt.getNks()
	ticks := []Tick{}
	for startTick <= max+eps {
		ticks = append(ticks, Tick{
			Value: startTick,
			Label: strconv.FormatFloat(startTick, 'f', nks, 64),
		})
		startTick += mt.stepWidth
	}

	return ticks
}

const ZEROS = "0000000000000000000000000000000000000000000000000000000000000000000000000"

func (mt *DenseTicks) checkTextWidth(size vg.Length, vks, nks int) bool {
	s := ZEROS[:vks]
	if nks > 0 {
		s += "." + ZEROS[:nks]
	}
	width := mt.stringSizer(s)
	return size > width
}

func (mt *DenseTicks) getPixels(width vg.Length) vg.Length {
	return width * vg.Length(mt.stepWidth*finer[mt.fineStep]/mt.delta)
}

func (mt *DenseTicks) getNks() int {
	nks := logCorr[mt.fineStep] - mt.log
	if nks < 0 {
		return 0
	}
	return nks
}

func (mt *DenseTicks) inc() {
	mt.fineStep++
	if mt.fineStep == len(finer) {
		mt.stepWidth /= 10
		mt.log--
		mt.fineStep = 0
	}
}

func (mt *DenseTicks) dec() {
	mt.fineStep--
	if mt.fineStep < 0 {
		mt.stepWidth *= 10
		mt.log++
		mt.fineStep = len(finer) - 1
	}
}

func exp10(log int) float64 {
	exp10 := 1.0
	if log < 0 {
		for i := 0; i < -log; i++ {
			exp10 /= 10
		}
	} else {
		for i := 0; i < log; i++ {
			exp10 *= 10
		}
	}
	return exp10
}
