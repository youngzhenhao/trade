#!/bin/bash
binaryName="Trade"
pid=$(pgrep -f $binaryName)
if [ -z "$pid" ]; then
  echo "No process named '$binaryName' was found, started a new one."
  nohup /root/trade/Trade >> /root/trade/nohup.log 2>&1 &
  pidNew=$(pgrep -f $binaryName)
  echo "$binaryName's new pid is $pidNew."
else
  # shellcheck disable=SC2086
  kill $pid
  nohup /root/trade/Trade >> /root/trade/nohup.log 2>&1 &
  echo "process '$binaryName' was restarted."
  pidNew=$(pgrep -f $binaryName)
  echo "$binaryName's new pid is $pidNew."
fi
