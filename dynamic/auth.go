package dynamic

// import (
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/golang-jwt/jwt/v4"
// )

// func InitAuthRoutes(router fiber.Router) {
// 	signin := router.Group("/sign-in")
// 	signin.Get("", signinGet)
// 	signin.Post("", signinPost)
// }

// func signinPost(c *fiber.Ctx) error {
// 	var input Administrator

// 	if err := c.BodyParser(&input); err != nil {
// 		return err
// 	}

// 	if input.Secret != "Fuck_me_for_the_win007" {
// 		return fiber.NewError(500, "Bad password")
// 	}

// jwtconf := JWT()

// claims := jwt.MapClaims{
// 	"iss":   input.Secret,
// 	"exp":   jwt.NewNumericDate(time.Now().Add(jwtconf.ExpireTime)),
// 	"iat":   jwt.NewNumericDate(time.Now()),
// 	"aud":   "arctic",
// 	"scope": "wewe",
// }

// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// mytoken, err := token.SignedString([]byte(jwtconf.SigningKey))
// if err != nil {
// 	return err
// }

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": mytoken})
// }

// func signinGet(c *fiber.Ctx) error {
// 	response := "Implement me"

// 	c.Send([]byte(response))

// 	return nil
// }

// type JWTConfig struct {
// 	SigningKey string
// 	ExpireTime time.Duration
// }

// func JWT() JWTConfig {
// 	expireHours := 24

// 	return JWTConfig{
// 		SigningKey: "private_key_x69s",
// 		ExpireTime: time.Duration(expireHours) * time.Hour,
// 	}
// }

// func AuthRequired(c *fiber.Ctx) error {
// 	token := c.Locals("token").(*jwt.Token)
// 	claims := token.Claims.(jwt.MapClaims)

// 	secret := claims["iss"]

// 	c.Locals("username", secret)

// 	if secret != "Fuck_me_for_the_win007" {
// 		return fiber.NewError(1, "bad secret")
// 	}

// 	return c.Next()
// }
