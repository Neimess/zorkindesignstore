package args

import "github.com/alexflint/go-arg"

type Args struct {
	Migrate    bool   `arg:"-m,--migrate" help:"Run database migrations"`
	Force      int    `arg:"--force"      help:"Force database to specific migration version before running --migrate (0 - skip)"`
	ConfigPath string `arg:"-c,--config"  help:"Path to the configuration file"`
}

func Parse() *Args {
	var args Args
	arg.MustParse(&args)
	return &args
}
