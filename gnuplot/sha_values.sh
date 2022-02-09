#!/bin/bash

PLOT=/c/Users/Administrator/Dev/gnuplot/bin/gnuplot.exe
PLOTFILE="$(pwd)/sha_values.plt"

(
	cd ../out

	for f in shavalues_be_round*.dat
	do
		d=$(echo ${f} | sed 's/\.dat/\.png/g')
    	#$PLOT -e "inputfile='${f}'; outputfile='${d}';" $PLOTFILE
    	$PLOT -e "inputfile='${f}'; outputfile='${d}';" $PLOTFILE

	done
)
