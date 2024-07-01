package api

import (
	"encoding/json"
	"github.com/vincent-petithory/dataurl"
	"os"
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

func NewMetaWithImageStr(description string, imageData string) *Meta {
	meta := Meta{
		Description: description,
		ImageData:   imageData,
	}
	return &meta
}

func NewMeta(description string) *Meta {
	meta := Meta{
		Description: description,
	}
	return &meta
}

func (m *Meta) LoadImage(file string) (bool, error) {
	if file != "" {
		image, err := os.ReadFile(file)
		if err != nil {
			return false, utils.AppendErrorInfo(err, "ReadFile")
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
	}
	err := json.Unmarshal([]byte(metaStr), m)
	if err != nil {
		m.Description = metaStr
	}
}

func (m *Meta) FetchAssetMeta(isHash bool, data string) string {
	response, err := fetchAssetMeta(isHash, data)
	if err != nil {
		return utils.MakeJsonResult(false, err.Error(), nil)
	}
	m.GetMetaFromStr(string(response.Data))
	return utils.MakeJsonResult(true, "", nil)
}
