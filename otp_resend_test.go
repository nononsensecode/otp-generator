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
			timeProvider:      timeProvider("25/08/2022 08:26:00"),
			generator:         generator("34567", nil),
			wantedErr:         nil,
			wantedData:        newOtpData("34567", "25/08/2022 08:25:00", 3*time.Minute, 2, 3, false),
		},

		"resendable returns stale error": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			duration:          3 * time.Minute,
			resendAttempts:    1,
			maxResendAttempts: 3,
			isStale:           true,
			timeProvider:      timeProvider("25/08/2022 08:29:00"),
			generator:         generator("34567", nil),
			wantedErr:         otp.ErrStale,
			wantedData:        newOtpData("12345", "25/08/2022 08:25:00", 3*time.Minute, 1, 3, true),
		},

		"resendable returns expiry error": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			duration:          3 * time.Minute,
			resendAttempts:    1,
			maxResendAttempts: 3,
			isStale:           false,
			timeProvider:      timeProvider("25/08/2022 08:29:00"),
			generator:         generator("34567", nil),
			wantedErr:         otp.ErrExpiry,
			wantedData:        newOtpData("12345", "25/08/2022 08:25:00", 3*time.Minute, 1, 3, true),
		},

		"resendable returns resend exceeds error": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			duration:          3 * time.Minute,
			resendAttempts:    3,
			maxResendAttempts: 3,
			isStale:           false,
			timeProvider:      timeProvider("25/08/2022 08:28:00"),
			generator:         generator("34567", nil),
			wantedErr:         otp.ErrResendExcceded,
			wantedData:        newOtpData("12345", "25/08/2022 08:25:00", 3*time.Minute, 3, 3, true),
		},

		"resendable returns otp generation error": {
			createdOn:         giveTime("25/08/2022 08:25:00"),
			otpVal:            "12345",
			duration:          3 * time.Minute,
			resendAttempts:    1,
			maxResendAttempts: 3,
			isStale:           false,
			timeProvider:      timeProvider("25/08/2022 08:26:00"),
			generator:         generator("", fmt.Errorf("i can't create otp")),
			wantedErr:         fmt.Errorf("i can't create otp"),
			wantedData:        newOtpData("12345", "25/08/2022 08:25:00", 3*time.Minute, 1, 3, false),
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
