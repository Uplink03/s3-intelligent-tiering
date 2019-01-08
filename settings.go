package main

import (
	"flag"
	"runtime"
)

type RunMode int

const (
	RunModeInvalid RunMode = iota
	RunModeListObjects
	RunModeChangeObjects
)

type Settings struct {
	WorkerCount int
	RunMode     RunMode
	Region      string
	Bucket      string
}

var settings Settings

func init() {
	flag.IntVar(&settings.WorkerCount, "workers", runtime.NumCPU(), "Number of workers to deploy")
	list := flag.Bool("list", false, "List objects and store in SQLite database")
	change := flag.Bool("change", false, "Change objects listed in SQLite database")

	flag.Parse()

	if *list == *change {
		// neither or both of list and change specified - that's an error
		exitErrorf("Must specify run mode: -list or -change")
	}

	switch true {
	case *list:
		settings.RunMode = RunModeListObjects
	case *change:
		settings.RunMode = RunModeChangeObjects
	}

	if settings.WorkerCount < 1 {
		settings.WorkerCount = 1
	}

	settings.Region = flag.Arg(0)
	settings.Bucket = flag.Arg(1)
}
