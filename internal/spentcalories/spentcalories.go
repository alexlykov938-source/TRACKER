package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format")
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid steps value")
	}

	trainingType := strings.TrimSpace(parts[1])

	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid data format")
	}

	return steps, trainingType, duration, nil
}


func distance(steps int, height float64) float64 {
	lengthStep := float64(steps) * (height * stepLengthCoefficient)
	return lengthStep / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration.Minutes() == 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("error parsing training data: %v", err)
	}

	switch trainingType {
	case "Бег", "Running":
		distance := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, duration)
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("error calculating running calories: %v", err)
		}
		return fmt.Sprintf("Тип тренировки: %s.\nДлительность: %s.\nДистанция: %.2f км.\nСкорость: %.2f км/ч.\nСожгли калорий: %.2f ккал.", trainingType, duration, distance, meanSpeed, calories), nil
	case "Ходьба", "Walking":
		distance := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, duration)
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("error calculating walking calories: %v", err)
		}
		return fmt.Sprintf("Тип тренировки: %s.\nДлительность: %s.\nДистанция: %.2f км.\nСкорость: %.2f км/ч.\nСожгли калорий: %.2f ккал.", trainingType, duration, distance, meanSpeed, calories), nil
	default:
		return "", fmt.Errorf("unknown training type: %s", trainingType)
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("invalid input values")
	}

	meanSpeed := meanSpeed(steps, height, duration)
	calories := (weight * meanSpeed * duration.Minutes()) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("invalid input values")
	}

	meanSpeed := meanSpeed(steps, height, duration)
	calories := (weight * meanSpeed * duration.Minutes()) / minInH * walkingCaloriesCoefficient
	return calories, nil
}
