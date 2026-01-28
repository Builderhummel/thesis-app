package protected_controller

import (
	"crypto/sha256"
	"encoding/hex"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	"github.com/gin-gonic/gin"
)

// HandleFileUpload handles file uploads for a thesis
func HandleFileUpload(c *gin.Context) {
	print("HandleFileUpload: Start of function")

	tuid := html.EscapeString(c.PostForm("tuid"))
	if tuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No thesis ID provided"})
		return
	}

	// Validate TUID
	tuidNum, err := strconv.Atoi(tuid)
	if err != nil || tuidNum < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid thesis ID"})
		return
	}

	// Get the uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Check file size
	if header.Size > Config.FileMaxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large. Maximum size is 50MB"})
		return
	}

	// Check filename length
	if len(header.Filename) > Config.FileMaxFilenameLen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename too long"})
		return
	}

	print("HandleFileUpload: after filename length check\n")

	// Get username from JWT token and convert to PDUID
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get PDUID from LoginHandle
	pduid, err := db_model.GetUserPuidFromUserId(user_id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
		return
	}

	// Get category from form (optional)
	category := html.EscapeString(c.PostForm("category"))

	// Generate unique filename using hash + timestamp
	ext := filepath.Ext(header.Filename)
	hash := sha256.New()
	hash.Write([]byte(header.Filename + time.Now().String()))
	uniqueName := hex.EncodeToString(hash.Sum(nil))[:16] + ext

	// Create directory for this thesis if it doesn't exist
	thesisDir := filepath.Join(Config.FileUploadDir, tuid)
	if err := os.MkdirAll(thesisDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create storage directory"})
		return
	}

	// Full path for the file
	filePath := filepath.Join(thesisDir, uniqueName)

	// Create the file on disk
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath) // Clean up on error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Insert file record into database
	_, err = db_model.InsertFileRecord(tuid, uniqueName, header.Filename, header.Size, pduid, category)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": header.Filename,
	})
}

// HandleFileDownload serves a file for download
func HandleFileDownload(c *gin.Context) {
	fuid := html.EscapeString(c.Query("fuid"))
	if fuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file ID provided"})
		return
	}

	// Get file info from database
	fileInfo, err := db_model.GetFileByID(fuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Construct file path
	filePath := filepath.Join(Config.FileUploadDir, fileInfo.TUID, fileInfo.FileName)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found on disk"})
		return
	}

	// Set headers for download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.OriginalFileName)
	c.Header("Content-Type", "application/octet-stream")

	// Serve the file
	c.File(filePath)
}

// HandleFileList returns a list of files for a thesis
func HandleFileList(c *gin.Context) {
	tuid := html.EscapeString(c.Query("tuid"))
	if tuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No thesis ID provided"})
		return
	}

	// Validate TUID
	tuidNum, err := strconv.Atoi(tuid)
	if err != nil || tuidNum < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid thesis ID"})
		return
	}

	// Get files from database
	files, err := db_model.GetFilesByThesis(tuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve files"})
		return
	}

	// Format file info for response
	type FileInfo struct {
		FUID             int64  `json:"fuid"`
		OriginalFileName string `json:"filename"`
		FileSize         int64  `json:"size"`
		UploadDate       string `json:"upload_date"`
		PDUID            string `json:"pduid"`
		Category         string `json:"category"`
	}

	var fileList []FileInfo
	for _, f := range files {
		fileList = append(fileList, FileInfo{
			FUID:             f.FUID,
			OriginalFileName: f.OriginalFileName,
			FileSize:         f.FileSize,
			UploadDate:       f.UploadDate.Format("2006-01-02 15:04:05"),
			PDUID:            f.PDUID,
			Category:         f.Category,
		})
	}

	c.JSON(http.StatusOK, gin.H{"files": fileList})
}

// HandleFileDelete deletes a file
func HandleFileDelete(c *gin.Context) {
	fuid := html.EscapeString(c.Query("fuid"))
	if fuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file ID provided"})
		return
	}

	// Get file info from database
	fileInfo, err := db_model.GetFileByID(fuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Construct file path
	filePath := filepath.Join(Config.FileUploadDir, fileInfo.TUID, fileInfo.FileName)

	// Delete file from disk
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from disk"})
		return
	}

	// Delete file record from database
	if err := db_model.DeleteFileRecord(fuid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
