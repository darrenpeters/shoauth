package shoauth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type expectedSuccessHandler struct{}

func (s *expectedSuccessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

type expectedFailureHandler struct{}

func (s *expectedFailureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

type unexpectedSuccessHandler struct {
	t *testing.T
}

func (s *unexpectedSuccessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.t.Errorf("Request should have failed, but succeeded instead!")
}

type unexpectedFailureHandler struct {
	t *testing.T
}

func (s *unexpectedFailureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.t.Errorf("Request should have succeeded, but failed instead!")
}

type testPersistence struct {
	exists  bool
	create  error
	session bool
}

func (t *testPersistence) InstallationExists(shopID string) bool {
	return t.exists
}

func (t *testPersistence) CreateInstallation(shopID string, accessToken string) error {
	return t.create
}

func (t *testPersistence) HasValidSession(r *http.Request) bool {
	return t.session
}

func TestValidSession(t *testing.T) {
	p := testPersistence{
		exists:  true,
		create:  errors.New("Installation already exists!"),
		session: true,
	}
	handler := NewShopifyOauthHandler(&expectedSuccessHandler{}, &unexpectedFailureHandler{t: t}, &p)
	testServer := httptest.NewServer(handler)
	defer testServer.Close()
	http.Get(testServer.URL + "/install")
}

func TestAlreadyInstalled(t *testing.T) {
	p := testPersistence{
		exists:  true,
		create:  errors.New("Installation already exists!"),
		session: false,
	}
	handler := NewShopifyOauthHandler(&unexpectedSuccessHandler{t: t}, &expectedFailureHandler{}, &p)
	testServer := httptest.NewServer(handler)
	defer testServer.Close()
	http.Get(testServer.URL + "/install")
}
