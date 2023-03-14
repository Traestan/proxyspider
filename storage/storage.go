package storage

// Storage
type Storage interface {
	CheckStorage() error            // Проверяет наличие storage
	WriteStorage(text string) error // Пишем в storage
	Stat() error                    // Статистика по добавленниям
}
