# prefix

🔴 prepend numbers, stats, dates, durations to streams

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/moul.io/prefix)
[![License](https://img.shields.io/badge/license-Apache--2.0%20%2F%20MIT-%2397ca00.svg)](https://github.com/moul/prefix/blob/master/COPYRIGHT)
[![GitHub release](https://img.shields.io/github/release/moul/prefix.svg)](https://github.com/moul/prefix/releases)
[![Docker Metrics](https://images.microbadger.com/badges/image/moul/prefix.svg)](https://microbadger.com/images/moul/prefix)
[![Made by Manfred Touron](https://img.shields.io/badge/made%20by-Manfred%20Touron-blue.svg?style=flat)](https://manfred.life/)

[![Go](https://github.com/moul/prefix/workflows/Go/badge.svg)](https://github.com/moul/prefix/actions?query=workflow%3AGo)
[![Release](https://github.com/moul/prefix/workflows/Release/badge.svg)](https://github.com/moul/prefix/actions?query=workflow%3ARelease)
[![PR](https://github.com/moul/prefix/workflows/PR/badge.svg)](https://github.com/moul/prefix/actions?query=workflow%3APR)
[![GolangCI](https://golangci.com/badges/github.com/moul/prefix.svg)](https://golangci.com/r/github.com/moul/prefix)
[![codecov](https://codecov.io/gh/moul/prefix/branch/master/graph/badge.svg)](https://codecov.io/gh/moul/prefix)
[![Go Report Card](https://goreportcard.com/badge/moul.io/prefix)](https://goreportcard.com/report/moul.io/prefix)
[![CodeFactor](https://www.codefactor.io/repository/github/moul/prefix/badge)](https://www.codefactor.io/repository/github/moul/prefix)


## Usage

[embedmd]:# (.tmp/usage.txt console)
```console
foo@bar:~$ prefix -h
USAGE
  prefix [flags] file

FLAGS
  -format string
    	format string (default "{{DEFAULT}} ")

SYNTAX
  {{.Duration | short_duration}}      {{.Duration}} displayed in a pretty & short format (len<=7)
  {{.Duration}}                       time since previous line was started
  {{.Format}}                         the value you set with -format
  {{.LineNumber3}}                    alias for {{printf "%-3d" .LineNumber}}
  {{.LineNumber4}}                    alias for {{printf "%-4d" .LineNumber}}
  {{.LineNumber5}}                    alias for {{printf "%-5d" .LineNumber}}
  {{.LineNumber}}                     display line number
  {{.ShortDuration}}                  alias for {{.Duration | short_duration}}
  {{.ShortUptime}}                    alias for {{.Uptime | short_duration}}
  {{.Uptime | short_duration}}        {{.Uptime}} displayed in a pretty & short format (len<=7)
  {{.Uptime}}                         time since the the prefixer was initialized
  {{env "USER"}}                      replace with content of the $USER env var
  {{now | unixEpoch}}                 current timestamp
  {{now}}                             current date (format: 2006-01-02 15:04:05.999999999 -0700 MST)
  {{uuidv4}}                          UUID of the v4 (randomly generated) type

  the following helpers are also available:
  - from the text/template library    https://golang.org/pkg/text/template/
  - from the sprig project            https://github.com/masterminds/sprig#usage

PRESETS
  {{DEFAULT}}          {{.LineNumber3}} up={{.ShortUptime}} d={{.ShortDuration}} |
  {{SHORT_DATE}}       {{now | date "06/02/01 15:04:05"}}
  {{SLOW_LINES}}       {{if (gt .Duration 1000000000)}}SLOW{{else}}    {{end}} {{.Duration | short_duration}} 

EXAMPLES
  prefix apache.log
  prefix -format=">>>" apache.log
  tail -f apache.log | prefix -
  my-cool-program 2>&1 | prefix -format="#{{.LineNumber5}} " -
```

[embedmd]:# (.tmp/example-1.txt console)
```console
foo@bar:~$ generate-fake-data | prefix -format="#{{.LineNumber3}} {{.ShortUptime}} {{.ShortDuration}} | "
#1   73.6µs  80.9µs  | At illum ut est sit soluta nulla numquam.
#2   112ms   111.9ms | Sunt quaerat ea dolores facere deleniti culpa numquam.
#3   327.1ms 215.1ms | Distinctio maxime consequatur est qui corporis sunt officia.
#4   605.3ms 278.2ms | Et quia odit molestias voluptas porro repellendus magnam.
#5   897.7ms 292.3ms | Corporis eos rem non hic esse optio quisquam.
#6   1.1s    211.7ms | Natus earum molestias iste architecto porro et blanditiis.
#7   1.3s    238.2ms | Eum repellendus nostrum qui eius suscipit fugit quia.
#8   1.4s    50.6ms  | Et nesciunt quod fuga ut vel pariatur libero.
#9   1.6s    209.6ms | Rerum omnis soluta facilis voluptatem possimus et voluptas.
#10  1.9s    274.7ms | Possimus harum voluptatibus aperiam voluptatibus qui autem quam.
```

[embedmd]:# (.tmp/example-2.txt console)
```console
foo@bar:~$ generate-fake-data | prefix -format="{{.LineNumber3}} "
1   At illum ut est sit soluta nulla numquam.
2   Nobis sunt quaerat ea dolores facere deleniti culpa.
3   Numquam ut distinctio maxime consequatur est qui corporis.
4   Sunt officia odit et quia odit molestias voluptas.
5   Porro repellendus magnam ipsa corporis eos rem non.
6   Hic esse optio quisquam hic natus earum molestias.
7   Iste architecto porro et blanditiis iste eum repellendus.
8   Nostrum qui eius suscipit fugit quia quo et.
9   Nesciunt quod fuga ut vel pariatur libero sequi.
10  Rerum omnis soluta facilis voluptatem possimus et voluptas.
```

[embedmd]:# (.tmp/example-3.txt console)
```console
foo@bar:~$ generate-fake-data | prefix -format=">>> "
>>> At illum ut est sit soluta nulla numquam.
>>> Nobis sunt quaerat ea dolores facere deleniti culpa.
>>> Numquam ut distinctio maxime consequatur est qui corporis.
>>> Sunt officia odit et quia odit molestias voluptas.
>>> Porro repellendus magnam ipsa corporis eos rem non.
>>> Hic esse optio quisquam hic natus earum molestias.
>>> Iste architecto porro et blanditiis iste eum repellendus.
>>> Nostrum qui eius suscipit fugit quia quo et.
>>> Nesciunt quod fuga ut vel pariatur libero sequi.
>>> Rerum omnis soluta facilis voluptatem possimus et voluptas.
```

[embedmd]:# (.tmp/example-4.txt console)
```console
foo@bar:~$ generate-fake-data | prefix -format="{{SLOW_LINES}} up={{.ShortUptime}} | "
     876.7µs  up=892.6µs | Rerum natus quo quo explicabo tempore et delectus.
SLOW 1s       up=1.1s    | Dolor blanditiis voluptas dolorum sint laudantium eveniet amet.
SLOW 1.3s     up=2.4s    | Qui asperiores molestiae est quia est eum omnis.
SLOW 1.3s     up=3.6s    | Illum explicabo aut illum iste pariatur aut laudantium.
     982.2ms  up=4.6s    | Quibusdam asperiores consequatur est dolores quas dolor ipsam.
     185.5ms  up=4.8s    | Possimus qui non rem qui cum sit temporibus.
     167.5ms  up=5s      | Ea debitis sit deleniti cum ut adipisci in.
     520.3ms  up=5.5s    | Eveniet molestias voluptatem voluptatem deserunt nisi tempora iusto.
     215ms    up=5.7s    | Fugiat minus quam eos voluptatem labore sit velit.
SLOW 1s       up=6.7s    | Enim aut autem tenetur fugit minima quo atque.
```

[embedmd]:# (.tmp/example-5.txt console)
```console
foo@bar:~$ generate-fake-data | prefix -format="{{SHORT_DATE}} "
20/19/09 00:02:43 At illum ut est sit soluta nulla numquam.
20/19/09 00:02:44 Sunt quaerat ea dolores facere deleniti culpa numquam.
20/19/09 00:02:45 Distinctio maxime consequatur est qui corporis sunt officia.
20/19/09 00:02:46 Et quia odit molestias voluptas porro repellendus magnam.
20/19/09 00:02:48 Corporis eos rem non hic esse optio quisquam.
20/19/09 00:02:49 Natus earum molestias iste architecto porro et blanditiis.
20/19/09 00:02:50 Eum repellendus nostrum qui eius suscipit fugit quia.
20/19/09 00:02:50 Et nesciunt quod fuga ut vel pariatur libero.
20/19/09 00:02:51 Rerum omnis soluta facilis voluptatem possimus et voluptas.
20/19/09 00:02:53 Possimus harum voluptatibus aperiam voluptatibus qui autem quam.
```


## Install

### Using go

```console
$ go get -u moul.io/prefix
```

### Releases

See https://github.com/moul/prefix/releases

## Contribute

![Contribute <3](https://raw.githubusercontent.com/moul/moul/master/contribute.gif)

I really welcome contributions. Your input is the most precious material. I'm well aware of that and I thank you in advance. Everyone is encouraged to look at what they can do on their own scale; no effort is too small.

Everything on contribution is sum up here: [CONTRIBUTING.md](./CONTRIBUTING.md)

### Contributors ✨

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-2-orange.svg)](#contributors)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="http://manfred.life"><img src="https://avatars1.githubusercontent.com/u/94029?v=4" width="100px;" alt=""/><br /><sub><b>Manfred Touron</b></sub></a><br /><a href="#maintenance-moul" title="Maintenance">🚧</a> <a href="https://github.com/moul/prefix/commits?author=moul" title="Documentation">📖</a> <a href="https://github.com/moul/prefix/commits?author=moul" title="Tests">⚠️</a> <a href="https://github.com/moul/prefix/commits?author=moul" title="Code">💻</a></td>
    <td align="center"><a href="https://manfred.life/moul-bot"><img src="https://avatars1.githubusercontent.com/u/41326314?v=4" width="100px;" alt=""/><br /><sub><b>moul-bot</b></sub></a><br /><a href="#maintenance-moul-bot" title="Maintenance">🚧</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!

### Stargazers over time

[![Stargazers over time](https://starchart.cc/moul/prefix.svg)](https://starchart.cc/moul/prefix)

## License

© 2020 [Manfred Touron](https://manfred.life)

Licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) ([`LICENSE-APACHE`](LICENSE-APACHE)) or the [MIT license](https://opensource.org/licenses/MIT) ([`LICENSE-MIT`](LICENSE-MIT)), at your option. See the [`COPYRIGHT`](COPYRIGHT) file for more details.

`SPDX-License-Identifier: (Apache-2.0 OR MIT)`
