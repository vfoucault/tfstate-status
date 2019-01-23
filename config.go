package main

import (
	"flag"
	"runtime"
)

type Config struct {
	ContainerName   string `json:"containerName"`
	Prefix          string `json:"prefix"`
	Verbose         bool   `json:"verbose"`
	Provider        string `json:"provider"`
	Threads         int    `json:"threads"`
	Empty           bool   `json:"empty_only"`
	FilterWorkspace string `json:"filter_by_name"`
}

// init the config with flag args
func (c *Config) Init() {
	flag.StringVar(&c.Provider, "provider", "", "s3 so far")
	flag.StringVar(&c.ContainerName, "container", "", "The ContainerName")
	flag.StringVar(&c.Prefix, "prefix", "", "Prefix")
	flag.BoolVar(&c.Verbose, "verbose", false, "Be Verbose")
	flag.IntVar(&c.Threads, "threads", runtime.NumCPU(), "Number of threads. Default to cores count")
	flag.BoolVar(&c.Empty, "list-empty", false, "List empty states only	")
	flag.StringVar(&c.FilterWorkspace, "filter-workspace", "", "Filter workspace by name")

	flag.Parse()
}
