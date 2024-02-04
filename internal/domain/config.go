package domain

type Config struct {
	App    App          `yaml:"app" env-prefix:"NOOLINGO_" json:"app"`
	Listen ListenConfig `yaml:"listen" env-prefix:"NOOLINGO_LISTEN_" json:"listen"`
	Log    Logger       `yaml:"log" env-prefix:"NOOLINGO_" json:"log"`
	Grpc   Grpc         `yaml:"grpc" env-prefix:"NOOLINGO_" json:"grpc"`
}

type App struct {
	Cors         bool              `yaml:"cors" env:"CORS"`
	RolesAccess  map[string]int    `yaml:"rolesAccess"`
	AccessMap    map[string]string `yaml:"accessMap"`
	AccessPrefix []string          `yaml:"accessPrefix"`
}

type ListenConfig struct {
	Ports portConfig `yaml:"ports" env-prefix:"PORTS_"`
	Host  string     `yaml:"host" env:"HOSTNAME" env-default:"0.0.0.0"`
}

type portConfig struct {
	Http string `yaml:"http" env:"HTTP"`
}

type Logger struct {
	Level map[string]string `yaml:"level" env:"LOG_LEVEL"`
}

type Grpc struct {
	Clients Clients `yaml:"clients" env-prefix:"GRPC_CLIENTS_"`
}

type Clients struct {
	UserService      string `yaml:"userservice" env:"USERSERVICE"`
	CardService      string `yaml:"cardservice" env:"CARDSERVICE"`
	DeckService      string `yaml:"deckservice" env:"DECKSERVICE"`
	StatisticService string `yaml:"statisticservice" env:"STATISTICSERVICE"`
}
