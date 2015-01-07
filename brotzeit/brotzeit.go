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

package brotzeit

import (
	"errors"
	"net/http"

	"github.com/SchumacherFM/wanderlust/Godeps/_workspace/src/github.com/juju/errgo"
	log "github.com/SchumacherFM/wanderlust/Godeps/_workspace/src/github.com/segmentio/go-log"
	"github.com/SchumacherFM/wanderlust/github.com/robfig/cron"
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/provisionerApi"
	"github.com/SchumacherFM/wanderlust/provisioners"
	"github.com/SchumacherFM/wanderlust/rucksack"
)

const (
	// dbConfig bucket/database contains mainly the cron configuration
	dbConfig = "bzcfg"
	// dbCronKeyPrefix is the key prefix for identifying a real cron schedule
	dbCronKeyPrefix = "cronExpression"
	// dbURLPrefix bucket/database contains millions of URLs
	dbURLPrefix = "bzcol"
)

var (
	ErrCronScheduleEmpty = errors.New("Cron Schedule is empty.")
	crond                = cron.New()
	// @todo implement http://talks.golang.org/2013/advconc.slide
	chanCronAdd = make(chan *BzConfig, 1)
	chanCronDel = make(chan *BzConfig, 1)
)

type (
	Brotzeit struct {
		l  *log.Logger
		bp rucksack.Backpacker
		pc provisionerApi.Collectioner
	}

	BzConfig struct {
		Route    string
		Name     string
		Icon     string
		Schedule string
		UrlCount int
	}
	BzConfigs struct {
		// C holds the Collection of BzConfig
		Collection []*BzConfig
	}
)

func New(l *log.Logger, bp rucksack.Backpacker) *Brotzeit {
	pc, _ := provisioners.GetAvailable()
	return &Brotzeit{
		l:  l.New("BZ"),
		bp: bp,
		pc: pc,
	}
}

// BootCron runs in a goroutine and starts the cron daemon only once!
func (b *Brotzeit) BootCron() {

	for _, p := range b.pc.Collection() {
		// get cron schedule
		s, _ := b.bp.FindOne(dbConfig, dbCronKeyPrefix+p.Api.Route())
		if nil == s {
			continue
		}

		ok, err := p.Api.ConfigComplete(b.bp)
		b.errCheck(err)
		if true == ok {
			js := string(s)
			if err := crond.AddFunc(js, p.Api.FetchURLs(b.bp, b.l)); nil != err {
				b.errCheck(err)
			} else {
				b.l.Debug("Cron added for: %s, Schedule: %s", p.Api.Route(), js)
			}
		}
	}
	startCron()
	b.l.Debug("Brotzeit cron daemon started!")
}

// BootCronNotifier runs in a goroutine and listens for the chanCronAdd channel. Once there is a save event
// it handles the starting and stopping of the cron jobs
func (b *Brotzeit) BootCronNotifier() {
	for {
		select {
		case bzAdd := <-chanCronAdd:
			b.l.Debug("BootCronNotifier Add: %#v", bzAdd)
			// get API
			// check config complete
			// if not remove cron, if running
			// if complete start cron
			// would be nice to use the websocket to notify the user about this steps here
		case bzDel := <-chanCronDel:
			b.l.Debug("BootCronNotifier Del: %#v", bzDel)
		}
	}

}

func (b *Brotzeit) errCheck(e error) {
	if nil != e {
		b.l.Warning(errgo.Details(errgo.New(e.Error())))
	}
}

// Close shuts down the cron job when the app receives a signal
func (b *Brotzeit) Close() error {
	if len(crond.Entries()) > 0 {
		crond.Stop()
	}
	// not really necessary but who knows ;-)
	close(chanCronAdd)
	return nil
}

// startCron starts the cron service if there are any entries
func startCron() {
	if len(crond.Entries()) > 0 && false == crond.IsRunning() {
		crond.Start()
	}
}

// GetCollection returns a collection containing all the provisioners with their brotzeit cron schedule
// and UrlCount
func GetCollection(pc provisionerApi.Collectioner, bp rucksack.Backpacker) (*BzConfigs, error) {

	var bzcc = &BzConfigs{
		Collection: make([]*BzConfig, len(pc.Collection())),
	}

	for i, p := range pc.Collection() {

		// possibility that database and key will not exists but we ignore that
		cd, _ := bp.FindOne(dbConfig, dbCronKeyPrefix+p.Api.Route())
		uc, _ := bp.Count(dbURLPrefix + p.Api.Route())
		bzc := &BzConfig{
			Route:    p.Api.Route(),
			Name:     p.Name,
			Icon:     p.Icon,
			Schedule: string(cd),
			UrlCount: uc,
		}
		bzcc.Collection[i] = bzc
	}

	return bzcc, nil
}

// SaveConfig stores the cron schedule in the backpacker. If the schedule is empty then it will be deleted from the DB
func SaveConfig(bp rucksack.Backpacker, r *http.Request) error {

	f := &BzConfig{}
	helpers.DecodeJSON(r, f)

	if "" == f.Route {
		return ErrCronScheduleEmpty
	}

	if "" == f.Schedule {
		chanCronDel <- f
		bp.Delete(dbConfig, dbCronKeyPrefix+f.Route)
		return nil
	}

	_, err := cron.Parse(f.Schedule)
	if nil != err {
		return err
	}
	// @todo onSave event trigger the cron: check if config complete, add job, and press start
	chanCronAdd <- f
	return bp.Insert(dbConfig, dbCronKeyPrefix+f.Route, []byte(f.Schedule))
}
