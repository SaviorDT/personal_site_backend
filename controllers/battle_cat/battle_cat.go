package battlecat

import (
	"fmt"
	"personal_site/models"
	"sort"
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
	Time	string `json:"time"`
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

	for i := range collections {
		sort.Slice(collections[i].Levels, func(a, b int) bool {
			return !LongerThan(collections[i].Levels[a].Time, collections[i].Levels[b].Time)
		})
	}

	c.JSON(200, collections)
}

func permuteFindLevels1(enemies []string, collections *[]levelCollection, allPossibleLevels []models.BattleCatLevel) {
	for _, enemy := range enemies {
		var collection levelCollection
		collection.Enemies = []string{enemy}

		for _, level := range allPossibleLevels {
			if containsEnemy(level.Enemies, enemy) {
				collection.Levels = append(collection.Levels, levelData{Level: level.Level, Time: calMaxTime(level.Enemies, []string{enemy}), Name: level.Name, HP: level.HP, Enemies: level.Enemies})
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

			for i, level := range allPossibleLevels {
				if containsEnemy(level.Enemies, enemy1) && containsEnemy(level.Enemies, enemy2) {
					allPossibleLevels[i].Enemies = ""
					collection.Levels = append(collection.Levels, levelData{Level: level.Level, Time: calMaxTime(level.Enemies, []string{enemy1, enemy2}), Name: level.Name, HP: level.HP, Enemies: level.Enemies})
				}
			}
			*collections = append(*collections, collection)
		}
	}
}

func permuteFindLevels3(enemies []string, collections *[]levelCollection, allPossibleLevels []models.BattleCatLevel) {
	var collection levelCollection
	collection.Enemies = enemies

	for i, level := range allPossibleLevels {
		if containsEnemy(level.Enemies, enemies[0]) && containsEnemy(level.Enemies, enemies[1]) && containsEnemy(level.Enemies, enemies[2]) {
			allPossibleLevels[i].Enemies = ""
			collection.Levels = append(collection.Levels, levelData{Level: level.Level, Time: calMaxTime(level.Enemies, enemies), Name: level.Name, HP: level.HP, Enemies: level.Enemies})
		}
	}
	*collections = append(*collections, collection)
}

func containsEnemy(enemies string, target string) bool {
	enemyList := strings.Split(enemies, "、")
	for _, enemy := range enemyList {
		if idx := strings.Index(enemy, "（"); idx != -1 {
			enemy = enemy[:idx]
		}
		enemy = strings.TrimSpace(enemy)
		if enemy == target {
			return true
		}
	}
	return false
}

func calMaxTime(levelEnemies string, targetEnemies []string) string {
	var result string = "101%"

	enemyList := strings.Split(levelEnemies, "、")
	for _, enemy := range enemyList {
		var time string = "101%"
		if idx := strings.Index(enemy, "（"); idx != -1 {
			original := enemy
			enemy = enemy[:idx]
			time = original[idx+3 : len(original)-3]
		}
		enemy = strings.TrimSpace(enemy)
		for _, target := range targetEnemies {
			if enemy == target {
				result = longerTime(result, time)
				break
			}
		}
	}
	return result
}

func longerTime(time1 string, time2 string) string {
	if LongerThan(time1, time2) {
		return time1
	}
	return time2
}

func LongerThan(time1 string, time2 string) bool {
	var val1, val2 int
	var format1, format2 byte
	fmt.Sscanf(time1[:len(time1)-1], "%d", &val1)
	fmt.Sscanf(time2[:len(time2)-1], "%d", &val2)
	format1 = time1[len(time1)-1]
	format2 = time2[len(time2)-1]

	if format1 != format2 {
		if format1 == 's' {
			return true
		}
		return false
	}

	if format1 == 's' {
		if val1 > val2 {
			return true
		}
		return false
	}

	// 10% longer than 20%
	if val1 > val2 {
		return false
	}
	return true
}