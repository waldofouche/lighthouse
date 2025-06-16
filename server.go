package lighthouse

type ServerConf struct {
	Port int     `yaml:"port"`
	TLS  tlsConf `yaml:"tls"`
	// Secure bool    `yaml:"-"`
	// Basepath       string       `yaml:"-"`
}

type tlsConf struct {
	Enabled      bool   `yaml:"enabled"`
	RedirectHTTP bool   `yaml:"redirect_http"`
	Cert         string `yaml:"cert"`
	Key          string `yaml:"key"`
}
