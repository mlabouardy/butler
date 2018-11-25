package main

import "testing"

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
