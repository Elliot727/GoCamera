package types

type Config struct {
	Server   bool   `yaml:"server"`
	Transfer bool   `yaml:"transfer"`
	Organise bool   `yaml:"organise"`
	Source   string `yaml:"source"`
	Dest     string `yaml:"dest"`
	Port     string `yaml:"port"`
}
