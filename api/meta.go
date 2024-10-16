package api

import (
	"encoding/json"
	"fmt"
	"github.com/vincent-petithory/dataurl"
	"os"
	"path/filepath"
	"strings"
	"trade/models"
	"trade/utils"
)

type Meta struct {
	Acronym     string `json:"acronym,omitempty"`
	Description string `json:"description,omitempty"`
	ImageData   string `json:"image_data,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	GroupName   string `json:"groupName,omitempty"`
}

func NewMeta(description string) *Meta {
	meta := Meta{
		Description: description,
	}
	return &meta
}

func NewMetaWithGroupName(description string, groupName string) *Meta {
	meta := Meta{
		Description: description,
		GroupName:   groupName,
	}
	return &meta
}

func NewMetaWithImageStr(description string, imageData string) *Meta {
	meta := Meta{
		Description: description,
		ImageData:   imageData,
	}
	return &meta
}

func (m *Meta) LoadImageByByte(image []byte) (bool, error) {
	if len(image) == 0 {
		fmt.Println("image data is nil")
		return false, fmt.Errorf("image data is nil")
	}
	imageStr := dataurl.EncodeBytes(image)
	m.ImageData = imageStr
	return true, nil
}

func (m *Meta) LoadImage(file string) (bool, error) {
	if file != "" {
		image, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("open image file is error:", err)
			return false, err
		}
		imageStr := dataurl.EncodeBytes(image)
		m.ImageData = imageStr
	}
	return true, nil
}

func (m *Meta) ToJsonStr() string {
	metastr, _ := json.Marshal(m)
	return string(metastr)
}

func (m *Meta) GetMetaFromStr(metaStr string) {
	if metaStr == "" {
		m.Description = "This asset has no meta."
		return
	}

	first := metaStr[:1]
	end := metaStr[len(metaStr)-1:]
	var s string
	if first == "\"" && end == "\"" {
		s = metaStr[1 : len(metaStr)-1]
	} else {
		s = metaStr
	}
	err := json.Unmarshal([]byte(s), m)
	if err != nil {
		m.Description = s
	}
}

func (m *Meta) SaveImage(dir string, name string) bool {
	if m.ImageData == "" {
		return false
	}
	dataUrl, err := dataurl.DecodeString(m.ImageData)
	if err != nil {
		return false
	}
	ContentType := dataUrl.MediaType.ContentType()
	datatype := strings.Split(ContentType, "/")
	if datatype[0] != "image" {
		fmt.Println("is not image dataurl")
		return false
	}
	formatName := strings.Split(name, ".")
	file := filepath.Join(dir, formatName[0]+"."+datatype[1])
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("create new image error:", err)
		return false
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	_, err = f.Write(dataUrl.Data)
	if err != nil {
		fmt.Println("Write data fail:", err)
		return false
	}
	return true
}

func (m *Meta) GetImage() []byte {
	if m.ImageData == "" {
		return nil
	}
	dataUrl, err := dataurl.DecodeString(m.ImageData)
	if err != nil {
		return nil
	}
	ContentType := dataUrl.MediaType.ContentType()
	datatype := strings.Split(ContentType, "/")
	if datatype[0] != "image" {
		fmt.Println("is not image dataurl")
		return nil
	}
	return dataUrl.Data
}

func (m *Meta) FetchAssetMeta(isHash bool, data string) string {
	response, err := fetchAssetMeta(isHash, data)
	if err != nil {
		return utils.MakeJsonErrorResult(models.FetchAssetMetaErr, err.Error(), nil)
	}
	m.GetMetaFromStr(string(response.Data))
	return utils.MakeJsonErrorResult(models.SUCCESS, "", nil)
}
