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

func TestJobList_WithoutFolders(t *testing.T) {

	tests := []struct {
		name   string
		fields JobList
		want   JobList
	}{
		{
			name: "No jobs given",
			fields: JobList{
				Jobs: []Job{},
			},
			want: JobList{
				Jobs: []Job{},
			},
		},
		{
			name: "Mixed",
			fields: JobList{
				Jobs: []Job{
					Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
					Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
					Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
					Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
					Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
					Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
					Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
				},
			},
			want: JobList{
				Jobs: []Job{
					Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
					Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
					Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobList := &JobList{
				Jobs: tt.fields.Jobs,
			}
			if got := jobList.WithoutFolders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JobList.WithoutFolders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJobList_GetSubfolders(t *testing.T) {

	tests := []struct {
		name   string
		fields JobList
		want   JobList
	}{
		{
			name: "No jobs given",
			fields: JobList{
				Jobs: []Job{},
			},
			want: JobList{
				Jobs: []Job{},
			},
		},
		{
			name: "Mixed",
			fields: JobList{
				Jobs: []Job{
					Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
					Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
					Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
					Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
					Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
					Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
					Job{Class: "org.jenkinsci.plugins.workflow.job.WorkflowJob"},
				},
			},
			want: JobList{
				Jobs: []Job{
					Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
					Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
					Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
					Job{Class: "com.cloudbees.hudson.plugins.folder.Folder"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobList := &JobList{
				Jobs: tt.fields.Jobs,
			}
			if got := jobList.GetSubfolders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JobList.GetSubfolders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_getFolderName(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "Root folder",
			url:  "https://jenkins.example.org",
			want: "",
		},
		{
			name: "One level down",
			url:  "https://jenkins.example.org/job/imafolder",
			want: "imafolder",
		},
		{
			name: "Two levels down",
			url:  "https://jenkins.example.org/job/imafolder/job/metoo",
			want: "imafolder/metoo",
		},
		{
			name: "Three levels down",
			url:  "https://jenkins.example.org/job/imafolder/job/metoo/job/gettingdarkhere",
			want: "imafolder/metoo/gettingdarkhere",
		},
		{
			name: "Two levels down and ends with /job",
			url:  "https://jenkins.example.org/job/imafolder/job/metoo/job",
			want: "imafolder/metoo",
		},
		{
			name: "Two levels down and ends with /job/",
			url:  "https://jenkins.example.org/job/imafolder/job/metoo/job/",
			want: "imafolder/metoo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := &Job{
				URL: tt.url,
			}
			if got := job.GetFolderName(); got != tt.want {
				t.Errorf("Job.GetFolderName() = %v, want %v", got, tt.want)
			}
		})
	}
}
