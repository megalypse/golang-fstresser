package main

import (
	"fmt"
	"time"

	"github.com/megalypse/golang-fstresser/internal/application/service/profile"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

func main() {
	asp := profile.AnomalyStressProfile{
		Req: entity.Request{
			Method:    "POST",
			Url:       "https://service.ci.aucn.io/resi/listing/v1/listings/minimal",
			BytesBody: []byte(`[1]`),
			Headers: map[string]string{
				"Content-Type":    "application/json",
				"Authorization":   fmt.Sprintf("Bearer %v", accessToken),
				"Accept-Encoding": "gzip, deflate, br",
				"Connection":      "keep-alive",
			},
		},
		Config: profile.Config{
			PeakRps:                 200,
			RampUpTime:              time.Minute * 20,
			BeginAnomalyAfter:       time.Minute * 20,
			AnomalyDuration:         time.Minute * 2,
			AnomalyRps:              300,
			HoldPeakAfterAnomalyFor: time.Minute * 5,
		},
	}

	asp.StartLoad()
}

var accessToken string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJsYXN0TmFtZSI6IkRpYXMiLCJ1c2VyX25hbWUiOiJiYWxsaWd1aWVyaUByZWRjbWFpbC5jb20iLCJhdXRob3JpdGllcyI6WyJJVEsiLCJBU0NNIiwiTUxDIiwiTVBHIiwiTVZNIiwiTVZUTSIsIkFGTSIsIk1WUiIsIkxOTSIsIk1WVCIsIk1QTSIsIkhCSSIsIlRSTSIsIlVDIiwiTURPQyIsIk1EIiwiTURNIiwiTUUiLCJQRE0iLCJVTSIsIk1CTSIsIk1MT00iLCJBU0wiLCJNUEtHIiwiT0JSIiwiTUwiLCJNTSIsIlNVIiwiTVAiLCJNTFNUIiwiTVEiLCJNUiIsIk9CVyIsIk1WIiwiR1NXIiwiTUxDTSIsIlRVTSIsIk1FTSIsIk1FVE0iLCJJVEtVIiwiQVRLIiwiTUVUIiwiT0NSIl0sImNsaWVudF9pZCI6ImF1Y3Rpb25FbmdpbmUiLCJmaXJzdE5hbWUiOiJCcnVubyIsImF1ZCI6WyJBREMiXSwicmVmcmVzaF9zaWduYXR1cmUiOiIyMDIzLTAxLTIzVDIyOjI3OjIwLjQzMjQ4MjE5N1oiLCJ1c2VyX2lkIjo5NTE5NTIxLCJpc19wcm9kIjpmYWxzZSwic2NvcGUiOlsicmVhZCIsIndyaXRlIl0sInBhcnR5X2lkIjoyNjEzMTU2NCwic2hvcnRfdGVybSI6MTY3NDUxODQwMCwiZXhwIjoxNjc0NTE2NDQwLCJqdGkiOiJjaWxFbzQ1RkJXdVZCT1Rkd2VQQlRoc09JcTAiLCJhY2NvdW50X2NvbmZpcm0iOmZhbHNlfQ.jUMttaBoISTCURFzKpWuceNIK2bvMP-eKnHsOMefp6AueYfPkM50biDh-5fiwvMb2L5PpwF7nCw4NmUK7CJ8Tq7Gy_Sx1BVWktjvyUFDBYxiSjdDFcoTzDe_RVne8xdRP2Ty1-rOqvOsgDzQZcuTFtgxDD6OFkSvJ64T5YxmB6BQXug9sIu0R3eh7bmeS3OBg93mecR-3fiGKyHiqLkFzreiUh9q61B5tk0sSvU9v80Cs7oiLmtUOuK75rrLYQoPxFTHLZrOwWYdTuR3Qm3gmynSPejarcBYtRhvtrZFCchF-P9_-cZWUDO42v37hKNg1urTJzBZkDn3rHzeFFpnNg"
