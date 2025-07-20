package redis

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	"codematic/pkg/environment"
	// "codematic/storage/redis/mock"
)

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
	redisStore KvStore
}

func (s *Suite) SetupSuite() {
	env := &environment.Env{}
	env.HelperForMocking(map[string]string{
		"REDIS_SERVER_ADDRESS": "rediss://username:password@host:6379",
	})
	dummyLog := zerolog.Nop()

	s.redisStore = *NewRedis(env, dummyLog, "rediss://username:password@host:6379")
}

func (s *Suite) AfterTest(_, _ string) {
}

// func (s *Suite) Test_GetValue() {
// 	ctrl := gomock.NewController(s.T())
// 	defer ctrl.Finish()

// 	foundKey := "FOUND_KEY"
// 	errorKey := "ERROR_KEY"
// 	m := mock.NewMockKvStore(ctrl)
// 	m.EXPECT().GetValue(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, key string, result *string) error {
// 		if strings.EqualFold(key, foundKey) {
// 			*result = foundKey
// 			return nil
// 		} else if strings.EqualFold(key, errorKey) {
// 			result = nil
// 			return ErrFailedToRetrieveValue
// 		}

// 		return nil // dont error but return nil for error and value retrieved
// 	}).AnyTimes()

// 	kv := KvStore(m)
// 	var fVal string
// 	fErr := kv.GetValue(context.Background(), foundKey, &fVal)
// 	require.NoError(s.T(), fErr)
// 	require.Equal(s.T(), fVal, foundKey)
// 	var eVal *string
// 	eErr := kv.GetValue(context.Background(), errorKey, eVal)
// 	require.Error(s.T(), eErr)
// 	require.Nil(s.T(), eVal)
// 	require.EqualError(s.T(), eErr, ErrFailedToRetrieveValue.Error())
// }
