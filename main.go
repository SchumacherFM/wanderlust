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
	"github.com/SchumacherFM/wanderlust/picnic"
	"log"
	"os"
	"runtime"
	"sync"
)

var wanderlustConfig *wanderlustApp

type wanderlustApp struct {
	cliContext *cli.Context
	waitGroup  sync.WaitGroup
	logger     *log.Logger
}

func init() {
	wanderlustConfig = &wanderlustApp{}
}

func (w *wanderlustApp) initLogger(logFile string) {
	if "" != logFile {
		logFilePointer, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		w.logger = log.New(logFilePointer, "", log.LstdFlags)
	} else {
		w.logger = log.New(os.Stderr, "", log.LstdFlags)
	}
}

// starts the HTTP server for the picnic web interface and runs it in a goroutine
func (w *wanderlustApp) bootPicnic() {
	w.waitGroup.Add(1)
	picnicApp := &picnic.PicnicApp{
		Port: uint(w.cliContext.Int("picnic-port")),
		Ip:   w.cliContext.String("picnic-ip"),
	}
	go func() {
		defer w.waitGroup.Done()
		picnicApp.Execute()
	}()
	w.logger.Printf("Picnic Running https://%s", picnicApp.GetListenAddress())
	w.waitGroup.Wait()
}

// mainAction will be executed when the CLI command run will be provided
func mainAction(c *cli.Context) {
	wanderlustConfig.cliContext = c
	wanderlustConfig.initLogger(c.String("logFile"))
	wanderlustConfig.bootPicnic()
}

func main() {
	if "" == os.Getenv("GOMAXPROCS") {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	app := cli.NewApp()
	app.Name = "Wanderlust"
	app.Version = "0.0.1"
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
					Name:  "picnic-ip,pip",
					Value: "",
					Usage: "IP address for picnic dashboard",
				},
				cli.IntFlag{
					Name:  "picnic-port,pp",
					Value: 3008,
					Usage: "Port for the picnic admin web interface",
				},
				cli.IntFlag{
					Name:  "rucksack-port,pr",
					Value: 3009,
					Usage: "Port for the rucksack JSON REST API @todo",
				},
				cli.StringFlag{
					Name:  "logFile,lf",
					Value: "",
					Usage: "Log to file or if empty to os.Stderr",
				},
				cli.StringFlag{
					Name:  "databaseDirectory,dd",
					Value: "",
					Usage: "Storage directory of the .db file",
				},
			},
		},
	}
	app.Run(os.Args)
}

func showHelp(c *cli.Context) {
	cli.ShowAppHelp(c)
}
