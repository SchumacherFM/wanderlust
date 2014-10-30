# Wanderlust - WIP

### Cache warmer with priorities

Wanderlust, the yearning to travel through Google Analytics (1), Piwik (1), 
REST endpoints (1), sitemap.xml, textarea input, etc to warm up the caches of your website.

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

Using [https://github.com/laher/goxc](https://github.com/laher/goxc)

Setup go/src for darwin, linux and windows [http://dave.cheney.net/2013/07/09/an-introduction-to-cross-compilation-with-go-1-1](http://dave.cheney.net/2013/07/09/an-introduction-to-cross-compilation-with-go-1-1)

Run `make build`. If you are interested in pre-compiled binaries, ping me.

Dependencies to external packages are all included in this repository. 
Yes there are apps like `godep` [here](https://coreos.com/blog/godep-for-end-user-go-projects/) 
to handle this but I decided 
to implement each package directly. Maybe some parts will be moved out of the 
repository, mainly those which are on github and keep those which needs hg or bzr.

Cross build inside docker container [https://github.com/docker-library/docs/tree/master/golang](https://github.com/docker-library/docs/tree/master/golang)

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

# License

Will change in the future.

General Public License *may change the license*

[http://www.gnu.org/copyleft/gpl.html](http://www.gnu.org/copyleft/gpl.html)

# Author

[Cyrill Schumacher](https://github.com/SchumacherFM) - [My pgp public key](http://www.schumacher.fm/cyrill.asc)

[@SchumacherFM](https://twitter.com/SchumacherFM)

Made in Sydney, Australia :-)
