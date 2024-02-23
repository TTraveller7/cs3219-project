package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/backend-common/common"
	"github.com/go-redis/redis/v8"
)

var (
	questionCache *redis.Client
	matchCache    *redis.Client
)

// Question Id list
// Key: 1-10, Value: qid

func initQidListAndPointer(ctx context.Context, matchId string) (bool, error) {

	difficulty, exist := getDifficulty(ctx, matchId)
	if !exist {
		log.Message("Fail to retrieve difficulty")
		return false, nil
	}

	list := questionIdList(difficulty)
	saveQuestionIdList(ctx, matchId, list)
	savePointer(ctx, matchId, int(INITIAL_POINTER))
	return true, nil
}

func saveQuestionIdList(ctx context.Context, matchId string, lst []uint) {
	var qidMap map[string]string = map[string]string{}
	for i, val := range lst {
		qidMap[strconv.FormatInt(int64(i+1), 10)] = strconv.FormatUint(uint64(val), 10)
	}
	log.Message("Qidmap created: " + fmt.Sprint(qidMap))

	if err := questionCache.HSet(ctx, QuestionListKey(matchId), qidMap).Err(); err != nil {
		log.Error("", err)
	}
}

func getQidByPointer(ctx context.Context, matchId string, ptr int) (uint, error) {
	qidMap, err := questionCache.HGetAll(ctx, QuestionListKey(matchId)).Result()
	if err != nil {
		return 0, err
	}

	if qidStr, ok := qidMap[strconv.FormatUint(uint64(ptr), 10)]; !ok {
		return 0, errors.New("pointer is out of range")
	} else {
		return common.StringToUint(qidStr)
	}
}

// Is the list of question ids empty?
func isQuestionIdListEmpty(ctx context.Context, matchId string) bool {
	if !existPointer(ctx, matchId) {
		log.Message("pointer does not exist")
		return true
	}
	pointer, _ := getPointer(ctx, matchId)

	log.Message(fmt.Sprintf("Pointer value is: %v", pointer))

	count, err := questionCache.Exists(ctx, QuestionListKey(matchId)).Result()
	log.Message(fmt.Sprintf("Count value is: %v", count))

	if err != nil {
		log.Error("", err)
		return true
	} else if count < 1 {
		return true
	}

	qidMapLength, err := questionCache.HLen(ctx, QuestionListKey(matchId)).Result()
	log.Message(fmt.Sprintf("QidMapLength value is: %v", qidMapLength))
	if err != nil {
		log.Error("", err)
		return true
	}

	return int(qidMapLength) < pointer
}

// Pointer

func savePointer(ctx context.Context, matchId string, pointer int) {
	// Save pointer
	questionCache.Set(ctx, PointerKey(matchId), pointer, 0)
}

func getPointer(ctx context.Context, matchId string) (int, error) {
	if pointerStr, err := questionCache.Get(ctx, PointerKey(matchId)).Result(); err != nil {
		return 0, err
	} else {
		return common.StringToInt(pointerStr)
	}
}

func existPointer(ctx context.Context, matchId string) bool {
	count, err := questionCache.Exists(ctx, PointerKey(matchId)).Result()
	if err != nil {
		log.Error("", err)
	}

	return (err == nil) && (count > 0)
}

// Lock

func tryLock(ctx context.Context, matchId string, retries int) bool {
	isLocked := false
	count := 0
	for !isLocked && count < retries {
		isLocked = lock(ctx, matchId)
		if !isLocked {
			time.Sleep(500 * time.Millisecond)
		}
		count++
	}

	return isLocked
}

func lock(ctx context.Context, matchId string) bool {
	res, err := questionCache.SetNX(ctx, LockKey(matchId), 1, 0).Result()
	if err != nil {
		log.Error("Fail to lock "+matchId+":", err)
		return false
	}

	return res
}

func unlock(ctx context.Context, matchId string) {
	if err := questionCache.Del(ctx, LockKey(matchId)).Err(); err != nil {
		log.Error("Fail to unlock "+matchId+":", err)
	}
}

// Difficulty

func getDifficulty(ctx context.Context, matchId string) (string, bool) {
	matchMap, err := matchCache.HGetAll(ctx, matchId).Result()
	if err != nil {
		log.Error("", err)
		return "", false
	} else {
		difficulty := matchMap[MATCH_DIFFICULTY_KEY]
		if len(difficulty) == 0 {
			log.Error("", errors.New("fail to retrieve difficulty: difficulty is empty for matchId "+matchId))
			return "", false
		}
		return difficulty, true
	}
}
