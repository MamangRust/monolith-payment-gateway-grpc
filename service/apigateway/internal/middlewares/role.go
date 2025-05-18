package middlewares

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"github.com/labstack/echo/v4"
)

func RoleMiddleware(allowedRoles []string, roleService pb.RoleServiceClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := c.Get("userID").(string)
			if !ok || userID == "" {
				return echo.NewHTTPError(http.StatusForbidden, "User ID not found in context")
			}

			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
			}

			roleResponse, err := roleService.FindByUserId(c.Request().Context(), &pb.FindByIdUserRoleRequest{
				UserId: int32(userIDInt),
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user information")
			}

			for _, userRole := range roleResponse.Data {
				for _, allowedRole := range allowedRoles {
					if userRole.Name == allowedRole {

						return next(c)
					}
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to access this resource")
		}
	}
}
