package storage

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type updateFolderRequest struct {
	Name string `json:"name"`
}

func CreateFolder(c *gin.Context) {
	folderPath, err := convertToDefaultStoragePath(c.Param("folder_name"), c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create directory"})
		return
	}

	err = mkDirIfNotExists(folderPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create directory"})
		return
	}

	c.JSON(200, gin.H{"message": "Directory created successfully"})
}

func ListFolder(c *gin.Context) {
	folderPath, err := convertToDefaultStoragePath(c.Param("folder_name"), c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list folder contents"})
		return
	}

	folderContent, err := getFolderContent(folderPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list folder contents"})
		return
	}

	c.JSON(200, folderContent)
}

func UpdateFolder(c *gin.Context) {
	folderPath, err := convertToDefaultStoragePath(c.Param("folder_name"), c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update folder"})
		return
	}

	updateReq := updateFolderRequest{}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	newFolderName := updateReq.Name
	if newFolderName != "" {
		err := renameFolder(folderPath, newFolderName)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to rename folder"})
			return
		}
	}

	c.JSON(200, gin.H{"message": "Folder updated successfully"})
}

func DeleteFolder(c *gin.Context) {
	folderPath, err := convertToDefaultStoragePath(c.Param("folder_name"), c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete folder"})
		return
	}

	err = rmdir(folderPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete folder"})
		return
	}

	c.JSON(200, gin.H{"message": "Folder deleted successfully"})
}

func renameFolder(oldFolderPath, newFolderName string) error {
	newFolderPath := filepath.Join(filepath.Dir(oldFolderPath), newFolderName)

	return move(oldFolderPath, newFolderPath)
}
