package reurl

import (
    "errors"
    "fmt"
    "net/http"
    "personal_site/controllers/utils"
    "personal_site/models"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// Request/response shapes
type CreateReurlRequest struct {
    Key       *string `json:"key"`
    TargetURL string  `json:"target_url" binding:"required"`
    ExpiresIn *string `json:"expires_in"` // one of: "1h","12h","1d","7d","30d"
}

type PatchReurlRequest struct {
    Key       *string `json:"key"`
    TargetURL *string `json:"target_url"`
    ExpiresIn *string `json:"expires_in"` // nil means not changing; empty string means remove expiry
}


// CreateReurl handles creating a new reurl mapping. Auth required.
func CreateReurl(c *gin.Context, db *gorm.DB) {
    var req CreateReurlRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
        return
    }

    // owner
    user, err := utils.GetTokenUser(c)
    if err != nil || user.ID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
        return
    }

    // Validate expires: default to 7d when not specified
    var expiresAt *time.Time
    if req.ExpiresIn == nil {
        // call parseExpiresIn with "7d"
        s := "7d"
        tptr, err := parseExpiresIn(&s)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error parsing default expiry", "details": err.Error()})
            return
        }
        expiresAt = tptr
    } else {
        // use helper to parse
        tptr, err := parseExpiresIn(req.ExpiresIn)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_in value"})
            return
        }
        expiresAt = tptr
    }

    // determine key
    finalKey := ""
    if req.Key != nil && *req.Key != "" {
        finalKey = *req.Key
    } else {
        k, genErr := generateKey(db)
        if genErr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "key generation failed", "details": genErr.Error()})
            return
        }
        finalKey = k
    }

    // ensure uniqueness
    var existing models.Reurl
    if err := db.Where("`key` = ?", finalKey).First(&existing).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "key already exists"})
        return
    } else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "details": err.Error()})
        return
    }

    reurl := models.Reurl{
        Key:       finalKey,
        TargetURL: req.TargetURL,
        ExpiresAt: expiresAt,
        OwnerID:   user.ID,
    }

	_ = ClearExpiredUrls(db, &user.ID, &finalKey, nil)
    if err := db.Create(&reurl).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reurl", "details": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"data": reurl})
}

// ListReurls lists mappings. Admins see all, regular users see their own.
func ListReurls(c *gin.Context, db *gorm.DB) {
    user, _ := utils.GetTokenUser(c) // AuthRequired ensures existence

    var results []models.Reurl
    if utils.IsAdminUser(c) {
		_ = ClearExpiredUrls(db, nil, nil, nil)
        if err := db.Preload("Owner").Find(&results).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "details": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"data": results})
        return
    }

	_ = ClearExpiredUrls(db, &user.ID, nil, nil)
    if err := db.Where("owner_id = ?", user.ID).Find(&results).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": results})
}

// GetReurl returns a single mapping by ID. Admins can view any; users can view only their own.
func GetReurl(c *gin.Context, db *gorm.DB) {
    idStr := c.Param("id")

    // Parse id to uint
    var id uint
    if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    // Clear expired urls for this id
    _ = ClearExpiredUrls(db, nil, nil, &id)

    var reurl models.Reurl
    if err := db.Preload("Owner").First(&reurl, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "details": err.Error()})
        return
    }

    if !utils.IsAdminUser(c) {
        user, _ := utils.GetTokenUser(c)
        if reurl.OwnerID != user.ID {
            c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{"data": reurl})
}

// PatchReurl updates fields that are provided. Admins can update any; users only their own.
func PatchReurl(c *gin.Context, db *gorm.DB) {
    idStr := c.Param("id")

    // Parse id to uint
    var id uint
    if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    // Clear expired urls for this id
    _ = ClearExpiredUrls(db, nil, nil, &id)

    var reurl models.Reurl
    if err := db.First(&reurl, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "details": err.Error()})
        return
    }

    if !utils.IsAdminUser(c) {
        user, _ := utils.GetTokenUser(c)
        if reurl.OwnerID != user.ID {
            c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
            return
        }
    }

    var req PatchReurlRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
        return
    }

    // Update key
    if req.Key != nil {
        newKey := *req.Key
        if newKey != reurl.Key {
            var existing models.Reurl
            if err := db.Where("`key` = ?", newKey).First(&existing).Error; err == nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "key already exists"})
                return
            } else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "details": err.Error()})
                return
            }
            reurl.Key = newKey
        }
    }

    if req.TargetURL != nil {
        reurl.TargetURL = *req.TargetURL
    }

    if req.ExpiresIn != nil {
        // empty string means remove expiry
        if *req.ExpiresIn == "" {
            reurl.ExpiresAt = nil
        } else {
            tptr, err := parseExpiresIn(req.ExpiresIn)
            if err != nil || tptr == nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_in value"})
                return
            }
            reurl.ExpiresAt = tptr
        }
    }

    if err := db.Save(&reurl).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": reurl})
}

// DeleteReurl deletes a mapping. Admins can delete any; users only their own.
func DeleteReurl(c *gin.Context, db *gorm.DB) {
    idStr := c.Param("id")

    // Parse id to uint
    var id uint
    if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    // Clear expired urls for this id
    _ = ClearExpiredUrls(db, nil, nil, &id)

    var reurl models.Reurl
    if err := db.First(&reurl, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "details": err.Error()})
        return
    }

    if !utils.IsAdminUser(c) {
        user, _ := utils.GetTokenUser(c)
        if reurl.OwnerID != user.ID {
            c.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
            return
        }
    }

    if err := db.Delete(&reurl).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true})
}

// Redirect looks up mapping by key (public) and performs HTTP redirect if found and not expired.
func Redirect(c *gin.Context, db *gorm.DB) {
    key := c.Param("key")

    // Clear expired urls for this key
    _ = ClearExpiredUrls(db, nil, &key, nil)

    var reurl models.Reurl
    if err := db.Where("`key` = ?", key).First(&reurl).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "details": err.Error()})
        return
    }

    if reurl.ExpiresAt != nil && time.Now().After(*reurl.ExpiresAt) {
        _ = db.Delete(&reurl)
        c.JSON(http.StatusGone, gin.H{"error": "expired"})
        return
    }

    // Use 302 Found (temporary) redirect
    c.Redirect(http.StatusFound, reurl.TargetURL)
}