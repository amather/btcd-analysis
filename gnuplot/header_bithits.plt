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
set y2range [0:30]
#set y2range [0:4]
set xtics 0, 16, 680 tc "yellow" 

set grid xtics


#set style rect fc lt -1 fs solid 0.15 noborder
#set style rect back fc rgb "#ff0000" fs solid 0.15 noborder
set obj rect from 0, graph 0 to 32, graph 1 behind fc rgb "#404040" fs solid noborder 
set obj rect from 32, graph 0 to 288, graph 1 behind fc rgb "#606060" fs solid noborder 
set obj rect from 288, graph 0 to 544, graph 1 behind fc rgb "#404040" fs solid noborder 
set obj rect from 544, graph 0 to 576, graph 1 behind fc rgb "#606060" fs solid noborder 
set obj rect from 576, graph 0 to 608, graph 1 behind fc rgb "#404040" fs solid noborder 
set obj rect from 608, graph 0 to 640, graph 1 behind fc rgb "#606060" fs solid noborder 


set style fill solid 1.0 noborder 

#plot inputfile using 1:3 with points pointtype 1 lc rgb "#ff0000" axis x1y2, \
#     inputfile using 1:2 with points pointtype 1 lc rgb "#00aa00"
plot inputfile using 1:2 with points pointtype 5 lc rgb "#00aa00",\
     inputfile using 1:3 with points pointtype 13 lc rgb "#ff0000" axis x1y2



# https://edg.uchicago.edu/tutorials/pretty_plots_with_gnuplot/
