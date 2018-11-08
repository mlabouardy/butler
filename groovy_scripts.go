package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const decryptScriptTemplate = `import groovy.json.JsonSlurperClassic
import groovy.json.JsonOutput


def json = """<<JSON HERE>>"""

def data = new JsonSlurperClassic().parseText(json)
data.userpass.each {
  it.password = hudson.util.Secret.fromString(it.password).getPlainText()
}

data.secretfile.each {
	it.secretBytes = new String(com.cloudbees.plugins.credentials.SecretBytes.fromString(it.secretBytes).getPlainData(), "ASCII")
}

println JsonOutput.toJson(data)`

func GetDecryptScriptForCredentials(credentials Credentials) string {
	marshalledCredentials, _ := json.Marshal(credentials)
	return strings.Replace(decryptScriptTemplate, "<<JSON HERE>>", string(marshalledCredentials), 1)
}

func ExecuteGroovyScriptOnJenkins(script string, rawUrl string, username string, password string) string {
	apiURL := fmt.Sprintf("%s/scriptText", rawUrl)
	data := url.Values{}
	data.Set("script", script)
	body := strings.NewReader(data.Encode())
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	crumb, err := GetCrumb(rawUrl, username, password)
	if err != nil {
		fmt.Errorf("No crumb issueing possible: %v", err)
	} else {
		req.Header.Set(crumb[0], crumb[1])
	}

	req.SetBasicAuth(username, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)
	return string(responseBody)
}
