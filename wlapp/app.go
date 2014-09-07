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

package wlapp

import (
	"github.com/SchumacherFM/wanderlust/github.com/HouzuoGuo/tiedot/db"
	"github.com/SchumacherFM/wanderlust/github.com/codegangsta/cli"
	"github.com/SchumacherFM/wanderlust/picnic"
	"github.com/SchumacherFM/wanderlust/rucksack"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

type WanderlustApp struct {
	CliContext *cli.Context
	waitGroup  sync.WaitGroup
	Logger     *log.Logger
	db         *db.DB
}

// final method to wait on all the goroutines which are running mostly the HTTP server or other daemons
func (w *WanderlustApp) Finalizer() {
	w.Logger.Printf("GOMAXPROCS is set to %d", runtime.NumCPU())
	w.catchSysCall()
	w.waitGroup.Wait()
}

func (w *WanderlustApp) InitLogger(logFile string) {
	if "" != logFile {
		logFilePointer, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		w.Logger = log.New(logFilePointer, "", log.LstdFlags)
	} else {
		w.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}
}

// starts the HTTP server for the picnic web interface and runs it in a goroutine
func (w *WanderlustApp) BootPicnic() {

	picnicApp, err := NewPicnicApp(
		w.CliContext.String("picnic-listen-address"),
		w.CliContext.String("picnic-pem-dir"),
		w.Logger,
	)

	if nil != err {
		w.Logger.Fatal(err)
	}

	if "" != picnicApp.GetListenAddress() { // don't start if empty
		w.waitGroup.Add(1)
		go func() {
			defer w.waitGroup.Done()
			err := picnicApp.Execute()
			if nil != err {
				w.Logger.Fatal(err)
			}
		}()
		w.Logger.Printf("Picnic Running https://%s", picnicApp.GetListenAddress())
	}
}

func (w *WanderlustApp) BootWanderer() {
	w.Logger.Print("Booting Wanderer ... @todo")
}

func (w *WanderlustApp) BootBrotzeit() {
	w.Logger.Print("Booting Brotzeit ... @todo")
}

// @todo remove this and add it to the picnic app as a feature to start and stop the DB backend via web panel
// inits the rucksack and boots on the default http mux
func (w *WanderlustApp) BootRucksack() {
	rucksackApp := &rucksack.RucksackApp{
		ListenAddress: w.CliContext.String("rucksack-listen-address"),
		DbDir:         w.CliContext.String("rucksack-dir"),
		Logger:        w.Logger,
	}
	rucksackApp.InitDb()
	w.db = rucksackApp.GetDb()
	if "" != rucksackApp.ListenAddress {
		w.waitGroup.Add(1)
		go func() {
			defer w.waitGroup.Done()
			rucksackApp.StartHttp()
		}()
		w.Logger.Printf("Rucksack Running http://%s", rucksackApp.GetListenAddress())
	}
}

// catchSysCall ends the program correctly when receiving a sys call
// @todo add things like remove PEM dir, DB dir when no CLI value has been provided
func (w *WanderlustApp) catchSysCall() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(
		signalChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for sig := range signalChannel {
			w.Logger.Printf("Received signal: %s\n", sig.String())

			os.Exit(0)
		}
	}()
}
