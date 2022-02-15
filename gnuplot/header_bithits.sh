#!/bin/bash

PLOT=/c/Users/Administrator/Dev/gnuplot/bin/gnuplot.exe
PLOTFILE="$(pwd)/header_bithits.plt"

(
	cd ../out

	for f in header_bithits_*.dat
	do
		d=$(echo ${f} | sed 's/\.dat/\.png/g')
    	#$PLOT -e "inputfile='${f}'; outputfile='${d}';" $PLOTFILE
    	$PLOT -e "inputfile='${f}'; outputfile='${d}';" $PLOTFILE

	done
)
