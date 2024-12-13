package api

import (
	"auth/api/types"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// Low level API functionalities
type Router struct {
	logger   *zap.SugaredLogger
	gin      *gin.Engine
	database *mongo.Database
	logs     bool
	sessions map[string]*types.Session
}

func (router *Router) Log(message string) {
	if router.logs {
		router.logger.Info(message)
	}
}

func (router *Router) ErrorLog(message string) {
	if router.logs {
		router.logger.Error(message)
	}
}

func (router *Router) valudateAuth(c *gin.Context) {
	token := c.GetHeader("HARMONY_MICRO_SERVICES")

	// The microservice only validate the main microservices key, it will be located only in a local environment using the local subsystems
	// It is not valid for an external use so the applications should use an API router instead of using the AUTH API directly
	if token != os.Getenv("HARMONY_MICRO_SERVICES_KEY") {
		router.ErrorLog("Invalid micro service token key")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Next()
}

// SECTION - Sessions

func (api *Router) createSession(user types.User) (*string, error) {
	validated := api.validateSessionTime(user.Email)

	if validated {
		api.updateSession(user.Email)
		return &api.sessions[user.Email].SessionCode, nil
	}

	sessionCode := uuid.New().String()

	newSession := types.Session{
		User:        user,
		SessionCode: sessionCode,
		Date:        time.Now().UnixMilli(),
	}

	api.sessions[user.Email] = &newSession

	return &sessionCode, nil
}

func (api *Router) validateSessionTime(email string) bool {
	session, ok := api.sessions[email]

	if !ok {
		return false
	}

	now := time.Now().UnixMilli()

	active := now < session.Date+SESSION_TIME

	if !active {
		delete(api.sessions, email)
	}

	return active
}

func (api *Router) validateFullSession(email string, sessionCode string) bool {
	session, ok := api.sessions[email]

	if !ok || session.SessionCode != sessionCode {
		return false
	}

	validTime := api.validateSessionTime(email)

	return validTime
}

func (api *Router) updateSession(email string) bool {
	session, ok := api.sessions[email]

	if !ok {
		return false
	}

	session.Date += SESSION_TIME

	return true
}

func (api *Router) getActiveSessions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": false, "data": api.sessions})
}

func (api *Router) middleWareValidateSession(c *gin.Context) {
	var request types.Logged

	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success := api.validateFullSession(request.FromUser, request.SessionCode)

	if !success {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Next()
}

// Default methods
// We can use more methods but I will reduce it to GET and POST, I know that I can use PUT, UPDATE...

func (router *Router) ValidatedGet(url string, callback func(c *gin.Context)) {
	router.gin.GET(url, router.valudateAuth, callback)
}

func (router *Router) ValidatedPost(url string, callback func(c *gin.Context)) {
	router.gin.POST(url, router.valudateAuth, callback)
}

func (router *Router) ValidatedWithSessionPost(url string, callback func(c *gin.Context)) {
	router.gin.POST(url, router.valudateAuth, router.middleWareValidateSession, callback)
}
