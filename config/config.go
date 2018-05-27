// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period        time.Duration `config:"period"`
	CrowdUrl      string        `config:"crowd_url"`
	CrowdUsername string        `config:"crowd_username"`
	CrowdPassword string        `config:"crowd_password"`
}

var DefaultConfig = Config{
	Period: 10 * time.Minute,
}
