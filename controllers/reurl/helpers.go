package reurl

import (
    "errors"
    "personal_site/models"
    "time"

    "gorm.io/gorm"
)

// Allowed expiration keywords and their durations
var allowedExpires = map[string]time.Duration{
    "1h":  time.Hour,
    "12h": 12 * time.Hour,
    "1d":  24 * time.Hour,
    "7d":  7 * 24 * time.Hour,
    "30d": 30 * 24 * time.Hour,
}

// parseExpiresIn parses the optional expires_in string and returns an *time.Time
// If input is nil, it means do not change. If input is empty string, it means clear expiry (nil pointer).
func parseExpiresIn(in *string) (*time.Time, error) {
    if in == nil {
        return nil, nil
    }
    if *in == "" {
        return nil, nil
    }
    if d, ok := allowedExpires[*in]; ok {
        t := time.Now().Add(d)
        return &t, nil
    }
    return nil, errors.New("invalid expires_in value")
}

// findReurlByID returns the Reurl record by numeric id (string param). It returns gorm.ErrRecordNotFound if not found.
func findReurlByID(db *gorm.DB, id string) (models.Reurl, error) {
    var r models.Reurl
    if err := db.First(&r, id).Error; err != nil {
        return models.Reurl{}, err
    }
    return r, nil
}

// ClearExpiredUrls deletes expired reurl records based on the provided filters.
// This is used to improve user experience by cleaning up expired records before operations.
// It should not be used for determining if records exist in the database.
func ClearExpiredUrls(db *gorm.DB, owner *uint, key *string, id *uint) error {
    query := db.Model(&models.Reurl{}).Where("expires_at IS NOT NULL AND expires_at < ?", time.Now())

    if owner != nil {
        query = query.Where("owner_id = ?", *owner)
    }

    if key != nil {
        query = query.Where("`key` = ?", *key)
    }

    if id != nil {
        query = query.Where("id = ?", *id)
    }

    return query.Delete(&models.Reurl{}).Error
}
