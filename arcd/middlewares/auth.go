/*
 * Arc - Copyleft of Simone 'evilsocket' Margaritelli.
 * evilsocket at protonmail dot com
 * https://www.evilsocket.net/
 *
 * See LICENSE.
 */
package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/evilsocket/arc/arcd/config"
	"github.com/evilsocket/arc/arcd/log"
	"github.com/evilsocket/arc/arcd/utils"
	"github.com/gin-gonic/gin"
	"regexp"
	"time"
)

var authTokenParser = regexp.MustCompile("^(?i)Bearer:\\s*(.+)$")

func GenerateToken(k []byte, userId string) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	claims := make(jwt.MapClaims)
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(config.Conf.TokenDuration)).Unix()
	token.Claims = claims
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(k)
	return tokenString, err
}

func ValidateToken(t string, k string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(k), nil
	})

	return token, err
}

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.GetString("user_id")
		// do we need to refresh session data?
		if user_id != config.Conf.Username {
			// Parse bearer token from Authorization header.
			authorization := c.Request.Header.Get("Authorization")
			m := authTokenParser.FindStringSubmatch(authorization)
			if len(m) != 2 {
				utils.Forbidden(c)
				return
			}
			// Validate token
			token := m[1]
			valid, err := ValidateToken(token, config.Conf.Secret)
			if err != nil {
				log.Api(log.WARNING, c, "Error while validating bearer token: %s", err)
				utils.Forbidden(c)
				return
			}

			// set session data
			c.Set("user_id", valid.Claims.(jwt.MapClaims)["user_id"])
		}

		c.Next()
	}
}