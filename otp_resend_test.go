package otp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nononsensecode/otp-generator"
	"github.com/stretchr/testify/assert"
)

func Test_Resendable(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		createdOn         time.Time
		otpVal            string
		resendAttempts    int
		maxResendAttempts int
		duration          time.Duration
		isStale           bool
		timeProvider      otp.CurrentTimeProvider
		generator         otp.OtpGenerator
		length            int
		wantedErr         error
		wantedData        otp.OtpData
	}{
		"resendable success": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			duration:          3 * time.Minute,
			resendAttempts:    1,
			maxResendAttempts: 3,
			isStale:           false,
			timeProvider: func() time.Time {
				t := giveTime("25/08/2022 08:26:00")
				return t
			},
			generator: func(length int) (string, error) { return "34567", nil },
			wantedErr: nil,
			wantedData: otp.OtpData{
				Otp:               "34567",
				CreatedOn:         giveTime("25/08/2022 08:25:00"),
				ExpiryDuration:    3 * time.Minute,
				ResendAttempts:    2,
				MaxResendAttempts: 3,
				Stale:             false,
			},
		},

		"resendable returns stale error": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			duration:          3 * time.Minute,
			resendAttempts:    1,
			maxResendAttempts: 3,
			isStale:           true,
			timeProvider: func() time.Time {
				t := giveTime("25/08/2022 08:29:00")
				return t
			},
			generator: func(length int) (string, error) { return "34567", nil },
			wantedErr: otp.ErrStale,
			wantedData: otp.OtpData{
				Otp:               "12345",
				CreatedOn:         giveTime("25/08/2022 08:25:00"),
				ExpiryDuration:    3 * time.Minute,
				ResendAttempts:    1,
				MaxResendAttempts: 3,
				Stale:             true,
			},
		},

		"resendable returns expiry error": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			duration:          3 * time.Minute,
			resendAttempts:    1,
			maxResendAttempts: 3,
			isStale:           false,
			timeProvider: func() time.Time {
				t := giveTime("25/08/2022 08:29:00")
				return t
			},
			generator: func(length int) (string, error) { return "34567", nil },
			wantedErr: otp.ErrExpiry,
			wantedData: otp.OtpData{
				Otp:               "12345",
				CreatedOn:         giveTime("25/08/2022 08:25:00"),
				ExpiryDuration:    3 * time.Minute,
				ResendAttempts:    1,
				MaxResendAttempts: 3,
				Stale:             true,
			},
		},

		"resendable returns resend exceeds error": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			duration:          3 * time.Minute,
			resendAttempts:    3,
			maxResendAttempts: 3,
			isStale:           false,
			timeProvider: func() time.Time {
				t := giveTime("25/08/2022 08:28:00")
				return t
			},
			generator: func(length int) (string, error) { return "34567", nil },
			wantedErr: otp.ErrResendExcceded,
			wantedData: otp.OtpData{
				Otp:               "12345",
				CreatedOn:         giveTime("25/08/2022 08:25:00"),
				ExpiryDuration:    3 * time.Minute,
				ResendAttempts:    3,
				MaxResendAttempts: 3,
				Stale:             true,
			},
		},

		"resendable returns otp generation error": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			duration:          3 * time.Minute,
			resendAttempts:    1,
			maxResendAttempts: 3,
			isStale:           false,
			timeProvider: func() time.Time {
				t := giveTime("25/08/2022 08:26:00")
				return t
			},
			generator: func(length int) (string, error) { return "", fmt.Errorf("i can't create otp") },
			wantedErr: fmt.Errorf("i can't create otp"),
			wantedData: otp.OtpData{
				Otp:               "12345",
				CreatedOn:         giveTime("25/08/2022 08:25:00"),
				ExpiryDuration:    3 * time.Minute,
				ResendAttempts:    1,
				MaxResendAttempts: 3,
				Stale:             false,
			},
		},
	}

	for tname, tcase := range tests {
		t.Run(tname, func(t *testing.T) {
			o := otp.FromPersistence(tcase.otpVal, tcase.createdOn,
				tcase.duration, tcase.resendAttempts, tcase.maxResendAttempts, tcase.isStale)
			err := o.Resendable(tcase.timeProvider, tcase.generator, tcase.length)
			assert.Equal(t, tcase.wantedErr, err)
			assert.Equal(t, tcase.wantedData, o.OtpData())
		})
	}
}
