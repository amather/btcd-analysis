#!/bin/bash

PLOT=/c/Users/Administrator/Dev/gnuplot/bin/gnuplot.exe
PLOTFILE="$(pwd)/bitcount_header.plt"

(
	cd ../out
	$PLOT -e "inputfile='bitcount_header.dat'; outputfile='bitcount_header.png';" $PLOTFILE
)
