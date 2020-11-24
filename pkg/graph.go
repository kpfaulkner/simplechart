package pkg

import (
"fmt"
"github.com/wcharczuk/go-chart"
"os"
"strconv"
"time"
)

type GraphPoint struct {
	Timestamp  time.Time
	Count int
}

func parseInt(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}

func parseFloat64(str string) float64 {
	v, _ := strconv.ParseFloat(str, 64)
	return v
}

func convertToSeparateDataSlices(data []GraphPoint) ([]time.Time, []float64) {
	var xvalues []time.Time
	var yvalues []float64

	for _, m := range data {
		xvalues = append(xvalues, m.Timestamp)
		yvalues = append(yvalues, float64(m.Count))
	}
	return xvalues, yvalues
}

// DrawChart takes a map of strings to GraphPoint array.
// The string is the name of the data being graphed.
// The GraphPoint array is the time/value pair.
func DrawChart(data map[string][]GraphPoint, filename string) {

	// list of stroke colours.
	//strokeColours := []drawing.Color {chart.ColorWhite,chart.ColorBlue,chart.ColorCyan,chart.ColorGreen,chart.ColorRed,chart.ColorOrange,chart.ColorYellow,chart.ColorBlack,chart.ColorLightGray,chart.ColorAlternateBlue,chart.ColorAlternateGreen,chart.ColorAlternateGray,chart.ColorAlternateYellow,chart.ColorAlternateLightGray}
	seriesList := []chart.Series{}
	maxY := 0.0

	strokeNumber := 0
	for k,v := range data {
		xValues, yValues := convertToSeparateDataSlices(v)

		// get maxY for graphing later.
		for _,y := range yValues {
			if y > maxY {
				maxY = y
			}
		}

		errorSeries := chart.TimeSeries{
			Name: k,
			/*
				Style: chart.Style{
					Show:        true,
					StrokeColor: strokeColours[strokeNumber],
					StrokeWidth: chart.Disabled,
					DotWidth:5,

					//FillColor:   chart.ColorBlue.WithAlpha(100),

				}, */
			XValues: xValues,
			YValues: yValues,

		}
		strokeNumber++
		if strokeNumber >= 14 {
			strokeNumber = 0
		}

		seriesList = append(seriesList, errorSeries)
	}

	// give it some headroom
	maxY += 20.0
	graph := chart.Chart{
		/*
				Width:  2280,
				Height: 720,

				Background: chart.Style{
					Padding: chart.Box{
						Top: 50,
					},
				},
				Canvas:chart.Style {
		  		  FillColor:chart.ColorBlack,
				},*/

		YAxis: chart.YAxis{
			Name:      "Count",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
			TickStyle: chart.Style{
				TextRotationDegrees: 45.0,
			},
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%d", int(v.(float64)))
			},
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: maxY,
			},
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: chart.TimeValueFormatterWithFormat("2006-01-02.15-04-05"),
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 1.0,
			},
			//GridLines: releases(),
		},
		Series: seriesList,
	}

	graph.Elements = []chart.Renderable{chart.LegendLeft(&graph)}
	//, chart.Style{ FillColor:drawing.Color{R: 50, G: 120, B: 203, A: 255}})}

	// output to file... ?

	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("unable to create file\n")
	}

	graph.Render(chart.PNG, f)
}

