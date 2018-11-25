package main

import "strings"

func GetFolderURL(url string, folderName string) string {
	if folderName == "" {
		return url
	}

	url = strings.TrimRight(url, "/")
	folderName = strings.TrimRight(folderName, "/")
	if !strings.HasPrefix(folderName, "/") {
		folderName = "/" + folderName
	}
	path := strings.Replace(folderName, "/", "/job/", -1)
	return url + path
}
