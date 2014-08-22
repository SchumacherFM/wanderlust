# Wanderlust - Work in Progress

### Cache warmer with priorities

Wanderlust, the yearning to travel through Google Analytics, Piwik, REST endpoints, sitemap.xml, etc to warm up the caches of your website.

### Cold Cache vs. Warm Cache ?

#### Cold Cache

There is analogy with cold engine and warm engine of the car. 

When the cache is empty or has irrelevant data, so that CPU needs to do a slower read from main memory for your program data requirement.

#### Warm Cache

When the cache contains relevant data, and all the reads for your program are satisfied from the cache itself.

## Architectural Overview

Please see the wanderlust.png image in the mindmap directory.

- Rucksack: Database [https://github.com/HouzuoGuo/tiedot](https://github.com/HouzuoGuo/tiedot)
- Picnic: Web Router: [https://github.com/codegangsta/negroni](https://github.com/codegangsta/negroni)

## Build

Using [https://github.com/laher/goxc](https://github.com/laher/goxc)

Setup go/src for darwin, linux and windows [http://dave.cheney.net/2013/07/09/an-introduction-to-cross-compilation-with-go-1-1](http://dave.cheney.net/2013/07/09/an-introduction-to-cross-compilation-with-go-1-1)

Run `make build`. If you are interested in pre-compiled binaries, ping me.

Using `godep` for dependency management.

# Contributing

Please do not hesitate to send pull requests.

# License

General Public License

[http://www.gnu.org/copyleft/gpl.html](http://www.gnu.org/copyleft/gpl.html)

# Author

[Cyrill Schumacher](https://github.com/SchumacherFM) - [My pgp public key](http://www.schumacher.fm/cyrill.asc)

[@SchumacherFM](https://twitter.com/SchumacherFM)

Made in Sydney, Australia :-)

If you consider a donation please contribute to: [http://www.seashepherd.org/](http://www.seashepherd.org/)

