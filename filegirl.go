package main

type FileGirl struct {
	Core struct {
		Version int `yaml:"version"`
	}
	Monitor struct {
		Types       []string `yaml:"types"`
		ExceptTypes []string `yaml:"exceptTypes"`
		IncludeDirs []string `yaml:"includeDirs"`
		ExceptDirs  []string `yaml:"exceptDirs"`
		Events      []string `yaml:"events"`
		// convert to
		TypesMap       map[string]bool `yaml:"-"`
		ExceptTypesMap map[string]bool `yaml:"-"`
		IncludeDirsMap map[string]bool `yaml:"-"`
		ExceptDirsMap  map[string]bool `yaml:"-"`
		DirsMap        map[string]bool `yaml:"-"`

		IncludeDirsRec map[string]bool `yaml:"-"`
	}
	Command struct {
		Exec            []string `yaml:"exec"`
		DelayMillSecond int      `yaml:"delayMillSecond"`
	}
	Notifier struct {
		CallUrl string `yaml:"callUrl"`
	}
	Instruction []string `yaml:"instruction"`

	// convert to
	InstructionMap map[string]bool `yaml:"-"`
}
