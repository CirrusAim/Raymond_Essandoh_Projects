#! /bin/bash
set -e
prefix="ABC"
three="localhost:50081,localhost:50082,localhost:50083"
five="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085"
server=$three
if [[ $3 == 5 ]]; then
    server=$five
fi
for i in $( seq 1 $1 )
do
    val=""
    for j in $( seq 1 $2 )
    do
        val="${val},${prefix}${i}${j}"
    done
    modified="${val:1}"
    ./paxosclient -addrs="${server}" -clientRequest="${modified}" -clientId "${i}" &
done