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
	"github.com/SchumacherFM/wanderlust/github.com/robfig/cron"
	log "github.com/SchumacherFM/wanderlust/github.com/segmentio/go-log"
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/provisionerApi"
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
	Logger               *log.Logger
	BackPacker           rucksack.Backpacker
	ErrCronScheduleEmpty = errors.New("Cron Schedule is empty.")
	crond                = cron.New()
	bootCronOne          sync.Once
)

// BootCron starts the cron daemon once
func BootCron() {
	bootCronOne.Do(bootCron)
}

// bootCron internal start function
func bootCron() {

	// retrieve schedule collection
	sc, err := BackPacker.FindAll(bpDbConfig)
	if nil != err {
		Logger.Notice("%s",err)
	}

	for i := 0; i < len(sc); i = i + 2 {
		jobName := sc[i]
		jobSchedule := sc[i+1]
		Logger.Debug("%s: %s\n",jobName, jobSchedule)
	}

	Logger.Debug("Brotzeit cron daemon started!")
}

type (
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

// GetCollection returns a collection containing all the provisioners with their brotzeit cron schedule
// and UrlCount
func GetCollection(pc *provisionerApi.Provisioners, bp rucksack.Backpacker) (*BzConfigs, error) {

	bzcc := &BzConfigs{}
	bzcc.Collection = make([]*BzConfig, pc.Length())
	for i, p := range pc.Collection {

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
	return bp.Insert(bpDbConfig, bpCronKeyPrefix+f.Route, []byte(f.Schedule))
}
