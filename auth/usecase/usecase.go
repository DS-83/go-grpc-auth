package usecase

import (
	"context"
	"errors"
	"example-grpc-auth/auth"
	e "example-grpc-auth/err"
	"example-grpc-auth/models"

	"time"

	pb "example-grpc-auth/api"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	userRepo  auth.UserRepo
	tokenRepo auth.TokenRepo
	jwtKey    []byte
}

type AuthClaims struct {
	User *models.User
	jwt.RegisteredClaims
}

func NewAuthServer(a auth.UserRepo, t auth.TokenRepo, b []byte) *AuthServer {
	return &AuthServer{
		userRepo:  a,
		tokenRepo: t,
		jwtKey:    b,
	}
}

func (s *AuthServer) SignUp(ctx context.Context, r *pb.SignUpRequest) (*pb.User, error) {
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8
	//  (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), 8)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.CreateUser(ctx, r.Username, string(hashedPassword))
	if err != nil {
		if errors.Is(err, e.ErrDupKey) {
			return nil, status.Error(codes.AlreadyExists, e.ErrDupKey.Error())
		}
		return nil, err
	}

	resp, err := s.userRepo.GetUser(ctx, r.Username, r.Password)
	if err != nil {
		return nil, err
	}

	return toPbUser(resp), nil

}

// Sign in user and get JWT string
func (s *AuthServer) SignIn(ctx context.Context, r *pb.SignInRequest) (*pb.SignInResponce, error) {
	user, err := s.userRepo.GetUser(ctx, r.Username, r.Password)
	if err != nil {
		switch err {
		case e.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, e.ErrUserNotFound.Error())
		case e.ErrInvalidCred:
			return nil, status.Error(codes.InvalidArgument, e.ErrInvalidCred.Error())
		}
		return nil, err
	}
	// Token expiration time:
	exp := time.Now().Add(86400 * time.Second)
	// Create the Claims
	claims := AuthClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	// Create jwt Token. Signing method HS256 uses a []byte key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Token string
	ts, err := token.SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}

	return &pb.SignInResponce{
		Token: ts,
	}, nil
}

func (s *AuthServer) Delete(ctx context.Context, r *pb.DelRequest) (*pb.Response, error) {
	user := toModelsUser(r.User)

	if err := s.userRepo.DeleteUser(ctx, user); err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	if err := s.tokenRepo.RevokeToken(ctx, r.Token); err != nil {
		return nil, err
	}
	return &pb.Response{
		Response: "Ok",
	}, nil
}

func (s *AuthServer) Update(ctx context.Context, r *pb.UpdRequest) (*pb.User, error) {
	filt := toModelsUser(r.Filtr)
	upd := toModelsUser(r.Upd)

	if upd.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(upd.Password), 8)
		if err != nil {
			return nil, err
		}
		upd.Password = string(hashedPassword)
	}

	user, err := s.userRepo.UpdateUser(ctx, filt, upd)
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	// Check if request is from sign-up case. Break if true.
	if r.SignUp {
		return toPbUser(user), nil
	}

	if err := s.tokenRepo.RevokeToken(ctx, r.Token); err != nil {
		return nil, err
	}

	return toPbUser(user), nil
}

func (s *AuthServer) ParseToken(ctx context.Context, r *pb.ParseRequest) (*pb.User, error) {
	token, err := jwt.ParseWithClaims(r.Token, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, e.ErrInvalidAccessToken.Error())
	}
	// Check for revoked token
	if ok, _ = s.tokenRepo.IsRevoked(ctx, r.Token); ok {
		return nil, status.Error(codes.Unauthenticated, e.ErrInvalidAccessToken.Error())
	}

	return toPbUser(claims.User), nil
}

func toModelsUser(u *pb.User) *models.User {
	return &models.User{
		ID:       u.Id,
		MysqlID:  int(u.MysqlId),
		Username: u.Username,
		Password: u.Password,
	}
}

func toPbUser(u *models.User) *pb.User {
	return &pb.User{
		Id:       u.ID,
		MysqlId:  int64(u.MysqlID),
		Username: u.Username,
		Password: u.Password,
	}
}
