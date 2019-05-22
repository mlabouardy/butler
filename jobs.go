package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"net/url"
)

type JobList struct {
	Jobs []Job `xml:"job"`
}

type JenkinsHTTPClient struct {
	BasicAuthSettings BasicAuthSettings
}

type BasicAuthSettings struct {
	Username string
	Password string
}

type Job struct {
	Class      string `xml:"_class,attr"`
	Name       string `xml:"name"`
	URL        string `xml:"url"`
	httpClient *JenkinsHTTPClient
}

func (job *Job) IsFolder() bool {
	return strings.HasSuffix(job.Class, "Folder")
}

func (job *Job) GetFolderName() (ret string) {

	path := strings.Split(job.URL, "/job")
	if len(path) == 1 {
		return
	}
	ret = strings.Join(path[1:], "/")
	ret = strings.Replace(ret, "//", "/", -1)
	ret = strings.TrimLeft(ret, "/")
	ret = strings.TrimRight(ret, "/")
	return
}

func (jobList *JobList) GetSubfolders() JobList {
	var subfolders JobList

	isFolder := func(job Job) bool { return job.IsFolder() }
	subfolders.Jobs = choose(jobList.Jobs, isFolder)
	return subfolders
}

func (jobList *JobList) GetSubfoldersRecursively() (JobList, error) {
	var subfolders = jobList.GetSubfolders()

	for _, innerfolder := range subfolders.Jobs {
		allJobsOfInnerFolder, err := innerfolder.GetJobs()
		if err != nil {
			return subfolders, err
		}
		recursiveSubfolders, err := allJobsOfInnerFolder.GetSubfoldersRecursively()
		if err != nil {
			return subfolders, err
		}
		subfolders.Jobs = append(subfolders.Jobs, recursiveSubfolders.Jobs...)
	}

	return subfolders, nil
}

func (job *Job) GetJobs() (JobList, error) {
	url := fmt.Sprintf("%s/api/xml", job.URL)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(job.httpClient.BasicAuthSettings.Username, job.httpClient.BasicAuthSettings.Password)
	if err != nil {
		return NewJobList(), err
	}

	resp, err := client.Do(req)
	if err != nil {
		return NewJobList(), err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return NewJobList(), errors.New("Unauthorized 401")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return NewJobList(), err
	}

	var jobList JobList
	xml.Unmarshal(data, &jobList)

	for i, _ := range jobList.Jobs {
		jobList.Jobs[i].httpClient = job.httpClient
	}
	return jobList, nil
}

func NewJob(serverUrl string, folderName string, httpClient *JenkinsHTTPClient) Job {
	return Job{
		URL:        GetFolderURL(serverUrl, folderName),
		httpClient: httpClient,
	}
}

func (jobList *JobList) WithoutFolders() JobList {
	isNoFolder := func(job Job) bool { return !job.IsFolder() }
	return JobList{Jobs: choose(jobList.Jobs, isNoFolder)}
}

func choose(jobs []Job, test func(Job) bool) (ret []Job) {
	ret = make([]Job, 0)
	for _, job := range jobs {
		if test(job) {
			ret = append(ret, job)
		}
	}
	return
}

func NewJobList() JobList {
	return JobList{Jobs: make([]Job, 0)}
}

func ListFolders(server string, folderName string, username string, password string, recursive bool) error {
	httpClient := &JenkinsHTTPClient{
		BasicAuthSettings: BasicAuthSettings{
			Username: username,
			Password: password,
		},
	}
	rootJob := NewJob(server, folderName, httpClient)
	var jobsList JobList
	jobsList, err := rootJob.GetJobs()
	if err != nil {
		return err
	}

	jobsList = jobsList.GetSubfolders()

	if recursive {
		jobsList, err = jobsList.GetSubfoldersRecursively()
		if err != nil {
			return err
		}
	}

	for _, folder := range jobsList.Jobs {
		fmt.Println(folder.GetFolderName())
	}
	return nil
}

func ExportJobs(server string, folderName string, username string, password string, skipFolder bool) error {
	httpClient := &JenkinsHTTPClient{
		BasicAuthSettings: BasicAuthSettings{
			Username: username,
			Password: password,
		},
	}
	rootJob := NewJob(server, folderName, httpClient)
	jobs, err := rootJob.GetJobs()
	if err != nil {
		return err
	}

	if skipFolder {
		jobs = jobs.WithoutFolders()
	}

	var directory = "jobs"
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, 0755)
	}

	for _, job := range jobs.Jobs {
		fmt.Printf("Exporting job: %s\n", job.Name)
		err := ExportJob(job, username, password)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExportJob(job Job, username string, password string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", job.URL+"/config.xml", nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return errors.New("Unauthorized 401")
	}

	if resp.StatusCode != 200 {
		return errors.New("Job couldn't not be exported")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var directory = "jobs/" + job.Name
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, 0755)
	}

	f, err := os.Create("jobs/" + job.Name + "/config.xml")
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "%s", data)
	if job.IsFolder() {
		fmt.Printf("\tJob is a folder.\n")
	}
	return nil
}

func GetCrumb(host string, username string, password string) ([]string, error) {
	crumbUrl := `%s/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,":",//crumb)`
	url := fmt.Sprintf(crumbUrl, host)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		return []string{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return []string{}, errors.New("Unauthorized 401")
	}

	if resp.StatusCode == 404 {
		return []string{}, errors.New("Not found 404")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(data), ":"), nil
}

func ImportJobs(server string, username string, password string, folder string) error {
	jobs, err := ioutil.ReadDir("jobs")
	if err != nil {
		return err
	}

	for _, job := range jobs {
		fmt.Printf("Import job: %s\n", job.Name())
		err := ImportJob(job.Name(), folder, server, username, password)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func ImportJob(name string, folderName string, server string, username string, password string) error {
	jsonStr, err := ioutil.ReadFile("jobs/" + name + "/config.xml")
	if err != nil {
		return err
	}

	// escape is mandatory if the job has special char in the name
	name = url.PathEscape(name)
	createURL := fmt.Sprintf("%s/createItem?name=%s", GetFolderURL(server, folderName), name)
	client := &http.Client{}
	req, err := http.NewRequest("POST", createURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.SetBasicAuth(username, password)

	crumb, err := GetCrumb(server, username, password)
	if err != nil {
		return err
	}

	req.Header.Set(crumb[0], crumb[1])
	req.Header.Set("Content-Type", "text/xml")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return errors.New("Unauthorized 401")
	}

	if resp.StatusCode != 200 {
		errorMsg := fmt.Sprintf(`Job couldn't not be imported (code %s).
		Check if it is already existing and verify all plugins used in this job are installed on the target jenkins instance`,
		strconv.Itoa(resp.StatusCode))
		b, errRead := ioutil.ReadAll(resp.Body)
		if errRead == nil {
			errorMsg = errorMsg + fmt.Sprintf("%s", b)
		}
		return errors.New(errorMsg)
	}

	return nil
}
