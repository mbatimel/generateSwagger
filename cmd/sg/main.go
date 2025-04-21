package main

import (
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/mbatimel/generateSwagger/pkg/generator"
	"github.com/mbatimel/generateSwagger/pkg/logger"
)

var (
	Version    = "v0.0.1"
	BuildStamp = time.Now().String()
)

var log = logger.Log.WithField("module", "sg")

func main() {

	app := cli.NewApp()
	app.Version = Version
	app.EnableBashCompletion = true
	app.Usage = "make generator easy"
	app.Name = "golang generator swagger"
	app.Compiled, _ = time.Parse(time.RFC3339, BuildStamp)

	app.Commands = []*cli.Command{
		{
			Name:   "swagger",
			Usage:  "generate swagger documentation by interfaces in 'service' package",
			Action: cmdSwagger,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "services",
					Value: "./pkg/someService/service",
					Usage: "path to services package",
				},
				&cli.StringFlag{
					Name:  "outFile",
					Usage: "path to output folder",
				},
				&cli.StringSliceFlag{
					Name:  "ifaces",
					Usage: "included interfaces",
				},
				&cli.StringFlag{
					Name:  "redoc",
					Usage: "path to output redoc bundle",
				},
			},

			UsageText:   "sg swagger --include firstIface --exclude secondIface",
			Description: "generate swagger documentation by interfaces",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
func cmdSwagger(c *cli.Context) (err error) {

	defer func() {
		if err == nil {
			log.Info("done")
		}
	}()

	var tr generator.Transport
	if tr, err = generator.NewTransport(log, Version, c.String("services"), c.StringSlice("ifaces")...); err != nil {
		return
	}

	outPath := path.Join(c.String("services"), "swagger.yaml")

	if c.String("outFile") != "" {
		outPath = c.String("outFile")
	}
	if err = tr.RenderSwagger(outPath, c.StringSlice("ifaces")...); err == nil {
		if c.String("redoc") != "" {
			var output []byte
			log.Infof("write to %s", c.String("redoc"))
			if output, err = exec.Command("redoc-cli", "bundle", outPath, "-o", c.String("redoc")).Output(); err != nil {
				log.WithError(err).Error(string(output))
			}
		}
	}
	return
}
