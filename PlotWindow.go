package DataVisual

import (
	"errors"
	"image"
	"image/color"
	
	"github.com/BurntSushi/xgb/xproto"
	
	
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/mousebind"
)

type PlotWindow struct {
	window *xwindow.Window
	waitingForPlotWindowCloseFlag bool
	plotWindowClosed chan bool
	// Settings for axis drawing and labeling. 
	// None of these are used currently, but we will definitely need them 
	// shortly so I've included them.
	xAxisVisible bool
	yAxisVisible bool
	xAxisLabel string
	yAxisLabel string
}

// The width of the plot window. This is the width of the UI Window.
func (plot *PlotWindow) WindowWidth() int {
	plot.window.Geometry()
	return plot.window.Geom.Width()
}

// The height of the plot window. This is the height of the UI Window.
func (plot *PlotWindow) WindowHeight() int {
	plot.window.Geometry()
	return plot.window.Geom.Height()
}

// Set the window title.
func (plot *PlotWindow) SetTitle(Title string) {
	ewmh.WmNameSet(system.X, plot.window.Id, Title)
}



// Creates and displays a new plot window.
// Defaults to an all white background with the upper left corner at (0,0).
// If title is "" then the title will be an auto incrementing "Figure - #"
func NewPlotWindow(width, height int, title string) (*PlotWindow, error) {
	
	if width <= 0 {
		err := errors.New("Width is negative or 0, this is invalid.")
		return nil, err
	}
	
	if height <= 0 {
		err := errors.New("Height is negative or 0, this is invalid.")
		return nil, err
	}
	
	system, err := Initialize()
	
	if system == nil {
		return nil, err
	}
	
	resultantPlotWindow := new(PlotWindow)
	
	win, err := xwindow.Generate(system.X)
	if err != nil {
		return nil, err
	}
	
	resultantPlotWindow.window = win
	xpos := 0
	ypos := 0
	var backgroundColor uint32 
	backgroundColor = 0xffffffff
	// Create window, checked because we want to fail if this doesn't work.
	err = win.CreateChecked(
		system.X.RootWin(), 
		xpos, 
		ypos,
		width,
		height,
		xproto.CwBackPixel|xproto.CwEventMask,
		backgroundColor,
		xproto.EventMaskButtonRelease)
	if err != nil {
		return nil, err
	}
	
	// Gracefully removes event handling system on window close.
	// In addition it removes the plotWindow from the DVSystem which
	// sends off all the required plotWindow closing messages.
	win.WMGracefulClose(
		func(w *xwindow.Window) {
			// Detach all event handlers.
			// This should always be done when a window can no longer
			// receive events.
			xevent.Detach(w.X, w.Id)
			mousebind.Detach(w.X, w.Id)
			w.Destroy()
			
			// We need to remove the plot from the DVSystem list since it's now
			// no longer visible.
			system.removePlotWindow(resultantPlotWindow)
		})
		
	if title == "" {
		title = system.getNextPlotWindowDefaultName()
	}
	
	resultantPlotWindow.SetTitle(title)
	system.addPlotWindow(resultantPlotWindow)
	
	// Show the underlying window
	win.Map()
	
	// Hand back the PlotWindow
	return resultantPlotWindow, nil
}

// Tells the DVSystem to remove this plot (which closes it) and lets it handle 
// the cases of waiting for all plots to close and sending the per plot close
// message and the all plot closed message.
func (plot *PlotWindow) Close() {
	system.removePlotWindow(plot)
}

// Sends a closed message on the closed channel when we are waiting on closed.
func (plot *PlotWindow) sendCloseMessage() {
	if plot.waitingForPlotWindowCloseFlag == true {
		plot.plotWindowClosed <- true
	}
}

// A setter so we can do deferment on the waiting flag.
func (plot *PlotWindow) setWaitingFlag(flag bool) {
	plot.waitingForPlotWindowCloseFlag = flag
}

// If the window exists it waits on the closed channel until it receives 
// anything. This is simply a wait/signal system but it should be true simply 
// for the sake of consistancy.
func (plot *PlotWindow) WaitForClosed() {
	if plot == nil {
		return
	}
	// We are waiting but we want to flip the waiting flag back when we are done
	plot.setWaitingFlag(true)
	defer plot.setWaitingFlag(false)

	// If this window is nil then it's obviously all ready closed.
	if plot.window == nil {
		return
	}
	<- plot.plotWindowClosed
}

// Draws a array of points on the window in the window (rather then plot) 
// coordinate system.
// This routine is faster then the single point system because it performs only
// one buffer flip per Point array rather then per point.
//
// This routine should be avoided. This is here as an assistance function rather
// then the intended routine to be used for plotting. Prefer the 'Plot' data 
// functions over this function in almost all circumstances.
func (plot *PlotWindow) DrawPoints(points []image.Point, pixelColor color.Color) {
	if plot == nil {
		return
	}
	ximg, err := xgraphics.NewDrawable(system.X, xproto.Drawable(plot.window.Id))
	if err != nil {
		return
	}
	defer ximg.Destroy()

	// 
	for index, _ := range points {
		ximg.Set(points[index].X, points[index].Y, pixelColor)
	}
	
	ximg.XSurfaceSet(plot.window.Id)
	
	ximg.XDraw()
	
	ximg.XPaint(plot.window.Id)
	
}

// Draws a single point on the window in the window (rather then plot) 
// coordinate system.
// This routine is slow because it is essentially performing a buffer flip for 
// each point.
//
// This routine should be avoided. This is here as an assistance function rather
// then the intended routine to be used for plotting. Prefer the 'Plot' data 
// functions over this function in almost all circumstances.
func (plot *PlotWindow) DrawPoint(x, y int, pixelColor color.Color) {
	if plot == nil {
		return
	}
	ximg, err := xgraphics.NewDrawable(system.X, xproto.Drawable(plot.window.Id))
	if err != nil {
		return
	}
	defer ximg.Destroy()	
	
	ximg.Set(x,y,pixelColor)
	
	ximg.XSurfaceSet(plot.window.Id)
	
	ximg.XDraw()
	
	ximg.XPaint(plot.window.Id)

}
