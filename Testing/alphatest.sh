#!/bin/sh
echo "{\"text\":\"How much potassium in a kiwano?\"}" > input
JSON2=`curl -v -X POST -d @input localhost:3001/alpha`
echo $JSON2

