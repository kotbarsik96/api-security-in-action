package main

import (
	"api-security-in-action/src/api/ctrlauth"
	"api-security-in-action/src/api/ctrlmessage"
	"api-security-in-action/src/api/ctrlspace"
	"api-security-in-action/src/db"
	"api-security-in-action/src/domain/audit"
	"api-security-in-action/src/domain/auth"
	"api-security-in-action/src/domain/message"
	"api-security-in-action/src/domain/space"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	gormDb, err := db.NewGormDB()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	RegisterControllers(router, gormDb)

	router.Run()
}

func RegisterControllers(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")

	authMdw := auth.MiddlewareAuthentication(db)
	auditMdw := audit.AuditMiddleware(db)

	// space
	spaceCtrl := ctrlspace.NewSpaceController(
		space.NewSpaceCreateService(db))

	spaceCtrl.RegisterRoutes(api, authMdw, auditMdw)

	// message
	messageCtrl := ctrlmessage.NewMessageController(
		message.NewMessageCreateService(db),
		message.NewMessagesRepository(db))

	messageCtrl.RegisterRoutes(api, authMdw, auditMdw)

	// auth
	authCtrl := ctrlauth.NewAuthController(
		auth.NewAuthService(db))

	authCtrl.RegisterRoutes(api)
}
