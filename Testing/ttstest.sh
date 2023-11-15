#!/bin/sh
echo "{\"text\":\"What is the circumference of the earth?\"}" > input
JSON2=`curl -s -X POST -d @input localhost:3003/tts`
echo $JSON2

