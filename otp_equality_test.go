package otp_test

import (
	"testing"
	"time"

	"github.com/nononsensecode/otp-generator"
	"github.com/stretchr/testify/assert"
)

func Test_Equality_With_Staling(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		createdOn         time.Time
		otpVal            string
		resendAttempts    int
		maxResendAttempts int
		duration          time.Duration
		isStale           bool
		timeProvider      otp.CurrentTimeProvider
		otpToCheck        string
		wantedEquality    bool
		wantedData        otp.OtpData
	}{
		"otp is equal": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			resendAttempts:    1,
			maxResendAttempts: 3,
			duration:          3 * time.Minute,
			isStale:           false,
			timeProvider: func() time.Time {
				t := giveTime("25/08/2022 08:26:00")
				return t
			},
			otpToCheck:     "12345",
			wantedEquality: true,
			wantedData: otp.OtpData{
				CreatedOn:         giveTime("25/08/2022 08:25:00"),
				Otp:               "12345",
				ExpiryDuration:    3 * time.Minute,
				ResendAttempts:    1,
				MaxResendAttempts: 3,
				Stale:             true,
			},
		},
		"otp is not equal": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			resendAttempts:    1,
			maxResendAttempts: 3,
			duration:          3 * time.Minute,
			isStale:           false,
			timeProvider: func() time.Time {
				t := giveTime("25/08/2022 08:26:00")
				return t
			},
			otpToCheck:     "34567",
			wantedEquality: false,
			wantedData: otp.OtpData{
				CreatedOn:         giveTime("25/08/2022 08:25:00"),
				Otp:               "12345",
				ExpiryDuration:    3 * time.Minute,
				ResendAttempts:    1,
				MaxResendAttempts: 3,
				Stale:             true,
			},
		},
		"otp is already stale": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			resendAttempts:    1,
			maxResendAttempts: 3,
			duration:          3 * time.Minute,
			isStale:           true,
			timeProvider: func() time.Time {
				t := giveTime("25/08/2022 08:26:00")
				return t
			},
			otpToCheck:     "12345",
			wantedEquality: false,
			wantedData: otp.OtpData{
				CreatedOn:         giveTime("25/08/2022 08:25:00"),
				Otp:               "12345",
				ExpiryDuration:    3 * time.Minute,
				ResendAttempts:    1,
				MaxResendAttempts: 3,
				Stale:             true,
			},
		},

		"otp is stale as it is expired": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			resendAttempts:    1,
			maxResendAttempts: 3,
			duration:          3 * time.Minute,
			isStale:           false,
			timeProvider: func() time.Time {
				t := giveTime("25/08/2022 08:29:00")
				return t
			},
			otpToCheck:     "12345",
			wantedEquality: false,
			wantedData: otp.OtpData{
				CreatedOn:         giveTime("25/08/2022 08:25:00"),
				Otp:               "12345",
				ExpiryDuration:    3 * time.Minute,
				ResendAttempts:    1,
				MaxResendAttempts: 3,
				Stale:             true,
			},
		},

		"otp is stale as resend attempt exceeded": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			resendAttempts:    3,
			maxResendAttempts: 3,
			duration:          3 * time.Minute,
			isStale:           false,
			timeProvider: func() time.Time {
				t := giveTime("25/08/2022 08:26:00")
				return t
			},
			otpToCheck:     "12345",
			wantedEquality: false,
			wantedData: otp.OtpData{
				CreatedOn:         giveTime("25/08/2022 08:25:00"),
				Otp:               "12345",
				ExpiryDuration:    3 * time.Minute,
				ResendAttempts:    3,
				MaxResendAttempts: 3,
				Stale:             true,
			},
		},
	}

	for tname, tcase := range tests {
		t.Run(tname, func(t *testing.T) {
			o := otp.FromPersistence(tcase.otpVal, tcase.createdOn, tcase.duration,
				tcase.resendAttempts, tcase.maxResendAttempts, tcase.isStale)
			isEqual := o.StaleMeAfterEqualityCheck(tcase.otpToCheck, tcase.timeProvider)
			assert.Equal(t, tcase.wantedEquality, isEqual)
			assert.Equal(t, tcase.wantedData, o.OtpData())
		})
	}

	for tname, tcase := range tests {
		t.Run(tname, func(t *testing.T) {
			o := otp.FromPersistence(tcase.otpVal, tcase.createdOn, tcase.duration,
				tcase.resendAttempts, tcase.maxResendAttempts, tcase.isStale)
			isEqual := o.StaleMeOnlyIfEqualsEqualityCheck(tcase.otpToCheck, tcase.timeProvider)
			if tname == "otp is not equal" {
				tcase.wantedData.Stale = false
			}
			assert.Equal(t, tcase.wantedEquality, isEqual)
			assert.Equal(t, tcase.wantedData, o.OtpData())
		})
	}
}
