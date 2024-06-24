package services

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	LndMainnetDataPath = "/root/mainnet-lit/.lnd/data"
)

func SnapshotToZipLast() {
	// lnd snapshot path
	filePaths := []string{
		LndMainnetDataPath + "/chain/bitcoin/mainnet/block_headers.bin",
		LndMainnetDataPath + "/chain/bitcoin/mainnet/neutrino.db",
		LndMainnetDataPath + "/chain/bitcoin/mainnet/reg_filter_headers.bin",
		LndMainnetDataPath + "/graph/mainnet/channel.db",
	}

	// zip path
	zipFilePath := "/root/neutrino/data.zip"
	//delete zip
	errRemove := os.Remove(zipFilePath)
	if errRemove != nil {
		log.Println("delete snapshot zip err:", errRemove)
	} else {
		log.Println("delete snapshot zip ok")
	}

	// create zip
	file, errzip := os.Create(zipFilePath)
	if errzip != nil {
		log.Println("Create Zip File err: ", errzip)
	}
	defer file.Close()

	// create zip write
	w := zip.NewWriter(file)
	defer w.Close()

	// lnd snapshot in zip
	for _, filePath := range filePaths {
		// get file name
		fileName := filepath.Base(filePath)
		// create file in zip
		f, err := w.Create(fileName)
		if err != nil {
			log.Println(" w.Create err: ", err)
		}

		// open file，in zip
		file, err := os.Open(filePath)
		if err != nil {
			log.Println("os.Open err: ", err)
		}
		defer file.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			log.Println("io.Copy err: ", err)
		}
	}
	log.Println("creat zip ok")
}
