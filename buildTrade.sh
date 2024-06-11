#!/bin/bash

start_time=$(date +%s)
currentPath=$(pwd)
currentFolderName="${currentPath##*/}"
specificString="trade"
configName="config.yaml"
tradeMarketPath="/root/trade"
binaryName="Trade"

if [ "$currentFolderName" = "$specificString" ]; then
    echo ""
    echo "Build Trade is in progress, please wait..."
#    if [ -d "$tradeMarketPath" ]; then
#        echo "$tradeMarketPath already exists."
#    else
#        echo "$tradeMarketPath does not exist, new folder will be created."
#        mkdir -p $tradeMarketPath
#    fi
    go build -o $tradeMarketPath/$binaryName ./main/main.go
    if [ -f "$tradeMarketPath/$configName" ]; then
        echo "$tradeMarketPath/$configName already exists."
    else
        echo "$tradeMarketPath/$configName does not exist, new example config file will be copied to that path."
        echo "Please edit config in $configName before running trade."
        cp ./main/config.yaml.example $tradeMarketPath/$configName
    fi
    end_time=$(date +%s)
    time_taken=$((end_time - start_time))
    echo "Build Trade finished, time cost: $time_taken seconds."
    echo ""
    echo "Tips:"
    echo "    1. go to $tradeMarketPath path:"
    echo "        cd $tradeMarketPath"
    echo ""
    echo "    2. Use this following command to start trade in daemon:"
    echo "        nohup $tradeMarketPath/$binaryName >> $tradeMarketPath/nohup.log 2>&1 &"
    echo ""
    echo "    3. Use this following command to check Trade is working properly:"
    echo "        lsof -i -n -P | grep LISTEN"
    echo ""
    echo "    4. Use this following command to stop the Trade:"
    echo "        kill -9 [PID]"
    echo ""
else
    echo "Wrong current directory, please run script in trade."
    # shellcheck disable=SC2162
    read -p "Press any key to continue..."
fi
