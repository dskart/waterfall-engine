package cache

type Config struct {
	// If given, an in-memory cache will be used. This is only recommended for dev and tests.
	InMemory bool `yaml:"InMemory"`

	// If given, Redis will be used for the store. This should be an IP address and port such as
	// "127.0.0.1:6379".
	RedisAddress string `yaml:"RedisAddress"`
}
