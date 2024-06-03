$start_time = Get-Date
$currentPath = Get-Location
$currentFolderName = Split-Path -Path $currentPath -Leaf
$specificString = "trade"
$configName = "config.yaml"
$tradeMarketPath = "/root/tradeMarket"
if ($currentFolderName -eq $specificString) {
    Write-Host "build Trade is in progress, please wait..."
    if (Test-Path ./$tradeMarketPath) {
        Write-Host "$tradeMarketPath already exists."
    } else {
        Write-Host "$tradeMarketPath does not exist, new folder will be created."
        mkdir $tradeMarketPath
    }
    go build -o /root/tradeMarket/Trade ./main/main.go
    if (Test-Path ./$configName) {
        Write-Host "$fileName already exists."
    } else {
        Write-Host "$fileName does not exist, new example config file will be copied to current path."
        Copy-Item ./main/config.yaml.example ./$configName
    }
    $end_time = Get-Date
    $time_taken = $end_time - $start_time
    Write-Host "Time cost: $($time_taken.TotalSeconds) seconds."
} else {
	Write-Output "Wrong current directory, please run script in trade."
    pause
}
