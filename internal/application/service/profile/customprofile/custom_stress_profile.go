package customprofile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
	"github.com/megalypse/golang-fstresser/internal/domain/usecase"
)

type CustomStressProfile struct {
	IsActive           bool
	Requests           []CustomProfileRequest
	Config             CustomProfileConfig
	MakeRequestUsecase usecase.MakeRequestUsecase
}

func (csp CustomStressProfile) StartLoad(ctx context.Context, cancelCtx context.CancelFunc) {
	if !csp.IsActive {
		return
	}

	deployCustomProfileOrchestrator(ctx, cancelCtx, &csp)
}

type CustomProfileRequest struct {
	Rate    int
	Method  string
	Url     string
	Body    CustomBody
	Headers map[string]string
}

func (cpr CustomProfileRequest) ToEntity() entity.Request {
	res, err := json.Marshal(cpr.Body.JsonBody)

	if err != nil {
		common.GetLogger().Log(err.Error())
	}

	return entity.Request{
		Method:    cpr.Method,
		Url:       cpr.Url,
		BytesBody: res,
		MapBody:   cpr.Body.FormBody,
		Headers:   cpr.Headers,
	}
}

type CustomBody struct {
	BodyType string
	JsonBody any
	FormBody map[string][]string
}

func (cb *CustomBody) UnmarshalJSON(rawJson []byte) error {
	holder := new(struct {
		BodyType string
		Body     any
	})

	json.Unmarshal(rawJson, &holder)

	switch holder.BodyType {
	case "JSON":
		cb.JsonBody = holder.Body
	case "POSTFORM":
		finalHolder := make(map[string][]string)
		json.Unmarshal(holder.Body.([]byte), &finalHolder)
		cb.FormBody = finalHolder
	default:
		return fmt.Errorf("body type not supported: %q", holder.BodyType)
	}

	return nil
}

type CustomProfileConfig struct {
	PeakRps             int
	RampUpTime          DurationInput
	CustomLoads         []CustomLoad
	ExecutionEndsAt     DurationInput
	RpsIncreaseInterval DurationInput
	ErrorThreshold      string
	LoadsInterval       DurationInput
	GlobalHeaders       map[string]string

	// This was created as a computed value to increase load tests flexibility
	RpsRampupPace float64
}

func (cpc *CustomProfileConfig) bootstrap() {
	rampUpsAmount := cpc.RampUpTime.Duration.Milliseconds() / cpc.RpsIncreaseInterval.Duration.Milliseconds()
	cpc.RpsRampupPace = float64(cpc.PeakRps) / float64(rampUpsAmount)
}

type CustomLoad struct {
	StartsAt DurationInput
	EndsAt   DurationInput
	Rps      int
	IsLogged bool

	unixStartsAt int64
	unixEndsAt   int64
}

func (cl *CustomLoad) calculateWindow(loadStartPoint time.Time) {
	cl.unixStartsAt = loadStartPoint.Add(cl.StartsAt.Duration).Unix()
	cl.unixEndsAt = loadStartPoint.Add(cl.EndsAt.Duration).Unix()
}

func (cl CustomLoad) GetStartPoint() int64 {
	return cl.unixStartsAt
}

func (cl CustomLoad) GetEndPoint() int64 {
	return cl.unixEndsAt
}

type DurationInput struct {
	Duration time.Duration
}

func (di *DurationInput) UnmarshalJSON(rawJson []byte) error {
	parseErr := errors.New("duration input must be on {int}:{sec|min|hour} format")
	rawDuration := string(rawJson)
	result, _ := strconv.Unquote(rawDuration)

	raw := strings.Split(result, ":")

	if len(raw) != 2 {
		return parseErr
	}

	rawNumber := raw[0]
	rawMetric := raw[1]

	contextlessNumber, err := strconv.Atoi(rawNumber)
	if err != nil {
		return err
	}

	if rawMetric == "" {
		return parseErr
	}

	durationMetric := func() time.Duration {
		switch rawMetric {
		case "sec":
			return time.Second
		case "min":
			return time.Minute
		case "hour":
			return time.Hour
		default:
			return time.Nanosecond
		}
	}()

	if durationMetric == time.Nanosecond {
		return parseErr
	}

	di.Duration = durationMetric * time.Duration(contextlessNumber)
	return nil
}

type RateCounter struct {
	maxValue     int
	currentValue int
}

func (rc *RateCounter) Next() int {
	if rc.currentValue+1 > rc.maxValue {
		rc.currentValue = 0
	} else {
		rc.currentValue++
	}

	return rc.currentValue
}

type CustomRequesterPayload struct {
	Request          *entity.Request
	CustomLoadConfig *CustomLoad
}

type DefaultRequesterPayload struct {
	Request *entity.Request
	Rps     int
}
