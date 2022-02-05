#!/bin/bash

PLOT=/c/Users/Administrator/Dev/gnuplot/bin/gnuplot.exe
PLOTFILE="$(pwd)/bitcount_hash.plt"

(
	cd ../out
	$PLOT -e "inputfile='bitcount_hash.dat'; outputfile='bitcount_hash.png';" $PLOTFILE
)
