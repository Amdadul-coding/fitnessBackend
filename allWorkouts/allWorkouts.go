package allworkouts

import "fitnessBackend/workouts"

func GetAbsWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("abs")
}

func GetAbductorsWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("abductors")
}

func GetAdductorsWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("adductors")
}

func GetBicepsWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("biceps")
}

func GetCalvesWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("calves")
}

func GetCardiovascularSystemWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("cardiovascular system")
}

func GetDeltsWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("delts")
}

func GetForearmsWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("forearms")
}

func GetGlutesWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("glutes")
}

func GetHamstringsWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("hamstrings")
}

func GetLatsWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("lats")
}

func GetLevatorScapulaeWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("levator scapulae")
}

func GetPectoralsWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("pectorals")
}

func GetQuadsWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("quads")
}

func GetSerratusAnteriorWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("serratus anterior")
}

func GetSpineWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("spine")
}

func GetTrapsWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("traps")
}

func GetTricepsWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("triceps")
}

func GetUpperBackWorkouts() workouts.Workouts {
	return workouts.AllWorkouts.ForTargetMuscle("upper back")
}

// GetChestWorkouts selects the dataset's pectorals target muscle.
func GetChestWorkouts() workouts.Workouts {
	return GetPectoralsWorkouts()
}
