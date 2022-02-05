#!/bin/bash

PLOT=/c/Users/Administrator/Dev/gnuplot/bin/gnuplot.exe
PLOTFILE="$(pwd)/nonce.plt"

#$PLOT -e "inputfile='nonce_le.dat'; outputfile='nonce_le.png';" nonce.gnuplot
#$PLOT -e "inputfile='nonce_be.dat'; outputfile='nonce_be.png';" nonce.gnuplot
(
	cd ../out
    $PLOT -e "inputfile='nonce_le.dat'; outputfile='nonce_le.png';" $PLOTFILE
    $PLOT -e "inputfile='nonce_be.dat'; outputfile='nonce_be.png';" $PLOTFILE
)
