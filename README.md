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
| --- | --- | --- |
| admin | admin123 | yes |
| sugriwa | user123 | no |
| subali | user123 | no |


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

- On  Windows:

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
### Linux:
- Fedora 37 Workstation

### Go: 
- go version go1.19.7 linux/amd64

### MySQL: 
- Server version: 10.5.18-MariaDB MariaDB Server


### Bahasa Indonesia
Saya sendang berjuang sembuh dari Stroke pada 03 Oktober 2018-hari ini. Coding ini dirilis 7 Januari 2020, ternyata coding sedikit lupa. Kata-katanya dari Bahasa Indonesia sedikit lupa dan Bahasa Inggris kayaknya sulit. Insya Allah, perlahan-lahan sembuh. Aamiin.

Allah itu baik. 🙂

### English (translate[.]google[.]co[.]id)
I'm struggling to recover from a stroke on October 03, 2018-today. This coding was released January 7, 2020, apparently the coding was a little forgotten. The words from Indonesian are a little forgotten and English seems difficult. Insya Allah, slowly recover. Aamiin.

Allah is good. 🙂

---

Copyright © 2020 by Ocki Bagus Pratama
