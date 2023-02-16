package customprofile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/megalypse/golang-fstresser/internal/domain/entity"
	"github.com/megalypse/golang-fstresser/internal/domain/usecase"
)

/*
CustomStressProfile implements the `StressProfile`, and is meant to be used to
translate json files to a valid stress profile to be used.
*/
type CustomStressProfile struct {
	// IsActive if false, the profile will be ignored and not executed
	IsActive bool

	// The HTTP requests to be used during the stressing
	Requests []CustomProfileRequest

	// The configuration for the stressing
	Config CustomProfileConfig

	// An instance of a struct that implements `usecase.MakeRequestUsecase`
	MakeRequestUsecase usecase.MakeRequestUsecase
}

func (csp CustomStressProfile) StartLoad(ctx context.Context, cancelCtx context.CancelFunc) {
	if !csp.IsActive {
		return
	}

	deployCustomProfileOrchestrator(ctx, cancelCtx, &csp)
}

type CustomProfileRequest struct {
	/*
		The rate this request will be used during the stressing.

		Example:
		Two `CustomProfileRequest` instances are created, each with `Rate` == 1.

		This means that if the stress have made 100k requests by the end of its execution,
		each `CustomProfileRequest` mentioned above will be responsible for 50k of the total amount.
	*/
	Rate int

	/*Method represents the HTTP method the request will be. (POST | GET | PUT | etc...)*/
	Method string

	Url     string
	Body    CustomBody
	Headers map[string]string
}

// ToEntity converts `CustomProfileRequest` to `entity.Request` to be used on `MakeLightweightRequest`.
func (cpr CustomProfileRequest) ToEntity() entity.Request {
	res, err := json.Marshal(cpr.Body.JsonBody)

	if err != nil {
		log.Fatal(err.Error())
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
	// BodyType represents if the request body will be a JSON or POSTFORM
	BodyType string

	// Fill this with any JSON structure if `BodyType` == "JSON"
	JsonBody any

	// Fill this with the bellow type if `BodyType` == "POSTFORM"
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
	// How many RPS the stresser should shot by the end of the ramp up time
	PeakRps int

	// How long the stresser should take before reaching the top RPS
	RampUpTime DurationInput

	// When the stresser should finish its execution
	ExecutionEndsAt DurationInput

	// The interval to be waited before each RPS increase
	RpsIncreaseInterval DurationInput

	// The error threshold to be tolerated before the software stops its execution
	ErrorThreshold string

	// The interval between each load, the default is 1s.
	LoadsInterval DurationInput
	CustomLoads   []CustomLoad
	GlobalHeaders map[string]string

	// This was created as a computed value to increase load tests flexibility
	RpsRampupPace float64
}

// bootstrap pre calculates values of interest for the stresser
func (cpc *CustomProfileConfig) bootstrap() {
	// How many times the RPS will increase
	rampUpsAmount := cpc.RampUpTime.Duration.Milliseconds() / cpc.RpsIncreaseInterval.Duration.Milliseconds()

	// How many requests will be added on each rampup
	cpc.RpsRampupPace = float64(cpc.PeakRps) / float64(rampUpsAmount)
}

// CustomLoad holds the necessary data for the end user be able to create custom loads during the stresser execution
type CustomLoad struct {
	/*
		Represents when the custom load will start.

		Example: "20:min".
		This means this custom load will start when the stresser runtime gets to 20 minutes.
	*/
	StartsAt DurationInput

	/*
		Represents when the custom load will end.

		Example: "21:min".
		This means this custom load will end when the stresser runtime gets to 21 minutes.
	*/
	EndsAt DurationInput

	// The RPS this custom load will generate
	Rps int

	// TODO: remove this unused field
	IsLogged bool

	// This field will hold a runtime computed value that represents `StartsAt` converted to Unix
	unixStartsAt int64

	// This field will hold a runtime computed value that represents `EndsAt` converted to Unix
	unixEndsAt int64
}

// calculateWindow calculates the unix points of when the custom load will start and finish.
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

/*
DurationInput purpose is to allow the enduser to configure the stresser timestamps using the following format:

{uint}:{sec|min|hour}
*/
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

	// If the duration metric == `time.Nanosecond` it means an invalid time metric was provided
	// and an error should be returned.
	if durationMetric == time.Nanosecond {
		return parseErr
	}

	di.Duration = durationMetric * time.Duration(contextlessNumber)
	return nil
}

/*
RateCounter is basically an interator. Its original purpose is to keep track of which
request should be used when making one during the stress.
*/
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

/*
CustomRequesterPayload is the payload to be sent on the custom requester channel.
*/
type CustomRequesterPayload struct {
	// The request to be sent.
	//
	// This field should be filled through the `RateCounter` iterator to keep the desired ratio.
	Request          *entity.Request
	CustomLoadConfig *CustomLoad
}

type DefaultRequesterPayload struct {
	// The request to be sent.
	//
	// This field should be filled through the `RateCounter` iterator to keep the desired ratio.
	Request *entity.Request
	Rps     int
}
