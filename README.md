# mds

A fork of [tlight's md](https://github.com/tlight/md) that uses the [Goldmark](https://github.com/yuin/goldmark) renderer instead of [Blackfriday](https://github.com/russross/blackfriday/). As is the original project:

> Zero configuration minimal markdown server for local rendering.

## Installation

```bash
go get github.com/dienakakim/mds
```

## Usage

```bash
mds --file=README.md [--port=8080] [--dark]
```

## License

Copyright &copy; 2020 Dien Tran. See LICENSE file. Enough parts of the code have been rewritten that I can safely pronounce it "original".

## P.S

`md` is a bad name for a project, as it clashes with the "make directory" command on Windows. This spurred me to initiate this fork, but along came various other features:

- Tailwind CSS styling support
- Dark theme (enabled with the flag `--dark`)
- Syntax highlighting, for your favorite language
