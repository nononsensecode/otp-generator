package otp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nononsensecode/otp-generator"
	"github.com/stretchr/testify/assert"
)

func giveTime(asString string) time.Time {
	t, _ := time.Parse("02/01/2006 15:04:05", asString)
	return t
}

func Test_New(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		timeProvider      otp.CurrentTimeProvider
		duration          time.Duration
		generator         otp.OtpGenerator
		length            int
		maxResendAttempts int
		wantedErr         error
		wantedData        otp.OtpData
	}{
		"new otp creation success": {
			timeProvider: func() time.Time {
				t := giveTime("25/08/2022 08:25:00")
				return t
			},
			duration: 3 * time.Minute,
			generator: func(length int) (string, error) {
				return "12345", nil
			},
			length:            5,
			maxResendAttempts: 3,
			wantedErr:         nil,
			wantedData: otp.OtpData{
				Otp:               "12345",
				CreatedOn:         giveTime("25/08/2022 08:25:00"),
				ExpiryDuration:    3 * time.Minute,
				ResendAttempts:    1,
				MaxResendAttempts: 3,
				Stale:             false,
			},
		},

		"new otp creation return otp generator error": {
			timeProvider: func() time.Time {
				t := giveTime("25/08/2022 08:25:00")
				return t
			},
			duration: 3 * time.Minute,
			generator: func(length int) (string, error) {
				return "", fmt.Errorf("unknown generator error")
			},
			length:            5,
			maxResendAttempts: 3,
			wantedErr:         fmt.Errorf("unknown generator error"),
			wantedData:        otp.OtpData{},
		},
	}

	for tname, tcase := range tests {
		t.Run(tname, func(t *testing.T) {
			o, err := otp.New(tcase.timeProvider, tcase.duration, tcase.generator,
				tcase.length, tcase.maxResendAttempts)
			assert.Equal(t, tcase.wantedErr, err)
			assert.Equal(t, tcase.wantedData, o.OtpData())
		})
	}
}
