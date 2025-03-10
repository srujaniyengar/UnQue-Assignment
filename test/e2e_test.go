package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"UnQue/configs"
	"UnQue/routes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

type Appointment struct {
	ID           string `json:"id"`
	Student      string `json:"student"`
	Professor    string `json:"professor"`
	Availability string `json:"availability"`
	Status       string `json:"status"`
}

func perform_req(router http.Handler, req *http.Request, v interface{}) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if v != nil {
		if err := json.NewDecoder(rr.Body).Decode(v); err != nil {
			panic(err)
		}
	}
	return rr
}

func E2E_test(t *testing.T) {
	// feat: init DB connection
	configs.ConnectDB()

	// feat: prepare cleanup for student1 appointments
	student_a1_obj, err := primitive.ObjectIDFromHex("67c96fe328a7c547e2e47822")
	if err != nil {
		t.Fatalf("student1: invalid objectid: %v", err)
	}
	// fix: cleanup student1 appointments
	_, err = configs.DB.Collection("appointments").DeleteMany(context.Background(), bson.M{"student": student_a1_obj})
	if err != nil {
		t.Fatalf("student1: cleanup failed: %v", err)
	}

	// feat: setup router
	router := routes.SetupRoutes()

	// feat: student1 login
	student_a1_login := map[string]string{"email": "student@example.com", "password": "password"}
	var student_a1_resp LoginResponse
	body, _ := json.Marshal(student_a1_login)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	perform_req(router, req, &student_a1_resp)
	studentA1Token := student_a1_resp.Token
	t.Logf("student1: token: %s", studentA1Token)

	// feat: professor login
	proff_login := map[string]string{"email": "professor@example.com", "password": "password"}
	var proff_resp LoginResponse
	body, _ = json.Marshal(proff_login)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	perform_req(router, req, &proff_resp)
	professorToken := proff_resp.Token
	t.Logf("professor: token: %s", professorToken)

	// feat: professor sets availability
	availabilityPayload := map[string]interface{}{"slots": []string{"2025-03-10T09:00:00Z", "2025-03-10T10:00:00Z"}}
	body, _ = json.Marshal(availabilityPayload)
	req, _ = http.NewRequest("POST", "/availability", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+professorToken)
	var availabilityResp []map[string]interface{}
	perform_req(router, req, &availabilityResp)
	if len(availabilityResp) < 2 {
		t.Fatalf("professor: expected 2 avail slots, got %d", len(availabilityResp))
	}
	t.Logf("professor: availability: %v", availabilityResp)

	// feat: student1 views availability for professor
	req, _ = http.NewRequest("GET", "/availability?professor_id=67c9705a28a7c547e2e47823", nil)
	req.Header.Set("Authorization", "Bearer "+studentA1Token)
	var slotsResp []map[string]interface{}
	perform_req(router, req, &slotsResp)
	if len(slotsResp) == 0 {
		t.Fatalf("student1: no avail slots found")
	}
	t.Logf("student1: sees slots: %v", slotsResp)

	// feat: student1 books appointment for time T1
	app_paylod := map[string]string{"professor_id": "67c9705a28a7c547e2e47823", "slot": "2025-03-10T09:00:00Z"}
	body, _ = json.Marshal(app_paylod)
	req, _ = http.NewRequest("POST", "/appointments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+studentA1Token)
	var appointmentResp Appointment
	perform_req(router, req, &appointmentResp)
	studentA1AppointmentID := appointmentResp.ID
	t.Logf("student1: appointment id: %s", studentA1AppointmentID)

	// feat: student2 login
	student_a2_login := map[string]string{"email": "student2@example.com", "password": "password"}
	var student_a2_resp LoginResponse
	body, _ = json.Marshal(student_a2_login)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	perform_req(router, req, &student_a2_resp)
	studentA2Token := student_a2_resp.Token
	t.Logf("student2: token: %s", studentA2Token)

	// feat: student2 books appointment for time T2
	app_paylod = map[string]string{"professor_id": "67c9705a28a7c547e2e47823", "slot": "2025-03-10T10:00:00Z"}
	body, _ = json.Marshal(app_paylod)
	req, _ = http.NewRequest("POST", "/appointments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+studentA2Token)
	var app_resp Appointment
	perform_req(router, req, &app_resp)
	t.Logf("student2: appointment id: %s", app_resp.ID)

	// feat: professor cancels student1's appointment
	req, _ = http.NewRequest("DELETE", "/appointments/"+studentA1AppointmentID, nil)
	req.Header.Set("Authorization", "Bearer "+professorToken)
	var deleteResp map[string]string
	perform_req(router, req, &deleteResp)
	if deleteResp["message"] != "Appointment cancelled" {
		t.Fatalf("professor: cancellation failed: %v", deleteResp)
	}
	t.Log("professor: appointment cancelled")

	// feat: student1 checks appointments (should be empty)
	req, _ = http.NewRequest("GET", "/appointments", nil)
	req.Header.Set("Authorization", "Bearer "+studentA1Token)
	var appointmentsCheck []Appointment
	perform_req(router, req, &appointmentsCheck)
	if len(appointmentsCheck) != 0 {
		t.Fatalf("student1: expected no appointments, got: %v", appointmentsCheck)
	}
	t.Log("student1: no pending appointments")
}
