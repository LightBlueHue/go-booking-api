# Welcome to Revel
[![Build](https://github.com/LightBlueHue/go-booking-api/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/LightBlueHue/go-booking-api/actions/workflows/go.yml)

A high-productivity web framework for the [Go language](http://www.golang.org/).

### Create new app:
revel new -a go-booking-api -r

### Start the web server:

   revel run go-booking-api

### Go to http://localhost:9000/ and you'll see:

    "It works"

## Code Layout

The directory structure of a generated Revel application:

    conf/             Configuration directory
        app.conf      Main app configuration file
        routes        Routes definition file

    app/              App sources
        init.go       Interceptor registration
        controllers/  App controllers go here
        views/        Templates directory

    messages/         Message files

    public/           Public static assets
        css/          CSS files
        js/           Javascript files
        images/       Image files

    tests/            Test suites


## Help

* The [Getting Started with Revel](http://revel.github.io/tutorial/gettingstarted.html).
* The [Revel guides](http://revel.github.io/manual/index.html).
* The [Revel sample apps](http://revel.github.io/examples/index.html).
* The [API documentation](https://godoc.org/github.com/revel/revel).

* [Json tags](https://drstearns.github.io/tutorials/gojson/)
* [Odata response](https://docs.oasis-open.org/odata/odata-json-format/v4.0/errata02/os/odata-json-format-v4.0-errata02-os-complete.html#_Toc403940655)
* [jwt](https://medium.com/wesionary-team/jwt-authentication-in-golang-with-gin-63dbc0816d55)
* https://stackoverflow.com/questions/44589854/is-it-possible-to-debug-go-revel-framework-from-visual-studio-code
* https://medium.com/learn-go/go-path-explained-cab31a0d90b9
* https://go.dev/ref/spec#Passing_arguments_to_..._parameters
* https://pkg.go.dev/syreclabs.com/go/faker#readme-name


## Docker
`docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres`

`docker run --name pgAdmin -p 4000:80 -e PGADMIN_DEFAULT_EMAIL="postgres@postgres.com" -e PGADMIN_DEFAULT_PASSWORD="postgres" dpage/pgadmin4`


PgAdmin should be available on http://localhost:4000/

To connect pgAdmin to your postgress server, run the following in cmd to find your container ip: docker network inspect bridge or type host.docker.internal as host
