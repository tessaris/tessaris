package config

type Config struct {
	// Port is the port that the server will listen on.
	Port int
	// TrimTrailingSlash determines whether or not trailing slashes should be removed from URLs.
	TrimTrailingSlash bool
	// InertiaView is the name of the Go HTML template file that will be used to render Inertia pages.
	InertiaView string
}

/* New creates a new Config instance with default values. */
func New() *Config {
	return &Config{
		Port:              3000,
		TrimTrailingSlash: true,
		InertiaView:       "app",
	}
}
