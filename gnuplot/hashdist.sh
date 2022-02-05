#!/bin/bash

PLOT=/c/Users/Administrator/Dev/gnuplot/bin/gnuplot.exe
PLOTFILE="$(pwd)/hashdist.plt"

(
	cd ../out

	for i in $(seq 1 8)
	do
    	$PLOT -e "inputfile='hashdist_le_${i}.dat'; outputfile='hashdist_le_${i}.png';" $PLOTFILE
    	$PLOT -e "inputfile='hashdist_be_${i}.dat'; outputfile='hashdist_be_${i}.png';" $PLOTFILE

	done
)
