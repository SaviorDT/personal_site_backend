package battlecat

import (
	"personal_site/models"
	"strings"

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
	var allPossibleLevels []models.BattleCatLevel
	db.Where("stage = ?", req.Stage).Find(&allPossibleLevels)

	if len(req.Enemies) >= 3 {
		permuteFindLevels3(req.Enemies, &collections, allPossibleLevels)
	}
	if len(req.Enemies) >= 2 {
		permuteFindLevels2(req.Enemies, &collections, allPossibleLevels)
	}
	permuteFindLevels1(req.Enemies, &collections, allPossibleLevels)

	c.JSON(200, collections)
}

func permuteFindLevels1(enemies []string, collections *[]levelCollection, allPossibleLevels []models.BattleCatLevel) {
	for _, enemy := range enemies {
		var collection levelCollection
		collection.Enemies = []string{enemy}

		for _, level := range allPossibleLevels {
			if containsEnemy(level.Enemies, enemy) {
				collection.Levels = append(collection.Levels, levelData{Level: level.Level, Name: level.Name, HP: level.HP, Enemies: level.Enemies})
			}
		}
		*collections = append(*collections, collection)
	}
}

func permuteFindLevels2(enemies []string, collections *[]levelCollection, allPossibleLevels []models.BattleCatLevel) {
	for i, enemy1 := range enemies {
		for _, enemy2 := range enemies[i+1:] {
			var collection levelCollection
			collection.Enemies = []string{enemy1, enemy2}

			for _, level := range allPossibleLevels {
				if containsEnemy(level.Enemies, enemy1) && containsEnemy(level.Enemies, enemy2) {
					collection.Levels = append(collection.Levels, levelData{Level: level.Level, Name: level.Name, HP: level.HP, Enemies: level.Enemies})
				}
			}
			*collections = append(*collections, collection)
		}
	}
}

func permuteFindLevels3(enemies []string, collections *[]levelCollection, allPossibleLevels []models.BattleCatLevel) {
	var collection levelCollection
	collection.Enemies = enemies

	for _, level := range allPossibleLevels {
		if containsEnemy(level.Enemies, enemies[0]) && containsEnemy(level.Enemies, enemies[1]) && containsEnemy(level.Enemies, enemies[2]) {
			collection.Levels = append(collection.Levels, levelData{Level: level.Level, Name: level.Name, HP: level.HP, Enemies: level.Enemies})
		}
	}
	*collections = append(*collections, collection)
}

func containsEnemy(enemies string, target string) bool {
	enemyList := strings.Split(enemies, "、")
	for _, enemy := range enemyList {
		if enemy == target {
			return true
		}
	}
	return false
}