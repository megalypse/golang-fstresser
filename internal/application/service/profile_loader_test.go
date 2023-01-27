package service

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

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
	loader := LocalProfileLoader{}

	result := loader.LoadProfile(cancelCtx, testResourcesPath+"/test_profile.json")
	request := result[0].Requests[0]
	config := result[0].Config

	assert.Equal(1, request.Rate)
	assert.Equal("POST", request.Method)
	assert.Equal("https://mock-url.com", request.Url)
	assert.Equal("application/json", request.Headers["Content-Type"])

	assert.Equal(200, config.PeakRps)
	assert.Equal(time.Hour, config.EndLoadAt.Duration)
	assert.Equal(time.Second*20, config.RampUpTime.Duration)
	assert.Equal(time.Hour, config.EndLoadAt.Duration)
	assert.Equal(time.Minute*20, config.CustomLoads[0].StartsAt.Duration)
	assert.Equal(time.Minute*23, config.CustomLoads[0].EndsAt.Duration)
	assert.Equal(300, config.CustomLoads[0].Rps)
}
