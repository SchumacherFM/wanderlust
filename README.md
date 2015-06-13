# Wanderlust - WIP

### Cache warmer/primer with priorities

Wanderlust, the yearning to travel through Google Analytics (1), Piwik (1), 
REST endpoints (1), sitemap.xml, textarea input, etc to warm up the caches of your website.

An app for priming your cache to help avoid a [cache stampede](http://en.wikipedia.org/wiki/Cache_stampede).

(1) available to buy.

### Cold Cache vs. Warm Cache ?

#### Cold Cache

There is analogy with cold engine and warm engine of the car. 

When the cache is empty or has irrelevant data, 
so that CPU needs to do a slower read from main memory 
for your program data requirement.

#### Warm Cache

When the cache contains relevant data, 
and all the reads for your program are satisfied from the cache itself.

## Architectural Overview

![mindmap](https://raw.githubusercontent.com/SchumacherFM/wanderlust/master/mindmap/wanderlust.png)

- Ideas: [https://github.com/GordonLesti/Lesti_Fpc/issues/25](https://github.com/GordonLesti/Lesti_Fpc/issues/25)

## Build

Install [http://getgb.io/](http://getgb.io/)

```
$ go get -u github.com/constabulary/gb/...
$ git clone --recursive git@github.com:SchumacherFM/wanderlust.git
$ gb build all
$ ./bin/wanderlust
```

@todo: Review special build process when including static assets.


## Profiling

To start the program with CPU profiling: `$ WL_PPROF_CPU=1 ./wanderlust [options] [arguments]`

To start the program with Memory profiling: `$ WL_PPROF_MEM=1 ./wanderlust [options] [arguments]`

To read the generated profile:

```
$ go tool pprof ./wanderlust /path/to/cpu.pprof
```

To generated fancy graphics to stdout:

```
$ go tool pprof [--pdf|svg|gif] ./wanderlust /path/to/cpu.pprof
```

# Contributing

Please do not hesitate to send pull requests or tweet me.

# Free vs. Commercial Version

### Idea

The free version includes the provisioners for loading URLs from the sitemap.xml and from a textarea.
The free version implements a lot of tracking via Piwik. The tracking data will be send to my
private Piwik server and of course not shared with others. There is a privacy statement in the app.

As soon as you buy (via in-app-purchase) one provisioner you'll get the commercial version. 
Each provisioner has a recurring monthly fee. Provisioners must be compiled into this source code so delivery 
of the new binary file may take some time. With the commercial version tracking is removed.

Available build-in free provisioners:

- Sitemap (up to XX URLs crawled randomly) 
- Textarea (up to XX URLs crawled randomly)

Available In-App provisioners for purchase:

- Google Analytics
- Piwik
- Any REST endpoint
- Magento, TYPO3, Drupal, WordPress, ... (2)
- Any other analytics service which has a data out API

(2) Module is open source but for Wanderlust integration a recurring fee is required.

##### Some thoughts

In-App-Purchases may take a while to deliver in the early days because the Go binaries needs to be recompiled.

Depending on how much you contribute you can get modules for free or even share in revenue.

Open Source ain’t Charity.

The shop system is provided via a Magento REST API hosted on my server.

Implement: https://github.com/headzoo/surf

# License

Will change in the future.

General Public License *may change the license*

[http://www.gnu.org/copyleft/gpl.html](http://www.gnu.org/copyleft/gpl.html)

# Author

[Cyrill Schumacher](https://github.com/SchumacherFM) - [My pgp public key](http://www.schumacher.fm/cyrill.asc)

[@SchumacherFM](https://twitter.com/SchumacherFM)

Made in Sydney, Australia :-)
