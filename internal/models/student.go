package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Student struct {
	StudentID  int
	Name       string
	LastName   string
	LessonTime string
	Weekday    int
	TeacherID  int
}

func ParseStudent(s *Student, str string) error {
	parts := strings.Split(str, " ")
	if len(parts) != 4 {
		return fmt.Errorf("неверный формат ввода")
	}

	name := strings.TrimSpace(parts[0])
	lastName := strings.TrimSpace(parts[1])
	weekday, err := strconv.Atoi(parts[2])
	if err != nil || weekday < 0 || weekday > 6 {
		return fmt.Errorf("неверный формат ввода дня недели")
	}
	lessonTime := strings.TrimSpace(parts[3])

	s.Name = name
	s.LastName = lastName
	s.Weekday = weekday
	s.LessonTime = lessonTime

	return nil
}

func ParseForDelete(s *Student, str string) error {
	parts := strings.Split(str, " ")
	if len(parts) != 2 {
		return fmt.Errorf("неверный формат ввода")
	}

	name := strings.TrimSpace(parts[0])
	lastName := strings.TrimSpace(parts[1])

	s.Name = name
	s.LastName = lastName

	return nil
}

func ParseForUpdateWeekday(s *Student, currentWeekday int, str string) error {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return fmt.Errorf("неверный формат ввода")
	}

	name := strings.TrimSpace(parts[0])
	lastName := strings.TrimSpace(parts[1])
	weekdayNew := strings.TrimSpace(parts[2])
	weekdayOld := strings.TrimSpace(parts[3])

	intWeekdayNew, err := strconv.Atoi(weekdayNew)
	if err != nil {
		return fmt.Errorf("неверный формат ввода")
	}
	intWeekdayOld, err := strconv.Atoi(weekdayOld)
	if err != nil {
		return fmt.Errorf("неверный формат ввода")
	}
	currentWeekday = intWeekdayOld
	s.Name = name
	s.LastName = lastName
	s.Weekday = intWeekdayNew

	return nil
}

func ParseForUpdateLessonTime(s *Student, currentLessonTime string, str string) error {
	parts := strings.Split(str, " ")
	if len(parts) != 4 {
		return fmt.Errorf("неверный формат ввода")
	}

	name := strings.TrimSpace(parts[0])
	lastName := strings.TrimSpace(parts[1])
	newLessonTime := strings.TrimSpace(parts[2])
	oldLessonTime := strings.TrimSpace(parts[3])

	currentLessonTime = oldLessonTime
	s.Name = name
	s.LastName = lastName
	s.LessonTime = newLessonTime

	return nil
}
