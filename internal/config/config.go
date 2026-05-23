package config

type Config struct {
	RootPath string
	Top      int
	MaxDepth int
	Workers  int
	ByGit    bool
	JSON     bool
	Excludes map[string]bool
}
