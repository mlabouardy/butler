package main

import "testing"

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
