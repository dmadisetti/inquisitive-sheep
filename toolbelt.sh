#!/bin/sh

show_help(){
    echo "
    Why Hello there! You must be looking for help\n\
    \n\
    The Flags: \n\
    r - run \n\
    t - test \n\
    d - deploy \n\
    s - setup\n\
    p - ci push
    c - clean
    \n\
    Chain em together as you see fit \n\
    "
}

setup(){
    export FILE=go_appengine_sdk_linux_amd64-$(curl https://appengine.google.com/api/updatecheck | grep release | grep -o '[0-9\.]*').zip
    curl -O https://storage.googleapis.com/appengine-sdks/featured/$FILE
    unzip -q $FILE
}

run(){
    ./go_appengine/goapp serve;
}

try(){
    ./go_appengine/goapp build ./shoop;
    ./go_appengine/goapp test ./tests;
}

deploy(){
    echo $PASSWORD | go_appengine/appcfg.py --email=dylan.madisetti@gmail.com --passin update ./
}

push(){
    try || exit 1;
    git branch | grep "\*\ [^(master)\]" || {
        deploy;
    }
}

clean(){
    rm -rf go_appengine*;
}

while getopts "h?rtpscdx:" opt; do
    case "$opt" in
    h|\?)
        show_help
        ;;
    s)  setup
        ;;
    d)  deploy
        ;;
    r)  run
        ;;
    t)  try
        ;;
    p)  push
        ;;
    c)  clean
        ;;
    esac
done
