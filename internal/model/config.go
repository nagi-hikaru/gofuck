package model

type RuleConfig struct {
	Enabled  bool `mapstructure:"enabled"`
	Priority int  `mapstructure:"priority"`
}

type Config struct {
	RulesConfig map[string]RuleConfig `mapstruct[ure:"rulesconfig"`
}
