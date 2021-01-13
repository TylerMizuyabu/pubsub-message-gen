package types

// Config struct
type Config struct {
	Topic           string `yaml:"topic"`
	MessageType		string `yaml:"messageType"`
	ProjectID       string `yaml:"projectID"`
	MessageInterval int64 `yaml:"messageInterval"`
}

//TODO Config validation