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
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
	"github.com/SchumacherFM/wanderlust/github.com/robfig/cron"
	log "github.com/SchumacherFM/wanderlust/github.com/segmentio/go-log"
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/provisionerApi"
	"github.com/SchumacherFM/wanderlust/provisioners"
	"github.com/SchumacherFM/wanderlust/rucksack"
	"net/http"
	"sync"
)

const (
	// bpDbConfig bucket/database contains mainly the cron configuration
	bpDbConfig = "bzcfg"
	// bpCronKeyPrefix is the key prefix for identifying a real cron schedule
	bpCronKeyPrefix = "cronExpression"
	// bpDbUrlsPrefix bucket/database contains millions of URLs
	bpDbUrlsPrefix = "bzcol"
)

var (
	ErrCronScheduleEmpty = errors.New("Cron Schedule is empty.")
	crond                = cron.New()
	bootCronOne          sync.Once
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

// BootCron starts the cron daemon once
func (b *Brotzeit) BootCron() {
	bootCronOne.Do(b.bootCron)
}

// bootCron internal start function to start the cronjobs
func (b *Brotzeit) bootCron() {
	for _, p := range b.pc.Collection() {
		s, _ := b.bp.FindOne(bpDbConfig, bpCronKeyPrefix+p.Api.Route())
		if nil == s {
			continue
		}

		ok, err := p.Api.ConfigComplete(b.bp)
		b.errWar(err)
		if true == ok {
			js := string(s)
			if err := crond.AddFunc(js, p.Api.FetchURLs(b.bp, b.l)); nil != err {
				b.errWar(err)
			} else {
				b.l.Debug("Cron added for: %s, Schedule: %s", p.Api.Route(), js)
			}
		}
	}

	if len(crond.Entries()) > 0 {
		crond.Start()
	}

	b.l.Debug("Brotzeit cron daemon started!")
}

func (b *Brotzeit) errWar(e error) {
	if nil != e {
		b.l.Warning(errgo.Details(errgo.New(e.Error())))
	}
}

// Close shuts down the cron job when the app receives a signal
func (b *Brotzeit) Close() error {
	if len(crond.Entries()) > 0 {
		crond.Stop()
	}
	return nil
}

// GetCollection returns a collection containing all the provisioners with their brotzeit cron schedule
// and UrlCount
func GetCollection(pc provisionerApi.Collectioner, bp rucksack.Backpacker) (*BzConfigs, error) {

	var bzcc = &BzConfigs{
		Collection: make([]*BzConfig, len(pc.Collection())),
	}

	for i, p := range pc.Collection() {

		// possibility that database and key will not exists but we ignore that
		cd, _ := bp.FindOne(bpDbConfig, bpCronKeyPrefix+p.Api.Route())
		uc, _ := bp.Count(bpDbUrlsPrefix + p.Api.Route())
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
		bp.Delete(bpDbConfig, bpCronKeyPrefix+f.Route)
		return nil
	}

	_, err := cron.Parse(f.Schedule)
	if nil != err {
		return err
	}
	// @todo onSave event trigger the cron: check if config complete, add job, and press start
	return bp.Insert(bpDbConfig, bpCronKeyPrefix+f.Route, []byte(f.Schedule))
}
