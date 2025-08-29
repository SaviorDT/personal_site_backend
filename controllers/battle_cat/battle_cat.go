package battlecat

import (
	"personal_site/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FilterLevelsRequest struct {
	Stage   string   `form:"stage" binding:"required,max=3"`
	Enemies []string `form:"enemy" binding:"required,max=3"`
}

type levelData struct {
	Level   string `json:"level"`
	Name    string `json:"name"`
	HP      uint16 `json:"hp"`
	Enemies string `json:"enemies"`
}

type levelCollection struct {
	Enemies []string    `json:"enemies"`
	Levels  []levelData `json:"levels"`
}

func FilterLevels(c *gin.Context, db *gorm.DB) {
	var req FilterLevelsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var collections []levelCollection
	stageDB := db.Where("stage = ?", req.Stage)
	if len(req.Enemies) >= 3 {
		permuteFindLevels3(req.Enemies, &collections, stageDB)
	}
	if len(req.Enemies) >= 2 {
		permuteFindLevels2(req.Enemies, &collections, stageDB)
	}
	permuteFindLevels1(req.Enemies, &collections, stageDB)

	c.JSON(200, collections)
}

func permuteFindLevels1(enemies []string, collections *[]levelCollection, db *gorm.DB) {
	var levels []models.BattleCatLevel

	err := db.Where("INSTR(enemies, ?) > 0", enemies[0]).Find(&levels).Error
	if err == nil {
		var collection levelCollection
		collection.Enemies = []string{enemies[0]}

		for _, level := range levels {
			collection.Levels = append(collection.Levels, levelData{Level: level.Level, Name: level.Name, HP: level.HP, Enemies: level.Enemies})
		}

		*collections = append(*collections, collection)
	}

	if len(enemies) < 2 {
		return
	}

	err = db.Where("INSTR(enemies, ?) > 0", enemies[1]).Find(&levels).Error
	if err == nil {
		var collection levelCollection
		collection.Enemies = []string{enemies[1]}

		for _, level := range levels {
			collection.Levels = append(collection.Levels, levelData{Level: level.Level, Name: level.Name, HP: level.HP, Enemies: level.Enemies})
		}

		*collections = append(*collections, collection)
	}

	if len(enemies) < 3 {
		return
	}

	err = db.Where("INSTR(enemies, ?) > 0", enemies[2]).Find(&levels).Error
	if err == nil {
		var collection levelCollection
		collection.Enemies = []string{enemies[2]}

		for _, level := range levels {
			collection.Levels = append(collection.Levels, levelData{Level: level.Level, Name: level.Name, HP: level.HP, Enemies: level.Enemies})
		}

		*collections = append(*collections, collection)
	}
}

func permuteFindLevels2(enemies []string, collections *[]levelCollection, db *gorm.DB) {
	var levels []models.BattleCatLevel

	err := db.Where("INSTR(enemies, ?) > 0 AND INSTR(enemies, ?) > 0", enemies[0], enemies[1]).Find(&levels).Error
	if err == nil {
		var collection levelCollection
		collection.Enemies = []string{enemies[0], enemies[1]}
		for _, level := range levels {
			collection.Levels = append(collection.Levels, levelData{Level: level.Level, Name: level.Name, HP: level.HP, Enemies: level.Enemies})
		}

		*collections = append(*collections, collection)
	}

	if len(enemies) < 3 {
		return
	}

	err = db.Where("INSTR(enemies, ?) > 0 AND INSTR(enemies, ?) > 0", enemies[0], enemies[2]).Find(&levels).Error
	if err == nil {
		var collection levelCollection
		collection.Enemies = []string{enemies[0], enemies[2]}

		for _, level := range levels {
			collection.Levels = append(collection.Levels, levelData{Level: level.Level, Name: level.Name, HP: level.HP, Enemies: level.Enemies})
		}

		*collections = append(*collections, collection)
	}

	err = db.Where("INSTR(enemies, ?) > 0 AND INSTR(enemies, ?) > 0", enemies[1], enemies[2]).Find(&levels).Error
	if err == nil {
		var collection levelCollection
		collection.Enemies = []string{enemies[1], enemies[2]}

		for _, level := range levels {
			collection.Levels = append(collection.Levels, levelData{Level: level.Level, Name: level.Name, HP: level.HP, Enemies: level.Enemies})
		}

		*collections = append(*collections, collection)
	}
}

func permuteFindLevels3(enemies []string, collections *[]levelCollection, db *gorm.DB) {
	var levels []models.BattleCatLevel

	err := db.Where("INSTR(enemies, ?) > 0 AND INSTR(enemies, ?) > 0 AND INSTR(enemies, ?)", enemies[0], enemies[1], enemies[2]).Find(&levels).Error
	if err == nil {
		var collection levelCollection
		collection.Enemies = []string{enemies[0], enemies[1], enemies[2]}

		for _, level := range levels {
			collection.Levels = append(collection.Levels, levelData{Level: level.Level, Name: level.Name, HP: level.HP, Enemies: level.Enemies})
		}

		*collections = append(*collections, collection)
	}
}
