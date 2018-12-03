package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
)

type Plugin struct {
	Name        string `json:"shortName"`
	Description string `json:"longName"`
	Version     string `json:"version"`
}

type PluginData struct {
	Class   string   `json:"_class"`
	Plugins []Plugin `json:"plugins"`
}

func GetPlugins(server string, username string, password string) ([]Plugin, error) {
	url := fmt.Sprintf("%s/pluginManager/api/json?depth=1", server)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		return []Plugin{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return []Plugin{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return []Plugin{}, errors.New("Unauthorized 401")
	}

	var data PluginData
	json.NewDecoder(resp.Body).Decode(&data)

	return data.Plugins, nil
}

func ExportPlugins(server string, username string, password string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Version", "Description"})

	plugins, err := GetPlugins(server, username, password)
	if err != nil {
		return err
	}

	file, err := os.Create("plugins.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	for _, plugin := range plugins {
		table.Append([]string{plugin.Name, plugin.Version, plugin.Description})
		file.WriteString(fmt.Sprintf("%s@%s\n", plugin.Name, plugin.Version))
	}

	table.Render()
	return nil
}

func ImportPlugins(server string, username string, password string) error {
	url := fmt.Sprintf("%s/pluginManager/installNecessaryPlugins", server)
	file, err := os.Open("plugins.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fmt.Printf("Installing %s\n", scanner.Text())
		reqBody := fmt.Sprintf(`<jenkins><install plugin="%s" /></jenkins>`, scanner.Text())
		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(reqBody)))
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
			return errors.New("Plugin cannot be installed")
		}
	}
	return nil
}
