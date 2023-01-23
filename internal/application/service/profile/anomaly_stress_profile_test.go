package profile

// import (
// 	"context"
// 	"log"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/megalypse/golang-fstresser/internal/domain/entity"
// )

// var asp AnomalyStressProfile
// var testMu sync.Mutex
// var mockReqService MockRequestService

// func init() {
// 	asp = AnomalyStressProfile{
// 		RequestService: &mockReqService,
// 		Req: entity.Request{
// 			Method:    "GET",
// 			Url:       "https://mock-url.com",
// 			BytesBody: []byte{},
// 			MapBody:   make(map[string][]string),
// 			Headers:   make(map[string]string),
// 		},
// 		Config: Config{
// 			PeakRps:                 100,
// 			RampUpTime:              time.Second * 5,
// 			BeginAnomalyAfter:       time.Second * 6,
// 			AnomalyDuration:         time.Second * 3,
// 			AnomalyRps:              2,
// 			HoldPeakAfterAnomalyFor: time.Second,
// 		},
// 		State: State{
// 			IsDefaultFluxActive: true,
// 			IsAnomalyFluxActive: false,
// 		},
// 	}
// }

// func Test_runtime(t *testing.T) {
// 	asp.StartLoad()

// 	log.Println(mockReqService.GetTotalCount())

// 	if asp.State.EffectiveRps != asp.Config.PeakRps {
// 		t.Errorf("\nPeak RPS not reached. \nExpected: %d\nActual: %d", asp.Config.PeakRps, asp.State.EffectiveRps)
// 	}

// 	if asp.State.Runtime != asp.Config.ExpectedExecutionTime {
// 		t.Errorf("\nExpected runtime not met. \nExpected: %d\nActual: %d", asp.Config.ExpectedExecutionTime, asp.State.Runtime)

// 	}
// }

// type MockRequestService struct {
// 	getCount      int
// 	postCount     int
// 	postFormCount int
// }

// func (mrs MockRequestService) GetTotalCount() int {
// 	return mrs.getCount + mrs.postCount + mrs.postFormCount
// }

// func (mrs *MockRequestService) Get(closeCtx context.CancelFunc, req *entity.Request) *entity.Response {
// 	testMu.Lock()
// 	mrs.getCount++
// 	testMu.Unlock()

// 	return &entity.Response{}
// }

// func (mrs *MockRequestService) Post(closeCtx context.CancelFunc, req *entity.Request) *entity.Response {
// 	testMu.Lock()
// 	mrs.postCount++
// 	testMu.Unlock()

// 	return &entity.Response{}
// }

// func (mrs *MockRequestService) PostForm(closeCtx context.CancelFunc, req *entity.Request) *entity.Response {
// 	testMu.Lock()
// 	mrs.postFormCount++
// 	testMu.Unlock()

// 	return &entity.Response{}
// }
