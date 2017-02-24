package owl

import (
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap"
	"github.com/urfave/cli"
)

// GetLogger returns zap logger configured for owl's needs
func GetLogger(level zapcore.Level) *zap.Logger {

	zapconfig := zap.NewProductionConfig()
	zapconfig.Encoding = "console"
	zapconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapconfig.DisableStacktrace = true
	zapconfig.DisableCaller = true
	zapconfig.Level.SetLevel(level)

	logger, err := zapconfig.Build()

	if err != nil {
		panic(err)
	}

	return logger
}

// GetApp returns configured urfave/cli application configured for owl's needs
func GetApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Owl"
	app.Usage = "filewatcher that runs bash command when file changes"

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "Version of owl",
	}

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "ignore, i",
			Usage: "All directories with name `IGNORE` are ignored",
		},
		cli.StringFlag{
			Name:  "run, r",
			Usage: "If is any file changed, run `RUN`",
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "verbose mode",
		},
		cli.StringFlag{
			Name:  "debounce, d",
			Usage: "Waiting time for executing in miliseconds",
		},
		cli.StringSliceFlag{
			Name:  "filter, f",
			Usage: "Files are filtered by expression",
		},
	}

	return app
}
