package config

type Version struct {
	Name string `koanf:"name"`
	URL  string `koanf:"url"`
}

type Config struct {
	Versions []Version `koanf:"versions"`
}
