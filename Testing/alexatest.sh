#!/bin/sh
echo "{\"speech\":\"`base64 -i question3.wav`\"}" > input
JSON2=`curl -v -X POST -d @input localhost:3000/alexa`
echo $JSON2 | cut -d '"' -f4 | base64 -d > answer.wav



