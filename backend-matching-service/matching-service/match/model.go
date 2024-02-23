package match

import (
	"fmt"
	"strconv"
	"time"

	"github.com/backend-common/common"
)

type MatchRequest struct {
	MatchRequestId string    `json:"matchRequestId"`
	Username       string    `json:"username" binding:"required"`
	Difficulty     string    `json:"difficulty" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
}

type Match struct {
	MatchId    string    `json:"matchId"`
	Difficulty string    `json:"difficulty"`
	UsernameA  string    `json:"usernameA"`
	UsernameB  string    `json:"usernameB"`
	IsEnded    bool      `json:"isEnded"`
	CreatedAt  time.Time `json:"createdAt"`
	EndedAt    time.Time `json:"endedAt"`
}

type User struct {
	Username           string
	Status             string
	LastMatchRequestId string
	MatchId            string
}

func getPendingStatus(difficulty string) string {
	return "pending_" + difficulty
}

func (u *User) getValuesMap() map[string]string {
	return map[string]string{
		"Status":             u.Status,
		"LastMatchRequestId": u.LastMatchRequestId,
		"MatchId":            u.MatchId,
	}
}

func (u *User) isPending() bool {
	return common.PendingStatusSet()[u.Status]
}

func valuesMapToUser(username string, valuesMap map[string]string) (*User, error) {
	return &User{
		Username:           username,
		Status:             valuesMap["Status"],
		LastMatchRequestId: valuesMap["LastMatchRequestId"],
		MatchId:            valuesMap["MatchId"],
	}, nil
}

// Match fucntions

func newMatch(difficulty string, usernameA string, usernameB string) *Match {
	return &Match{
		MatchId:    getNewUuid(),
		Difficulty: difficulty,
		UsernameA:  usernameA,
		UsernameB:  usernameB,
		IsEnded:    false,
		CreatedAt:  time.Now(),
		EndedAt:    time.Time{},
	}
}

func (m *Match) toString() string {
	return fmt.Sprintf("%#v", *m)
}

func (m *Match) getValuesMap() map[string]string {
	return map[string]string{
		"Difficulty": m.Difficulty,
		"UsernameA":  m.UsernameA,
		"UsernameB":  m.UsernameB,
		"IsEnded":    strconv.FormatBool(m.IsEnded),
		"CreatedAt":  m.CreatedAt.Format(TIME_FORMAT),
		"EndedAt":    m.EndedAt.Format(TIME_FORMAT),
	}
}

func valueMapToMatch(matchId string, valueMap map[string]string) (*Match, error) {
	createdAt, err := time.Parse(TIME_FORMAT, valueMap["CreatedAt"])
	if err != nil {
		return &Match{}, err
	}

	endedAt, err := time.Parse(TIME_FORMAT, valueMap["EndedAt"])
	if err != nil {
		return &Match{}, err
	}

	isEnded, err := strconv.ParseBool(valueMap["IsEnded"])
	if err != nil {
		return &Match{}, err
	}

	return &Match{
		MatchId:    matchId,
		Difficulty: valueMap["Difficulty"],
		UsernameA:  valueMap["UsernameA"],
		UsernameB:  valueMap["UsernameB"],
		IsEnded:    isEnded,
		CreatedAt:  createdAt,
		EndedAt:    endedAt,
	}, nil
}
