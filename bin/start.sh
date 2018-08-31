#!/bin/bash

trap "kill 0" EXIT

cd ~/go/src/github.com/agradl/golang-exercise && go build && ./golang-exercise &

listOptions(){
    echo "Choose from the following commands: "
    echo " hash <password>"
    echo " stats"
    echo " shutdown"
    echo " hash/<int>"
}

listOptions
while true; do 
    read cmd arg
    case $cmd in 
        hash )
            curl http://localhost:8080/hash --data "password=$arg";printf "\n";;
        stats )
            curl "http://localhost:8080/stats";printf "\n";;
        shutdown )
            curl "http://localhost:8080/shutdown";printf "\n";;
        hash/* )
            curl "http://localhost:8080/$cmd";printf "\n";;
        *)
            echo "$cmd is not a valid options, see below";listOptions;;
        esac
        sleep .2
        kill -0 $! &> /dev/null
        if [ "$?" != "0" ]; then
            exit 1
        fi
done