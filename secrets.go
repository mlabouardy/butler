package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type JenkinsFolder struct {
	XMLName     xml.Name `xml:"com.cloudbees.hudson.plugins.folder.Folder"`
	Plugin      string   `xml:"plugin,attr"`
	Actions     string   `xml:"actions"`
	DisplayName string   `xml:"displayName"`
	Properties  struct {
		CredentialProperty struct {
			DomainCredentials struct {
				Class string `xml:"class,attr"`
				Entry struct {
					Credentials Credentials `xml:"java.util.concurrent.CopyOnWriteArrayList"`
				} `xml:"entry"`
			} `xml:"domainCredentialsMap"`
		} `xml:"com.cloudbees.hudson.plugins.folder.properties.FolderCredentialsProvider_-FolderCredentialsProperty"`
	} `xml:"properties"`
	FolderViews struct {
		Text  string `xml:",chardata"`
		Class string `xml:"class,attr"`
		Views struct {
			Text               string `xml:",chardata"`
			HudsonModelAllView struct {
				Text  string `xml:",chardata"`
				Owner struct {
					Text      string `xml:",chardata"`
					Class     string `xml:"class,attr"`
					Reference string `xml:"reference,attr"`
				} `xml:"owner"`
				Name            string `xml:"name"`
				FilterExecutors string `xml:"filterExecutors"`
				FilterQueue     string `xml:"filterQueue"`
				Properties      struct {
					Text  string `xml:",chardata"`
					Class string `xml:"class,attr"`
				} `xml:"properties"`
			} `xml:"hudson.model.AllView"`
		} `xml:"views"`
		PrimaryView string `xml:"primaryView"`
		TabBar      struct {
			Text  string `xml:",chardata"`
			Class string `xml:"class,attr"`
		} `xml:"tabBar"`
	} `xml:"folderViews"`
	HealthMetrics struct {
		Text                                                        string `xml:",chardata"`
		ComCloudbeesHudsonPluginsFolderHealthWorstChildHealthMetric struct {
			Text         string `xml:",chardata"`
			NonRecursive string `xml:"nonRecursive"`
		} `xml:"com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric"`
	} `xml:"healthMetrics"`
	Icon struct {
		Text  string `xml:",chardata"`
		Class string `xml:"class,attr"`
	} `xml:"icon"`
}

type Credentials struct {
	UsernamePassword []UsernamePasswordCredential `xml:"com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl" json:"userpass"`
	SecretFile       []struct {
		Plugin      string `xml:"plugin,attr"`
		ID          string `xml:"id" json:"id"`
		Description string `xml:"description" json:"description"`
		FileName    string `xml:"fileName" json:"fileName"`
		SecretBytes string `xml:"secretBytes" json:"secretBytes"`
	} `xml:"org.jenkinsci.plugins.plaincredentials.impl.FileCredentialsImpl"`
}

type UsernamePasswordCredential struct {
	Plugin      string `xml:"plugin,attr" json:"plugin"`
	ID          string `xml:"id" json:"id"`
	Description string `xml:"description" json:"description"`
	Username    string `xml:"username" json:"username"`
	Password    string `xml:"password" json:"password"`
}

func ParseJenkinsFolder(xmlInput []byte) JenkinsFolder {
	var folder JenkinsFolder
	xmlInput = []byte(strings.Trim(string(xmlInput), "<?xml version='1.1' encoding='UTF-8'?>"))
	xml.Unmarshal(xmlInput, &folder)
	return folder
}

func (folder *JenkinsFolder) GetCredentials() Credentials {
	return folder.Properties.CredentialProperty.DomainCredentials.Entry.Credentials
}

const testxml = `<?xml version='1.1' encoding='UTF-8'?>
<com.cloudbees.hudson.plugins.folder.Folder plugin="cloudbees-folder@6.6">
  <actions/>
  <displayName>Ansible</displayName>
  <properties>
    <com.cloudbees.hudson.plugins.folder.properties.FolderCredentialsProvider_-FolderCredentialsProperty>
      <domainCredentialsMap class="hudson.util.CopyOnWriteMap$Hash">
        <entry>
          <com.cloudbees.plugins.credentials.domains.Domain plugin="credentials@2.1.18">
            <specifications/>
          </com.cloudbees.plugins.credentials.domains.Domain>
          <java.util.concurrent.CopyOnWriteArrayList>
            <com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl plugin="credentials@2.1.18">
              <id>test</id>
              <description>Test</description>
              <username>Test</username>
              <password>{AQAAABAAAAAQ7EzV5N/fXZEKM9HyG+1T66P67iqU+tptVCNuvNX1TM0=}</password>
            </com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl>
            <org.jenkinsci.plugins.plaincredentials.impl.FileCredentialsImpl plugin="plain-credentials@1.4">
              <id>deploy-key-file</id>
              <description>blub</description>
              <fileName>accessKeys.csv</fileName>
              <secretBytes>{bEtRRJ+hCoQHgEAmcGhAOlKFx6J5tVuKmwdBVSgdq4zkktsLwG1zHO6swI3mQ5z9UhbgRRHDf2W8oSHlfmno8+KHWKWKyNmQUL5cv6/8n5JnmvsMGx+DT4KJL2XDVl33nuNbDpkcJEDGBWqb2hA47iRtW6h4mxlbNja5E12eUMs=}</secretBytes>
            </org.jenkinsci.plugins.plaincredentials.impl.FileCredentialsImpl>
          </java.util.concurrent.CopyOnWriteArrayList>
        </entry>
      </domainCredentialsMap>
    </com.cloudbees.hudson.plugins.folder.properties.FolderCredentialsProvider_-FolderCredentialsProperty>
  </properties>
  <folderViews class="com.cloudbees.hudson.plugins.folder.views.DefaultFolderViewHolder">
    <views>
      <hudson.model.AllView>
        <owner class="com.cloudbees.hudson.plugins.folder.Folder" reference="../../../.."/>
        <name>all</name>
        <filterExecutors>false</filterExecutors>
        <filterQueue>false</filterQueue>
        <properties class="hudson.model.View$PropertyList"/>
      </hudson.model.AllView>
    </views>
    <primaryView>all</primaryView>
    <tabBar class="hudson.views.DefaultViewsTabBar"/>
  </folderViews>
  <healthMetrics>
    <com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric>
      <nonRecursive>false</nonRecursive>
    </com.cloudbees.hudson.plugins.folder.health.WorstChildHealthMetric>
  </healthMetrics>
  <icon class="com.cloudbees.hudson.plugins.folder.icons.StockFolderIcon"/>
</com.cloudbees.hudson.plugins.folder.Folder>`

func DecryptFolder(url string, folderName string, username string, password string) error {
	folder, _ := GetFolder(GetFolderURL(url, folderName), username, password)
	credentials := folder.GetCredentials()
	b, _ := json.Marshal(credentials)
	fmt.Printf("%v\n", string(b))
	return nil
}

func GetFolderURL(url string, folderName string) string {
	url = strings.TrimRight(url, "/")
	folderName = strings.TrimRight(folderName, "/")
	if !strings.HasPrefix(folderName, "/") {
		folderName = "/" + folderName
	}
	path := strings.Replace(folderName, "/", "/job/", -1)
	return url + path
}

func GetFolder(baseUrlOfFolder string, username string, password string) (JenkinsFolder, error) {
	url := fmt.Sprintf("%s/config.xml", baseUrlOfFolder)

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

	return ParseJenkinsFolder(data), nil
}
