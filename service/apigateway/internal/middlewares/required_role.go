package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RequireRoles is an Echo middleware that checks if the value of the key "role_names"
// (set by the RoleValidator middleware) contains any of the given allowedRoles. If
// it does, the request is allowed to proceed. Otherwise, it returns a 403 status code
// with an error message indicating that the role is not permitted.
//
// Example:
//
//	e.GET("/admin", RequireRoles("admin", "superadmin"), func(c echo.Context) error {
//		// only users with role "admin" or "superadmin" can access this route
//		return c.String(http.StatusOK, "Hello, "+c.Get("user_id").(string))
//	})
func RequireRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			raw, ok := c.Get("role_names").([]string)
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "Roles not found")
			}

			for _, r := range raw {
				for _, allowed := range allowedRoles {
					if r == allowed {
						return next(c)
					}
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Role not permitted")
		}
	}
}
