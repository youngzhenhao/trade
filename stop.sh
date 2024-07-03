#!/bin/bash
binaryName="Trade"
pid=$(pgrep -x $binaryName)
kill "$pid"
echo "process '$binaryName' was killed."
