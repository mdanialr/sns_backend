package otp_service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mdanialr/sns_backend/internal/core/repository/otp_repository"
	req "github.com/mdanialr/sns_backend/internal/requests"
	"github.com/mdanialr/sns_backend/pkg/logger"
	"github.com/mdanialr/sns_backend/pkg/twofa"
	"github.com/spf13/viper"
)

type otpSvc struct {
	v    *viper.Viper
	log  logger.Writer
	repo otp_repository.IRepository
}

// New return new service that can be used to exchange otp code with a JWT that
// has validity based on the config.
func New(v *viper.Viper, l logger.Writer, repo otp_repository.IRepository) IService {
	return &otpSvc{v, l, repo}
}

func (o *otpSvc) ValidateOTP(ctx context.Context, req *req.OTP) bool {
	ot, err := twofa.InitOTPWithConfig(o.v)
	if err != nil {
		o.log.Err("failed to ini otp with config from app:", err)
		return false
	}
	valid, err := ot.VerifyCode(req.Code)
	if err != nil {
		o.log.Err("failed to verify otp with code", req.Code, "and error:", err)
	}
	// if it's valid, then make sure that it's not been used before
	if valid {
		if ro, _ := o.repo.GetByCode(ctx, req.Code); ro != nil {
			// return false if it's exist in db
			if ro.ID != 0 {
				return false
			}
			// delete all past records
			if err = o.repo.DeleteAll(ctx); err != nil {
				o.log.Err("failed to delete all records of RegisteredCode:", err)
				return false
			}
			// then save the recent one
			if _, err = o.repo.Create(ctx, req.Code); err != nil {
				o.log.Err("failed to save new RegisteredCode:", err)
				return false
			}
			return valid
		}
	}
	return false
}

func (o *otpSvc) GetJWT(_ context.Context) (string, error) {
	dur, _ := time.ParseDuration(o.v.GetString("jwt.duration") + "m")
	claims := jwt.MapClaims{
		"user": o.v.GetString("jwt.secret"),
		"exp":  time.Now().Add(dur).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(o.v.GetString("jwt.secret")))
}
