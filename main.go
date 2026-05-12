package main

import (
	"api-security-in-action/src/api/controllers"
	"api-security-in-action/src/db"
	"api-security-in-action/src/db/gorm_services/audit"
	"api-security-in-action/src/db/gorm_services/auth"
	"api-security-in-action/src/db/gorm_services/message"
	"api-security-in-action/src/db/gorm_services/space"

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
	spaceCtrl := controllers.NewSpaceController(
		space.NewSpaceCreateService(db))

	spaceCtrl.RegisterRoutes(api, authMdw, auditMdw)

	// message
	messageCtrl := controllers.NewMessageController(
		message.NewMessageCreateService(db),
		message.NewMessagesRepository(db))

	messageCtrl.RegisterRoutes(api, authMdw, auditMdw)

	// auth
	authCtrl := controllers.NewAuthController(
		auth.NewAuthService(db))

	authCtrl.RegisterRoutes(api)

	// audit
	auditCtrl := controllers.NewAuditController(
		audit.NewAuditRepository(db))

	auditCtrl.RegisterRoutes(api, authMdw, auditMdw)
}
