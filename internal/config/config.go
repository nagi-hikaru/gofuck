package config

type RuleConfig struct {
	Enabled  bool `mapstructure:"enabled"`
	Priority int  `mapstructure:"priority"`
}

type ProviderConfig struct {
	BaseURL string   `mapstructure:"base_url"`
	APIKey  string   `mapstructure:"api_key"`
	Models  []string `mapstructure:"models"`
}

type Settings struct {
	RulesConfig     map[string]RuleConfig     `mapstructure:"rules"`
	ProvidersConfig map[string]ProviderConfig `mapstructure:"providers"`
}

var Config Settings
