package service

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/megalypse/golang-fstresser/internal/application/service/profile/customprofile"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

var ctx context.Context
var cancelCtx context.CancelFunc
var testResourcesPath string

func init() {
	rawCtx := context.Background()
	ctx, cancelCtx = context.WithCancel(rawCtx)

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	testResourcesPath = basepath + "/../../../test_resources"
}

func TestLocalProfileLoading(t *testing.T) {
	assert := assert.New(t)
	loader := LocalProfileLoader{
		MakeRequestUsecase: MockMakeRequest{},
	}

	result := loader.LoadProfile(ctx, cancelCtx, testResourcesPath+"/test_profile.json")
	request := result[0].(customprofile.CustomStressProfile).Requests[0]
	config := result[0].(customprofile.CustomStressProfile).Config

	assert.Equal(1, request.Rate)
	assert.Equal("POST", request.Method)
	assert.Equal("https://mock-url.com", request.Url)

	assert.Equal(200, config.PeakRps)
	assert.Equal(time.Hour, config.ExecutionEndsAt.Duration)
	assert.Equal(time.Second*20, config.RampUpTime.Duration)
	assert.Equal(time.Hour, config.ExecutionEndsAt.Duration)
	assert.Equal(time.Minute*20, config.CustomLoads[0].StartsAt.Duration)
	assert.Equal(time.Minute*23, config.CustomLoads[0].EndsAt.Duration)
	assert.Equal("application/json", config.GlobalHeaders["Content-Type"])
	assert.Equal(300, config.CustomLoads[0].Rps)
}

type MockMakeRequest struct{}

func (MockMakeRequest) Request(context.Context, context.CancelFunc, *entity.Request, map[string]string) *entity.Response {
	return &entity.Response{}
}
