package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImdyaW1tIn0.G2xNcdMlunsSioRXkcR3NbBNnTPMF94UdSs-Ke_kV74"

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HelloWorldHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"hello": world}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRollDiceHandler(t *testing.T) {
	go BroadcastRoll()
	req, err := http.NewRequest("GET", "/roll?value=1d20", nil)
	if err != nil {
		t.Fatal(err)
	}
	jwtMiddleware := setUpJwtMiddleware()

	rr := httptest.NewRecorder()
	handler := jwtMiddleware.Handler(RollDiceHandler)

	req.Header.Set("Authorization", "Bearer "+token)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var roll DiceResponse
	json.Unmarshal(rr.Body.Bytes(), &roll)

	log.Printf("roll %v", roll)

	if roll.User.Username != "grimm" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			roll.User.Username, "grimm")
	}

	if len(roll.DiceRoll.DiceList) != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(roll.DiceRoll.DiceList), 1)
	}
	if roll.DiceRoll.OverallRollValue > 20 || roll.DiceRoll.OverallRollValue < 1 {
		t.Errorf("handler returned unexpected body: got %v wanted between %v and %v ",
			len(roll.DiceRoll.DiceList), 1, 20)
	}
}


func TestInitiativeDiceHandler(t *testing.T) {
	go HandleInitiative()
	req, err := http.NewRequest("POST", "/init", nil)
	if err != nil {
		t.Fatal(err)
	}
	jwtMiddleware := setUpJwtMiddleware()

	rr := httptest.NewRecorder()
	handler := jwtMiddleware.Handler(RollDiceHandler)

	req.Header.Set("Authorization", "Bearer "+token)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func setUpJwtMiddleware() *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})
}
