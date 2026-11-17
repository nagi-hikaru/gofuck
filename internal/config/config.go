package config

type RuleConfig struct {
	Enabled  bool `mapstructure:"enabled"`
	Priority int  `mapstructure:"priority"`
}

type Settings struct {
	RulesConfig map[string]RuleConfig
}

var Config Settings
