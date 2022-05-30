package controller

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"samplegoapp.com/db"
	"samplegoapp.com/models"
)

//TODO: add authorization handler, JWT secret
const JWT_SECRET = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	// check if email is already in use
	var existingUser models.User
	db.DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID != 0 {
		c.Status(fiber.StatusConflict)
		return c.JSON(fiber.Map{"message": "Email already in use"})
	} else {
		db.DB.Create(&user)
		return c.JSON(user)
	}
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	var user models.User

	db.DB.Where("email = ?", data["email"]).First(&user)

	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"message": "User not found"})
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		return c.JSON(fiber.Map{"message": "Invalid credentials"})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // one day 24hrs
	})

	token, err := claims.SignedString([]byte(JWT_SECRET))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)

		return c.JSON(fiber.Map{"message": "Error signing token"})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"message": "Login successful"})
}

func GetUsers(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"message": "Unauthorized"})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User
	db.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().AddDate(0, 0, -1),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"message": "Logout successful"})

}

func MakePost(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"message": "Unauthorized"})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User
	db.DB.Where("id = ?", claims.Issuer).First(&user)

	var post models.Post
	if err := c.BodyParser(&post); err != nil {
		return err
	}

	post.AuthorID = user.ID
	post.CreatedAt = time.Now().Unix()
	post.UpdatedAt = time.Now().Unix()

	db.DB.Create(&post)

	return c.JSON(post)
}

func GetPostsByUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"message": "Unauthorized"})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User
	db.DB.Where("id = ?", claims.Issuer).First(&user)

	// return posts in order of creation


	var posts []models.Post
	db.DB.Where("author_id = ?", user.ID).Find(&posts)

	return c.JSON(posts)
}

// get post by user and by param id
func GetPostByUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"message": "Unauthorized"})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User
	db.DB.Where("id = ?", claims.Issuer).First(&user)

	var post models.Post
	id, _ := strconv.Atoi(c.Params("id"))

	// try to find post by id if not found then return 404
	db.DB.Where("id = ? AND author_id = ?", id, user.ID).First(&post)

	if post.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"message": "Post not found"})
	}
	return c.JSON(post)
}

func UpdatePost(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"message": "Unauthorized"})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User
	db.DB.Where("id = ?", claims.Issuer).First(&user)

	var post models.Post
	id, _ := strconv.Atoi(c.Params("id"))

	// try to find post by id if not found then return 404
	db.DB.Where("id = ? AND author_id = ?", id, user.ID).First(&post)

	if post.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"message": "Post not found"})
	}

	if err := c.BodyParser(&post); err != nil {
		return err
	}

	post.UpdatedAt = time.Now().Unix()

	db.DB.Save(&post)

	return c.JSON(post)
}

func DeletePost(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"message": "Unauthorized"})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User
	db.DB.Where("id = ?", claims.Issuer).First(&user)

	var post models.Post
	id, _ := strconv.Atoi(c.Params("id"))

	// try to find post by id if not found then return 404
	db.DB.Where("id = ? AND author_id = ?", id, user.ID).First(&post)

	if post.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"message": "Post not found"})
	}

	db.DB.Delete(&post)

	return c.JSON(fiber.Map{"message": "Post deleted"})
}
