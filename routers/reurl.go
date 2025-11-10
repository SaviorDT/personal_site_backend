package routers

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    reurlController "personal_site/controllers/reurl"
    "personal_site/middlewares"
)

// reurlRouter registers RESTful endpoints to manage redirect mappings.
// Routes are mounted under the API prefix + `/reurl`.
type reurlRouter struct{}

func (reurlRouter) RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
    // List all mappings
    r.GET("/", middlewares.AuthRequired(), func(c *gin.Context) {
        reurlController.ListReurls(c, db)
    })

    // Create a new mapping (protected)
    r.POST("/", middlewares.AuthRequired(), func(c *gin.Context) {
        reurlController.CreateReurl(c, db)
    })

    // Get a mapping by ID
    r.GET("/:id", middlewares.AuthRequired(), func(c *gin.Context) {
        reurlController.GetReurl(c, db)
    })

    // Patch a mapping by ID (protected)
    r.PATCH("/:id", middlewares.AuthRequired(), func(c *gin.Context) {
        reurlController.PatchReurl(c, db)
    })

    // Delete a mapping by ID (protected)
    r.DELETE("/:id", middlewares.AuthRequired(), func(c *gin.Context) {
        reurlController.DeleteReurl(c, db)
    })

    // Public redirect endpoint. Example: GET /reurl/redirect/:key
    // This will perform an HTTP redirect to the configured URL if present and not expired.
    r.GET("/redirect/:key", func(c *gin.Context) {
        reurlController.Redirect(c, db)
    })
}
