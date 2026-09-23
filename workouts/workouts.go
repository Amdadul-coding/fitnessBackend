package workouts

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Workout struct {
	ExerciseID       string   `json:"exerciseId"`
	Name             string   `json:"name"`
	GIFURL           string   `json:"gifUrl"`
	BodyParts        []string `json:"bodyParts"`
	Equipments       []string `json:"equipments"`
	TargetMuscles    []string `json:"targetMuscles"`
	SecondaryMuscles []string `json:"secondaryMuscles"`
	Instructions     []string `json:"instructions"`
}

type Workouts []Workout

// WorkoutType classifies body-weight-tagged exercises as home and all others as gym.
func (workout Workout) WorkoutType() string {
	if contains(workout.Equipments, "body weight") {
		return "home"
	}
	return "gym"
}

//go:embed workoutsData.json
var workoutsJSON []byte

var AllWorkouts = mustLoadWorkouts()

func mustLoadWorkouts() Workouts {
	var workouts Workouts
	if err := json.Unmarshal(workoutsJSON, &workouts); err != nil {
		panic(fmt.Errorf("parse embedded workoutsData.json: %w", err))
	}
	return workouts
}

// ForTargetMuscle selects exercises that primarily target the requested muscle.
func (workouts Workouts) ForTargetMuscle(muscle string) Workouts {
	result := make(Workouts, 0)
	for _, workout := range workouts {
		if contains(workout.TargetMuscles, muscle) {
			result = append(result, workout)
		}
	}
	return result
}

func (workouts Workouts) WithEquipment(equipment string) Workouts {
	result := make(Workouts, 0)
	for _, workout := range workouts {
		if contains(workout.Equipments, equipment) {
			result = append(result, workout)
		}
	}
	return result
}

func (workouts Workouts) Bodyweight() Workouts {
	return workouts.WithEquipment("body weight")
}

func (workouts Workouts) MuscleGroups() []string {
	seen := make(map[string]bool)
	for _, workout := range workouts {
		for _, muscle := range workout.TargetMuscles {
			seen[muscle] = true
		}
		for _, muscle := range workout.SecondaryMuscles {
			seen[muscle] = true
		}
	}
	muscles := make([]string, 0, len(seen))
	for muscle := range seen {
		muscles = append(muscles, muscle)
	}
	sort.Strings(muscles)
	return muscles
}

func contains(values []string, query string) bool {
	for _, value := range values {
		if strings.EqualFold(value, strings.TrimSpace(query)) {
			return true
		}
	}
	return false
}
