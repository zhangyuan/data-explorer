package conf

type ConnectionsConfiguration struct {
	Connections []Connection `yaml:"connections"`
}

type Connection struct {
	Id  string `yaml:"id"`
	DSN string `yaml:"dsn"`
}

func LoadConnection(path string) (*ConnectionsConfiguration, error) {
	return LoadFromYAML[ConnectionsConfiguration](path)
}
