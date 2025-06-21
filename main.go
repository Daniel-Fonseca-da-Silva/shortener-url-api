package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var (
	urlStore    = make(map[string]string)
	secretKey   = []byte("12345678901234567890123456789012")
	mu          sync.Mutex
	lettersRune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	sugar       *zap.SugaredLogger
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	sugar = logger.Sugar()
	sugar.Info("URL Shortener service starting on port 8080")

	http.HandleFunc("/shorten", shortenUrl)
	http.HandleFunc("/", redirectHandler)

	sugar.Info("Server is running on port 8080")
	sugar.Fatal(http.ListenAndServe(":8080", nil))
}

func encrypt(orignalUrl string) (result string) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		sugar.Fatal("Failed to create cipher block", "error", err)
	}

	plainText := []byte(orignalUrl)
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]

	if _, err := rand.Read(iv); err != nil {
		sugar.Fatal("Failed to generate IV", "error", err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	result = hex.EncodeToString(cipherText)
	sugar.Debug("URL encrypted successfully", "originalLength", len(orignalUrl))
	return
}

// generateShortId Take a number and convert to base 64 to get a random letter or number
func generateShortId() (result string) {
	b := make([]rune, 6)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(lettersRune))))
		if err != nil {
			sugar.Fatal("Failed to generate random number", "error", err)
		}
		b[i] = lettersRune[num.Int64()]
	}
	result = string(b)
	sugar.Debug("Short ID generated", "id", result)
	return
}

func shortenUrl(w http.ResponseWriter, r *http.Request) {
	orinalUrl := r.URL.Query().Get("url")
	if orinalUrl == "" {
		sugar.Warn("Missing URL parameter in request")
		http.Error(w, "URL parameter in query is required", http.StatusBadRequest)
		return
	}

	if !(strings.HasPrefix(orinalUrl, "https://") || strings.HasPrefix(orinalUrl, "http://")) {
		sugar.Warn("Invalid URL format", "url", orinalUrl)
		http.Error(w, "URL parameter must have the value https:// or http://", http.StatusBadRequest)
		return
	}

	encryptedUrl := encrypt(orinalUrl)
	shortId := generateShortId()
	mu.Lock()
	urlStore[shortId] = encryptedUrl
	mu.Unlock()

	shortUrl := fmt.Sprintf("http://localhost:8080/%s", shortId)
	sugar.Info("URL shortened successfully", "originalUrl", orinalUrl, "shortId", shortId, "shortUrl", shortUrl)
	fmt.Fprintf(w, "The shortened url is: %s", shortUrl)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortId := r.URL.Path[1:]

	mu.Lock()
	encryptedUrl, ok := urlStore[shortId]
	mu.Unlock()

	if !ok {
		sugar.Warn("Short ID not found", "shortId", shortId)
		http.Error(w, "This url does not exist in our project", http.StatusNotFound)
		return
	}

	decryptedUrl := decrypt(encryptedUrl)
	sugar.Info("Redirecting to original URL", "shortId", shortId, "originalUrl", decryptedUrl)
	http.Redirect(w, r, decryptedUrl, http.StatusFound)
}

func decrypt(encryptedUrl string) (result string) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		sugar.Fatal("Failed to create cipher block for decryption", "error", err)
	}

	cipherText, err := hex.DecodeString(encryptedUrl)
	if err != nil {
		sugar.Fatal("Failed to decode hex string", "error", err)
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	result = string(cipherText)
	sugar.Debug("URL decrypted successfully", "decryptedLength", len(result))
	return
}
