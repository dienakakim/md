# mds

A fork of [tlight's md](https://github.com/tlight/md) that uses the [goldmark](https://github.com/yuin/goldmark) renderer instead of [blackfriday](https://github.com/russross/blackfriday/). As is the original project:

> Zero configuration minimal markdown server for local rendering.

## Installation

```bash
go get github.com/dienakakim/mds
```

## Usage

```bash
mds README.md
```

## License

Copyright (c) 2018 Tim, (c) 2020 Dien Tran. See LICENSE file.

## P.S

`md` is a bad name for a project, as it clashes with the "make directory" command on Windows. This spurred me to initiate this fork, but that is not the only fix this fork makes: upcoming features include syntax highlighting and a dark themed render.
