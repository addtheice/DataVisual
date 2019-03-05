Discontinued working on this. SVGo solved all of my problems and looks much 
nicer than this project which was, admittedly, more about playing around in
Go than solving a problem. Leaving this up in case others want to play with 
it or want to send me hints about better Go style.

DataVisual
==========

X11 based Go library for creating basic plots, graphs, and general 
data visualization.

The goals of this library is to create 'quick and dirty' plots for 'what does my 
data look like?' moments.

Easy of use is emphasized over performance or optional capabilities.

Think of this as you would fmt.Print() when doing fast value checks while 
developing, only for graphs and plots.

The current plan is to use reflection to allow for a single Plot() function 
which just takes the data and displays it in the most reasonable default manner.
Eventually there will be a system for adding new plot stylings, default 
overrides, and features like zooming, saving of plots, and overlays/multiplots. 

Ease of use will always take preference over extensibility.

Current Status
==============

Done:
-----

1. PlotWindows can be created and removed asynchronously.
2. PlotWindows can be waited on.
3. Hook in for displaying plots.

TODO:
-----

1. Static point plotting. "Here are some points. Plot them."
2. Static complex point plotting. "Here are some complex numbers. Plot them."
3. Static time series plotting. "Here is some points on a time line. Plot them."
4. Dynamic plotting of the above. "Here is a channel. Listen to it and plot the 
stuff I give you, updating the plot as it comes in."
