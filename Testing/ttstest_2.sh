#!/bin/sh
echo "{\"text\":\"What is the area of Germany?\"}" > input
JSON2=`curl -s -X POST -d @input localhost:3003/tts`
echo $JSON2 | cut -d '"' -f4 | base64 -d > question3.wav

