package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Eucastan/eucastanpay/common/pkg/auth"
	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/user/config"
	"github.com/Eucastan/eucastanpay/services/user/internal/domain"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/request"
	"github.com/Eucastan/eucastanpay/services/user/internal/dto/response"
	"github.com/Eucastan/eucastanpay/services/user/internal/infra/redis"
	"github.com/Eucastan/eucastanpay/services/user/internal/repository"
	"github.com/Eucastan/eucastanpay/services/user/internal/usecase"
	"github.com/Eucastan/eucastanpay/services/user/internal/util/password"
	"github.com/Eucastan/eucastanpay/services/user/internal/worker"
)

type UserUseCase struct {
	User      repository.UserRepository
	Auth      repository.AuthRepository
	telemetry *telemetry.Telemetry
	cfg       *config.Config
	Email     usecase.EmailSender
	Redis     *redis.RedisClient
	Publisher *worker.PublishUserRegistration
}

func NewUserUseCase(
	user repository.UserRepository,
	auth repository.AuthRepository,
	telemetry *telemetry.Telemetry,
	cfg *config.Config,
	email usecase.EmailSender,
	redis *redis.RedisClient,
	publisher *worker.PublishUserRegistration,
) *UserUseCase {

	return &UserUseCase{
		User:      user,
		Auth:      auth,
		telemetry: telemetry,
		cfg:       cfg,
		Email:     email,
		Redis:     redis,
		Publisher: publisher,
	}
}

func (u *UserUseCase) Register(ctx context.Context, input *request.RegisterRequest) (*response.UserResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.Register")
	defer span.End()

	log.Println("1. checking email")

	if _, err := u.User.FindByEmail(ctx, input.Email); err == nil {
		log.Printf("7. failed while checking duplicate email: %v", err)
		span.RecordError(err)
		return nil, errmessage.ErrDuplicateEmail
	}

	log.Println("2. hashing password")

	hashPass, err := password.GeneratePassHash(input.Password)
	if err != nil {
		log.Printf("7. failed hash password: %v", err)
		span.RecordError(err)
		return nil, err
	}

	user := &domain.User{
		ID:            uuid.NewString(),
		Email:         input.Email,
		Phone:         input.Phone,
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Password:      hashPass,
		DateOfBirth:   input.DateOfBirth,
		Role:          "user",
		Status:        domain.StatusPending,
		EmailVerified: false,
		CreatedAt:     time.Now(),
	}

	log.Println("3. creating user")

	if err := u.User.Create(ctx, user); err != nil {
		log.Printf("7. failed to register user: %v", err)
		span.RecordError(err)
		return nil, err
	}

	log.Println("4. generating access token")

	token, err := auth.GenerateAccessToken(user.ID, user.Email, string(domain.EmailToken), u.cfg.SharedCfg.JWTSecret)
	if err != nil {
		log.Printf("7. failed access token generation: %v", err)
		span.RecordError(err)
		return nil, err
	}

	log.Println("5. saving token")

	err = u.Auth.Create(ctx, &domain.Token{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     token,
		TokenType: domain.EmailToken,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		log.Printf("7. failed to save token: %v", err)
		span.RecordError(err)
		return nil, err
	}

	log.Println("6. generating email token")

	link := "http://localhost:8080/verify?token=" + token
	if err := u.Email.SendVerificationEmail(user.Email, link); err != nil {
		log.Printf("7. failed to send email: %v", err)
		span.RecordError(err)
		return nil, err
	}

	response := response.ToUserResponse(user)
	log.Println("7. publishing event")

	if err := u.Publisher.OnUserRegistration(ctx, &response); err != nil {
		log.Printf("7. failed event: %v", err)
		span.RecordError(err)
		return nil, err
	}

	return &response, nil
}

func (u *UserUseCase) VerifyEmail(ctx context.Context, token string) error {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.VerifyEmail")
	defer span.End()

	claims, err := auth.ValidateToken(token, u.cfg.SharedCfg.JWTSecret)
	if err != nil {
		span.RecordError(err)
		return err
	}

	tkn, err := u.Auth.FindToken(ctx, token, string(domain.EmailToken))
	if err != nil {
		span.RecordError(err)
		return err
	}

	if tkn.Revoked || time.Now().After(tkn.ExpiredAt) {
		return errmessage.ErrInvalidToken
	}

	user, err := u.User.FindByID(ctx, claims.UserID)
	if err != nil {
		span.RecordError(err)
		return err
	}

	if user.EmailVerified {
		return nil
	}

	user.EmailVerified = true
	user.Status = domain.StatusActive

	if err := u.User.Update(ctx, user); err != nil {
		span.RecordError(err)
		return err
	}

	return u.Auth.Revoked(ctx, token)
}

func (u *UserUseCase) Login(ctx context.Context, input *request.LoginRequest) (*response.AuthResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.Login")
	defer span.End()

	user, err := u.User.FindByEmail(ctx, input.Email)
	if err != nil {
		span.RecordError(err)
		return nil, errmessage.ErrInvalidCredentials
	}

	if err := password.IsMatch(user.Password, input.Password); err != nil {
		span.RecordError(err)
		return nil, errmessage.ErrPasswordNotConfirmed
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Email, user.Role, u.cfg.SharedCfg.JWTSecret)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	refreshToken, err := auth.RefreshToken(user.ID, user.Email, user.Role, u.cfg.SharedCfg.JWTSecret)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	tkn := &domain.Token{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     refreshToken,
		TokenType: domain.RefreshToken,
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err = u.Auth.Create(ctx, tkn); err != nil {
		span.RecordError(err)
		return nil, err
	}

	resp := response.UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Phone:         user.Phone,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Status:        string(domain.StatusActive),
		EmailVerified: true,
		CreatedAt:     user.CreatedAt,
	}

	return &response.AuthResponse{
		User:         resp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *UserUseCase) GetAllUsers(ctx context.Context) ([]response.UserResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.GetAllUsers")
	defer span.End()

	acc, err := u.User.FindAll(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	var resp []response.UserResponse
	for _, v := range acc {
		resp = append(resp, response.UserResponse{
			ID:            v.ID,
			Email:         v.Email,
			Phone:         v.Phone,
			FirstName:     v.FirstName,
			LastName:      v.LastName,
			Status:        string(v.Status),
			EmailVerified: v.EmailVerified,
			CreatedAt:     v.CreatedAt,
		})
	}

	return resp, err
}

func (u *UserUseCase) GetUserByID(ctx context.Context, id string) (*response.UserResponse, error) {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.GetUserByID")
	defer span.End()

	user, err := u.User.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	resp := response.ToUserResponse(user)
	return &resp, nil

}

func (u *UserUseCase) UserCurrentStatus(ctx context.Context, id, status string) (string, error) {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.UserCurrentStatus")
	defer span.End()

	user, err := u.User.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	// Only allow status changes if the account is currently Active
	if user.Status != domain.StatusActive {
		return "", fmt.Errorf("only active accounts can change status. current status: %s", user.Status)
	}

	var message string

	switch status {
	case "suspended":
		user.Status = domain.StatusSuspended
		message = fmt.Sprintf("This account %s has been suspended", user.Email)

	case "closed":
		user.Status = domain.StatusClosed
		message = fmt.Sprintf("This account %s has been closed", user.ID)

	case "pending":
		user.Status = domain.StatusPending
		message = "This account is Pending"

	default:
		return "", fmt.Errorf("invalid status: %s. allowed: suspended, closed, pending", status)
	}

	if err := u.User.Update(ctx, user); err != nil {
		span.RecordError(err)
		return "", err
	}

	return message, nil
}

func (u *UserUseCase) RefreshToken(ctx context.Context, oldToken string) (string, string, error) {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.RefreshToken")
	defer span.End()

	_, err := auth.ValidateToken(oldToken, u.cfg.SharedCfg.JWTSecret)
	if err != nil {
		span.RecordError(err)
		return "", "", err
	}

	// Get token from DB
	storedToken, err := u.Auth.FindToken(ctx, oldToken, string(domain.RefreshToken))
	if err != nil {
		span.RecordError(err)
		return "", "", err
	}

	// Reuse detection
	if storedToken.Revoked {
		_ = u.Auth.RevokeAllByUser(ctx, storedToken.UserID)
		return "", "", errmessage.ErrInvalidToken
	}

	// Check expiration
	if time.Now().After(storedToken.ExpiredAt) {
		return "", "", errmessage.ErrExpiredToken
	}

	// Revoke old token
	if err := u.Auth.Revoked(ctx, oldToken); err != nil {
		span.RecordError(err)
		return "", "", err
	}

	user, err := u.User.FindByID(ctx, storedToken.UserID)
	if err != nil {
		span.RecordError(err)
		return "", "", err
	}

	newRefresh, err := auth.RefreshToken(
		user.ID,
		user.Email,
		user.Role,
		u.cfg.SharedCfg.JWTSecret,
	)
	if err != nil {
		span.RecordError(err)
		return "", "", err
	}

	// Store new token (link it to old one)
	newToken := &domain.Token{
		ID:        uuid.NewString(),
		UserID:    storedToken.UserID,
		Token:     newRefresh,
		TokenType: domain.RefreshToken,
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour),
		ParentID:  &storedToken.ID,
	}

	if err := u.Auth.Create(ctx, newToken); err != nil {
		span.RecordError(err)
		return "", "", err
	}

	// Generate new access token
	newAccess, err := auth.GenerateAccessToken(
		user.ID, user.Email, user.Role, u.cfg.SharedCfg.JWTSecret,
	)
	if err != nil {
		span.RecordError(err)
		return "", "", err
	}

	return newAccess, newRefresh, nil
}

func (u *UserUseCase) LogoutAllUsers(ctx context.Context, userID string) error {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.LogoutAllUsers")
	defer span.End()

	return u.Auth.RevokeAllByUser(ctx, userID)
}

func (u *UserUseCase) Logout(ctx context.Context, refreshToken string) error {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.Logout")
	defer span.End()

	return u.Auth.Revoked(ctx, refreshToken)
}

func (u *UserUseCase) ForgotPassword(ctx context.Context, input *request.ForgotPasswordRequest) error {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.ForgottenPassword")
	defer span.End()

	exists, err := u.User.FindByEmail(ctx, input.Email)
	if err != nil {
		span.RecordError(err)
		logrus.Errorf("user with email %s not found: %v", input.Email, err)
		return nil
	}

	resetToken, err := auth.GenerateAccessToken(exists.ID, exists.Email, exists.Role, u.cfg.SharedCfg.JWTSecret)
	if err != nil {
		span.RecordError(err)
		return err
	}

	tkn := &domain.Token{
		ID:        uuid.NewString(),
		UserID:    exists.ID,
		Token:     resetToken,
		TokenType: domain.PasswordResetToken,
		ExpiredAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := u.Auth.Create(ctx, tkn); err != nil {
		span.RecordError(err)
		return err
	}

	// Send an email
	link := "https://localhost:8080/reset-password?token=" + resetToken
	return u.Email.SendResetPasswordEmail(exists.Email, link)
}

func (u *UserUseCase) ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.ResetPassword")
	defer span.End()

	claims, err := auth.ValidateToken(req.Token, u.cfg.SharedCfg.JWTSecret)
	if err != nil {
		span.RecordError(err)
		return err
	}

	tkn, err := u.Auth.FindToken(ctx, req.Token, string(domain.PasswordResetToken))
	if err != nil || tkn.Revoked {
		span.RecordError(err)
		return err
	}

	if tkn.Revoked || time.Now().After(tkn.ExpiredAt) {
		return errmessage.ErrInvalidToken
	}

	user, err := u.User.FindByID(ctx, claims.UserID)
	if err != nil {
		span.RecordError(err)
		return err
	}

	if req.Password != req.ConfirmPass {
		return errmessage.ErrPasswordNotConfirmed
	}

	hashed, err := password.GeneratePassHash(req.Password)
	if err != nil {
		span.RecordError(err)
		return err
	}

	user.Password = hashed

	if err := u.User.Update(ctx, user); err != nil {
		span.RecordError(err)
		return err
	}

	if err := u.Auth.RevokeAllByUser(ctx, user.ID); err != nil {
		span.RecordError(err)
		return err
	}

	return u.Auth.Revoked(ctx, req.Token)
}

func (u *UserUseCase) Update(ctx context.Context, userID string, input *request.UpdateRequest) error {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.Update")
	defer span.End()

	updateUser := &domain.User{
		ID:            userID,
		Password:      input.Password,
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Status:        domain.UserStatus(input.Status),
		EmailVerified: input.EmailVerified,
	}

	return u.User.Update(ctx, updateUser)
}

func (u *UserUseCase) DeleteUser(ctx context.Context, userID string) error {
	ctx, span := u.telemetry.Start(ctx, "UserUseCase.DeleteUser")
	defer span.End()

	return u.User.Delete(ctx, userID)
}
