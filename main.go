package main

import (
	"encoding/json"
	"fitnessBackend/allWorkouts"
	"fitnessBackend/workouts"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type workoutRequest struct {
	BodyParts string `json:"bodyParts"`
}

type workoutResponse struct {
	BodyParts     []string `json:"bodyParts"`
	GIFURL        string   `json:"gifUrl"`
	TargetMuscles []string `json:"targetMuscles"`
	Instructions  []string `json:"instructions"`
}

var targetMuscleGetters = map[string]func() workouts.Workouts{
	"abs":                   allworkouts.GetAbsWorkouts,
	"abductors":             allworkouts.GetAbductorsWorkouts,
	"adductors":             allworkouts.GetAdductorsWorkouts,
	"biceps":                allworkouts.GetBicepsWorkouts,
	"calves":                allworkouts.GetCalvesWorkouts,
	"cardiovascular system": allworkouts.GetCardiovascularSystemWorkouts,
	"delts":                 allworkouts.GetDeltsWorkouts,
	"forearms":              allworkouts.GetForearmsWorkouts,
	"glutes":                allworkouts.GetGlutesWorkouts,
	"hamstrings":            allworkouts.GetHamstringsWorkouts,
	"lats":                  allworkouts.GetLatsWorkouts,
	"levator scapulae":      allworkouts.GetLevatorScapulaeWorkouts,
	"pectorals":             allworkouts.GetPectoralsWorkouts,
	"quads":                 allworkouts.GetQuadsWorkouts,
	"serratus anterior":     allworkouts.GetSerratusAnteriorWorkouts,
	"spine":                 allworkouts.GetSpineWorkouts,
	"traps":                 allworkouts.GetTrapsWorkouts,
	"triceps":               allworkouts.GetTricepsWorkouts,
	"upper back":            allworkouts.GetUpperBackWorkouts,
	"chest":                 allworkouts.GetChestWorkouts,
}

func main() {
	server := &http.Server{
		Addr:              ":8080",
		Handler:           newHandler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Println("Workout API listening on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}

func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/workouts", handleWorkouts)
	return mux
}

func handleWorkouts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "use POST /workouts"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request workoutRequest
	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": `expected a JSON object such as {"bodyParts":"chest"}`})
		return
	}
	if err := decoder.Decode(new(interface{})); err != io.EOF {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request must contain exactly one JSON object"})
		return
	}

	targetMuscle := strings.ToLower(strings.TrimSpace(request.BodyParts))
	getter, ok := targetMuscleGetters[targetMuscle]
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bodyParts must name a target muscle: " + strings.Join(supportedTargetMuscles(), ", ")})
		return
	}

	selected := getter()
	response := make([]workoutResponse, 0, len(selected))
	for _, workout := range selected {
		response = append(response, workoutResponse{
			BodyParts: workout.BodyParts, GIFURL: workout.GIFURL,
			TargetMuscles: workout.TargetMuscles, Instructions: workout.Instructions,
		})
	}
	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}

func supportedTargetMuscles() []string {
	names := make([]string, 0, len(targetMuscleGetters))
	for name := range targetMuscleGetters {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
