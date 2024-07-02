package services

import (
	"os/exec"
)

func GenerateBlocks(block int) (string, error) {
	cmd := exec.Command("/bin/bash", "/root/bitcoin-reg/autogen.sh")
	out, err := cmd.Output()
	if err != nil {
		return "", nil
	}
	return string(out), nil
}
