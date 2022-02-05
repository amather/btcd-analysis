#!/usr/bin/env gnuplot

#set terminal png size 3000,700 background rgb "#202020"
set terminal png size 1920,700 background rgb "#202020"
set output outputfile
set key off
set datafile separator " "
set format y "%x"

set print "-"
print inputfile, " --> ", outputfile 

set border lc "yellow"
set xtics tc "yellow"

# 1073741824 == 2^32/4
#set ytics 0, 1073741824
# 536870912 == 2^32/8
#set ytics 0, 536870912  
# 268435456 == 2^32/16
set ytics 0, 268435456  tc "yellow"
# 67108864 == /64
#set mytics 0, 67108864 
set mytics 4


# Y-minor tics
set style line 81 lt 1 lc rgb "blue"
set grid mytics back ls 81 lw 0.5

# Y-major tics
set style line 80 lt 1 lc rgb "#555555"
set grid ytics behind ls 80


#plot "nonce.dat"  with points pointtype 0 lc "#00aa00"
plot inputfile with points pointtype 0 lc "#00aa00"


# https://edg.uchicago.edu/tutorials/pretty_plots_with_gnuplot/
