$start_time = Get-Date
$currentPath = Get-Location
$currentFolderName = Split-Path -Path $currentPath -Leaf
$specificString = "trade"
$configName = "config.yaml"
$tradeMarketPath = "/root/trade"
$binaryName = "Trade"
if ($currentFolderName -eq $specificString) {
    Write-Host ""
    Write-Host "Build Trade is in progress, please wait..."
#    if (Test-Path $tradeMarketPath) {
#        Write-Host "$tradeMarketPath already exists."
#    } else {
#        Write-Host "$tradeMarketPath does not exist, new folder will be created."
#        mkdir $tradeMarketPath
#    }
    go build -o $tradeMarketPath/$binaryName ./main/main.go
    if (Test-Path $tradeMarketPath/$configName) {
        Write-Host "$tradeMarketPath/$configName already exists."
    } else {
        Write-Host "$tradeMarketPath/$configName does not exist, new example config file will be copied to that path."
        Write-Host "Please edit config in $configName before running trade."
        Copy-Item ./main/config.yaml.example $tradeMarketPath/$configName
    }
    $end_time = Get-Date
    $time_taken = $end_time - $start_time
    Write-Host "Build Trade finished, time cost: $($time_taken.TotalSeconds) seconds."
    Write-Host ""
    Write-Host "Tips:"
    Write-Host "    1.go to $tradeMarketPath path:"
    Write-Host "        cd $tradeMarketPath"
    Write-Host ""
    Write-Host "    2.Use this fllowing command to start trade in daemon:"
    Write-Host "        nohup $tradeMarketPath/$binaryName >> $tradeMarketPath/nohup.log 2>&1 &"
    Write-Host ""
    Write-Host "    3.Use this fllowing command to check Trade is working properly:"
    Write-Host "        lsof -i -n -P | grep LISTEN"
    Write-Host ""
    Write-Host "    4.Use this fllowing command to stop the Trade:"
    Write-Host "        kill -9 [PID]"
    Write-Host ""
} else {
    Write-Output "Wrong current directory, please run script in trade."
    pause
}

