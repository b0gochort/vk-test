package repo

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/b0gochort/vk-test/internal/bot"
	"github.com/b0gochort/vk-test/internal/models"
)

func NewTeacher(t *models.Teacher, db *sql.DB) (int, error) {
	query := `INSERT INTO teachers (teacher_id, name, last_name) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, t.TeacherID, t.Name, t.LastName)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return t.TeacherID, nil
}

func TurnOffTeacherNotifyStatus(t *models.Teacher, db *sql.DB) error {
	query := `UPDATE teachers SET status_notified = false WHERE teacher_id = $1`
	_, err := db.Exec(query, t.TeacherID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func TurnOnTeacherNotifyStatus(t *models.Teacher, db *sql.DB) error {
	query := `UPDATE teachers SET status_notified = true WHERE teacher_id = $1`
	_, err := db.Exec(query, t.TeacherID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func DeleteTeacher(t *models.Teacher, db *sql.DB) error {
	// Удаляем всех учеников учителя
	_, err := db.Exec("DELETE FROM students WHERE teacher_id = $1", t.TeacherID)
	if err != nil {
		return err
	}

	// Удаляем учителя
	_, err = db.Exec("DELETE FROM teachers WHERE id = $1", t.TeacherID)
	if err != nil {
		return err
	}
	return nil
}

func NewStudent(s *models.Student, db *sql.DB) error {
	query := `INSERT INTO students (name, last_name, lesson_time, weekday, teacher_id) 
		VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, s.Name, s.LastName, s.LessonTime, s.Weekday, s.TeacherID)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Student with ID %d has been created", s.StudentID)
	return nil
}

func DeleteStudent(s *models.Student, db *sql.DB) error {
	query := `DELETE FROM students WHERE name = $1 AND last_name = $2;
				WHERE teacher_id = $3`
	_, err := db.Exec(query, s.Name, s.LastName)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Student with %s %s has been deleted", s.Name, s.LastName)
	return nil
}

func UpdateWeekDayStudent(s *models.Student, currentWeekDay int, db *sql.DB) error {
	query := `UPDATE students SET weekday = $1 WHERE name = $2 AND last_name = $3 AND weekday = $4`
	_, err := db.Exec(query, s.Weekday, s.Name, s.LastName, currentWeekDay)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Student with %s %s has been update", s.Name, s.LastName)
	return nil
}

func UpdateLessonTimeStudent(s *models.Student, currentLessonTime string, db *sql.DB) error {
	query := `UPDATE students SET lesson_time = $1 WHERE name = $2 AND last_name = $3 AND lesson_time = $4`
	lessonTime, err := time.Parse("15:04:05", s.LessonTime)
	if err != nil {
		return err
	}
	_, err = db.Exec(query, lessonTime, s.Name, s.LastName, currentLessonTime)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Student with %s %s has been update", s.Name, s.LastName)
	return nil
}

func CheckLessonTime(db *sql.DB, vk *api.VK) {
	for {
		day := int(time.Now().Weekday())
		// Получаем всех студентов, чье занятие начинается в течение следующей минуты
		query := `SELECT s.name AS student_name, s.last_name AS student_last_name, s.student_id AS student_id, t.name AS teacher_name, t.last_name AS teacher_last_name,s.teacher_id AS teacher_id
			FROM students AS s
			JOIN teachers AS t ON s.teacher_id = t.teacher_id
			WHERE EXTRACT(EPOCH FROM (s.lesson_time - NOW()::time without time zone))/60 <= 15 AND s.lesson_time > NOW()::time without time zone AND s.notified = false AND s.weekday = $1 AND t.status_notified = true;
		`
		rows, err := db.Query(query, day)
		if err != nil {
			log.Println(err)
			continue
		}

		// Обрабатываем каждую запись
		for rows.Next() {
			var s models.Student
			var t models.Teacher

			err := rows.Scan(&s.Name, &s.LastName, &s.StudentID, &t.Name, &t.LastName, &t.TeacherID)
			if err != nil {
				log.Println(err)
				continue
			}

			fmt.Printf("У студента %s %s занятие через 15 минут с преподавателем %s %s\n", s.Name, s.LastName, t.Name, t.LastName)

			// Отправляем уведомление пользователю
			err = bot.SendMessage(vk, t.TeacherID, fmt.Sprintf("У студента %s %s занятие через 15 минут\n", s.Name, s.LastName))
			if err != nil {
				log.Println(err)
			}

			// Обновляем поле notified для данного студента
			queryUpdateNotify := `UPDATE students SET notified = true WHERE student_id = $1`
			_, err = db.Exec(queryUpdateNotify, s.StudentID)
			if err != nil {
				log.Println(err)
			}
		}

		// Закрываем результаты запроса
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}

		// Ожидаем одну минуту перед следующей проверкой
		time.Sleep(1 * time.Minute)
	}
}

func UpdateNotifyStatus(db *sql.DB) {
	for {
		// Получаем всех студентов, чье занятие начинается в течение следующей минуты
		query := `SELECT s.student_id
		FROM students AS s
		JOIN teachers AS t ON s.teacher_id = t.teacher_id
		WHERE s.lesson_time < NOW()::time without time zone AND s.notified = true
		`
		rows, err := db.Query(query)
		if err != nil {
			log.Println(err)
			continue
		}

		// Обрабатываем каждую запись
		for rows.Next() {
			var s models.Student

			err := rows.Scan(&s.StudentID)
			if err != nil {
				log.Println(err)
				continue
			}

			// Обновляем поле notified для данного студента
			queryUpdateNotify := `UPDATE students SET notified = false WHERE student_id = $1`
			_, err = db.Exec(queryUpdateNotify, s.StudentID)
			if err != nil {
				log.Println(err)
			}
		}

		// Закрываем результаты запроса
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}

		// Ожидаем одну минуту перед следующей проверкой
		time.Sleep(1 * time.Minute)
	}
}

func GetStudentsForToday(teacherID int, db *sql.DB) (string, error) {

	today := int(time.Now().Weekday())
	var ans string

	// Формируем SQL-запрос
	query := `SELECT name, last_name, lesson_time
			  FROM students WHERE teacher_id=$1 AND weekday=$2`

	// Выполняем запрос
	rows, err := db.Query(query, teacherID, today)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Обрабатываем результаты запроса
	students := []string{}
	for rows.Next() {
		student := models.Student{}
		err := rows.Scan(&student.Name, &student.LastName, &student.LessonTime)
		if err != nil {
			return "", err
		}
		time, err := bot.GetTimeFromString(student.LessonTime)

		if err != nil {
			return "", err
		}

		students = append(students, fmt.Sprintf("%s %s %s", student.Name, student.LastName, time))
		ans = bot.ConcatenateStrings(students)
	}

	return ans, nil
}
