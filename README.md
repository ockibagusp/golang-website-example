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

## Setting MySQL database

### Database 
file: golang-website-example.sql -> new database: golang-website-example

### Testing
file: golang-website-example.sql -> new database: golang-website-example_test

## Users Table

| Username | Password | is Admin |
| --- | --- | --- |
| admin | admin123 | yes |
| sugriwa | user123 | no |
| subali | user123 | no |


## Running app

### Compile and run Go program
```
$ go run main.go
```

or,

### Build compiles the packages

```
$ go build
```

- On Linux or Mac:

    ```
    $ ./golang-website-example
    ```

- On  Windows:

    ```
    $ golang-website-example.exe
    ```


## Router
This using [router](https://github.com/ockibagusp/golang-website-example/blob/master/router/router.go).


## httpexpect: Test for Debug
This using [debug](https://github.com/ockibagusp/golang-website-example/blob/master/tests/main_test.go).

Optional. Default value @debug: {true} or {1}.

1. function debug (bool)

    @function debug: {true} or {false}

2. os.Setenv("debug", ...)

    - @debug: {true} or {1}

        ```
        os.Setenv("debug", "true") 
        ```
        or,
        ```
        os.Setenv("debug", "1")
        ```

    - @debug: {false} or {0}
        ```
        os.Setenv("debug", "false") 
        ```
        or,
        ```
        os.Setenv("debug", "0")
        ```

### Test the packages

```
$ go test github.com/ockibagusp/golang-website-example/tests 
```

or, verbose output (-v) flag:

```
$ go test github.com/ockibagusp/golang-website-example/tests -v
```

## Air: cloud Live reload for Go apps
[Air](https://github.com/cosmtrek/air) is yet another live-reloading command line utility for Go applications in development.


## TODO List
- Admin Users: delete table User `deleted_at` @route: /users/admin/delete

    button:
    Restore | Delete Permanently

- Admin user button: delete not for admin
- too much

## Operating System (with me)
### Linux:
- Fedora 36 Workstation

### Go: 
- go1.18.5 linux/amd64

### MySQL: 
- mysql  Ver 8.0.27 for Linux on x86_64 (Source distribution)


### Bahasa Indonesia
Saya sendang berjuang sembuh dari Stroke pada 03 Oktober 2018-hari ini. Coding ini dirilis 7 Januari 2020, ternyata coding sedikit lupa. Kata-katanya dari Bahasa Indonesia sedikit lupa dan Bahasa Inggris kayaknya sulit. Insya Allah, perlahan-lahan sembuh. Aamiin.

Allah itu baik. 🙂

### English (translate[.]google[.]co[.]id)
I'm struggling to recover from a stroke on October 03, 2018-today. This coding was released January 7, 2020, apparently the coding was a little forgotten. The words from Indonesian are a little forgotten and English seems difficult. Insya Allah, slowly recover. Aamiin.

Allah is good. 🙂

---

Copyright © 2020 by Ocki Bagus Pratama
