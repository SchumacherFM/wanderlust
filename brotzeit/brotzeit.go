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
)

const (
	// bpConfigName bucket/database contains mainly the cron configuration
	bpConfigName           = "bzcfg"
	// bpCronKeyPrefix is the key prefix for identifying a real cron schedule
	bpCronKeyPrefix        = "cronExpression"
	// bpCollectionNamePrefix bucket/database contains millions of URLs
	bpCollectionNamePrefix = "bzcol"
)

var (
	Logger               *log.Logger
	ErrCronScheduleEmpty = errors.New("Cron Schedule is empty.")
)

func BootMe() error {

	// job/worker

	//	pc, err := provisioners.GetAvailable()
	//	if nil != err {
	//		return err
	//	}
	//	for _, prov := range pc.Collection {
	//
	//		go prov.Api.ProduceUrls()
	//
	//	}
	return nil
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
		cd, _ := bp.FindOne(bpConfigName, bpCronKeyPrefix+p.Api.Route())

		uc, _ := bp.Count(bpCollectionNamePrefix + p.Api.Route())

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
		bp.Delete(bpConfigName, bpCronKeyPrefix+f.Route)
		return nil
	}

	_, err := cron.Parse(f.Schedule)
	if nil != err {
		return err
	}
	return bp.Insert(bpConfigName, bpCronKeyPrefix+f.Route, []byte(f.Schedule))
}
