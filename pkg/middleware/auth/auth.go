package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/unification-com/unode-onboard-api/pkg/models"
)

func Middleware(dbClient *models.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			// logrus.Errorln("No token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized. No token"})
		}

		tokenString = tokenString[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil {
			logrus.Errorf("Unable to parse token: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized: Invalid token"})
		}

		if !token.Valid {
			logrus.Errorln("Invalid token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		walletAddressClaim := claims["wallet"]
		walletAddress := convertInterfaceToString(walletAddressClaim)

		// checking if walletAddress still exists in DB
		var account models.Account
		if err := dbClient.Model(&account).Where("payment_wallet = ?", walletAddress).Select(); err != nil {
			logrus.Errorf("Wallet address %s does not exist in DB. Err: %v", walletAddress, err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
		}

		if account.AccountID == 0 || account.PaymentWallet == "" {
			logrus.Errorln("Wallet address or account ID does not exist")
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token compromised"})
		}

		c.Locals("wallet", account.PaymentWallet)
		c.Locals("accountid", account.AccountID)

		return c.Next()
	}
}
