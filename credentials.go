package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func DecryptFolderCredentials(url string, folderName string, username string, password string) error {
	folder, _ := GetFolder(url, folderName, username, password)
	credentials := folder.GetCredentials()
	script := GetDecryptScriptForCredentials(credentials)
	response := ExecuteGroovyScriptOnJenkins(script, url, username, password)
	fmt.Println(response)
	return nil
}

func ApplyFolderCredentials(url string, folderName string, username string, password string) error {
	var credentials Credentials

	err := json.NewDecoder(os.Stdin).Decode(&credentials)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	script := GetApplyScriptForCredentials(credentials, folderName)
	response := ExecuteGroovyScriptOnJenkins(script, url, username, password)
	fmt.Println(response)

	return nil
}
