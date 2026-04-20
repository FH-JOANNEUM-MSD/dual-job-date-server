package repository

import (
	"bytes"
	"dual-job-date-server/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	defaultCompanyLogoBucket  = "company-logos"
	defaultCompanyImageBucket = "company-images"
	defaultLogoMaxBytes       = 5 * 1024 * 1024
)

type CompanyLogoUploadResult struct {
	CompanyID int    `json:"company_id"`
	Bucket    string `json:"bucket"`
	ObjectKey string `json:"object_key"`
	LogoURL   string `json:"logo_url"`
}

type CompanyImageUploadResult struct {
	CompanyID int    `json:"company_id"`
	Bucket    string `json:"bucket"`
	ObjectKey string `json:"object_key"`
	ImageURL  string `json:"image_url"`
}

type storageBucketInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Public bool   `json:"public"`
}

func UploadCompanyLogo(companyID int, originalName, contentType string, fileData []byte) (CompanyLogoUploadResult, error) {
	if len(fileData) == 0 {
		return CompanyLogoUploadResult{}, fmt.Errorf("leere datei")
	}

	company, err := GetCompanyByID(companyID)
	if err != nil {
		return CompanyLogoUploadResult{}, err
	}
	if company.ID == 0 {
		return CompanyLogoUploadResult{}, fmt.Errorf("company nicht gefunden")
	}

	normalizedContentType := normalizeContentType(contentType)
	if !isSupportedLogoContentType(normalizedContentType) {
		return CompanyLogoUploadResult{}, fmt.Errorf("nur png, jpg, jpeg oder webp erlaubt")
	}

	maxBytes := getLogoMaxBytes()
	if len(fileData) > maxBytes {
		return CompanyLogoUploadResult{}, fmt.Errorf("datei zu gross: max %d bytes", maxBytes)
	}

	bucket := getCompanyLogoBucket()
	if err := ensurePublicBucket(bucket); err != nil {
		return CompanyLogoUploadResult{}, err
	}

	objectKey := buildCompanyLogoObjectKey(company, originalName, normalizedContentType)
	if err := uploadToBucket(bucket, objectKey, normalizedContentType, fileData); err != nil {
		return CompanyLogoUploadResult{}, err
	}

	logoURL := buildPublicObjectURL(bucket, objectKey)
	if err := UpdateCompanyLogoURL(companyID, logoURL); err != nil {
		return CompanyLogoUploadResult{}, err
	}

	return CompanyLogoUploadResult{
		CompanyID: companyID,
		Bucket:    bucket,
		ObjectKey: objectKey,
		LogoURL:   logoURL,
	}, nil
}

func UploadCompanyImage(companyID int, originalName, contentType string, fileData []byte) (CompanyImageUploadResult, error) {
	if len(fileData) == 0 {
		return CompanyImageUploadResult{}, fmt.Errorf("leere datei")
	}

	company, err := GetCompanyByID(companyID)
	if err != nil {
		return CompanyImageUploadResult{}, err
	}
	if company.ID == 0 {
		return CompanyImageUploadResult{}, fmt.Errorf("company nicht gefunden")
	}

	normalizedContentType := normalizeContentType(contentType)
	if !isSupportedLogoContentType(normalizedContentType) {
		return CompanyImageUploadResult{}, fmt.Errorf("nur png, jpg, jpeg oder webp erlaubt")
	}

	maxBytes := getLogoMaxBytes()
	if len(fileData) > maxBytes {
		return CompanyImageUploadResult{}, fmt.Errorf("datei zu gross: max %d bytes", maxBytes)
	}

	bucket := getCompanyImageBucket()
	if err := ensurePublicBucket(bucket); err != nil {
		return CompanyImageUploadResult{}, err
	}

	objectKey := buildCompanyLogoObjectKey(company, originalName, normalizedContentType)
	if err := uploadToBucket(bucket, objectKey, normalizedContentType, fileData); err != nil {
		return CompanyImageUploadResult{}, err
	}

	imageURL := buildPublicObjectURL(bucket, objectKey)
	if err := AddCompanyImageURL(companyID, imageURL); err != nil {
		return CompanyImageUploadResult{}, err
	}

	return CompanyImageUploadResult{
		CompanyID: companyID,
		Bucket:    bucket,
		ObjectKey: objectKey,
		ImageURL:  imageURL,
	}, nil
}

func ensurePublicBucket(bucket string) error {
	info, found, err := getBucketByID(bucket)
	if err != nil {
		return err
	}

	if !found {
		return createPublicBucket(bucket)
	}

	if !info.Public {
		return setBucketPublic(bucket, info.Name)
	}

	return nil
}

func getBucketByID(bucket string) (storageBucketInfo, bool, error) {
	url := fmt.Sprintf("%s/storage/v1/bucket/%s", getSupabaseURL(), bucket)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return storageBucketInfo{}, false, err
	}
	attachSupabaseAuthHeaders(req)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return storageBucketInfo{}, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return storageBucketInfo{}, false, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		bodyText := strings.ToLower(string(body))
		// Supabase can return 400 with a payload-level 404 when bucket is missing.
		if strings.Contains(bodyText, "bucket not found") || strings.Contains(bodyText, `"statuscode":"404"`) {
			return storageBucketInfo{}, false, nil
		}
		return storageBucketInfo{}, false, fmt.Errorf("bucket lookup fehlgeschlagen (%d): %s", resp.StatusCode, string(body))
	}

	var info storageBucketInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return storageBucketInfo{}, false, err
	}

	return info, true, nil
}

func createPublicBucket(bucket string) error {
	payload := map[string]interface{}{
		"id":     bucket,
		"name":   bucket,
		"public": true,
	}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/storage/v1/bucket", getSupabaseURL())
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	attachSupabaseAuthHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	resBody, _ := io.ReadAll(resp.Body)
	// If bucket already exists (race/parallel), continue safely.
	if strings.Contains(strings.ToLower(string(resBody)), "already exists") || resp.StatusCode == http.StatusConflict {
		return nil
	}

	return fmt.Errorf("bucket erstellen fehlgeschlagen (%d): %s", resp.StatusCode, string(resBody))
}

func setBucketPublic(bucket, name string) error {
	if name == "" {
		name = bucket
	}

	payload := map[string]interface{}{
		"id":     bucket,
		"name":   name,
		"public": true,
	}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/storage/v1/bucket/%s", getSupabaseURL(), bucket)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	attachSupabaseAuthHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bucket update fehlgeschlagen (%d): %s", resp.StatusCode, string(resBody))
	}

	return nil
}

func uploadToBucket(bucket, objectKey, contentType string, fileData []byte) error {
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", getSupabaseURL(), bucket, objectKey)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(fileData))
	if err != nil {
		return err
	}
	attachSupabaseAuthHeaders(req)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload fehlgeschlagen (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func deleteObjectFromBucket(bucket, objectKey string) error {
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", getSupabaseURL(), bucket, objectKey)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	attachSupabaseAuthHeaders(req)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete fehlgeschlagen (%d): %s", resp.StatusCode, string(body))
	}
	return nil
}

func buildCompanyLogoObjectKey(company models.Company, originalName, contentType string) string {
	ext := strings.ToLower(filepath.Ext(originalName))
	if ext == "" {
		extensions, _ := mime.ExtensionsByType(contentType)
		if len(extensions) > 0 {
			ext = extensions[0]
		}
	}
	if ext == "" {
		ext = ".png"
	}

	companyPart := sanitizeFilePart(company.Name)
	if companyPart == "" {
		companyPart = "company"
	}

	return fmt.Sprintf("company-%d/%d-%s%s", company.ID, time.Now().UTC().Unix(), companyPart, ext)
}

func buildPublicObjectURL(bucket, objectKey string) string {
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", getSupabaseURL(), bucket, objectKey)
}

func DeleteCompanyImageObjectByURL(rawURL string) error {
	bucket, objectKey, ok := parsePublicObjectURL(rawURL)
	if !ok {
		return nil
	}
	if bucket != getCompanyImageBucket() {
		return nil
	}
	return deleteObjectFromBucket(bucket, objectKey)
}

func parsePublicObjectURL(raw string) (bucket, objectKey string, ok bool) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", "", false
	}
	base, err := url.Parse(getSupabaseURL())
	if err != nil {
		return "", "", false
	}
	if !strings.EqualFold(parsed.Scheme, base.Scheme) || !strings.EqualFold(parsed.Host, base.Host) {
		return "", "", false
	}

	path := parsed.EscapedPath()
	const prefix = "/storage/v1/object/public/"
	if !strings.HasPrefix(path, prefix) {
		return "", "", false
	}
	rest := strings.TrimPrefix(path, prefix)
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}

	unescapedKey, err := url.PathUnescape(parts[1])
	if err != nil {
		return "", "", false
	}
	return parts[0], unescapedKey, true
}

func attachSupabaseAuthHeaders(req *http.Request) {
	key := os.Getenv("SUPABASE_KEY")
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("apikey", key)
}

func getSupabaseURL() string {
	return strings.TrimRight(os.Getenv("SUPABASE_URL"), "/")
}

func getCompanyLogoBucket() string {
	if bucket := strings.TrimSpace(os.Getenv("COMPANY_LOGO_BUCKET")); bucket != "" {
		return bucket
	}
	return defaultCompanyLogoBucket
}

func getCompanyImageBucket() string {
	if bucket := strings.TrimSpace(os.Getenv("COMPANY_IMAGE_BUCKET")); bucket != "" {
		return bucket
	}
	return defaultCompanyImageBucket
}

func getLogoMaxBytes() int {
	raw := strings.TrimSpace(os.Getenv("COMPANY_LOGO_MAX_BYTES"))
	if raw == "" {
		return defaultLogoMaxBytes
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return defaultLogoMaxBytes
	}
	return parsed
}

func sanitizeFilePart(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " ", "-")

	var builder strings.Builder
	for _, ch := range value {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' {
			builder.WriteRune(ch)
		}
	}

	result := builder.String()
	result = strings.Trim(result, "-_")
	return result
}

func normalizeContentType(contentType string) string {
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	if idx := strings.Index(contentType, ";"); idx >= 0 {
		contentType = strings.TrimSpace(contentType[:idx])
	}
	return contentType
}

func isSupportedLogoContentType(contentType string) bool {
	switch contentType {
	case "image/png", "image/jpeg", "image/jpg", "image/webp":
		return true
	default:
		return false
	}
}
