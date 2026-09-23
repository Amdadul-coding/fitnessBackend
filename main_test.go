package main

import (
	"encoding/json"
	"fitnessBackend/allWorkouts"
	"fitnessBackend/workouts"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestChestWorkoutsResponse(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/workouts", strings.NewReader(`{"bodyParts":" Chest "}`))
	w := httptest.NewRecorder()
	newHandler().ServeHTTP(w, r)
	if w.Code != http.StatusOK || w.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected response: %d %s", w.Code, w.Body.String())
	}
	var rows []map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &rows); err != nil {
		t.Fatal(err)
	}
	expected := allworkouts.GetChestWorkouts()
	if len(rows) != len(expected) || len(rows) == 0 {
		t.Fatalf("got %d workouts, want %d", len(rows), len(expected))
	}
	for i, row := range rows {
		if len(row) != 6 {
			t.Fatalf("expected exactly six response fields, got %v", row)
		}
		for key, value := range map[string]interface{}{
			"workoutType": expected[i].WorkoutType(),
			"name":        expected[i].Name,
			"bodyParts":   expected[i].BodyParts, "gifUrl": expected[i].GIFURL,
			"targetMuscles": expected[i].TargetMuscles, "instructions": expected[i].Instructions,
		} {
			encoded, err := json.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual([]byte(row[key]), encoded) {
				t.Errorf("exercise %d: incorrect %s", i, key)
			}
		}
	}
}

func TestTricepsWorkoutsResponse(t *testing.T) {
	w := httptest.NewRecorder()
	newHandler().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/workouts", strings.NewReader(`{"bodyParts":" TRICEPS "}`)))
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected response: %d %s", w.Code, w.Body.String())
	}
	var rows []workoutResponse
	if err := json.Unmarshal(w.Body.Bytes(), &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) == 0 || len(rows) != len(allworkouts.GetTricepsWorkouts()) {
		t.Fatalf("unexpected number of triceps exercises: %d", len(rows))
	}
	for _, row := range rows {
		found := false
		for _, muscle := range row.TargetMuscles {
			if muscle == "triceps" {
				found = true
			}
		}
		if !found {
			t.Fatalf("returned exercise does not target triceps: %+v", row)
		}
	}
}

func TestWorkoutRequestValidation(t *testing.T) {
	for _, body := range []string{
		``, `{`, `{}`, `null`, `{"bodyParts":null}`, `{"bodyParts":"unknown"}`,
		`{"bodyParts":["chest"]}`, `{"bodyParts":"chest","extra":true}`,
		`{"bodyParts":"chest","workoutType":"outdoors"}`, `{"bodyParts":"chest","workoutType":123}`,
		`{"bodyParts":"upper arms"}`, `{"bodyParts":"lower arms"}`, `{"bodyParts":"lower legs"}`,
		`{"bodyParts":"chest"} {}`, `{"bodyParts":"chest"} trailing`,
		`{"bodyParts":"` + strings.Repeat("x", 4096) + `"}`,
	} {
		t.Run(body[:minLength(len(body), 60)], func(t *testing.T) {
			w := httptest.NewRecorder()
			newHandler().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/workouts", strings.NewReader(body)))
			if w.Code != http.StatusBadRequest {
				t.Fatalf("got status %d, want 400", w.Code)
			}
		})
	}
	w := httptest.NewRecorder()
	newHandler().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/workouts", nil))
	if w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") != "POST" {
		t.Fatalf("expected 405 with Allow: POST, got %d", w.Code)
	}
}

func TestEveryDatasetTargetMuscle(t *testing.T) {
	expected := make(map[string][]workoutResponse)
	for _, workout := range workouts.AllWorkouts {
		for _, muscle := range workout.TargetMuscles {
			expected[muscle] = append(expected[muscle], workoutResponse{
				WorkoutType: workout.WorkoutType(),
				Name:        workout.Name,
				BodyParts:   workout.BodyParts, GIFURL: workout.GIFURL,
				TargetMuscles: workout.TargetMuscles, Instructions: workout.Instructions,
			})
		}
	}
	// Chest is a friendly alias for the dataset's pectorals label.
	expected["chest"] = expected["pectorals"]
	for muscle, want := range expected {
		t.Run(muscle, func(t *testing.T) {
			payload, err := json.Marshal(map[string]string{"bodyParts": " " + strings.ToUpper(muscle) + " "})
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			newHandler().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/workouts", strings.NewReader(string(payload))))
			if w.Code != http.StatusOK {
				t.Fatalf("unexpected response: %d %s", w.Code, w.Body.String())
			}
			var got []workoutResponse
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("response does not match target-muscle selection for %s: got %d exercises, want %d", muscle, len(got), len(want))
			}
		})
	}
}

func minLength(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestWorkoutTypeFiltering(t *testing.T) {
	for _, muscle := range supportedTargetMuscles() {
		for _, kind := range []string{"home", "gym"} {
			t.Run(muscle+"/"+kind, func(t *testing.T) {
				want := make([]workoutResponse, 0)
				for _, workout := range targetMuscleGetters[muscle]() {
					isHome := false
					for _, equipment := range workout.Equipments {
						if equipment == "body weight" {
							isHome = true
						}
					}
					if isHome != (kind == "home") {
						continue
					}
					want = append(want, workoutResponse{
						Name: workout.Name, WorkoutType: kind, BodyParts: workout.BodyParts,
						GIFURL: workout.GIFURL, TargetMuscles: workout.TargetMuscles, Instructions: workout.Instructions,
					})
				}
				payload, err := json.Marshal(map[string]string{"bodyParts": muscle, "workoutType": " " + strings.ToUpper(kind) + " "})
				if err != nil {
					t.Fatal(err)
				}
				w := httptest.NewRecorder()
				newHandler().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/workouts", strings.NewReader(string(payload))))
				if w.Code != http.StatusOK {
					t.Fatalf("unexpected response: %d %s", w.Code, w.Body.String())
				}
				var got []workoutResponse
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("incorrect %s exercises for %s: got %d, want %d", kind, muscle, len(got), len(want))
				}
			})
		}
	}
}
