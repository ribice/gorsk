# GORSK - GO(lang) Restful Starter Kit

[![Build Status](https://travis-ci.org/ribice/gorsk.svg?branch=master)](https://travis-ci.org/ribice/gorsk)
[![codecov](https://codecov.io/gh/ribice/gorsk/branch/master/graph/badge.svg)](https://codecov.io/gh/ribice/gorsk)
[![Go Report Card](https://goreportcard.com/badge/github.com/ribice/gorsk)](https://goreportcard.com/report/github.com/ribice/gorsk)
[![Maintainability](https://api.codeclimate.com/v1/badges/c3cb09dbc0bc43186464/maintainability)](https://codeclimate.com/github/ribice/gorsk/maintainability)

**[Gorsk V2 is released - read more about it](https://www.ribice.ba/refactoring-gorsk/)**

Gorsk is a Golang starter kit for developing RESTful services. It is designed to help you kickstart your project, skipping the 'setting-up part' and jumping straight to writing business logic.

Previously Gorsk was built using [Gin](https://github.com/gin-gonic/gin). Gorsk using Gin is available [HERE](https://github.com/ribice/gorsk-gin).

Gorsk follows SOLID principles, with package design being inspired by several package designs, including Ben Johnson's [Standard Package Layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1), [Go Standard Package Layout](https://github.com/golang-standards/project-layout) with my own ideas applied to both. The idea for building this project and its readme structure was inspired by [this](https://github.com/qiangxue/golang-restful-starter-kit).

This starter kit currently provides:

* Fully featured RESTful endpoints for authentication, changing password and CRUD operations on the user entity
* JWT authentication and session
* Application configuration via config file (yaml)
* RBAC (role-based access control)
* Structured logging
* Great performance
* Request marshaling and data validation
* API Docs using SwaggerUI
* Mocking using stdlib
* Complete test coverage
* Containerized database query tests

The following dependencies are used in this project (generated using [Glice](https://github.com/ribice/glice)):

```bash
|-------------------------------------|--------------------------------------------|--------------|
|             DEPENDENCY              |                  REPOURL                   |   LICENSE    |
|-------------------------------------|--------------------------------------------|--------------|
| github.com/labstack/echo            | https://github.com/labstack/echo           | MIT          |
| github.com/go-pg/pg                 | https://github.com/go-pg/pg                | bsd-2-clause |
| github.com/dgrijalva/jwt-go         | https://github.com/dgrijalva/jwt-go        | MIT          |
| github.com/rs/zerolog               | https://github.com/rs/zerolog              | MIT          |
| golang.org/x/crypto/bcrypt          | https://github.com/golang/crypto           |              |
| gopkg.in/yaml.v2                    | https://github.com/go-yaml/yaml            |              |
| gopkg.in/go-playground/validator.v8 | https://github.com/go-playground/validator | MIT          |
| github.com/lib/pq                   | https://github.com/lib/pq                  | Other        |
| github.com/nbutton23/zxcvbn-go      | https://github.com/nbutton23/zxcvbn-go     | MIT          |
| github.com/fortytw2/dockertest      | https://github.com/fortytw2/dockertest     | MIT          |
| github.com/stretchr/testify         | https://github.com/stretchr/testify        | Other        |
|-------------------------------------|--------------------------------------------|--------------|
```

1. Echo - HTTP 'framework'.
2. Go-Pg - PostgreSQL ORM
3. JWT-Go - JWT Authentication
4. Zerolog - Structured logging
5. Bcrypt - Password hashing
6. Yaml - Unmarshalling YAML config file
7. Validator - Request validation.
8. lib/pq - PostgreSQL driver
9. zxcvbn-go - Password strength checker
10. DockerTest - Testing database queries
11. Testify/Assert - Asserting test results

Most of these can easily be replaced with your own choices since their usage is abstracted and localized.

## Getting started

Using Gorsk requires having Go 1.7 or above. Once you downloaded Gorsk (either using Git or go get) you need to configure the following:

1. To use Gorsk as a starting point of a real project whose package name is something like `github.com/author/project`, move the directory `$GOPATH/github.com/ribice/gorsk` to `$GOPATH/github.com/author/project`. We use [task](https://github.com/go-task/task) as a task runner. Run ``task relocate`` to rename the packages, e.g. ``task relocate TARGET=github.com/author/project``

2. Change the configuration file according to your needs, or create a new one.

3. Set the ("ENVIRONMENT_NAME") environment variable, either using terminal or os.Setenv("ENVIRONMENT_NAME","dev").

4. Set the JWT secret env var ("JWT_SECRET")

5. In cmd/migration/main.go set up psn variable and then run it (go run main.go). It will create all tables, and necessery data, with a new account username/password admin/admin.

6. Run the app using:

```bash
go run cmd/api/main.go
```

The application runs as an HTTP server at port 8080. It provides the following RESTful endpoints:

* `POST /login`: accepts username/passwords and returns jwt token and refresh token
* `GET /refresh/:token`: refreshes sessions and returns jwt token
* `GET /me`: returns info about currently logged in user
* `GET /swaggerui/` (with trailing slash): launches swaggerui in browser
* `GET /v1/users`: returns list of users
* `GET /v1/users/:id`: returns single user
* `POST /v1/users`: creates a new user
* `PATCH /v1/password/:id`: changes password for a user
* `DELETE /v1/users/:id`: deletes a user

You can log in as admin to the application by sending a post request to localhost:8080/login with username `admin` and password `admin` in JSON body.

### Implementing CRUD of another table

Let's say you have a table named 'cars' that handles employee's cars. To implement CRUD on this table you need:

1. Inside `pkg/utl/model` create a new file named `car.go`. Inside put your entity (struct), and methods on the struct if you need them.

2. Create a new `car` folder in the (micro)service where your service will be located, most probably inside `api`. Inside create a file/service named car.go and test file for it (`car/car.go` and `car/car_test.go`). You can test your code without writing a single query by mocking the database logic inside /mock/mockdb folder. If you have complex queries interfering with other entities, you can create in this folder other files such as car_users.go or car_templates.go for example.

3. Inside car folder, create folders named `platform`, `transport` and `logging`.

4. Code for interacting with a platform like database (postgresql) should be placed under `car/platform/pgsql`. (`pkg/api/car/platform/pgsql/car.go`)

5. In `pkg/api/car/transport` create a new file named `http.go`. This is where your handlers are located. Under the same location create http_test.go to test your API.

6. In logging directory create a file named `car.go` and copy the logic from another service. This serves as request/response logging.

6. In `pkg/api/api.go` wire up all the logic, by instantiating car service, passing it to the logging and transport service afterwards.

### Implementing other platforms

Similarly to implementing APIs relying only on a database, you can implement other platforms by:

1. In the service package, in car.go add interface that corresponds to the platform, for example, Indexer or Reporter.

2. Rest of the procedure is same, except that in `/platform` you would create a new folder for your platform, for example, `elastic`.

3. Once the new platform logic is implemented, create an instance of it in main.go (for example `elastic.Client`) and pass it as an argument to car service (`pkg/api/car/car.go`).

### Running database queries in transaction

To use a transaction, before interacting with db create a new transaction:

```go
err := s.db.RunInTransaction(func (tx *pg.Tx) error{
    // Application service here
})
````

Instead of passing database client as `s.db` , inside this function pass it as `tx`. Handle the error accordingly.

## Project Structure

1. Root directory contains things not related to code directly, e.g. docker-compose, CI/CD, readme, bash scripts etc. It should also contain vendor folder, Gopkg.toml and Gopkg.lock if dep is being used.

2. Cmd package contains code for starting applications (main packages). The directory name for each application should match the name of the executable you want to have. Gorsk is structured as a monolith application but can be easily restructured to contain multiple microservices. An application may produce multiple binaries, therefore Gorsk uses the Go convention of placing main package as a subdirectory of the cmd package. As an example, scheduler application's binary would be located under cmd/cron. It also loads the necessery configuration and passes it to the service initializers.

3. Rest of the code is located under /pkg. The pkg directory contains `utl` and 'microservice' directories.

4. Microservice directories, like api (naming corresponds to `cmd/` folder naming) contains multiple folders for each domain it interacts with, for example: user, car, appointment etc.

5. Domain directories, like user, contain all application/business logic and two additional directories: platform and transport.

6. Platform folder contains various packages that provide support for things like databases, authentication or even marshaling. Most of the packages located under platform are decoupled by using interfaces. Every platform has its own package, for example, postgres, elastic, redis, memcache etc.

7. Transport package contains HTTP handlers. The package receives the requests, marshals, validates then passes it to the corresponding service.

8. Utl directory contains helper packages and models. Packages such as mock, middleware, configuration, server are located here.

## License

gorsk is licensed under the MIT license. Check the [LICENSE](LICENSE) file for details.

## Author

[Emir Ribic](https://ribice.ba)
