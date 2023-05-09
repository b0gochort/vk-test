package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/b0gochort/vk-test/internal/bot"
	"github.com/b0gochort/vk-test/internal/models"
	"github.com/b0gochort/vk-test/internal/repo"
	"github.com/b0gochort/vk-test/pkg"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	db, err := pkg.NewPostgres()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	token := os.Getenv("TOKEN")
	vk := api.NewVK(token)
	log.Println("get token and create newVk")
	go repo.CheckLessonTime(db, vk)
	go repo.UpdateNotifyStatus(db)

	// get information about the group
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("get group information")

	// Initializing Long Poll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("init poll")

	// New message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)

		b := params.NewMessagesSendBuilder()

		if strings.ToLower(obj.Message.Text) == "старт" {
			user := obj.Message.FromID
			userInfo, err := vk.UsersGet(api.Params{
				"user_ids": user,
				"fields":   "first_name,last_name",
			})
			if err != nil {
				log.Fatal(err)
			}
			var t models.Teacher
			t.TeacherID = user
			t.Name = userInfo[0].FirstName
			t.LastName = userInfo[0].LastName

			repo.NewTeacher(&t, db)

			b.Message("Вы зарегистрированы!\n Выберите соответствующий пункт меню:")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int())
			b.PeerID(obj.Message.PeerID)
			kb := bot.MainKeyBoard()
			b.Keyboard(kb)
			_, err = vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		}

		switch obj.Message.Text {

		case "Ученики":

			b.Message("Выберите соответствующий пункт меню:")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int()) // замените 123456 на любое уникальное число
			b.PeerID(obj.Message.PeerID)
			kb := bot.KeyboardStudents()
			b.Keyboard(kb)
			_, err := vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}

		case "Добавить ученика":
			b.Message("Введите данные студента в формате 'Имя Фамилия день недели урока(0 - воскресенье, 6 - суббота) ВремяУрока' (например, 'Тимофей Иванов 1 19:00:00')")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int()) // замените 123456 на любое уникальное число
			b.PeerID(obj.Message.PeerID)
			_, err := vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
			student := &models.Student{}
			student.TeacherID = obj.Message.FromID
			for {
				// Retrieve the user's response
				history, err := vk.MessagesGetHistory(api.Params{
					"peer_id": obj.Message.PeerID,
					"count":   1,
				})
				if err != nil {
					log.Fatal(err)
				}
				if len(history.Items) == 0 {
					continue
				}

				err = models.ParseStudent(student, history.Items[0].Text)
				if err != nil {
					b.Message("Неверный формат ввода. Попробуйте снова.")
					_, err = vk.MessagesSend(b.Params)
					if err != nil {
						log.Fatal(err)
					}
					continue
				}

				break

			}

			err = repo.NewStudent(student, db)
			if err != nil {
				log.Fatal(err)
			}
			b.Message("Студент был успешно добавлен")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int()) // замените 123456 на любое уникальное число
			b.PeerID(obj.Message.PeerID)
			kb := bot.KeyboardStudents()
			b.Keyboard(kb)
			_, err = vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}

		case "Удалить ученика":
			b.Message("Введите данные студента, которого собираетесь удалить в формате Имя Фамилия (например: Иван Иванов)")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int()) // замените 123456 на любое уникальное число
			b.PeerID(obj.Message.PeerID)
			_, err := vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
			student := &models.Student{}
			student.TeacherID = obj.Message.PeerID
			for {
				// Retrieve the user's response
				history, err := vk.MessagesGetHistory(api.Params{
					"peer_id": obj.Message.PeerID,
					"count":   1,
				})
				if err != nil {
					log.Fatal(err)
				}
				if len(history.Items) == 0 {
					continue
				}

				err = models.ParseForDelete(student, history.Items[0].Text)
				if err != nil {
					b.Message("Неверный формат ввода. Попробуйте снова.")
					_, err = vk.MessagesSend(b.Params)
					if err != nil {
						log.Fatal(err)
					}
					continue
				}

				break

			}
			repo.DeleteStudent(student, db)
			b.Message("Студент был успешно удалён")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int()) // замените 123456 на любое уникальное число
			b.PeerID(obj.Message.PeerID)
			kb := bot.KeyboardStudents()
			b.Keyboard(kb)
			_, err = vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		case "Уроки":

			b.Message("Выберите соответствующий пункт меню:")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int()) // замените 123456 на любое уникальное число
			b.PeerID(obj.Message.PeerID)
			kb := bot.KeyboardLessons()
			b.Keyboard(kb)
			_, err := vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		case "Расписание на сегодня":

			lessons, err := repo.GetStudentsForToday(obj.Message.PeerID, db)
			if err != nil {
				log.Println("Can't get lessons today", err)
			}

			b.Message(fmt.Sprintf("Все уроки на сегодня:\n %s", lessons))
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int()) // замените 123456 на любое уникальное число
			b.PeerID(obj.Message.PeerID)
			kb := bot.KeyboardLessons()
			b.Keyboard(kb)
			_, err = vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}

		case "Расписание на неделю":

			b.Message("TBD")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int()) // замените 123456 на любое уникальное число
			b.PeerID(obj.Message.PeerID)
			kb := bot.KeyboardLessons()
			b.Keyboard(kb)
			_, err = vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}

		case "Назад":
			b.Message("Выберите соответствующий пункт меню:")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int())
			b.PeerID(obj.Message.PeerID)
			kb := bot.MainKeyBoard()
			b.Keyboard(kb)
			_, err = vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}

		case "Изменить":
			b.Message("Выберите соответствующий пункт меню:")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int())
			b.PeerID(obj.Message.PeerID)
			kb := bot.KeyboardChange()
			b.Keyboard(kb)
			_, err = vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		case "Изменить день":
			b.Message("Введите данные студента,у которого собираетесь изменить день недели в формате Имя Фамилия НомерДняНовый НомерДняСтарый(0-воскресенье, 6-Суббота)(например: Иван Иванов 3 2)")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int()) // замените 123456 на любое уникальное число
			b.PeerID(obj.Message.PeerID)
			_, err := vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
			var currentWeekday int
			student := &models.Student{}
			student.TeacherID = obj.Message.PeerID
			for {
				// Retrieve the user's response
				history, err := vk.MessagesGetHistory(api.Params{
					"peer_id": obj.Message.PeerID,
					"count":   1,
				})
				if err != nil {
					log.Fatal(err)
				}
				if len(history.Items) == 0 {
					continue
				}

				err = models.ParseForUpdateWeekday(student, currentWeekday, history.Items[0].Text)
				if err != nil {
					b.Message("Неверный формат ввода. Попробуйте снова.")
					_, err = vk.MessagesSend(b.Params)
					if err != nil {
						log.Fatal(err)
					}
					continue
				}

				break

			}
			repo.UpdateWeekDayStudent(student, currentWeekday, db)
			b.Message("День недели был успешно изменён")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int()) // замените 123456 на любое уникальное число
			b.PeerID(obj.Message.PeerID)
			kb := bot.KeyboardStudents()
			b.Keyboard(kb)
			_, err = vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}

		case "Изменить время":
			// b.Message("Введите данные студента,у которого собираетесь изменить день недели в формате Имя Фамилия НовоеВремя СтароеВремя (например: Иван Иванов 19:00:00 21:00:00)")
			b.Message("TBD")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int()) // замените 123456 на любое уникальное число
			b.PeerID(obj.Message.PeerID)
			_, err := vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
			// var oldLessonTime string
			// student := &models.Student{}
			// student.TeacherID = obj.Message.PeerID
			// for {
			// 	// Retrieve the user's response
			// 	history, err := vk.MessagesGetHistory(api.Params{
			// 		"peer_id": obj.Message.PeerID,
			// 		"count":   1,
			// 	})
			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	if len(history.Items) == 0 {
			// 		continue
			// 	}

			// 	err = models.ParseForUpdateLessonTime(student, oldLessonTime, history.Items[0].Text)
			// 	if err != nil {
			// 		b.Message("Неверный формат ввода. Попробуйте снова.")
			// 		_, err = vk.MessagesSend(b.Params)
			// 		if err != nil {
			// 			log.Fatal(err)
			// 		}
			// 		continue
			// 	}

			// 	break

			// }
			// repo.UpdateLessonTimeStudent(student, oldLessonTime, db)
			// b.Message("Время урока было успешно изменено")
			// rand.Seed(time.Now().UnixNano())
			// b.RandomID(rand.Int()) // замените 123456 на любое уникальное число
			// b.PeerID(obj.Message.PeerID)
			// kb := bot.Keyboard_students()
			// b.Keyboard(kb)
			// _, err = vk.MessagesSend(b.Params)
			// if err != nil {
			// 	log.Fatal(err)
			// }

		case "Уведомления":

			b.Message("Выберите соответствующий пункт меню:")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int())
			b.PeerID(obj.Message.PeerID)
			kb := bot.KeyboardDown()
			b.Keyboard(kb)
			_, err = vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}

		case "Выключить уведомления":
			user := obj.Message.FromID
			if err != nil {
				log.Fatal(err)
			}
			var t models.Teacher
			t.TeacherID = user

			repo.TurnOffTeacherNotifyStatus(&t, db)

			b.Message("Уведомления успешно выключены!")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int())
			b.PeerID(obj.Message.PeerID)
			kb := bot.KeyboardDown()
			b.Keyboard(kb)
			_, err = vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		case "Включить уведомления":
			user := obj.Message.FromID
			if err != nil {
				log.Fatal(err)
			}
			var t models.Teacher
			t.TeacherID = user

			repo.TurnOnTeacherNotifyStatus(&t, db)

			b.Message("Уведомления успешно выключены!")
			rand.Seed(time.Now().UnixNano())
			b.RandomID(rand.Int())
			b.PeerID(obj.Message.PeerID)
			kb := bot.KeyboardDown()
			b.Keyboard(kb)
			_, err = vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		}

	})
	// Run Bots Long Poll
	log.Println("Start Long Poll")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}

	// Ожидаем, чтобы функция CheckLessonTime продолжала работу
	// Если вы не добавите эту строку, программа закончит свою работу
	// сразу после запуска функции CheckLessonTime
	select {}
}
