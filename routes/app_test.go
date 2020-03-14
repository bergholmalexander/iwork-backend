package routes

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go-api/db/postgres"
	"go-api/utils"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type testRouteConfig struct {
	Method    string
	URL       string
	Body      io.Reader
	URLParams map[string]string
	Handler   func(http.ResponseWriter, *http.Request)
	Headers   map[string]string
}

func executeReq(t *testing.T, config *testRouteConfig) *httptest.ResponseRecorder {
	req, err := http.NewRequest(config.Method, config.URL, config.Body)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, config.URLParams)
	for hKey, hVal := range config.Headers {
		req.Header.Add(hKey, hVal)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(config.Handler)
	handler.ServeHTTP(rr, req)
	return rr
}

type AppTestSuite struct {
	suite.Suite
	app *App
}

func (suite *AppTestSuite) SetupSuite() {
	dbUrl := os.Getenv("TEST_DB_URL")
	fixtures := []string{
		"../resources/tables.sql",
		"../test-fixtures/floors.sql",
		"../test-fixtures/users.sql",
		"../test-fixtures/workspaces.sql",
		"../test-fixtures/book_offer.sql",
	}
	if err := utils.RunFixturesOnDB(dbUrl, fixtures); err != nil {
		suite.FailNow("failed to create test db")
	}
	suite.app = NewTestApp()
}

func TestApp(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

// MOCKS
type mockDrive struct {
	mock.Mock
}

func (m *mockDrive) UploadFloorPlan(name string, content io.Reader) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}

func NewTestApp() *App {
	dbUrl := os.Getenv("TEST_DB_URL")
	store, err := postgres.NewPostgresDataStore(dbUrl)
	if err != nil {
		log.Println("Failed to connect to database")
		log.Fatal(err)
	}
	return &App{
		router: mux.NewRouter().StrictSlash(true),
		store:  store,
		gDrive: nil,
	}
}