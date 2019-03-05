package main

import (
	"fmt"
	"math/rand"
	"image/color"
	"image"
	"github.com/scriptandcompile/DataVisual"
)

func main() {
	plot, err := DataVisual.NewPlotWindow(400,400,"")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	// Draw random points in another go routine.
	go randomDraw(plot)
	
	// wait to close.
	DataVisual.WaitForAllPlotsClosed()
}

// Draws 10,000 points inside the plot window.
// This is intended to demonstrate how to draw some points on top of a plot, not
// as the way to draw a plot. The 'plot' data series of functions are intended
// for drawing plots. 
func randomDraw(plot *DataVisual.PlotWindow) {
	var pixelPoints [10000]image.Point
	
	
	for {
		// Get the current UI window width/height.
		// Note that this is *not* the width/height of the plot.
		// we do this outside the loop since getting the width and height
		// of the ui window is a relatively slow process.
		width := plot.WindowWidth()
		height := plot.WindowHeight()
		
		// Set 10,000 random points
		for index, _ := range pixelPoints {
			
			pixelPoints[index].X = rand.Intn(width)
			pixelPoints[index].Y = rand.Intn(height)
		}
		
		// pick a draw color for the 10,000 points.
		color := color.RGBA {
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 0xff,
		}
		
		// draw all 10,000 points.
		plot.DrawPoints(pixelPoints[:],color)
	}
}
