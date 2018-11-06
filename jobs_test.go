package main

import (
	"reflect"
	"testing"
)

func TestJob_IsFolder(t *testing.T) {
	tests := []struct {
		name string
		job  Job
		want bool
	}{
		{"Should be a folder", Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"}, true},
		{"Should be no folder", Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.job.IsFolder(); got != tt.want {
				t.Errorf("Job.IsFolder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterOutFolders(t *testing.T) {
	tests := []struct {
		name       string
		unfiltered []Job
		want       []Job
	}{
		{
			"Filter out folder",
			[]Job{
				Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
				Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
			},
			[]Job{
				Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
			},
		},
		{
			"Filter out many folders",
			[]Job{
				Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
				Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
				Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
				Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
				Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
				Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
				Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
				Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
				Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
			},
			[]Job{
				Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterOutFolders(tt.unfiltered); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterOutFolders() = %v, want %v", got, tt.want)
			}
		})
	}
}
