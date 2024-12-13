package api

import (
	"auth/api/database"
	"auth/api/types"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type API struct {
	Router
	state bool
}

// SECTION - Restrictions

func (api *API) validateUserRole(user string, needRole string) (bool, error) {
	role, err := database.GetUserRoleByEmail(api.database, user)

	if err != nil {
		return false, err
	}

	return role == needRole, nil
}

// SECTION - Roles

func (api *API) getRoles(c *gin.Context) {
	roles, err := database.GetRoles(api.database)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "data": roles})
}

func (api *API) addRole(c *gin.Context) {
	var request types.APIAddRoleRequest

	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := api.validateUserRole(request.FromUser, "ADMIN")

	if err != nil || !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	success, err = database.AddRole(api.database, request.Name)

	if err != nil || !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (api *API) configRole(c *gin.Context) {
	var request types.APIConfigRoleRequest

	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := api.validateUserRole(request.FromUser, "ADMIN")

	if err != nil || !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	success, err = database.UpdateRoleConfig(api.database, request.Id, request.Config)

	if err != nil || !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

// SECTION - Users

func (api *API) createUser(c *gin.Context) {
	var request types.APIAddUser

	err := c.BindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pass, err := HashPassword(request.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := database.AddUser(api.database, types.User{
		Email:    request.Email,
		Name:     request.Name,
		Password: pass,
		Role:     request.Role,
	})

	if err != nil || !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (api *API) hasRole(c *gin.Context) {
	var email = c.Param("email")
	var role = c.Param("role")

	user, err := database.GetUserByEmail(api.database, email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "data": user.Role == role})
}

func (api *API) updateUserRole(c *gin.Context) {
	var request types.APIUpdateUserRole

	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := api.validateUserRole(request.Email, "ADMIN")

	if err != nil || !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	success, err = database.UpdateUserRoleByEmail(api.database, request.Email, request.Role)

	if err != nil || !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (api *API) getUserRole(c *gin.Context) {
	var email = c.Param("email")

	user, err := database.GetUserByEmail(api.database, email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "data": user.Role})
}

func (api *API) getUsers(c *gin.Context) {
	users, err := database.GetUsers(api.database)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "data": users})
}

func (api *API) login(c *gin.Context) {
	var request types.APILogin

	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := database.GetUserByEmail(api.database, request.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success := VerifyPassword(request.Password, user.Password)

	if !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	user.Password = ""
	session, err := api.createSession(*user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "data": gin.H{"session": session, "user": user}})
}

func (api *API) validateSession(c *gin.Context) {
	var request types.APIValidateSession

	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isValidSession := api.validateFullSession(request.FromUser, request.SessionCode)

	if !isValidSession {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (api *API) createRoutes() {
	// API - GET - /status
	// return the state of the API, if it fails internally must returns false
	api.ValidatedGet("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"state": api.state})
	})

	// API - GET - /sessions
	// get active session on the system
	api.ValidatedGet("/sessions", api.getActiveSessions)

	// API - GET - /roles
	// return all possible roles on the ecosystem
	api.ValidatedGet("/roles", api.getRoles)

	// API - GET - /has-role/:role/:userId
	// return if the user has a specific role, use it to validate the role
	api.ValidatedGet("/has-role/:role/:email", api.hasRole)

	// API - POST - /add-role
	// add new role to the ecosystem
	api.ValidatedWithSessionPost("/add-role", api.addRole)

	// API - POST - /config-role
	// update role configuration
	api.ValidatedWithSessionPost("/config-role", api.configRole)

	// API - POST - /users
	// get all users from database
	api.ValidatedGet("/users", api.getUsers)

	// API - POST - /set-role
	// add role to one user
	api.ValidatedWithSessionPost("/set-role", api.updateUserRole)

	// API - GET - /roles/:userId
	// return the roles of one user by id
	api.ValidatedGet("/roles/:email", api.getUserRole)

	// API - POST - /login
	// login with username and ID
	api.ValidatedPost("/login", api.login)

	// API - POST - /create-user
	// add new user to the database
	api.ValidatedPost("/create-user", api.createUser)

	// API - POST - /validate-session
	// validate if the session is active
	api.ValidatedPost("/validate-session", api.validateSession)

	api.logger.Info("Initilizing routes")
}

func (api *API) Initialize(logger *zap.SugaredLogger, database *mongo.Database) {
	api.logger = logger
	api.database = database

	api.state = true
	api.logs = true

	api.logger.Info("Created sessions map")
	api.sessions = make(map[string]*types.Session)

	api.gin = gin.Default()
	api.gin.Use(CorsMiddleware())

	api.createRoutes()

	api.gin.Run(":" + os.Getenv("PORT"))

}
