#!/usr/bin/env gnuplot

set terminal png size 1920,500 background rgb "#202020"
set output outputfile
set key off
set datafile separator " "
#set format y "%d"

set print "-"
print inputfile, " --> ", outputfile 


set border lc "yellow"
set yrange [0:30]
#set ytics 0, 500
#set mytics 2
set xtics 0, 32, 680 tc "yellow" 

#set boxwidth 0.05 relative
#set boxwidth 0.1 absolute
#set boxwidth -2
set style fill solid 1.0 noborder 
#set style fill solid 0.5 noborder 

#set arrow from graph 0,0.5 to graph 1,0.5 nohead lt rgb "blue" 

#plot inputfile using 2 with boxes lc rgb "#00aa00"
plot inputfile using 1:2 with points pointtype 1 lc rgb "#00aa00"


# https://edg.uchicago.edu/tutorials/pretty_plots_with_gnuplot/
