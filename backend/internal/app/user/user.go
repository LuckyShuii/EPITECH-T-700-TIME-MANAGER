package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// For now all in one file, then start the structure in this folder once basic concepts are understood
// such as services, repositories, handlers, models, auth, etc.

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

type User struct {
	ID     	  string  `json:"id"`
	Username  string  `json:"username"`
	Email 	  string  `json:"email"`
}

var fakeUsers = []User{
	{ID: "1", Username: "username&1", Email: "username@1.com"},
	{ID: "2", Username: "username&2", Email: "username@2.com"},
	{ID: "3", Username: "username&3", Email: "username@3.com"},
}

func ReturnFakeUsers(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, fakeUsers)
}

func ReturnFakeUserById(c *gin.Context) {
	// Param from /:id in the route
	id := c.Param("id")

    // Loop over the list of albums, looking for
    // a user whose ID value matches the parameter.
    for _, a := range fakeUsers {
        if a.ID == id {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }

    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func AddUser(c *gin.Context) {
    var newUser User

    // Call BindJSON to bind the received JSON to
    // newUser.
    if err := c.BindJSON(&newUser); err != nil {
        return
    }

    // Add the new album to the slice.
    fakeUsers = append(fakeUsers, newUser)
    c.IndentedJSON(http.StatusCreated, newUser)
}