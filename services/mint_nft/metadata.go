package mint_nft

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"trade/api"
	"trade/utils"
)

type MetaData struct {
	ID   string `json:"id"`
	Meta struct {
		Name       string          `json:"name"`
		Attributes []api.Attribute `json:"attributes"`
	} `json:"meta"`
}

func GetAttributesFromFile(path string) (attributes []api.Attribute, err error) {
	var file *os.File
	file, err = os.OpenFile(path, os.O_RDONLY, 0664)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "OpenFile")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("close file error:", err)
		}
	}(file)
	var content []byte
	content, err = io.ReadAll(file)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "ReadAll")
	}
	var metaData MetaData
	err = json.Unmarshal(content, &metaData)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "Unmarshal")
	}
	return metaData.Meta.Attributes, nil
}
