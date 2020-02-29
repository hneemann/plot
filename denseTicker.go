package plot

import (
	"gonum.org/v1/plot/vg"
	"math"
	"strconv"
	"time"
)

var _ SizeTicker = &DenseTicks{}

type stringSizer func(string) vg.Length

// DenseTicks creates tick marks as dense as possible
type DenseTicks struct {
	vks       int
	delta     float64
	fineStep  int
	stepWidth float64
	log       int

	axisSize    vg.Length
	stringSizer stringSizer
}

func (mt *DenseTicks) SetAxis(axis Axis, axisLen vg.Length, orientation orientation) {
	mt.axisSize = axisLen
	mt.stringSizer = createStringer(axis, orientation)
}

func createStringer(axis Axis, orientation orientation) stringSizer {
	if orientation == horizontal {
		return func(s string) vg.Length {
			return axis.Tick.Label.Font.Width(s + "M")
		}
	} else {
		return func(s string) vg.Length {
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

var _ SizeTicker = &DenseTimeTicks{}

// DenseTimeTicks creates tick marks as dense as possible
type DenseTimeTicks struct {
	// Format is used to format the date
	Format string

	// Time takes a float64 value and converts it into a time.Time.
	// If nil, UTC unix time is used.
	Time func(t float64) time.Time

	// Float takes a time.Time value and converts it into a float64
	// Must be the inverse of Time
	Float func(t time.Time) float64

	// Axis is required to transform the values from
	// the data coordinate system to the graphic coordinate system
	Axis Axis

	axisSize    vg.Length
	stringSizer stringSizer
}

func (t *DenseTimeTicks) SetAxis(axis Axis, axisLen vg.Length, orientation orientation) {
	t.axisSize = axisLen
	t.stringSizer = createStringer(axis, orientation)
	t.Axis = axis
}

type dateModifier func(time time.Time) time.Time

type incrementer struct {
	incr, norm dateModifier
}

var incrementerList = []incrementer{
	{daily(1), normTime},
	{daily(2), normDay},
	{weekly, normDay},
	{twoWeekly, normDay},
	{monthly(1), normDay},
	{monthly(2), normMonth},
	{monthly(3), normMonth},
	{monthly(4), normMonth},
	{monthly(6), normMonth},
	{yearly(1), normMonth},
	{yearly(2), normYear(2)},
	{yearly(5), normYear(5)},
	{yearly(10), normYear(10)},
	{yearly(20), normYear(20)},
}

func normYear(i int) dateModifier {
	return func(t time.Time) time.Time {
		y := (t.Year() / i) * i
		return time.Date(y, 1, 1, 0, 0, 0, 0, t.Location())
	}
}

func normMonth(t time.Time) time.Time {
	y := t.Year()
	return time.Date(y, 1, 1, 0, 0, 0, 0, t.Location())
}

func normDay(t time.Time) time.Time {
	y := t.Year()
	m := t.Month()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

func normTime(t time.Time) time.Time {
	y := t.Year()
	m := t.Month()
	d := t.Day()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func daily(days int) dateModifier {
	return func(t time.Time) time.Time {
		y := t.Year()
		m := t.Month()
		d := t.Day()
		return time.Date(y, m, d+days, 0, 0, 0, 0, t.Location())
	}
}
func weekly(t time.Time) time.Time {
	y := t.Year()
	m := t.Month()
	d := t.Day() + 7
	if d > 28 {
		d = 1
		m += 1
	}
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func twoWeekly(t time.Time) time.Time {
	y := t.Year()
	m := t.Month()
	d := t.Day() + 14
	if d > 28 {
		d = 1
		m += 1
	}
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func monthly(month time.Month) dateModifier {
	return func(t time.Time) time.Time {
		y := t.Year()
		m := t.Month() + month
		d := 1
		return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	}
}

func yearly(years int) func(t time.Time) time.Time {
	return func(t time.Time) time.Time {
		y := t.Year() + years
		return time.Date(y, 1, 1, 0, 0, 0, 0, t.Location())
	}
}

func (t *DenseTimeTicks) Ticks(min, max float64) []Tick {
	if t.Time == nil || t.Float == nil {
		t.Time = func(t float64) time.Time {
			return time.Unix(int64(t), 0).In(time.UTC)
		}
		t.Float = func(t time.Time) float64 {
			return float64(t.Unix())
		}
	}

	minTime := t.Time(min)
	size := t.stringSizer(minTime.Format(t.Format))

	index := 0
	for {
		t0 := incrementerList[index].norm(minTime)
		t1 := incrementerList[index].incr(t0)

		space := vg.Length(t.Axis.Norm(t.Float(t1))-t.Axis.Norm(t.Float(t0))) * t.axisSize

		if space > size || index == len(incrementerList)-1 {
			break
		}
		index++
	}

	incrementer := incrementerList[index]
	tickTime := incrementer.norm(minTime)

	for t.Float(tickTime) < min {
		tickTime = incrementer.incr(tickTime)
	}

	var ticker []Tick
	for {
		v := t.Float(tickTime)
		if v > max {
			break
		}
		ticker = append(ticker, Tick{
			Value: v,
			Label: t.Time(v).Format(t.Format),
		})
		tickTime = incrementer.incr(tickTime)
	}

	return ticker
}
