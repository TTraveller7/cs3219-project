package match

import (
	"context"
	"time"

	"github.com/backend-common/common"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const TIME_FORMAT = time.RFC3339Nano

var (
	log         *common.Logger
	mainContext context.Context
)

func getNewUuid() string {
	// Let's beg for probability
	// TODO: check error and duplicate
	id, _ := uuid.NewRandom()
	return id.String()
}

func createMatchIdCookie(matchId string) string {
	// Generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": matchId,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Jwt expires after 30 days
		"iat": time.Now().Unix(),                          // issue time
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte(common.SECRET)) // SECRET here should be a env variable e.g. []byte(os.Getenv("SECRET"))
	return tokenString
}

func Setup(ctx context.Context, l *common.Logger) {
	mainContext = ctx
	log = l
}

func ConnectCacheAndDb() {
	pool = common.CreateDbPool(getPostgresUrl(), log)
	userCache = common.CreateCache(getUserRedisUrl(), log)
	matchCache = common.CreateCache(getMatchRedisUrl(), log)
}

func CloseCacheAndDb() {
	common.CloseDbpool(pool)
}
