#!/usr/bin/env gnuplot

set terminal png size 1920,900 background rgb "#202020"
set output outputfile
set key off
set datafile separator " "

set print "-"
print inputfile, " --> ", outputfile 


set border lc "yellow"
set yrange [32:224]
#set yrange [0:257]
#set xrange [620000:630000]
set ytics 0, 32
#set mytics 2
#set xtics 0, 32, 256 tc "yellow" 


set arrow from graph 0,0.5 to graph 1,0.5 nohead front lt rgb "blue" 

plot inputfile using 1:2 with points pointtype 0 lc rgb "#00aa00"
#plot inputfile using 1:2 with points lc rgb "#00aa00"


# https://edg.uchicago.edu/tutorials/pretty_plots_with_gnuplot/
