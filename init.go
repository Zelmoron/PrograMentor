package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// init функция автоматически вызывается при запуске пакета main
func init() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка при загрузке файла .env, продолжаем без него:", err)
	}

	// Для проверки можно вывести переменную
	// например, чтобы убедиться, что DB_CONNECTION загружена
	if env := os.Getenv("DB_CONNECTION"); env == "" {
		log.Println("Переменная DB_CONNECTION не установлена")
	} else {
		log.Println("DB_CONNECTION загружена успешно")
	}
}
