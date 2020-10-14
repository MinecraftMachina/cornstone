package curseforge

import (
	"cornstone/util"
	"fmt"
)

var url = "https://addons-ecs.forgesvc.net/api/v2/addon/"
var client = util.DefaultClient.New()

func QueryAddon(id int) (*Addon, error) {
	addon := Addon{}
	if _, err := client.Get(fmt.Sprintf("%s%d", url, id)).ReceiveSuccess(&addon); err != nil {
		return nil, err
	}
	return &addon, nil
}

func GetAddonFileDownloadUrl(addonId int, fileId int) (string, error) {
	result := make([]byte, 0)
	if _, err := client.Get(fmt.Sprintf("%s%d/file/%d/download-url", url, addonId, fileId)).
		ByteResponse().
		ReceiveSuccess(&result); err != nil {
		return "", err
	}
	return string(result), nil
}
