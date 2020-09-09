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

TODO
[embedmd]:# (.tmp/usage.txt console)
```console
foo@bar:~$ prefix -h
Usage of prefix:
  -format string
    	format string (default "{{.LineNumber3}} upt={{.Uptime}} dur={{.Duration}} ")
```

[embedmd]:# (.tmp/example-1.txt console)
```console
foo@bar:~$ generate-fake-data | prefix -format="#{{.LineNumber3}} {{.Uptime}} {{.Duration}} | "
#1   61.4µs  69.5µs  | At illum ut est sit soluta nulla numquam.
#2   113.8ms 113.7ms | Sunt quaerat ea dolores facere deleniti culpa numquam.
#3   328.8ms 215.1ms | Distinctio maxime consequatur est qui corporis sunt officia.
#4   607.2ms 278.3ms | Et quia odit molestias voluptas porro repellendus magnam.
#5   899.5ms 292.3ms | Corporis eos rem non hic esse optio quisquam.
#6   1.1s    211.7ms | Natus earum molestias iste architecto porro et blanditiis.
#7   1.3s    238.1ms | Eum repellendus nostrum qui eius suscipit fugit quia.
#8   1.4s    50.6ms  | Et nesciunt quod fuga ut vel pariatur libero.
#9   1.6s    209.7ms | Rerum omnis soluta facilis voluptatem possimus et voluptas.
#10  1.9s    274.8ms | Possimus harum voluptatibus aperiam voluptatibus qui autem quam.
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
