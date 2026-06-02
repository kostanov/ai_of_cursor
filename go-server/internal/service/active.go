package service

import "sync"

// ActiveUsers — потокобезопасный список активных id.
type ActiveUsers struct {
	mu      sync.Mutex
	users   []int
	maxSize int
}

// NewActiveUsers создаёт сервис с лимитом длины списка.
func NewActiveUsers(maxSize int) *ActiveUsers {
	return &ActiveUsers{maxSize: maxSize}
}

// Add добавляет id и возвращает копию текущего списка.
func (a *ActiveUsers) Add(userID int) []int {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.users = append(a.users, userID)
	if len(a.users) > a.maxSize {
		overflow := len(a.users) - a.maxSize
		a.users = a.users[overflow:]
	}
	return append([]int(nil), a.users...)
}
