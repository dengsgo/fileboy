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
		//BeforeExec string   `yaml:"beforeExec"`
		Exec []string `yaml:"exec"`
		//AfterExec  string   `yaml:"afterExec"`
	}
}
