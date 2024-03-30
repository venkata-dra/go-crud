package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/venkata-dra/crud-go/initializers"
	"github.com/venkata-dra/crud-go/models"
	"golang.org/x/crypto/bcrypt"
)

func PostsCreate(c *gin.Context) {

	var body struct {
		Body  string `json:"body" binding:"required"`
		Title string `json:"title" binding:"required"`
	}
	// log.Fatal(body.Body)
	c.Bind(&body)
	// if err := c.Bind(&body); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	// 	return
	// }

	post := models.Post{Title: body.Title, Body: body.Body}
	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.Status(http.StatusInternalServerError)
		log.Fatal("Error occured while sending post request", result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func PostIndex(c *gin.Context) {
	var posts []models.Post
	initializers.DB.Find(&posts)
	c.JSON(http.StatusOK, gin.H{
		"post": posts,
	})
}

func PostsShow(c *gin.Context) {

	var posts models.Post

	id := c.Param("id")
	initializers.DB.First(&posts, id)

	c.JSON(http.StatusOK, gin.H{
		"post": posts,
	})
}

func PostsUpdate(c *gin.Context) {
	// get the record using id
	id := c.Param("id")

	// read the data of the request body
	var body struct {
		Body  string `json:"body" binding:"required"`
		Title string `json:"title" binding:"required"`
	}
	//find the post we are updating
	c.Bind(&body)
	var post models.Post
	initializers.DB.First(&post, id)
	//update it

	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Title,
		Body:  body.Body,
	})
	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})

}

func PostDelete(c *gin.Context) {
	id := c.Param("id")

	initializers.DB.Delete(&models.Post{}, id)

	c.Status(200)
}

func SignUp(c *gin.Context) {
	//get the email/pass off the request body

	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{"error": "failed to read body"})
		return
	}
	//hash the password
	hash, error := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if error != nil {
		c.JSON(400, gin.H{"error": "failed to hash password"})
		return
	}
	//create the user
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "failed to create user"})
		return
	}

	//respond
	c.JSON(200, gin.H{"error": "User Created"})

}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{"error": "failed to read body"})
		return
	}
	//get the email and pass off the req body
	var user models.User
	initializers.DB.First(&user, "email =? ", body.Email)

	if user.ID == 0 {
		c.JSON(400, gin.H{"error": "failed to get user"})
		return
	}
	// lookup the request user

	// compare sent in password with saved user pass hash

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(400, gin.H{"error": "invalid email or password"})
		return
	}

	// generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"time": time.Now().Add(time.Hour).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	//send it back
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to create token"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", true, true)
	c.JSON(200, gin.H{
		"msg": "login successful",
	})
}

func Validate(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "login successful",
	})

}
