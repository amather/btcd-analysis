#!/bin/bash

PLOT=/c/Users/Administrator/Dev/gnuplot/bin/gnuplot.exe
PLOTFILE="$(pwd)/bitsum.plt"

(
	cd ../out
	$PLOT -e "inputfile='bitsum_hash.dat'; outputfile='bitsum_hash.png';" $PLOTFILE
)
