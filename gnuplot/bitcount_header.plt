#!/usr/bin/env gnuplot

set terminal png size 1920,900 background rgb "#202020"
set output outputfile
set key off
set datafile separator " "
#set format y "%d"

set print "-"
print inputfile, " --> ", outputfile 


set border lc "yellow"
set yrange [206:434]
#set yrange [0:640]
#set xrange [620000:630000]
set ytics 0,8 
#set mytics 2
#set xtics 0, 32, 256 tc "yellow" 


set arrow from graph 0,0.5 to graph 1,0.5 nohead front lt rgb "blue" 

plot inputfile using 1:2 with points pointtype 0 lc rgb "#00aa00",\
     inputfile using 1:3 with points pointtype 0 lt rgb "#DA6A1A",\
     inputfile using 1:4 with points pointtype 0 lt rgb "#aaaaaa",\
     inputfile using 1:5 with points pointtype 0 lt rgb "#aaaaaa",\

#plot inputfile using 1:2 with points lc rgb "#00aa00"


# https://edg.uchicago.edu/tutorials/pretty_plots_with_gnuplot/
