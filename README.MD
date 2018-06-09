# GORSK - GO(lang) Restful Starter Kit

[![Build Status](https://travis-ci.org/ribice/gorsk.svg?branch=master)](https://travis-ci.org/ribice/gorsk)
[![codecov](https://codecov.io/gh/ribice/gorsk/branch/master/graph/badge.svg)](https://codecov.io/gh/ribice/gorsk)
[![Go Report Card](https://goreportcard.com/badge/github.com/ribice/gorsk)](https://goreportcard.com/report/github.com/ribice/gorsk)
[![Maintainability](https://api.codeclimate.com/v1/badges/c3cb09dbc0bc43186464/maintainability)](https://codeclimate.com/github/ribice/gorsk/maintainability)

Gorsk is a Golang starter kit for developing RESTful services. It is designed to help you kickstart your project, skipping the 'setting-up part' and jumping straight to writing business logic.

Previously Gorsk was built using [Gin](https://github.com/gin-gonic/gin). Gorsk using Gin is available [HERE](https://github.com/ribice/gorsk-gin).

Gorsk follows SOLID principles, with package design being inspired by Ben Johnson's [Standard Package Layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1). The idea for building this project and its readme structure was inspired by [this](https://github.com/qiangxue/golang-restful-starter-kit).

This starter kit currently provides:

* Fully featured RESTful endpoints for authentication, password reset and CRUD operations on the user entity
* Session handling
* JWT Based authentication
* Application configuration via config file (yaml)
* RBAC (role-based access control)
* Structured (error) logging
* Great performance
* Partial update (PATCH) handling using reflection
* Request marshaling and data validation
* API Docs using SwaggerUI
* Mocking using stdlib
* Full test coverage
* Containerized DB query tests

The following dependencies are used in this project (generated using [Glice](https://github.com/ribice/glice)):

```bash
|-------------------------------------|--------------------------------------------|--------------|
|             DEPENDENCY              |                  REPOURL                   |   LICENSE    |
|-------------------------------------|--------------------------------------------|--------------|
| github.com/labstack/echo            | https://github.com/labstack/echo           | MIT          |
| github.com/go-pg/pg                 | https://github.com/go-pg/pg                | bsd-2-clause |
| github.com/dgrijalva/jwt-go         | https://github.com/dgrijalva/jwt-go        | MIT          |
| github.com/rs/xid                   | https://github.com/rs/xid                  | MIT          |
| golang.org/x/crypto/bcrypt          | https://github.com/golang/crypto           |              |
| github.com/facebookgo/grace         | https://github.com/facebookgo/grace        | Other        |
| gopkg.in/yaml.v2                    | https://github.com/go-yaml/yaml            |              |
| gopkg.in/go-playground/validator.v8 | https://github.com/go-playground/validator | MIT          |
| github.com/lib/pq                   | https://github.com/lib/pq                  | Other        |
| github.com/fortytw2/dockertest      | https://github.com/fortytw2/dockertest     | MIT          |
| github.com/stretchr/testify         | https://github.com/stretchr/testify        | Other        |
|-------------------------------------|--------------------------------------------|--------------|
```

1. Echo - HTTP 'framework'.
2. Go-Pg - PostgreSQL ORM
3. JWT-Go - JWT Authentication
4. XID - Generating refresh tokens
5. Bcrypt - Password hashing
6. Yaml - Unmarshalling YAML config file
7. Validator - Request validation.
8. lib/pq - Postgres driver
9. DockerTest - Testing database queries
10. Testify/Assert - Asserting test results

Most of these can easily be replaced with your own choices since their usage is abstracted and localized.

## Getting started

Using Gorsk requires having Go 1.7 or above. Once you downloaded Gorsk (either using Git or go get) you need to configure the following:

1. To use Gorsk as a starting point of a real project whose package name is something like `github.com/author/project`, move the directory `$GOPATH/github.com/ribice/gorsk` to `$GOPATH/github.com/author/project` and do a global replacement of the string `github.com/ribice/gorsk` with `github.com/author/project`.

2. Change the configuration file according to your needs, or create a new one.

3. Set the ("ENVIRONMENT_NAME") environment variable, either using terminal or os.Setenv("ENVIRONMENT_NAME","dev").

4. In cmd/migration/main.go set up psn variable and then run it (go run main.go). It will create all tables, and necessery data, with a new account username/password admin/admin.

5. Run the app using:

```bash
go run cmd/api/main.go
```

The application runs as an HTTP server at port 8080. It provides the following RESTful endpoints:

* `POST /login`: accepts username/passwords and returns jwt token and refresh token
* `GET /refresh/:token`: refreshes sessions and returns jwt token
* `GET /me`: returns info about currently logged in user
* `GET /swaggerui/`: launches swaggerui in browser
* `GET /v1/users`: returns list of users
* `GET /v1/users/:id`: returns single user
* `POST /v1/users`: creates a new user
* `PATCH /v1/users/:id/password`: changes password for a user
* `DELETE /v1/users/:id`: deletes a user

You can log in as admin to the application by sending a post request to localhost:8080/login with username `admin` and password `admin` in JSON body.

### Implementing CRUD of another table

Let's say you have a table named 'cars' that handles employee's cars. To implement CRUD on this table you need:

1. In project's root location create a new file named `car.go`. Inside put your entity (struct), and methods on the struct if you need them. Create interfaces for all platforms you have, for example, database, indexing, reporting etc. An example for this is UserDB interface in user.go

2. Create a new folder in root named car, and inside create a file/service named car.go and tests for it car_test.go (`car/car.go` and `car/car_test.go`). You can test your code without writing a single query by mocking the database logic inside /mock/mockdb folder. If you have complex queries interfering with other entities, you can create in this folder other files such as car_users.go or car_templates.go for example.

3. Database access code can be found under platform/postgres folder. (`platform/postgres/car.go` and `platform/postgres/car.go`)

4. In `cmd/api/service` create a new file named `car.go`. This is where your handlers are located. To handle the request data, create a new file inside `cmd/api/request` that will handle validation and request marshaling. Under the same location create car_test.go to test your API.

5. In `cmd/api/swagger` create car.go containing small structs to generate swagger API docs.

6. In `cmd/api/main.go` wire up all the logic, by creating a CarDB instance, passing it to car service and then the service to the handler.

### Implementing other platforms

Similarly to implementing APIs relying only on a database, you can implement other platforms by:

1. In the model package, in car.go add struct that corresponds to the platform, for example, CarElastic, or CarReporting. Create an interface that will handle communication to the platform, for example, CarIndexer interface.

2. Rest of the procedure is same, except that in `/platform` you would create a new folder for your platform, for example, `elastic`.

3. Once the new platform logic is implemented, create an instance of it in main.go (for example `elastic.client`) and pass it as an argument to car service (`car/car.go`).

## Project Structure

The project structure is explained in details in [THIS](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1) blog post, with some small changes taken from [Package Oriented Design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html). and GoDD. Primarily the package design is architectured by [tonto](https://github.com/tonto) with some small changes done by me.

1. Root package contains things not related to code directly, e.g. docker-compose, CI/CD, readme, bash scripts etc. It should also contain vendor folder, Gopkg.toml and Gopkg.lock if dep is being used.

2. Main package ties together dependencies. Located in cmd/api/main.go, the main package instantiates routing, all configurations, connections etc and injects them as dependencies to services. Gorsk is structured as a monolith application but can be easily restructured to contain multiple microservices. An application may produce multiple binaries, therefore Gorsk uses the Go convention of placing main package as a subdirectory of the cmd package. As an example, scheduler application's binary would be located under cmd/cron.

3. Rest of the code is located under /internal. Internal root contains domain types, e.g. User, Car, Company. This package only contains simple data types like User struct for holding user data or a UserService interface for fetching or saving user data.

4. Domain packages are located under the internal directory, in folders named after the domain. For example, car package is located in internal/car. All application/business logic is found here, and it connects to other 'platforms' such as a database, reporting, indexing, etc.

5. Platform folder contains various packages that provide support for things like databases, authentication or even marshaling. Most of the packages located under platform are decoupled by using interfaces. Every platform has its own package, for example, postgres, elastic, redis, memcache etc.

6. Mock package is shared across the entire application. Specific platforms have subpackages such as mockdb, mockindex etc. Since most of the dependencies are injected, it is trivial to mock them and pass the mock service as an argument.

7. Service package contains HTTP handlers. They only receive the requests, marshal and validate them using a helper package (request) and pass it to the corresponding services.

8. MW Package contains middleware related implementation, such as JWT authentication, CORS, perhaps request/error logging, security etc.

9. Config package contains application configurations, as well as the implementation to read the config files.

10. Swagger is a helper package containing go-swagger annotations for SwaggerUI generation.

## License

gorsk is licensed under the MIT license. Check the [LICENSE](LICENSE.md) file for details.

## Author

[Emir Ribic](https://ribice.ba)
