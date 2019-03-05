// package DataVisual provides default plotting windows and functions for go.
// Uses a X11 binding in order to provide window handling which limits to X11
// compatible systems.
package DataVisual

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"strconv"
	"container/list"
)

type DVSystem struct {
	X *xgbutil.XUtil
	plotWindows *list.List
	plotCounter uint
	waitingForAllPlotWindowCloseFlag bool
	allPlotWindowsClosed chan bool
}

var system *DVSystem = nil

// Used when creating a plotWindows with a "" for the title we
// want to create a default name "Figure - #" where the number
// increments for each plot created.
func (system *DVSystem) getNextPlotWindowDefaultName() string {
	return "Figure - " + strconv.FormatUint(uint64(system.plotCounter),10)
}

// Add a plot and increments the figure number so that we still get figure
// incrementing even if we create a custom titled plot.
func (system *DVSystem) addPlotWindow(plot *PlotWindow)  {
	system.plotWindows.PushFront(plot)
	system.plotCounter++
}

// Removes a plot, closing it's underlying window in the process if it isn't 
// closed. 
func (system *DVSystem) removePlotWindow(plot *PlotWindow)  {
	plot.window.Destroy()
	plot.window = nil
	// Find this plot from the plot list and remove it.
	// I can't help but feel that generics would make this a lot nicer looking.
    for e := system.plotWindows.Front(); e != nil; e = e.Next() {
        if e.Value == plot {
        	system.plotWindows.Remove(e)
        	break
        }
    }
    
    // Send off the closed message for per plot close waiting.
    plot.sendCloseMessage()
    
    
    // If we have closed the last Plot send off the allPlotsClosed message.
	if system.plotWindows.Len() == 0 {
		sendAllPlotCloseMessage()	
	}
}

// If we are waiting on the plots to be closed then we signal that they have
// all been closed.
func sendAllPlotCloseMessage() {
	if system.waitingForAllPlotWindowCloseFlag == true {
		system.allPlotWindowsClosed <- true
	}
}

// DVSystem uses the (scary!) package global singleton system variable.
// Normally we won't be creating a DVSystem since we want to just single line
// create static or dynamic plots. But incase someone needs to fiddle with the
// internal DataVisual X11 binding components this will allow them to create
// as they need it instead of whenever the first plot is created.
func Initialize() (*DVSystem, error) {
	if system == nil {
		system = new(DVSystem)
		system.plotCounter = 1
		X, err := xgbutil.NewConn()
		if X == nil {
			return nil, err
		}
		system.X = X
		system.plotWindows = list.New()
		system.allPlotWindowsClosed = make(chan bool, 0)
		go system.pump()
	} 
	return system, nil
}


func CloseAllPlots() {

	// The package global system variable should actually exist by this point
	// since trying to close all plots when no plots exist makes little sense.
	// But just incase someone decides to do just that this will just return...
	// since they don't exist.
	if system == nil {
		return
	}
	
	// Close all of the plotWindows
	for element := system.plotWindows.Front(); element != nil; element = element.Next() {
		individualPlotWindow := element.Value.(*PlotWindow)
		individualPlotWindow.Close()
	}
	
	// send of the allPlotsClosed message
	if system.waitingForAllPlotWindowCloseFlag == true {
		system.allPlotWindowsClosed <- true
	}
	
}

func (system *DVSystem) pump () {
	xevent.Main(system.X)
}

func (system *DVSystem) setWaitingFlag(flag bool) {
	system.waitingForAllPlotWindowCloseFlag = flag
}

func WaitForAllPlotsClosed() {
	// The package global system variable should actually exist by this point
	// since waiting for all plots to close when no plots have been created 
	// makes little sense.
	// But just incase someone decides to do just that this will just return...
	// since they don't exist.
	if system == nil {
		return
	}
	system.setWaitingFlag(true)
	defer system.setWaitingFlag(false)

	if system.plotWindows.Len() == 0 {
		return
	}
	<- system.allPlotWindowsClosed
	
}



