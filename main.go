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
	"runtime"
	"os"
	"github.com/SchumacherFM/wanderlust/github.com/codegangsta/cli"
)

func showHelp(c *cli.Context) {
	cli.ShowAppHelp(c)
}


func mainAction(c *cli.Context) {

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
			Flags:     []cli.Flag{
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
