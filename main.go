package main

import (
	"api-security-in-action/src/api/controllers"
	"api-security-in-action/src/db"
	"api-security-in-action/src/db/gorm_services/audit"
	"api-security-in-action/src/db/gorm_services/auth"
	"api-security-in-action/src/db/gorm_services/message"
	"api-security-in-action/src/db/gorm_services/permission"
	"api-security-in-action/src/db/gorm_services/space"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}

	gormDb, err := db.NewGormDB()
	if err != nil {
		log.Fatalf("Error loading database: %v", err)
	}

	router := CreateRouter()

	RegisterControllers(router, gormDb)

	router.Run()
}

func CreateRouter() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	r.Use(sessions.Sessions(os.Getenv("APP_NAME"), store))

	return r
}

func RegisterControllers(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")

	authMdw := auth.MiddlewareAuthentication(db)
	auditMdw := audit.AuditMiddleware(db)

	csrfService := auth.NewCsrfService([]byte(os.Getenv("HMAC_SECRET")))
	csrfMdw := auth.MiddlewareCSRF(csrfService)

	permissionGuard := permission.NewPermissionGuard(db)

	// space
	spaceCtrl := controllers.NewSpaceController(
		space.NewSpaceCreateService(db),
		permission.NewPermissionService(db))

	spaceCtrl.RegisterRoutes(api, authMdw, auditMdw, csrfMdw)

	// message
	messageCtrl := controllers.NewMessageController(
		message.NewMessageCreateService(db),
		message.NewMessagesRepository(db),
		permissionGuard)

	messageCtrl.RegisterRoutes(api, authMdw, auditMdw, csrfMdw)

	// auth
	authCtrl := controllers.NewAuthController(
		auth.NewAuthService(db, auth.NewBcryptPasswordService()),
		csrfService)

	authCtrl.RegisterRoutes(api, csrfMdw)

	// audit
	auditCtrl := controllers.NewAuditController(
		audit.NewAuditRepository(db))

	auditCtrl.RegisterRoutes(api, authMdw, auditMdw)
}
