package config

var ZEPTO_TEST_CONFIG = Config{
	App: AppConfig{
		Name:    "zepto-test",
		Version: "1.0.0",
		Session: SessionConfig{
			Name:   "zsid",
			Secret: "test",
		},
	},
	Server: ServerConfig{
		Host:         "0.0.0.0",
		Port:         8000,
		ReadTimeout:  15000,
		WriteTimeout: 15000,
	},
	Logger: LoggerConfig{
		Level:     "error",
		Colors:    true,
		Timestamp: true,
	},
}
