# golang-website-example

Golang Echo and html template.

> move GitHub repository for [hello](https://github.com/ockibagusp/hello) to [golang-website-example](#)

### Visual Studio Code

#### Run and Debug: [launch.json](https://github.com/ockibagusp/golang-website-example/blob/master/.vscode/launch.json).

## Getting Started

First, clone the repo:

```bash
$ git clone https://github.com/ockibagusp/golang-website-example.git
```

Repo dep

```bash
$ make dep
```

## Setting MySQL database

### Database

file: golang-website-example.sql -> new database: golang-website-example

### Testing

file: golang-website-example.sql -> new database: golang-website-example-test

## Users Table

| Username | Password | is Admin |
| -------- | -------- | -------- |
| admin    | admin123 | yes      |
| sugriwa  | user123  | no       |
| subali   | user123  | no       |

## Running app

### Compile and run Go program

```
$ make run
```

or,

### Build compiles the packages

```
$ make build
```

- On Linux or Mac:

  ```
  $ ./main
  ```

- On Windows:

  ```
  $ main.exe
  ```

## Router

This using [router](https://github.com/ockibagusp/golang-website-example/blob/master/api/main/router/router.go).

### Test the packages

Example,

#### test all flag

```bash
$ make test
```

#### test controller flag

```bash
$ make test-ctrl
```

or, verbose output (-v) flag:

#### test verbose all flag

```bash
$ make test-verbose
```

#### test verbose controller flag

```bash
$ make test-verbose-ctrl
```

or, cover all flag

#### test cover

```bash
$ make cover
```

#### cover show flag:

```bash
$ make cover-show
```

#### cover html flag:

```bash
$ make cover-html
```

#### cover select function (-func) flag:

```bash
$ make cover-func
```

## Air: cloud Live reload for Go apps

[Air](https://github.com/cosmtrek/air) is yet another live-reloading command line utility for Go applications in development.

## Live Server (npm)

[Live Server](https://www.npmjs.com/package/live-server) this is a little development server with live reload capability. Use it for hacking your HTML/JavaScript/CSS files, but not for deploying the final site.

## TODO List

- too much

## Operating System (with me)

### macOS:

- macOS Ventura 13.3.1 (a)

### Go:

- go version go1.20.3 darwin/arm64

### MySQL:

- Active Instance: MySQL 8.0.32-arm64

- MySQL Workbench Community (GPL) for macOS version 8.0.33 CE build 2947366 (64 bit)

### Bahasa Indonesia

Saya sendang berjuang sembuh dari Stroke pada 03 Oktober 2018-hari ini. Saya dirilis 7 Januari 2020 meng-coding sedikit lupa. Mulai dari Bahasa Indonesia nol sampai saat ini; sekarang sedikit lupa. Sama, Bahasa Inggris nol sampai saat ini; sekarang banyak sulit. Belajar lagi. Insya Allah, perlahan-lahan sembuh. Aamiin.

Allah itu baik. 🙂

### English (translate[.]google[.]co[.]id)

I am currently struggling to recover from a stroke on October 3, 2018-today. I released January 7 2020 coding was a little forgotten. Starting from Indonesian zero until now; now a little forgot. Same, English zero so far; now much difficult. Study more. Insya Allah, slowly healed. Aamiin.

Allah is good. 🙂

---

Copyright © 2020 by Ocki Bagus Pratama
