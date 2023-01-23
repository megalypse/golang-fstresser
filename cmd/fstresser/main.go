package main

import (
	"fmt"
	"time"

	"github.com/megalypse/golang-fstresser/internal/application/service"
	"github.com/megalypse/golang-fstresser/internal/application/service/profile"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

func main() {
	asp := profile.AnomalyStressProfile{
		RequestService: &service.HttpService{},
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

var accessToken string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJsYXN0TmFtZSI6IkRpYXMiLCJ1c2VyX25hbWUiOiJiYWxsaWd1aWVyaUByZWRjbWFpbC5jb20iLCJhdXRob3JpdGllcyI6WyJJVEsiLCJBU0NNIiwiTUxDIiwiTVBHIiwiTVZNIiwiTVZUTSIsIkFGTSIsIk1WUiIsIkxOTSIsIk1WVCIsIk1QTSIsIkhCSSIsIlRSTSIsIlVDIiwiTURPQyIsIk1EIiwiTURNIiwiTUUiLCJQRE0iLCJVTSIsIk1CTSIsIk1MT00iLCJBU0wiLCJNUEtHIiwiT0JSIiwiTUwiLCJNTSIsIlNVIiwiTVAiLCJNTFNUIiwiTVEiLCJNUiIsIk9CVyIsIk1WIiwiR1NXIiwiTUxDTSIsIlRVTSIsIk1FTSIsIk1FVE0iLCJJVEtVIiwiQVRLIiwiTUVUIiwiT0NSIl0sImNsaWVudF9pZCI6ImF1Y3Rpb25FbmdpbmUiLCJmaXJzdE5hbWUiOiJCcnVubyIsImF1ZCI6WyJBREMiXSwicmVmcmVzaF9zaWduYXR1cmUiOiIyMDIzLTAxLTIzVDE3OjQ3OjIzLjQxMjUxNzU3MloiLCJ1c2VyX2lkIjo5NTE5NTIxLCJpc19wcm9kIjpmYWxzZSwic2NvcGUiOlsicmVhZCIsIndyaXRlIl0sInBhcnR5X2lkIjoyNjEzMTU2NCwic2hvcnRfdGVybSI6MTY3NDUxODQwMCwiZXhwIjoxNjc0NDk5NjQzLCJqdGkiOiJSMFk1NG1Hd3hzTzRsbjRHVjNOX0Z4a0MtT28iLCJhY2NvdW50X2NvbmZpcm0iOmZhbHNlfQ.NfMXGHPv6GCAo7jyzFlvQ7UGSsc-lV5OMYWj8mlqetM2zCkdAFz98g1Hvjw-j25qIEAnsPymwa8O7nrDkowE12lyf2ys7IdLHDZL1jOR6PIwZ8NQVT4IyfruVeDwAa8S6CP5LRMB0CEuksbuB3PE-imSDSZjLSq595ro7K6r4iNCCO4ZSrVaio8rTqlAZTMNkY1AWU6KauVh5avibY2ROqJRWDi2FG_gf9DsraDHcyiYAtK4IONh4z_pDbdzd3Xo5vjINpZAcgOBk63vAxaxssa1DoKKccNyIlhh4QU2r4MA2wNz4o3DMB812y81-EIrxvkFj5RHR8YhrBwG-S056g"
