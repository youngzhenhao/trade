package services

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strconv"
	"trade/utils"
)

func GenerateBlocks(block int) (string, error) {
	cmd := exec.Command("/bin/bash", "/root/bitcoin-reg/autogen.sh")
	out, err := cmd.Output()
	if err != nil {
		return "", nil
	}
	return string(out), nil
}

func FaucetTransferBtc(address string, value float64) (string, error) {
	if address == "" || value == 0 {
		return "", errors.New("address or value is empty")
	}
	request := "'{\"" + address + "\":" + strconv.FormatFloat(value, 'f', -1, 64) + "}'"
	utils.CreateFile(filepath.Join("/root/bitcoin-reg/", "faucet.sh"), "bitcoin-cli --conf=/root/bitcoin-reg/bitcoin.conf -rpcwallet=wlt send "+request)
	cmd := exec.Command("/bin/bash", "/root/bitcoin-reg/faucet.sh")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
