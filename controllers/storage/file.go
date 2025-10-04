package storage

import (
	"io"
	"log"
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

	c.JSON(201, gin.H{"message": "File uploaded successfully"})
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

	chunkPath := filepath.Join(tmpDir, chunkIndexStr)
	err = writeMultipartFile(chunkPath, file)
	if err != nil {
		return err
	}

	if chunkIndex+1 == totalChunks {
		filePath, err := convertToStoragePath(c.Param("file_path"), c)
		if err != nil {
			return err
		}
		if err := mkDirIfNotExists(filepath.Dir(filePath)); err != nil {
			return err
		}

		// Merge in background to avoid blocking the request
		go func(tmpDir, filePath string, totalChunks int) {
			defer os.RemoveAll(tmpDir)

			finalOut, err := os.Create(filePath)
			if err != nil {
				log.Println("merge create final file error:", err)
				return
			}
			defer finalOut.Close()

			// Merge all chunks 0..totalChunks-1
			for i := 0; i < totalChunks; i++ {
				chunkFilePath := filepath.Join(tmpDir, strconv.Itoa(i))
				chunkFile, err := os.Open(chunkFilePath)
				if err != nil {
					log.Println("open chunk error:", err, "path:", chunkFilePath)
					return
				}
				if _, err = io.Copy(finalOut, chunkFile); err != nil {
					chunkFile.Close()
					log.Println("copy chunk error:", err, "path:", chunkFilePath)
					return
				}
				chunkFile.Close()
			}

			// Cleanup tmp dir
			if err := os.RemoveAll(tmpDir); err != nil {
				log.Println("cleanup tmp dir error:", err, "dir:", tmpDir)
			}
		}(tmpDir, filePath, totalChunks)

		return nil
	}

	return nil
}
