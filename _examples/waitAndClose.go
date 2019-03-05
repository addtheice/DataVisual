package main

import (
	"github.com/scriptandcompile/DataVisual"
	"time"
	"fmt"
)

// Demonstrates how to create a basic empty plot window. This is probably
// not how you want to create a plot window. You will probably want one of the
// static or dynamic plot windows instead (which uses the PlotWindow structure)
// but I've included this since this is what I created at the beginning.
func main() {
	
	// Creates a new plot window and displays it.
	// The plot window will be at (0,0) and be 300 wide by 400 tall.
	// The background will be entirely white. These are the default settings.
	// Since we haven't specified a title of the plot it will default to 
	// "Figure - 1". Each subsequent untitled plot will increment. 
	// "Figure - 2", "Figure - 3", etc.
	plot, err := DataVisual.NewPlotWindow(300,400, "")
	if plot == nil {
		fmt.Print(err)
		return
	}
	
	// delay for a few seconds then close the plots.
	go delayCloseAll()
	
	// Wait for the plots to close.
	// This will just wait synchronously for all the plots to close *while* 
	// event handling still occurs. 
	//
	// Normal uses for this will be to open a plot with some information then
	// wait for the user to close it before continuing on with further 
	// processing.
	//
	// since this is channel safe it can also be used with the dynamic plots
	// to display a plot and then wait for it to close while another channel
	// sends data to the plot to display.
	DataVisual.WaitForAllPlotsClosed()
}

// just delay for a few seconds and then close all the plots
func delayCloseAll() {
	time.Sleep(time.Second*15)
	DataVisual.CloseAllPlots()
}
	
