package regexp

type Config struct {
	CacheEnabled bool `mapstructure:"cacheenabled"`

	CacheMaxNumEntries int `mapstructure:"cachemaxnumentries"`
}
