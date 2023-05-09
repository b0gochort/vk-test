package bot

import "github.com/SevereCloud/vksdk/v2/object"

func MainKeyBoard() *object.MessagesKeyboard {
	kb := object.NewMessagesKeyboardInline()
	kb.AddRow()
	kb.AddTextButton("Ученики", "student", "primary")
	kb.AddTextButton("Уроки", "lessons", "primary")
	kb.AddTextButton("Изменить", "change", "primary")
	kb.AddTextButton("Уведомления", "down", "negative")
	return kb
}

func KeyboardStudents() *object.MessagesKeyboard {
	kb := object.NewMessagesKeyboardInline()
	kb.AddRow()
	kb.AddTextButton("Добавить ученика", "add_student", "primary")
	kb.AddTextButton("Удалить ученика", "delete_student", "primary")
	kb.AddRow()
	kb.AddTextButton("Назад", "main_keyboard", "default")
	return kb
}

func KeyboardLessons() *object.MessagesKeyboard {
	kb := object.NewMessagesKeyboardInline()
	kb.AddRow()
	kb.AddTextButton("Расписание на сегодня", "today_lessons", "primary")
	kb.AddTextButton("Расписание на неделю", "week_lessons", "primary")
	kb.AddRow()
	kb.AddTextButton("Назад", "main_keyboard", "default")
	return kb
}

func KeyboardChange() *object.MessagesKeyboard {
	kb := object.NewMessagesKeyboardInline()
	kb.AddRow()
	kb.AddTextButton("Изменить день", "change_weekday_of_student", "primary")
	kb.AddTextButton("Изменить время", "change_lesson_time_of_student", "primary")
	kb.AddRow()
	kb.AddTextButton("Назад", "main_keyboard", "default")
	return kb
}

func KeyboardDown() *object.MessagesKeyboard {
	kb := object.NewMessagesKeyboardInline()
	kb.AddRow()
	kb.AddTextButton("Выключить уведомления", "turn_off_notify", "primary")
	kb.AddTextButton("Включить уведомления", "turn_on_notify", "primary")
	// kb.AddTextButton("Отключить бота", "delete_user", "primary") TBD
	kb.AddRow()
	kb.AddTextButton("Назад", "main_keyboard", "default")
	return kb
}
