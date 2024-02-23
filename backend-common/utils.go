package common

import "strconv"

var (
	difficulties = map[string]bool{
		DIFFICULTY_EASY:   true,
		DIFFICULTY_MEDIUM: true,
		DIFFICULTY_HARD:   true,
	}
	userStatuses = map[string]bool{
		USER_STATUS_IDLE:           true,
		USER_STATUS_PENDING_EASY:   true,
		USER_STATUS_PENDING_MEDIUM: true,
		USER_STATUS_PENDING_HARD:   true,
		USER_STATUS_MATCHED:        true,
		USER_STATUS_LEAVING:        true,
	}
	userPendingStatuses = map[string]bool{
		USER_STATUS_PENDING_EASY:   true,
		USER_STATUS_PENDING_MEDIUM: true,
		USER_STATUS_PENDING_HARD:   true,
	}
)

func ResponseBodyWithMessage(msg string, dataName string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"message": msg,
		dataName:  data,
	}
}

func ResponseBodyWithError(errorMsg string) map[string]interface{} {
	return map[string]interface{}{
		"error": errorMsg,
	}
}

func DifficultySet() map[string]bool {
	return difficulties
}

func UserStatusSet() map[string]bool {
	return userStatuses
}

func PendingStatusSet() map[string]bool {
	return userPendingStatuses
}

func StringToUint(s string) (uint, error) {
	num, err := strconv.ParseUint(s, 10, 64)
	return uint(num), err
}

func StringToInt(s string) (int, error) {
	num, err := strconv.ParseInt(s, 10, 64)
	return int(num), err
}
