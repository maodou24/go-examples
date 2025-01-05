package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	uploadDir = "./uploads"
	tempDir   = "./uploads/temp" // 用于存储临时分片
)

func main() {
	// 创建必要的目录
	for _, dir := range []string{uploadDir, tempDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}
	}

	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(gin.Logger())

	// 修正模板路径
	e.LoadHTMLFiles("index.html") // 改用 LoadHTMLFiles，直接加载 index.html

	e.GET("/", serveHTML)
	e.GET("/status", fileStatus)
	e.POST("/upload", handleUpload)

	e.Run(":8080")
}

type UploadInfo struct {
	FilePath     string
	Uploaded     map[int]bool
	TotalChunks  int
	FileSize     int64
	LastModified int64
	mu           sync.Mutex
}

var (
	uploadInfos = make(map[string]UploadInfo)
	uploadMu    sync.RWMutex
)

func fileStatus(ctx *gin.Context) {
	filename := ctx.Query("filename")
	totalChunks, _ := strconv.Atoi(ctx.Query("totalChunks"))
	fileSize, _ := strconv.ParseInt(ctx.Query("fileSize"), 10, 64)
	lastModified, _ := strconv.ParseInt(ctx.Query("lastModified"), 10, 64)

	uploadMu.Lock()
	info, exists := uploadInfos[filename]
	if exists {
		// 检查文件信息是否匹配
		if info.FileSize != fileSize || info.LastModified != lastModified {
			// 文件已更改，清除旧的上传记录
			delete(uploadInfos, filename)
			exists = false
		}
	}

	if !exists {
		info = UploadInfo{
			FilePath:     path.Join(tempDir, filename),
			Uploaded:     make(map[int]bool),
			TotalChunks:  totalChunks,
			FileSize:     fileSize,
			LastModified: lastModified,
		}
		uploadInfos[filename] = info
	}
	uploadMu.Unlock()

	info.mu.Lock()
	uploadedChunks := make([]int, 0)
	for i := 0; i < info.TotalChunks; i++ {
		if info.Uploaded[i] {
			uploadedChunks = append(uploadedChunks, i)
		}
	}
	info.mu.Unlock()

	ctx.JSON(http.StatusOK, gin.H{
		"uploadedChunks": uploadedChunks,
		"completed":      len(uploadedChunks) == totalChunks,
	})
}

func handleUpload(ctx *gin.Context) {
	chunkNumber, _ := strconv.Atoi(ctx.PostForm("chunkNumber"))
	totalChunks, _ := strconv.Atoi(ctx.PostForm("totalChunks"))
	filename := ctx.PostForm("filename")

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无法获取上传文件"})
		return
	}

	uploadMu.RLock()
	info, exists := uploadInfos[filename]
	uploadMu.RUnlock()

	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "未找到文件信息"})
		return
	}

	// 确保分片目录存在
	chunkPath := path.Join(tempDir, fmt.Sprintf("%s_chunk_%d", filename, chunkNumber))

	// 保存分片
	if err := ctx.SaveUploadedFile(file, chunkPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存分片失败"})
		return
	}

	// 更新上传状态
	info.mu.Lock()
	info.Uploaded[chunkNumber] = true
	uploadedCount := len(info.Uploaded)
	isComplete := uploadedCount == totalChunks
	info.mu.Unlock()

	// 如果所有分片都已上传，则合并文件
	if isComplete {
		if err := mergeChunks(filename, totalChunks); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "合并文件失败"})
			return
		}
		// 清理分片和上传信息
		cleanupChunks(filename, totalChunks)
		uploadMu.Lock()
		delete(uploadInfos, filename)
		uploadMu.Unlock()
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"completed": isComplete,
	})
}

// 添加合并文件的函数
func mergeChunks(filename string, totalChunks int) error {
	destPath := path.Join(uploadDir, filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for i := 0; i < totalChunks; i++ {
		chunkPath := path.Join(tempDir, fmt.Sprintf("%s_chunk_%d", filename, i))
		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			return err
		}
		if _, err := destFile.Write(chunkData); err != nil {
			return err
		}
	}
	return nil
}

// 添加清理分片的函数
func cleanupChunks(filename string, totalChunks int) {
	for i := 0; i < totalChunks; i++ {
		chunkPath := path.Join(tempDir, fmt.Sprintf("%s_chunk_%d", filename, i))
		os.Remove(chunkPath)
	}
}

func serveHTML(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"index": "index.html",
	})
}