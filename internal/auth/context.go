package auth

import "github.com/gin-gonic/gin"

const ctxKey = "auth.claims"

type Claims struct {
	Token  string  `json:"-"`
	UserID *string `json:"user_id,omitempty"`
	Role   string  `json:"role"` // админ или юзер
}

func setClaims(c *gin.Context, cl Claims) { c.Set(ctxKey, cl) }

func GetClaims(c *gin.Context) (Claims, bool) {
	v, ok := c.Get(ctxKey)
	if !ok {
		return Claims{}, false
	}
	cl, _ := v.(Claims)
	return cl, true
}
