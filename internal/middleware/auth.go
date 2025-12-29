package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/k/iRegistro/internal/application/auth"
	"github.com/k/iRegistro/internal/domain"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("[DEBUG] AuthMiddleware Header: %q\n", authHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			fmt.Println("[DEBUG] Bearer prefix missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &auth.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			fmt.Printf("AuthMiddleware: Invalid token: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": fmt.Sprintf("%v", err)})
			return
		}

		claims, ok := token.Claims.(*auth.CustomClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Set context variables
		c.Set("userID", claims.UserID)
		c.Set("schoolID", claims.SchoolID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RBACMiddleware(requiredRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Role not found in context"})
			return
		}

		userRole := roleVal.(domain.Role)

		// SuperAdmin bypass
		if userRole == domain.RoleSuperAdmin {
			c.Next()
			return
		}

		for _, role := range requiredRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		fmt.Printf("RBACMiddleware: Insufficient permissions for role %v. Required one of: %v\n", userRole, requiredRoles)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions", "your_role": userRole, "required_roles": requiredRoles})
	}
}
