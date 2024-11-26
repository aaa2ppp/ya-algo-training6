#!/bin/sh

pref='20241[01]*'
prev=

for d in $(ls -d ${pref} | sort); do
    for i in 1 2 3 4; do
        if ! [ -f $d/less$i.csv ] && [ -f $d/less$i.html ]; then
            ./bin/parse $d/less$i.html > $d/less$i.csv
        fi

        if ! [ -f $d/asterisk$i ] && [ -f $prev/asterisk$i ]; then
            cp $prev/asterisk$i $d/asterisk$i
        fi 

        if ! [ -f $d/less$i.csv ] && [ -f $prev/less$i.csv ]; then
            cp $prev/less$i.csv $d/less$i.csv
        fi
    done
    prev=$d
done

target=${1:-$prev}
echo target=$target
(cd $target && ../bin/stat)
