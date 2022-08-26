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

func newOtpData(otpVal string, createdOn string, duration time.Duration,
	resendAttempts, maxResendAttempts int, isStale bool) otp.OtpData {
	return otp.OtpData{
		Otp:               otpVal,
		CreatedOn:         giveTime(createdOn),
		ExpiryDuration:    duration,
		ResendAttempts:    resendAttempts,
		MaxResendAttempts: maxResendAttempts,
		Stale:             isStale,
	}
}

func timeProvider(t string) func() time.Time {
	return func() time.Time {
		return giveTime(t)
	}
}

func generator(otpVal string, err error) func(int) (string, error) {
	return func(length int) (string, error) {
		return otpVal, err
	}
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
			timeProvider:      timeProvider("25/08/2022 08:25:00"),
			duration:          3 * time.Minute,
			generator:         generator("12345", nil),
			length:            5,
			maxResendAttempts: 3,
			wantedErr:         nil,
			wantedData:        newOtpData("12345", "25/08/2022 08:25:00", 3*time.Minute, 1, 3, false),
		},

		"new otp creation return otp generator error": {
			timeProvider:      timeProvider("25/08/2022 08:25:00"),
			duration:          3 * time.Minute,
			generator:         generator("", fmt.Errorf("unknown generator error")),
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
