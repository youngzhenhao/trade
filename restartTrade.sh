#!/bin/bash
binaryName="Trade"
pid=$(pgrep -f $binaryName)
if [ -z "$pid" ]; then
  echo "No process named '$binaryName' was found, started a new one."
  nohup /root/trade/Trade >> /root/trade/nohup.log 2>&1 &
else
  kill $pid
  sleep 2
  nohup /root/trade/Trade >> /root/trade/nohup.log 2>&1 &
  echo "process '$binaryName' was restarted."
fi
