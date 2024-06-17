package midlewares

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gofiber/fiber/v2"
	"todolistserver.com/test/authenticator"
	"todolistserver.com/test/types"
)

type CustomClaims struct {
	Scope string `json:"scope"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func RouteMilewareAuth(auth *authenticator.Authenticator) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authorization := c.Get("Authorization")

		if authorization == "" {
			log.Println("Error Auth: Bad Request")

			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"err_type": types.ERR_TYPE_MESSAGE,
				"msg":      "You do not have the appropriate permissions to access.",
			})
		}

		if ok := strings.Contains(authorization, "Bearer"); !ok {
			log.Println("Error Auth: Bad Request")

			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"err_type": types.ERR_TYPE_MESSAGE,
				"msg":      "The provided authorization token is invalid",
			})
		}

		authstrs := strings.Split(authorization, " ")

		if len(authstrs) != 2 {
			log.Println("Error Auth: Bad Request")

			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"err_type": types.ERR_TYPE_MESSAGE,
				"msg":      "The provided authorization token is invalid",
			})
		}

		token := authstrs[1]

		issuerURL, issuerURLErr := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")

		if issuerURLErr != nil {
			log.Println("Error Auth:", issuerURLErr.Error())

			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"err_type": types.ERR_TYPE_MESSAGE,
				"msg":      "An unexpected error has occurred on the server",
			})
		}

		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

		jwtValidator, jwtValidatorErr := validator.New(
			provider.KeyFunc,
			validator.RS256,
			issuerURL.String(),
			[]string{os.Getenv("AUTH0_AUDIENCE")},
			validator.WithCustomClaims(
				func() validator.CustomClaims {
					return &CustomClaims{}
				},
			),
			validator.WithAllowedClockSkew(time.Minute),
		)

		if jwtValidatorErr != nil {
			log.Println("Error Auth:", jwtValidatorErr.Error())

			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"err_type": types.ERR_TYPE_MESSAGE,
				"msg":      "An unexpected error has occurred on the server",
			})
		}

		claims, claimsErr := jwtValidator.ValidateToken(c.UserContext(), token)

		if claimsErr != nil {
			log.Println("Error Auth:", claimsErr.Error())

			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"err_type": types.ERR_TYPE_MESSAGE,
				"msg":      "Invalid token",
			})
		}

		dic := claims.(*validator.ValidatedClaims)

		c.Locals("user", dic.RegisteredClaims.Subject)

		return c.Next()
	}
}
