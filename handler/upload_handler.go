package handler

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const maxUploadBytes = 5 << 20 // 5 MB

// Allowed image types, detected from the file content (not its name).
var imageExt = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// UploadHandler stores admin-uploaded images on disk. Files are served at
// /uploads/... and referenced by that relative path.
type UploadHandler struct {
	dir string
}

func NewUploadHandler(dir string) *UploadHandler {
	return &UploadHandler{dir: dir}
}

// LocationImage stores a location photo under /uploads/locations.
func (h *UploadHandler) LocationImage(c *gin.Context) { h.save(c, "locations") }

// TrailImage stores a Jalur Kisah cover photo under /uploads/trails.
func (h *UploadHandler) TrailImage(c *gin.Context) { h.save(c, "trails") }

// save accepts multipart field "file" and returns {"url": "/uploads/<folder>/<id>.<ext>"}.
func (h *UploadHandler) save(c *gin.Context, folder string) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadBytes+1<<10)
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file gambar wajib diunggah (maksimal 5 MB)"})
		return
	}
	if fh.Size > maxUploadBytes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ukuran gambar maksimal 5 MB"})
		return
	}
	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file tidak bisa dibaca"})
		return
	}
	defer f.Close()

	head := make([]byte, 512)
	n, _ := io.ReadFull(f, head)
	ext, ok := imageExt[http.DetectContentType(head[:n])]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format gambar harus JPG, PNG, atau WebP"})
		return
	}

	sub := filepath.Join(h.dir, folder)
	if err := os.MkdirAll(sub, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan gambar"})
		return
	}
	b := make([]byte, 16)
	rand.Read(b)
	name := hex.EncodeToString(b) + ext

	out, err := os.Create(filepath.Join(sub, name))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan gambar"})
		return
	}
	defer out.Close()
	if _, err := out.Write(head[:n]); err == nil {
		_, err = io.Copy(out, f)
	}
	if err != nil {
		os.Remove(out.Name())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan gambar"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"url": "/uploads/" + folder + "/" + name})
}
