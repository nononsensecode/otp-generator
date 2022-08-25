package otp

import (
	"fmt"
	"time"
)

var (
	ErrExpiry         = fmt.Errorf("otp is expired")
	ErrResendExcceded = fmt.Errorf("resend attempts exceeded")
	ErrStale          = fmt.Errorf("otp is already stale")
)

type OtpGenerator func(length int) (string, error)

type CurrentTimeProvider func() time.Time

type Otp struct {
	otp               string
	createdOn         time.Time
	expiryDuration    time.Duration
	resendAttempts    int
	maxResendAttempts int
	stale             bool
}

type OtpData struct {
	Otp               string
	CreatedOn         time.Time
	ExpiryDuration    time.Duration
	ResendAttempts    int
	MaxResendAttempts int
	Stale             bool
}

func New(timeProvider CurrentTimeProvider, duration time.Duration,
	generator OtpGenerator, length, maxResendAttempts int) (o Otp, err error) {
	otpVal, err := generator(length)
	if err != nil {
		return
	}

	o = Otp{
		otp:               otpVal,
		createdOn:         timeProvider(),
		expiryDuration:    duration,
		resendAttempts:    1,
		maxResendAttempts: maxResendAttempts,
	}

	return
}

func FromPersistence(otpVal string, createdOn time.Time,
	duration time.Duration, resendAttempts, maxResendAttempts int, isStale bool) (o Otp) {
	return Otp{
		otp:               otpVal,
		createdOn:         createdOn,
		expiryDuration:    duration,
		resendAttempts:    resendAttempts,
		maxResendAttempts: maxResendAttempts,
		stale:             isStale,
	}
}

func (o Otp) OtpData() OtpData {
	return OtpData{
		Otp:               o.otp,
		CreatedOn:         o.createdOn,
		ExpiryDuration:    o.expiryDuration,
		ResendAttempts:    o.resendAttempts,
		MaxResendAttempts: o.maxResendAttempts,
		Stale:             o.stale,
	}
}

func (o *Otp) StaleMe() {
	o.stale = true
}

func (o *Otp) Resendable(timeProvider CurrentTimeProvider, generator OtpGenerator, length int) (err error) {
	err = o.staleMeIfInvalid(timeProvider)
	if err != nil {
		return
	}

	otpVal, err := generator(length)
	if err != nil {
		return
	}

	o.otp = otpVal
	o.resendAttempts++
	return
}

func (o *Otp) validate(timeProvider CurrentTimeProvider) (err error) {
	if o.stale {
		err = ErrStale
		return
	}

	currTime := timeProvider()
	if o.createdOn.Add(o.expiryDuration).Before(currTime) {
		err = ErrExpiry
		o.StaleMe()
		return
	}

	if o.resendAttempts == o.maxResendAttempts {
		err = ErrResendExcceded
		o.StaleMe()
		return
	}

	return
}

func (o *Otp) staleMeIfInvalid(timeProvider CurrentTimeProvider) (err error) {
	err = o.validate(timeProvider)
	if err != nil {
		o.StaleMe()
		return
	}

	return
}

func (o *Otp) StaleMeAfterEqualityCheck(otpVal string, timeProvider CurrentTimeProvider) (isEqual bool) {
	err := o.staleMeIfInvalid(timeProvider)
	if err != nil {
		return
	}

	o.StaleMe()
	if o.otp == otpVal {
		isEqual = true
	}

	return
}

func (o *Otp) StaleMeOnlyIfEqualsEqualityCheck(otpVal string, timeProvider CurrentTimeProvider) (isEqual bool) {
	err := o.staleMeIfInvalid(timeProvider)
	if err != nil {
		return
	}

	if o.otp == otpVal {
		o.StaleMe()
		isEqual = true
		return
	}

	return
}
