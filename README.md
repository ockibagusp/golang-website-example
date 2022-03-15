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

### Setting MySQL database

#### Database 
file: golang-website-example.sql -> new database: golang-website-example

#### Testing
file: golang-website-example.sql -> new database: golang-website-example_test

## Users Table

| Username | Password | is Admin |
| --- | --- | --- |
| admin | admin123 | yes |
| sugriwa | user123 | no |
| subali | user123 | no |


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

### Running app

#### Compile and run Go program
```
$ go run main.go
```

or,

#### Build compiles the packages

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

#### Test the packages

```
$ go test github.com/ockibagusp/golang-website-example/tests 
```

or, verbose output (-v) flag:

```
$ go test github.com/ockibagusp/golang-website-example/tests -v
```


## TODO List
- Admin Users: delete table User `deleted_at` @route: /users/admin/delete

    button:
    Restore | Delete Permanently

- Admin user button: delete not for admin
- Admin user search
- mock unit test
- list pagination with next, previous, first and last
- moves files function Server and NewServer, etc.
- Mutex: BankAccount
- docker
- too much

## Operating System (with me)
### Linux:
- Fedora 35 Workstation

### Go: 
- go1.16.11 linux/amd64

### MySQL: 
- mysql  Ver 8.0.27 for Linux on x86_64 (Source distribution)


### Bahasa Indonesia
Der Schlaganfall 03.10.2018-heute. Dirilis 7 Januari 2020. Coding ini sedikit lupa. Pun, ini Bahasa Inggris lupa lagi. Perlahan-lahan dari stroke. Aamiin.

### English (translate[.]google[.]co[.]id)
Stroke: 03 10 2018-today. Released January 7, 2020. This coding is a little forgotten. This is English forgot again. Little by little from stroke. Aamiin.

---

Copyright © 2020 by Ocki Bagus Pratama
