// Copyright 2017 Emir Ribic. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// GORSK - Go(lang) restful starter kit
//
// API Docs for GORSK v1
//
// 	 Terms Of Service:  N/A
//     Schemes: http
//     Version: 1.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Emir Ribic <ribice@gmail.com> https://ribice.ba
//     Host: localhost:8080
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer: []
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package main

import (
	"io/ioutil"
	"log"

	"github.com/gin-contrib/cors"

	"github.com/ribice/gorsk/internal/platform/postgres"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/ribice/gorsk/cmd/api/config"
	"github.com/ribice/gorsk/cmd/api/mw"
	"github.com/ribice/gorsk/cmd/api/service"
	_ "github.com/ribice/gorsk/cmd/api/swagger"
	"github.com/ribice/gorsk/internal/account"
	"github.com/ribice/gorsk/internal/auth"
	"github.com/ribice/gorsk/internal/rbac"
	"github.com/ribice/gorsk/internal/user"
	"go.uber.org/zap"
)

func main() {

	r := gin.Default()
	mw.Add(r, cors.Default(), mw.SecureHeaders())

	cfg, err := config.Load("dev")
	checkErr(err)

	db, err := pgsql.New(cfg.DB)
	checkErr(err)

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	addV1Services(cfg, r, db, logger)
	r.Run()
}

func addV1Services(cfg *config.Configuration, r *gin.Engine, db *pg.DB, log *zap.Logger) {
	userDB := pgsql.NewUserDB(db, log)
	jwt := mw.NewJWT(cfg.JWT)
	authSvc := auth.New(userDB, jwt)
	service.NewAuth(authSvc, r)

	rbacSvc := rbac.New(userDB)

	v1Router := r.Group("/v1")
	v1Router.GET("/swagger", docHandler)
	v1Router.Use(jwt.MWFunc())

	accDB := pgsql.NewAccountDB(db, log)
	service.NewAccount(account.New(accDB, userDB, rbacSvc), v1Router)

	service.NewUser(user.New(userDB, rbacSvc, authSvc), v1Router)
}

func docHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	data, _ := ioutil.ReadFile("./cmd/api/swagger.json")
	c.Writer.Write(data)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
