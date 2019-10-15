package main

import (
	"fmt"
	"github.com/mbndr/logo"
	"github.com/olekukonko/tablewriter"
	"github.com/vfoucault/tfstate-status/models"
	"github.com/vfoucault/tfstate-status/providers"
	"os"
	"strconv"
	"sync"
)

var (
	log    = logo.NewSimpleLogger(os.Stderr, logo.INFO, "tfstate-status ", true)
	config = new(Config)
)

func main() {
	// Parse flags
	config.Init()

	// Debug ?
	if config.Verbose {
		log.SetLevel(logo.DEBUG)
	}

	// Getting file list over the provider.
	// s3
	if config.Provider != "aws" {
		log.Fatal("Only AWS s3 supported for far")
	}
	provider, _ := providers.NewProviderFactory(config.Provider, config.ContainerName, config.Prefix)

	files, _ := provider.ListFiles()

	// setup routines
	var wg sync.WaitGroup
	fileChan := make(chan providers.ObjectFile, len(files))
	stateChan := make(chan *models.TfState, len(files))
	done := make(chan bool, 1)

	for _, i := range files {
		fileChan <- i
	}
	close(fileChan)

	wg.Add(config.Threads)

	for i := 0; i < config.Threads; i++ {
		go func() {
			defer wg.Done()
			for stateFile := range fileChan {
				state, err := provider.ProcessState(stateFile)
				if err != nil {
					log.Errorf("Something went wrong processing %v. Err=%v", stateFile.Key, err.Error())
				}
				log.Debugf("Processed %v in workspace %v", state.Name, state.Workspace)
				stateChan <- state
			}
		}()
	}
	var wg2 sync.WaitGroup
	wg2.Add(1)
	var states models.ListTfStates
	go func() {
		for {
			select {
			case state := <-stateChan:
				states = append(states, state)
			case <-done:
				defer wg2.Done()
				close(stateChan)
				if config.Empty {
					states = states.Empty()
				}
				if config.FilterWorkspace != "" {
					states = states.FilterWorkSpace(config.FilterWorkspace)
				}
				TableWriter(states)
				return
			}
		}
	}()

	wg.Wait()
	done <- true
	wg2.Wait()
}

func TableWriter(data []*models.TfState) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(Header())

	for _, v := range PrepareTabulare(data) {
		table.Append(v)
	}
	table.Render() // Send output

}

func Header() []string {
	return []string{"File", "workspace", "Last Modified", "Resources"}
}

func PrepareTabulare(data []*models.TfState) [][]string {
	var output [][]string
	for _, item := range data {
		add := []string{item.FileName, item.Workspace, fmt.Sprint(item.LastModified), strconv.Itoa(item.State.NbsResources())}
		output = append(output, add)
	}
	return output
}
