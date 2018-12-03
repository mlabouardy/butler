package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFolderURL(t *testing.T) {
	type args struct {
		url        string
		folderName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Without preceeding slash",
			args: args{
				url:        "https://sample-jenkins/",
				folderName: "BLA/BLUB",
			},
			want: "https://sample-jenkins/job/BLA/job/BLUB",
		},
		{
			name: "With preceeding slash",
			args: args{
				url:        "https://sample-jenkins/",
				folderName: "/BLA/BLUB",
			},
			want: "https://sample-jenkins/job/BLA/job/BLUB",
		},
		{
			name: "With no trailing slash in url",
			args: args{
				url:        "https://sample-jenkins",
				folderName: "BLA/BLUB",
			},
			want: "https://sample-jenkins/job/BLA/job/BLUB",
		},
		{
			name: "With no trailing slash in url but preceeding slash in folder",
			args: args{
				url:        "https://sample-jenkins",
				folderName: "/BLA/BLUB",
			},
			want: "https://sample-jenkins/job/BLA/job/BLUB",
		},
		{
			name: "With no trailing slash in folder name",
			args: args{
				url:        "https://sample-jenkins",
				folderName: "/BLA/BLUB/",
			},
			want: "https://sample-jenkins/job/BLA/job/BLUB",
		},
		{
			name: "With empty folder name",
			args: args{
				url:        "https://sample-jenkins",
				folderName: "",
			},
			want: "https://sample-jenkins",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFolderURL(tt.args.url, tt.args.folderName); got != tt.want {
				t.Errorf("GetFolderURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJenkinsFolder(t *testing.T) {
	assert := assert.New(t)
	xml := `<?xml version='1.1' encoding='UTF-8'?>
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

	got := parseJenkinsFolder([]byte(xml))

	assert.NotNil(got, "Should not be nil.")
	assert.Equal(got.GetCredentials().UsernamePassword[0].ID, "test")
}
