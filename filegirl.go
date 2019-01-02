package main

type FileGirl struct {
	Core struct {
		Version int `yaml:"version"`
	}
	Monitor struct {
		Types       []string `yaml:"types"`
		IncludeDirs []string `yaml:"includeDirs"`
		ExceptDirs  []string `yaml:"exceptDirs"`
	}
	Command struct {
		Exec            []string `yaml:"exec"`
		DelayMillSecond int      `yaml:"delayMillSecond"`
	}
	Notifier struct {
		CallUrl string `yaml:"callUrl"`
	}
}
