package tasks

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"personal_site/controllers/storage"
)

// ClearTmpStorage 每小時自動清理 storageRoot/tmp 下超過一小時未修改的資料夾
func ClearTmpStorage() {
	go func() {
		for {
			storageRoot, err := storage.GetStorageRoot()
			if err != nil {
				log.Println("[ClearTmpStorage] get storage root error:", err)
				time.Sleep(time.Hour)
				continue
			}
			tmpDir := filepath.Join(storageRoot, "tmp")
			entries, err := os.ReadDir(tmpDir)
			if err != nil {
				log.Println("[ClearTmpStorage] read tmp dir error:", err)
				time.Sleep(time.Hour)
				continue
			}
			now := time.Now()
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				folderPath := filepath.Join(tmpDir, entry.Name())
				// 先處理第一層資料夾
				info, err := os.Stat(folderPath)
				if err != nil {
					log.Println("[ClearTmpStorage] stat error:", err, folderPath)
					continue
				}
				if now.Sub(info.ModTime()) > time.Hour {
					if err := os.RemoveAll(folderPath); err != nil {
						log.Println("[ClearTmpStorage] remove error:", err, folderPath)
					} else {
						log.Println("[ClearTmpStorage] removed:", folderPath)
					}
					continue
				}
				// 再處理第二層
				subEntries, err := os.ReadDir(folderPath)
				if err != nil {
					log.Println("[ClearTmpStorage] read subdir error:", err, folderPath)
					continue
				}
				for _, subEntry := range subEntries {
					if !subEntry.IsDir() {
						continue
					}
					subFolderPath := filepath.Join(folderPath, subEntry.Name())
					subInfo, err := os.Stat(subFolderPath)
					if err != nil {
						log.Println("[ClearTmpStorage] stat error:", err, subFolderPath)
						continue
					}
					if now.Sub(subInfo.ModTime()) > time.Hour {
						if err := os.RemoveAll(subFolderPath); err != nil {
							log.Println("[ClearTmpStorage] remove error:", err, subFolderPath)
						} else {
							log.Println("[ClearTmpStorage] removed:", subFolderPath)
						}
					}
				}
			}
			time.Sleep(time.Hour)
		}
	}()
}
