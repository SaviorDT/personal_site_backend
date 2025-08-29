package storage

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type updateFileRequest struct {
	Path string `json:"path"`
}

func GetFile(c *gin.Context) {
	filePath, err := convertToStoragePath(c.Param("file_path"), c)
	if err != nil {
		c.JSON(400, gin.H{"error": "Cannot get file"})
		return
	}

	c.File(filePath)
}

func DeleteFile(c *gin.Context) {
	filePath, err := convertToStoragePath(c.Param("file_path"), c)
	if err != nil {
		c.JSON(400, gin.H{"error": "Cannot delete file"})
		return
	}

	err = remove(filePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(200, gin.H{"message": "File deleted successfully"})
}

func UploadFile(c *gin.Context) {
	err := saveFile(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully"})
}

func UpdateFile(c *gin.Context) {
	filePath, err := convertToStoragePath(c.Param("file_path"), c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update file"})
		return
	}

	updateReq := updateFileRequest{}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	newFilePath := updateReq.Path
	if newFilePath != "" {
		newFilePath, err = convertToStoragePath(newFilePath, c)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid new file path"})
			return
		}
		err := move(filePath, newFilePath)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to move file"})
			return
		}
	}

	c.JSON(200, gin.H{"message": "File updated successfully"})
}

func saveFile(c *gin.Context) error {
	fileID := c.PostForm("file_id")
	chunkIndexStr := c.PostForm("chunk_index")
	totalChunksStr := c.PostForm("total_chunks")

	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		return err
	}
	totalChunks, err := strconv.Atoi(totalChunksStr)
	if err != nil {
		return err
	}

	file, _, err := c.Request.FormFile("chunk_data")
	if err != nil {
		return err
	}
	defer file.Close()

	// 暫存目錄
	tmpDir, err := convertToTmpDataPath(fileID, c)
	if err != nil {
		return err
	}
	err = mkDirIfNotExists(tmpDir)
	if err != nil {
		return err
	}

	if chunkIndex+1 < totalChunks {
		chunkPath, err := convertToTmpDataPath(filepath.Join(fileID, chunkIndexStr), c)
		if err != nil {
			return err
		}
		err = writeMultipartFile(chunkPath, file)
		if err != nil {
			return err
		}
	} else {
		filePath, err := convertToStoragePath(c.Param("file_path"), c)
		if err != nil {
			return err
		}
		err = mkDirIfNotExists(filepath.Dir(filePath))
		if err != nil {
			return err
		}

		finalOut, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer finalOut.Close()

		// 合併所有已儲存的 chunk（0 到 totalChunks-2）
		for i := range totalChunks - 1 {
			chunkFilePath, err := convertToTmpDataPath(filepath.Join(fileID, strconv.Itoa(i)), c)
			if err != nil {
				return err
			}
			chunkFile, err := os.Open(chunkFilePath)
			if err != nil {
				return err
			}
			_, err = io.Copy(finalOut, chunkFile)
			chunkFile.Close()
			if err != nil {
				return err
			}
		}

		// last chunk
		_, err = io.Copy(finalOut, file)
		if err != nil {
			return err
		}

		rmdir(tmpDir)
	}

	return nil
}
