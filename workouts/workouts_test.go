package workouts_test

import (
	"fitnessBackend/allWorkouts"
	"fitnessBackend/workouts"
	"reflect"
	"testing"
)

func TestEmbeddedWorkouts(t *testing.T) {
	if len(workouts.AllWorkouts) == 0 {
		t.Fatal("embedded dataset is empty")
	}
	for _, workout := range workouts.AllWorkouts {
		if workout.ExerciseID == "" || workout.Name == "" || workout.GIFURL == "" || len(workout.Instructions) == 0 {
			t.Fatalf("missing exercise information: %+v", workout)
		}
	}

	if len(allworkouts.GetChestWorkouts().Bodyweight()) == 0 {
		t.Fatal("expected bodyweight chest workouts")
	}

}

func TestTargetMuscleFilter(t *testing.T) {
	primary := workouts.Workout{ExerciseID: "primary", TargetMuscles: []string{"triceps"}}
	dataset := workouts.Workouts{
		primary,
		{ExerciseID: "secondary", TargetMuscles: []string{"pectorals"}, SecondaryMuscles: []string{"triceps"}},
	}
	if got := dataset.ForTargetMuscle(" TRICEPS "); !reflect.DeepEqual(got, workouts.Workouts{primary}) {
		t.Fatalf("expected only the primary triceps exercise, got %+v", got)
	}
	if len(dataset.ForTargetMuscle("unknown")) != 0 {
		t.Fatal("unknown muscle should return no exercises")
	}
}

func TestFiltersAndMuscles(t *testing.T) {
	workout := workouts.Workout{
		ExerciseID: "example", Name: "Example exercise", GIFURL: "example.gif",
		BodyParts: []string{"upper legs", "lower legs"}, Equipments: []string{"body weight"},
		TargetMuscles: []string{"quads"}, SecondaryMuscles: []string{"calves", "quads"},
		Instructions: []string{"Example instruction"},
	}
	dataset := workouts.Workouts{workout, {BodyParts: []string{"chest"}, Equipments: []string{"dumbbell"}}}
	selection := dataset.ForTargetMuscle(" QUADS ").Bodyweight()
	if !reflect.DeepEqual(selection, workouts.Workouts{workout}) {
		t.Fatalf("expected one complete exercise, got %+v", selection)
	}
	if got := selection.MuscleGroups(); !reflect.DeepEqual(got, []string{"calves", "quads"}) {
		t.Fatalf("unexpected muscle groups: %v", got)
	}
	if len(dataset.ForTargetMuscle("unknown")) != 0 {
		t.Fatal("unknown target muscle should return no exercises")
	}
}
