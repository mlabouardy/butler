package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type JenkinsFolder struct {
	Properties struct {
		CredentialProperty struct {
			DomainCredentials struct {
				Class string `xml:"class,attr"`
				Entry struct {
					Credentials Credentials `xml:"java.util.concurrent.CopyOnWriteArrayList"`
				} `xml:"entry"`
			} `xml:"domainCredentialsMap"`
		} `xml:"com.cloudbees.hudson.plugins.folder.properties.FolderCredentialsProvider_-FolderCredentialsProperty"`
	} `xml:"properties"`
}

type Credentials struct {
	UsernamePassword []UsernamePasswordCredential `xml:"com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl" json:"userpass"`
	SecretFile       []SecretFileCredential       `xml:"org.jenkinsci.plugins.plaincredentials.impl.FileCredentialsImpl" json:"secretfile"`
}

type UsernamePasswordCredential struct {
	Plugin      string `xml:"plugin,attr" json:"plugin"`
	ID          string `xml:"id" json:"id"`
	Description string `xml:"description" json:"description"`
	Username    string `xml:"username" json:"username"`
	Password    string `xml:"password" json:"password"`
}

type SecretFileCredential struct {
	Plugin             string `xml:"plugin,attr"`
	ID                 string `xml:"id" json:"id"`
	Description        string `xml:"description" json:"description"`
	FileName           string `xml:"fileName" json:"fileName"`
	SecretBytes        string `xml:"secretBytes" json:"secretBytes"`
	EncodedSecretBytes string `xml:"encodedSecretBytes" json:"encodedSecretBytes"`
}

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

func parseJenkinsFolder(xmlInput []byte) JenkinsFolder {
	var folder JenkinsFolder
	xmlInput = []byte(strings.Trim(string(xmlInput), "<?xml version='1.1' encoding='UTF-8'?>"))
	xml.Unmarshal(xmlInput, &folder)
	return folder
}

func (folder *JenkinsFolder) GetCredentials() Credentials {
	return folder.Properties.CredentialProperty.DomainCredentials.Entry.Credentials
}

func GetFolder(url string, folderName string, username string, password string) (JenkinsFolder, error) {
	url = fmt.Sprintf("%s/config.xml", GetFolderURL(url, folderName))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		return JenkinsFolder{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return JenkinsFolder{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return JenkinsFolder{}, errors.New("Unauthorized 401")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return JenkinsFolder{}, err
	}

	return parseJenkinsFolder(data), nil
}
