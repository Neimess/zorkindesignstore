package args

import "github.com/alexflint/go-arg"

type Args struct {
	ConfigPath string `arg:"-c,--config"  help:"Path to the configuration file"`
}

func Parse() *Args {
	var args Args
	arg.MustParse(&args)
	return &args
}
