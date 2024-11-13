#!/bin/sh

pref=202411
prev=

for d in $(ls -d ${pref}*); do
    for i in 1 2 3; do
        if ! [ -f $d/less$i.csv ] && [ -f $d/less$i.html ]; then
            ./bin/parse $d/less$i.html > $d/less$i.csv
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
