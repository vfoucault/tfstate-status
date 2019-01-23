package models

import (
	"strings"
	"time"
)

type ListTfStates []*TfState

func (lt ListTfStates) Empty() (data ListTfStates) {
	for _, s := range lt {
		if s.State.NbsResources() == 0 {
			data = append(data, s)
		}
	}
	return
}

func (lt ListTfStates) FilterWorkSpace(workspace string) (data ListTfStates) {
	for _, s := range lt {
		if strings.Contains(s.Workspace, workspace) {
			data = append(data, s)
		}
	}
	return
}

type TfState struct {
	Name         string    `json:"name"`
	FileName     string    `json:"fileName"`
	LastModified time.Time `json:"lastModified"`
	Workspace    string
	State        *TerraformState `json:"state"`
}

type TerraformState struct {
	Version          int     `json:"version"`
	TerraformVersion string  `json:"terraform_version"`
	Serial           int     `json:"serial"`
	Lineage          string  `json:"lineage"`
	Modules          Modules `json:"modules"`
}

func (t *TerraformState) NbsResources() (count int) {
	for _, i := range t.Modules {
		count += len(i.Resources)
	}
	return
}

type TerraformStateModule struct {
	Path      []string               `json:"path"`
	Outputs   map[string]interface{} `json:"outputs"`
	Resources map[string]interface{} `json:"resources"`
	DependsOn []string               `json:"depends_on"`
}

type Modules []TerraformStateModule
