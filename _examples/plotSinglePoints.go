package main

import (
	"fmt"
	"math/rand"
	"image/color"
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

// Draws a single point at a time inside the plot window.
// This is intended to demonstrate how to draw individual points on top of a 
// plot, not as the way to draw a plot. The 'plot' data series of functions are 
// intended for drawing plots. 
func randomDraw(plot *DataVisual.PlotWindow) {
	for {
	
		// Get the current UI window width/height.
		// Note that this is *not* the width/height of the plot.
		// we do this outside the loop since getting the width and height
		// of the ui window is a relatively slow process. Getting the 
		// width/height every 10 points is reasonable.
		width := plot.WindowWidth()
		height := plot.WindowHeight()
		
		// draw 10 points.
		for index:=0;index<10;index++ {
			
			// pick a draw color for the 10,000 points.
			color := color.RGBA {
				R: uint8(rand.Intn(256)),
				G: uint8(rand.Intn(256)),
				B: uint8(rand.Intn(256)),
				A: 0xff,
			}
			
			plot.DrawPoint(rand.Intn(width),rand.Intn(height),color)
			
		}
		
		
	}
}
