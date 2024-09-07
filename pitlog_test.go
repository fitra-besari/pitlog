package pitlog_test

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"pitlog"
	"testing"
)

// Mock Logger untuk testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args...)
}

// Helper function untuk setup logger menggunakan MockLogger
func setupLogger() (*pitlog.Pitlog_base, *MockLogger) {
	mockLogger := new(MockLogger)

	// Buat directory log jika belum ada
	_ = os.Mkdir("./logs", os.ModePerm)

	// Gunakan New_pitlog untuk membuat instance pitlog_base
	pitlogInstance, err := pitlog.New_pitlog("TestApp", "1.0.0", "development", "./logs", "true", "false", "true")
	if err != nil {
		panic(err) // Jika gagal, hentikan pengujian
	}

	// Di sini kita tidak dapat mengganti 'dedicated' secara langsung karena field privat, jadi logger internal akan digunakan
	// Kode tetap berjalan sesuai log yang dihasilkan oleh logger internal

	return pitlogInstance.(*pitlog.Pitlog_base), mockLogger
}

func TestMakeLogString(t *testing.T) {
	DI, mockLogger := setupLogger()

	// Set expectation pada mock logger
	expectedMessage := "\nINFO request_id - test-id - Test Title -=> Test Message"
	mockLogger.On("Info", expectedMessage).Return()

	// Jalankan fungsi
	DI.Make_log_string("test-id", 1, "Test Title", "Test Message")

	// Verifikasi apakah log dicatat dengan benar
	mockLogger.AssertCalled(t, "Info", expectedMessage)
}

func TestMakeLogObject(t *testing.T) {
	DI, mockLogger := setupLogger()

	// Object yang akan dilog
	testObj := map[string]string{"key": "value"}

	// Set expectation pada mock logger
	jsonData, _ := json.Marshal(testObj)
	expectedMessage := fmt.Sprintf("\nINFO request_id - test-id - Test Title -=> %s", string(jsonData))
	mockLogger.On("Info", expectedMessage).Return()

	// Jalankan fungsi
	DI.Make_log_object("test-id", 1, "Test Title", testObj)

	// Verifikasi apakah log dicatat dengan benar
	mockLogger.AssertCalled(t, "Info", expectedMessage)
}

func TestApiLogMiddleware(t *testing.T) {
	DI, mockLogger := setupLogger()

	// Buat instance Echo
	e := echo.New()

	// Pasang middleware dan route
	DI.Api_log_middleware(e, []string{"password"})
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Buat request dan response untuk pengujian
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Jalankan request melalui middleware dan handler
	e.ServeHTTP(rec, req)

	// Pastikan response sesuai
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())

	// Verifikasi bahwa logger mencatat log
	mockLogger.AssertCalled(t, "Info", mock.AnythingOfType("string"))
}
