package match

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

var (
	userCache  *redis.Client
	matchCache *redis.Client
)

func saveUser(ctx context.Context, u *User) error {
	return userCache.HSet(ctx, u.Username, u.getValuesMap()).Err()
}

func isUserInCache(ctx context.Context, username string) bool {
	userCount, err := userCache.Exists(ctx, username).Result()
	if err != nil {
		log.Error("Fail to check if user exists", err)
		return false
	}
	return userCount > 0
}

func getUser(ctx context.Context, username string) (*User, error) {
	if !isUserInCache(ctx, username) {
		return &User{}, errors.New("User with name " + username + "does not exist")
	}

	userMap, err := userCache.HGetAll(ctx, username).Result()
	if err != nil {
		log.Error("Fail to get user", err)
		return &User{}, err
	}

	u, err := valuesMapToUser(username, userMap)
	if err != nil {
		log.Error("Fail to restore user from values map", err)
		return &User{}, err
	}

	return u, err
}

func getUserField(ctx context.Context, username string, field string) (string, error) {
	return userCache.HGet(ctx, username, field).Result()
}

func setUser(ctx context.Context, username string, values ...string) error {
	return userCache.HSet(ctx, username, values).Err()
}

func updateUserStatus(ctx context.Context, username string, matchId string, status string) error {
	userMap, err := matchCache.HGetAll(ctx, username).Result()
	if err != nil {
		return err
	}
	userMap["MatchId"] = matchId
	userMap["Status"] = status
	return userCache.HSet(ctx, username, userMap).Err()
}

func saveMatch(ctx context.Context, m *Match) error {
	return matchCache.HSet(ctx, m.MatchId, m.getValuesMap()).Err()
}

func getMatch(ctx context.Context, matchId string) (*Match, error) {
	if !isMatchInCache(ctx, matchId) {
		return &Match{}, errors.New("Match with id " + matchId + "does not exist")
	}

	matchMap, err := matchCache.HGetAll(ctx, matchId).Result()
	if err != nil {
		return &Match{}, err
	}

	match, err := valueMapToMatch(matchId, matchMap)
	if err != nil {
		return &Match{}, err
	}

	return match, nil
}

func isMatchInCache(ctx context.Context, matchId string) bool {
	count, err := matchCache.Exists(ctx, matchId).Result()
	if err != nil {
		log.Error("Fail to check if match exists:", err)
		return false
	}
	return count > 0
}

func isMatchRequestValid(ctx context.Context, matchRequest MatchRequest) (bool, error) {
	userLatestMatchRequestId, err := getUserField(ctx, matchRequest.Username, "LastMatchRequestId")
	if err != nil {
		return false, err
	}

	if user, err := getUser(ctx, matchRequest.Username); err != nil {
		return false, err
	} else {
		return user.isPending() && matchRequest.MatchRequestId == userLatestMatchRequestId, nil
	}
}
