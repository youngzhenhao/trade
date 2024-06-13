#!/bin/bash
binaryName="Trade"
pid=$(pgrep -f $binaryName)
kill "$pid"
echo "process '$binaryName' was killed."
