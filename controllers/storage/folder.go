package storage

import (
	"github.com/gin-gonic/gin"
)

type updateFolderRequest struct {
	Path string `json:"path"`
}

func CreateFolder(c *gin.Context) {
	folderPath, err := convertToStoragePath(c.Param("folder_path"), c)
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
	folderPath, err := convertToStoragePath(c.Param("folder_path"), c)
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
	folderPath, err := convertToStoragePath(c.Param("folder_path"), c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update folder"})
		return
	}

	updateReq := updateFolderRequest{}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	newFolderPath := updateReq.Path
	if newFolderPath != "" {
		newFolderPath, err = convertToStoragePath(newFolderPath, c)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid new folder path"})
			return
		}
		err := move(folderPath, newFolderPath)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to move folder"})
			return
		}
	}

	c.JSON(200, gin.H{"message": "Folder updated successfully"})
}

func DeleteFolder(c *gin.Context) {
	folderPath, err := convertToStoragePath(c.Param("folder_path"), c)
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
