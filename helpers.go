package main

import (
	"sync"
)

// Mutex to ensure safe access to shared state
	var mu sync.Mutex

	// Track selected plans for each user
	var userStates = make(map[int64]*UserState)

	var getUserState = func(userID int64) *UserState {
		mu.Lock()
		defer mu.Unlock()
		if state, exists := userStates[userID]; exists {
			return state
		}
		newState := &UserState{HasSelectedPlan: false}
		userStates[userID] = newState
		return newState
	}

