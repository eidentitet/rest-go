package jwt

import (
	"context"
	"crypto/rsa"
	"errors"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	configuration "github.com/eidentitet/rest-go/config"
	"github.com/eidentitet/rest-go/http/middleware"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
)

type JwtMiddleware struct {
}

func init() {
	middleware.Register("auth", &JwtMiddleware{})
}

func (jwt JwtMiddleware) Next() func(handler http.Handler) http.Handler {
	config := configuration.GetAppConfig()
	return GetJwtValidator(config).CheckJWT
}

func GetJwtValidator(appConfig *configuration.AppConfiguration) *jwtmiddleware.JWTMiddleware {
	issuerURL, err := url.Parse(appConfig.OpenID.Issuer)
	if err != nil {
		log.Fatalln("failed to parse the issuer url: %v", err)
	}

	keyFunc := func(ctx context.Context) (interface{}, error) {
		verifyBytes, err := os.ReadFile(appConfig.OpenID.KeyPath)
		if err != nil {
			log.Errorln(err)
			return &rsa.PublicKey{}, err
		}
		publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return publicKey, nil
	}

	// Set up the validator.
	jwtValidator, err := validator.New(
		keyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{appConfig.OpenID.Audience},
	)
	if err != nil {
		log.Fatalln("failed to set up the validator: %v", err)
	}

	// Set up the middleware.
	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithCredentialsOptional(false),
		jwtmiddleware.WithErrorHandler(ErrorHandler),
	)
	return middleware
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, jwtmiddleware.ErrJWTMissing):
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error_code": 401, "error_message":"JWT is missing."}`))
	case errors.Is(err, jwtmiddleware.ErrJWTInvalid):
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error_code": 401, "error_message":"JWT is invalid."}`))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error_code": 401, "error_message":"Something went wrong while checking the JWT."}`))
	}
}
