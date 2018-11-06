package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type AllView struct {
	XMLName xml.Name `xml:"allView"`
	Jobs    []Job    `xml:"job"`
}

type Job struct {
	XMLName xml.Name `xml:"job"`
	Name    string   `xml:"name"`
	URL     string   `xml:"url"`
}

func ExportJobs(server string, username string, password string) error {
	jobs, err := GetJobs(server, username, password)
	if err != nil {
		return err
	}

	var directory = "jobs"
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, 0755)
	}

	for _, job := range jobs {
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
	return nil
}

func GetJobs(server string, username string, password string) ([]Job, error) {
	url := fmt.Sprintf("%s/view/all/api/xml", server)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		return make([]Job, 0), err
	}

	resp, err := client.Do(req)
	if err != nil {
		return make([]Job, 0), err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return make([]Job, 0), errors.New("Unauthorized 401")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make([]Job, 0), err
	}

	var view AllView
	xml.Unmarshal(data, &view)

	return view.Jobs, nil
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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(data), ":"), nil
}

func ImportJobs(server string, username string, password string) error {
	jobs, err := ioutil.ReadDir("jobs")
	if err != nil {
		return err
	}

	for _, job := range jobs {
		fmt.Printf("Import job: %s\n", job.Name())
		err := ImportJob(job.Name(), server, username, password)
		if err != nil {
			return err
		}
	}

	return nil
}

func ImportJob(name string, server string, username string, password string) error {
	jsonStr, err := ioutil.ReadFile("jobs/" + name + "/config.xml")
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/createItem?name=%s", server, name)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
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
		return errors.New("Job couldn't not be imported: verify all plugins used in this job are installed on the target jenkins instance")
	}

	return nil
}
