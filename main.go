// Copyright (C) Cyrill@Schumacher.fm @SchumacherFM Twitter/GitHub
// Wanderlust - a cache warmer for your web app with priorities
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"github.com/SchumacherFM/wanderlust/github.com/codegangsta/cli"
	"github.com/SchumacherFM/wanderlust/github.com/davecheney/profile"
	"github.com/SchumacherFM/wanderlust/wlapp"
	"os"
	"runtime"
)

// The following vars will be injected during the build process via -ldflags.
var (
	Version string
	GitSHA  string
)

// mainAction will be executed when the CLI command run will be provided
func mainAction(c *cli.Context) {
	wlapp.CliContext = c
	wlapp.Boot()
}

func main() {
	setMaxParallelism()

	if "" != os.Getenv("WL_PPROF_CPU") {
		defer profile.Start(profile.CPUProfile).Stop()
	}
	if "" != os.Getenv("WL_PPROF_MEM") {
		defer profile.Start(profile.MemProfile).Stop()
	}

	app := cli.NewApp()
	app.Name = "Wanderlust"
	app.Version = Version + " [GitSHA: " + GitSHA + "]"
	app.Usage = "Wanderlust - a cache warmer for your web app with priorities"
	app.Action = showHelp
	app.Commands = []cli.Command{
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "Run the wanderlust app. `help run` for more information",
			Action:    mainAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "picnic-listen-address,pla",
					Value: "127.0.0.1:3008",
					Usage: "IP:Port address for picnic dashboard to listen on",
				},
				cli.StringFlag{
					Name:  "picnic-pem-dir,ppd",
					Value: "",
					Usage: "Directory to store the .pem certificates. If empty will throw it somewhere in the system. If provided file names must be cert.pem and key.pem!",
				},
				cli.StringFlag{
					Name:  "rucksack-dir,rd",
					Value: "",
					Usage: "Storage directory of the rucksack files. If empty then /tmp/random directory will be used.",
				},
				cli.StringFlag{
					Name:  "logFile,lf",
					Value: "",
					Usage: "Log to file or if empty to os.Stderr",
				},
				cli.StringFlag{
					Name:  "logLevel,ll",
					Value: "",
					Usage: "Log level: debug, info, notice, warning, error, critical, alert, emergency. Default: debug",
				},
			},
		},
	}
	app.Run(os.Args)
}

func showHelp(c *cli.Context) {
	cli.ShowAppHelp(c)
}

func setMaxParallelism() {
	if "" != os.Getenv("GOMAXPROCS") {
		return
	}
	maxProcs := runtime.GOMAXPROCS(0)
	mp := runtime.NumCPU()
	if maxProcs < mp {
		mp = maxProcs
	}
	runtime.GOMAXPROCS(mp)
}
