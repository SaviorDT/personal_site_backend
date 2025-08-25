package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"personal_site/controllers/utils"

	"github.com/gin-gonic/gin"
)

// mkDirIfNotExists 創建目錄（如果不存在）
// 參數 dir 應該是相對於專案根目錄的路徑
// 例如: "storage/uploads", "tmp/cache", "logs" 等
// 函式會自動將相對路徑轉換為基於專案根目錄的絕對路徑
func mkDirIfNotExists(folderPath string) error {
	// 檢查目錄是否存在
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 創建目錄，包括所有必要的父目錄
		err = os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func getFolderContent(folderPath string) ([]map[string]any, error) {
	folderContent := make([]map[string]any, 0)

	// 獲取目錄內容
	contents, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	// 將目錄內容轉換為 map 格式
	for _, entry := range contents {
		entryInfo, err := entry.Info()

		var size int64
		if err == nil {
			size = entryInfo.Size()
		} else {
			size = 0
		}

		item := map[string]any{
			"name":   entry.Name(),
			"is_dir": entry.IsDir(),
			"size":   size,
		}
		folderContent = append(folderContent, item)
	}

	return folderContent, nil
}

func move(oldPath, newPath string) error {
	// 檢查目錄是否存在
	if _, err := os.Stat(oldPath); err != nil {
		return err
	}

	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("destination directory already exists: %s", newPath)
	}

	// 移動目錄
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}

	return nil
}

func rmdir(folderPath string) error {
	// 刪除目錄及其內容
	return os.RemoveAll(folderPath)
}

func remove(filePath string) error {
	// 刪除檔案
	return os.Remove(filePath)
}

func writeMultipartFile(filePath string, file multipart.File) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}

	return nil
}

func convertToTmpDataPath(path string, c *gin.Context) (string, error) {
	userID := utils.GetUserID(c)
	tmpDataPath := filepath.Join("tmp", fmt.Sprintf("%d", userID), path)
	storageRoot, err := GetStorageRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(storageRoot, tmpDataPath), nil
}

// func convertToMetadataPath(path string, c *gin.Context) (string, error) {
// 	userID := utils.GetUserID(c)
// 	metadataPath := filepath.Join("data", fmt.Sprintf("%d", userID), "metadata", path)
// 	storageRoot, err := GetStorageRoot()
// 	if err != nil {
// 		return "", err
// 	}
// 	return filepath.Join(storageRoot, metadataPath), nil
// }

func convertToStoragePath(path string, c *gin.Context) (string, error) {
	userID := utils.GetUserID(c)
	userNickname := utils.GetUserNickname(c)
	storagePath := filepath.Join("data", fmt.Sprintf("%d", userID), userNickname, path)
	storageRoot, err := GetStorageRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(storageRoot, storagePath), nil
}

func GetStorageRoot() (string, error) {
	// 獲取專案根目錄路徑
	projectRoot, err := getProjectRoot()
	if err != nil {
		return "", err
	}

	// 返回專案根目錄的 storage 子目錄
	return filepath.Join(projectRoot, "storage"), nil
}

// getProjectRoot 獲取專案根目錄路徑
// 通過尋找 go.mod 檔案來確定專案根目錄
func getProjectRoot() (string, error) {
	// 從當前工作目錄開始
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 向上尋找包含 go.mod 的目錄
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir, nil
		}

		// 移動到父目錄
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// 已經到達根目錄，但沒有找到 go.mod
			break
		}
		currentDir = parentDir
	}

	// 如果找不到 go.mod，返回當前工作目錄
	return os.Getwd()
}
