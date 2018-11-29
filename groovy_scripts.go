package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func GetDecryptScriptForCredentials(credentials Credentials) string {
	marshalledCredentials, _ := json.Marshal(credentials)
	return strings.Replace(decryptScriptTemplate, "<<JSON HERE>>", string(marshalledCredentials), 1)
}

const decryptScriptTemplate = `import groovy.json.JsonSlurperClassic
import groovy.json.JsonOutput
import java.util.Base64

def json = """<<JSON HERE>>"""

def data = new JsonSlurperClassic().parseText(json)
data.userpass.each {
  it.password = hudson.util.Secret.fromString(it.password).getPlainText()
}

data.secretfile.each {
	it.rawString = new String(com.cloudbees.plugins.credentials.SecretBytes.fromString(it.secretBytes).getPlainData(), "ASCII")
	it.encodedSecretBytes = new String(Base64.encoder.encode(com.cloudbees.plugins.credentials.SecretBytes.fromString(it.secretBytes).getPlainData()))
}

println JsonOutput.toJson(data)`

func GetApplyScriptForCredentials(credentials Credentials, folderPath string) string {
	marshalledCredentials, _ := json.Marshal(credentials)
	templated := strings.Replace(createOrUpdateCredentialsTemplate, "<<JSON HERE>>", string(marshalledCredentials), 1)
	templated = strings.Replace(templated, "<<FOLDER HERE>>", folderPath, 1)
	return templated
}

const createOrUpdateCredentialsTemplate = `import com.cloudbees.hudson.plugins.folder.properties.FolderCredentialsProvider.FolderCredentialsProperty
import com.cloudbees.hudson.plugins.folder.AbstractFolder
import com.cloudbees.hudson.plugins.folder.Folder
import jenkins.model.*
import com.cloudbees.plugins.credentials.*
import com.cloudbees.plugins.credentials.common.*
import com.cloudbees.plugins.credentials.domains.*
import com.cloudbees.plugins.credentials.impl.*
import org.jenkinsci.plugins.plaincredentials.*
import org.jenkinsci.plugins.plaincredentials.impl.*
import org.apache.commons.fileupload.FileItem
import groovy.json.JsonSlurperClassic
import groovy.json.JsonOutput
import java.util.Base64


def createOrUpdateCredential(credentialStore, newCredential, existingCredentials) {
    def existingCredential = existingCredentials.find{c -> c.getId() == newCredential.getId()}
    if (existingCredential)
        credentialStore.updateCredentials(Domain.global(), existingCredential, newCredential)
    else
        credentialStore.addCredentials(Domain.global(), newCredential)
}

def json = """<<JSON HERE>>"""
String folderPath = "<<FOLDER HERE>>"
def data = new JsonSlurperClassic().parseText(json)

Jenkins.instance.getAllItems(Folder.class)
    .findAll{it.fullName.equals(folderPath)}
    .each{
        AbstractFolder<?> folderAbs = AbstractFolder.class.cast(it)
      	println "We're at ${folderAbs.fullName}"
        FolderCredentialsProperty property = folderAbs.getProperties().get(FolderCredentialsProperty.class)
        if(property == null){
            property = new FolderCredentialsProperty()
            folderAbs.addProperty(property)
        }

        def store = property.getStore()
        def existingCredentials = property.getCredentials()
        data.userpass.each {
            Credentials c = new UsernamePasswordCredentialsImpl(CredentialsScope.GLOBAL, it.id, it.description, it.username, it.password)
            createOrUpdateCredential(store, c, existingCredentials)
		}
		data.secretfile.each {
			def rawSecretFile = it
			fileItem = [ getName: { return rawSecretFile.fileName},  get: { return Base64.decoder.decode(rawSecretFile.encodedSecretBytes) } ] as FileItem
			secretFile = new FileCredentialsImpl(
				CredentialsScope.GLOBAL,
				rawSecretFile.id,
				rawSecretFile.description,
				fileItem, // Don't use FileItem
				null,
				"")
			createOrUpdateCredential(store, secretFile, existingCredentials)
		}
        println existingCredentials.toString()
}`

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
