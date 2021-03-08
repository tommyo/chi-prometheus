package chiprometheus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Test_Logger(t *testing.T) {
	recorder := httptest.NewRecorder()

	n := chi.NewRouter()
	m := NewMiddleware("test")
	n.Use(m)

	n.Handle("/metrics", promhttp.Handler())

	n.Get(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	n.Get(`/users/{firstName}`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	req1path := "/ok"
	req1, err := http.NewRequest("GET", "http://localhost:3000"+req1path, nil)
	if err != nil {
		t.Error(err)
	}

	req2path := "/users/JoeBob"
	req2, err := http.NewRequest("GET", "http://localhost:3000"+req2path, nil)
	if err != nil {
		t.Error(err)
	}

	req3path := "/users/Misty"
	req3, err := http.NewRequest("GET", "http://localhost:3000"+req3path, nil)
	if err != nil {
		t.Error(err)
	}

	req4path := "/metrics"
	req4, err := http.NewRequest("GET", "http://localhost:3000"+req4path, nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req1)
	n.ServeHTTP(recorder, req2)
	n.ServeHTTP(recorder, req3)
	n.ServeHTTP(recorder, req4)

	body := recorder.Body.String()

	fmt.Println(body)

	if !strings.Contains(body, reqsName) {
		t.Errorf("body does not contain request total entry '%s'", reqsName)
	}
	if !strings.Contains(body, latencyName) {
		t.Errorf("body does not contain request duration entry '%s'", latencyName)
	}
	if !strings.Contains(body, `chi_request_duration_milliseconds_count{code="OK",method="GET",path="`+req1path+`",service="test"} 1`) {
		t.Errorf("body does not contain req1 count summary '%s'", req1path)
	}
	if strings.Contains(body, `chi_request_duration_milliseconds_count{code="OK",method="GET",path="`+req2path+`",service="test"} 1`) {
		t.Errorf("body should not contain req2 count summary '%s'", req2path)
	}
	if strings.Contains(body, `chi_request_duration_milliseconds_count{code="OK",method="GET",path="`+req3path+`",service="test"} 1`) {
		t.Errorf("body should not contain req3 count summary '%s'", req3path)
	}
	if !strings.Contains(body, `chi_request_duration_milliseconds_count{code="OK",method="GET",path="/users/{firstName}",service="test"} 2`) {
		t.Errorf("body does not contain patterned url count summary '%s'", "/users/{firstName}")
	}
}
