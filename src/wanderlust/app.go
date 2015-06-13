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
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/codegangsta/cli"
	log "github.com/segmentio/go-log"
	"brotzeit"
	"picnic"
	"rucksack"
)

var (
	CliContext *cli.Context
	wg         sync.WaitGroup
	logger     *log.Logger
	bp         *rucksack.Rucksack
	bz         *brotzeit.Brotzeit
)

func Boot() {
	initLogger()
	BootRucksack()
	BootPicnic() // depends on the rucksack
	BootBrotzeit()
	BootWanderer()
	Finalizer()
}

func initLogger() {
	logFile := CliContext.String("logFile")
	logLevel := CliContext.String("logLevel")
	if "" != logFile {
		logFilePointer, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		logger = log.New(logFilePointer, log.DEBUG, "WL")
	} else {
		logger = log.New(os.Stderr, log.DEBUG, "WL")
	}
	if "" == logLevel {
		logLevel = "debug"
	}
	if lerr := logger.SetLevelString(logLevel); nil != lerr {
		panic(lerr)
	}
}

// BootRucksack inits the rucksack database and starts the background jobs
func BootRucksack() {
	var err error
	bp, err = rucksack.New(
		CliContext.String("rucksack-file"),
		logger,
	)
	logger.Check(err)

	wg.Add(1)
	// here can be added more services
	go func() {
		defer wg.Done()
		bp.Writer()
	}()
	logger.Notice("DB Background Services started")
}

// starts the HTTP server for the picnic web interface and runs it in a goroutine
func BootPicnic() {

	app, err := picnic.New(
		CliContext.String("picnic-listen-address"),
		CliContext.String("picnic-pem-dir"),
		logger,
		bp,
	)

	if nil != err {
		logger.Check(err)
	}

	if "" != app.GetListenAddress() { // don't start if empty
		wg.Add(1)
		go func() {
			defer wg.Done()
			logger.Check(app.Execute())
		}()
		url := "https://" + app.GetListenAddress()
		logger.Notice("Picnic Running %s", url)
		if true == CliContext.Bool("browser") {
			_, err := exec.Command("which", "open").Output()
			if err == nil {
				exec.Command("open", url).Output()
			}
		}
	}
}

func BootBrotzeit() {

	bz = brotzeit.New(logger, bp)
	wg.Add(1)
	go func() {
		defer wg.Done()
		bz.BootCron()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		bz.BootCronNotifier()
	}()
	logger.Notice("Booting Brotzeit ... ")
}

func BootWanderer() {
	logger.Notice("Booting Wanderer ... @todo")
}

// final method to wait on all the goroutines which are running mostly the HTTP server or other daemons
func Finalizer() {
	logger.Notice("GOMAXPROCS is set to %d", runtime.NumCPU())
	catchSysCall(bp, bz)
	wg.Wait()
}

// catchSysCall ends the program correctly when receiving a sys call
// @todo add things like remove PEM dir, DB dir when no CLI value has been provided
func catchSysCall(cl ...io.Closer) {
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
			logger.Notice("Received signal: %s. Closing %d app/s ...", sig.String(), len(cl))
			for _, c := range cl {
				if err := c.Close(); nil != err {
					logger.Check(err)
				} else {
					logger.Notice("Successful closed %T!", c)
				}
			}
			os.Exit(0)
		}
	}()
}

func UpdateAdminPassword() {
	done := make(chan bool)
	rucksack.WriterDone = func() {
		done <- true
	}
	rf := CliContext.String("rucksack-file")
	if _, err := os.Stat(rf); os.IsNotExist(err) || "" == rf {
		print("Rucksack File not found!\n") // @todo also check if file has a lock resp. wanderlust is running
		os.Exit(255)
	}
	initLogger()
	BootRucksack()
	admin := picnic.NewUserModel(bp, picnic.USER_ROOT)
	f, err := admin.FindMe()
	logger.Check(err)
	if false == f {
		logger.Emergency("User %s not found!", picnic.USER_ROOT)
		os.Exit(1)
	}
	if err := admin.GeneratePassword(); nil != err {
		logger.Check(err)
	}
	logger.Emergency("Changed password for user %s: %s", picnic.USER_ROOT, admin.Password)
	if err := admin.EncryptPassword(); nil != err {
		logger.Check(err)
	}
	ae, err := admin.Encode()
	logger.Check(err)
	bp.Insert(picnic.USER_DB_COLLECTION_NAME, admin.GetId(), ae)
	<-done
	bp.Close()
	os.Exit(0)
}
