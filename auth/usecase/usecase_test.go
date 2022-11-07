package usecase

import (
	"context"
	pb "example-grpc-auth/api"
	"example-grpc-auth/auth/repo/mock"
	"example-grpc-auth/models"
	"reflect"
	"testing"

	mc "github.com/stretchr/testify/mock"
)

var (
	userRepo  = new(mock.UserRepoMock)
	tokenRepo = new(mock.TokenRepoMock)
	server    = new(pb.UnimplementedAuthServiceServer)
)

func TestAuthServer_ParseToken(t *testing.T) {
	type fields struct {
		UnimplementedAuthServiceServer *pb.UnimplementedAuthServiceServer
		userRepo                       *mock.UserRepoMock
		tokenRepo                      *mock.TokenRepoMock
		jwtKey                         []byte
	}
	type args struct {
		ctx context.Context
		r   *pb.ParseRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.User
		wantErr bool
	}{{
		name: "Valid token",
		fields: fields{
			UnimplementedAuthServiceServer: server,
			userRepo:                       userRepo,
			tokenRepo:                      tokenRepo,
			jwtKey:                         []byte("34989fdf3df"),
		},
		args: args{
			ctx: context.Background(),
			r: &pb.ParseRequest{
				Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7IklEIjoiNjM2ODI0NjVkNzczMTNhMDZhZDU0NWM2IiwiTXlzcWxJRCI6MzAsIlVzZXJuYW1lIjoidXNlcjQ5IiwiUGFzc3dvcmQiOiIkMmEkMDgkd1U3aHY3MVloMnBWamQ5V2VvejNqT0JNMmpaRktLZGd3bkxxVVJVdTRmUTdJaDk5QlQyQlMifSwiZXhwIjoxNjY3ODU4MTA5fQ.pm6YN0eP6Ez_esd4REXihXYuyUf8ABQ2a8_bFEOeu0g",
			},
		},
		want: &pb.User{
			Id:       "63682465d77313a06ad545c6",
			MysqlId:  30,
			Username: "user49",
			Password: "$2a$08$wU7hv71Yh2pVjd9Weoz3jOBM2jZFKKdgwnLqURUu4fQ7Ih99BT2BS",
		},
		wantErr: false,
	}, {
		name: "invalid token",
		fields: fields{
			UnimplementedAuthServiceServer: server,
			userRepo:                       userRepo,
			tokenRepo:                      tokenRepo,
			jwtKey:                         []byte("34989fdf3df"),
		},
		args: args{
			ctx: context.Background(),
			r: &pb.ParseRequest{
				Token: "eyJhbGciOiJIUzt1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7IklEIjoiNjM2ODr0NjVkNzczMTNhMDZhZDU0NWM2IiwiTXlzcWxJRCI6MzAsIlVzZXJuYW1lIjoidXNlcjQ5IiwiUGFzc3dvcmQiOiIkMmEdMDgkd1U3aHY3MVloMnBWamQ5V2VvejNqT0JNMmpaRktLZGd3bkxxVVJVdTRmUTdJaDk5QlQyQlMifSwiZXhwIjoxNjY35DU4MTA5fQ.pm6YN04P6Ez_esd4REXihXYuyUf8ABQ2a8_bFEOeu0g",
			},
		},
		want:    nil,
		wantErr: true,
	},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthServer{
				UnimplementedAuthServiceServer: *tt.fields.UnimplementedAuthServiceServer,
				userRepo:                       tt.fields.userRepo,
				tokenRepo:                      tt.fields.tokenRepo,
				jwtKey:                         tt.fields.jwtKey,
			}

			tt.fields.tokenRepo.On("IsRevoked", tt.args.r.Token).Return(false, nil)
			got, err := s.ParseToken(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthServer.ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthServer.ParseToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthServer_SignUp(t *testing.T) {
	type fields struct {
		UnimplementedAuthServiceServer *pb.UnimplementedAuthServiceServer
		userRepo                       *mock.UserRepoMock
		tokenRepo                      *mock.TokenRepoMock
		jwtKey                         []byte
	}
	type args struct {
		ctx context.Context
		r   *pb.SignUpRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.User
		wantErr bool
	}{{
		name: "valid input",
		fields: fields{
			UnimplementedAuthServiceServer: server,
			userRepo:                       userRepo,
			tokenRepo:                      tokenRepo,
			jwtKey:                         []byte("34989fdf3df"),
		},
		args: args{
			ctx: context.Background(),
			r: &pb.SignUpRequest{
				Username: "test",
				Password: mc.Anything,
			},
		},
		want: &pb.User{
			Username: "test",
			Password: mc.Anything,
		},
		wantErr: false,
	},
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthServer{
				UnimplementedAuthServiceServer: *tt.fields.UnimplementedAuthServiceServer,
				userRepo:                       tt.fields.userRepo,
				tokenRepo:                      tt.fields.tokenRepo,
				jwtKey:                         tt.fields.jwtKey,
			}

			tt.fields.userRepo.On("CreateUser", tt.args.r.Username, tt.args.r.Password).Return(nil)
			tt.fields.userRepo.On("GetUser", tt.args.r.Username, tt.args.r.Password).Return(&models.User{
				Username: "test",
				Password: mc.Anything,
			}, nil)

			got, err := s.SignUp(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthServer.SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthServer.SignUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthServer_SignIn(t *testing.T) {
	type fields struct {
		UnimplementedAuthServiceServer pb.UnimplementedAuthServiceServer
		userRepo                       *mock.UserRepoMock
		tokenRepo                      *mock.TokenRepoMock
		jwtKey                         []byte
	}
	type args struct {
		ctx context.Context
		r   *pb.SignInRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.User
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "valid input",
			fields: fields{
				UnimplementedAuthServiceServer: *server,
				userRepo:                       userRepo,
				tokenRepo:                      tokenRepo,
				jwtKey:                         []byte("123"),
			},
			args: args{
				ctx: context.Background(),
				r: &pb.SignInRequest{
					Username: "test",
					Password: mc.Anything,
				},
			},
			want: &pb.User{
				Id:       "",
				MysqlId:  0,
				Username: "test",
				Password: mc.Anything,
			},

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthServer{
				UnimplementedAuthServiceServer: tt.fields.UnimplementedAuthServiceServer,
				userRepo:                       tt.fields.userRepo,
				tokenRepo:                      tt.fields.tokenRepo,
				jwtKey:                         tt.fields.jwtKey,
			}
			tt.fields.userRepo.On("GetUser", tt.args.r.Username, tt.args.r.Password).Return(&models.User{
				Username: "test",
				Password: mc.Anything,
			}, nil)
			// request token
			got1, err := s.SignIn(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthServer.SignIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// parse token
			tt.fields.tokenRepo.On("IsRevoked", got1.Token).Return(false, nil)
			pr := &pb.ParseRequest{
				Token: got1.Token,
			}
			got2, err := s.ParseToken(tt.args.ctx, pr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthServer.SignIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got2, tt.want) {
				t.Errorf("AuthServer.SignIn() = %v, want %v", got2, tt.want)
			}
		})
	}
}

func TestAuthServer_Delete(t *testing.T) {
	type fields struct {
		UnimplementedAuthServiceServer pb.UnimplementedAuthServiceServer
		userRepo                       *mock.UserRepoMock
		tokenRepo                      *mock.TokenRepoMock
		jwtKey                         []byte
	}
	type args struct {
		ctx context.Context
		r   *pb.DelRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.Response
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "valid input",
			fields: fields{
				UnimplementedAuthServiceServer: *server,
				userRepo:                       userRepo,
				tokenRepo:                      tokenRepo,
				jwtKey:                         []byte("123"),
			},
			args: args{
				ctx: context.Background(),
				r: &pb.DelRequest{
					User: &pb.User{
						Username: "test",
						Password: mc.Anything,
					},
					Token: mc.Anything,
				},
			},
			want: &pb.Response{
				Response: "Ok",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthServer{
				UnimplementedAuthServiceServer: tt.fields.UnimplementedAuthServiceServer,
				userRepo:                       tt.fields.userRepo,
				tokenRepo:                      tt.fields.tokenRepo,
				jwtKey:                         tt.fields.jwtKey,
			}

			tt.fields.userRepo.On("DeleteUser", toModelsUser(tt.args.r.User)).Return(nil)
			tt.fields.tokenRepo.On("RevokeToken", tt.args.r.Token).Return(nil)

			got, err := s.Delete(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthServer.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthServer.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthServer_Update(t *testing.T) {
	type fields struct {
		UnimplementedAuthServiceServer *pb.UnimplementedAuthServiceServer
		userRepo                       *mock.UserRepoMock
		tokenRepo                      *mock.TokenRepoMock
		jwtKey                         []byte
	}
	type args struct {
		ctx context.Context
		r   *pb.UpdRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.User
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "valid input",
			fields: fields{
				UnimplementedAuthServiceServer: server,
				userRepo:                       userRepo,
				tokenRepo:                      tokenRepo,
				jwtKey:                         []byte{},
			},
			args: args{
				ctx: context.Background(),
				r: &pb.UpdRequest{
					Filtr: &pb.User{
						Id:       "1",
						MysqlId:  0,
						Username: "test",
						Password: mc.Anything,
					},
					Upd: &pb.User{
						Id:       "1",
						MysqlId:  1,
						Username: "test1",
						Password: "",
					},
					Token:  mc.Anything,
					SignUp: false,
				},
			},
			want: &pb.User{
				Id:       "1",
				MysqlId:  1,
				Username: "test1",
				Password: mc.Anything,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthServer{
				UnimplementedAuthServiceServer: *tt.fields.UnimplementedAuthServiceServer,
				userRepo:                       tt.fields.userRepo,
				tokenRepo:                      tt.fields.tokenRepo,
				jwtKey:                         tt.fields.jwtKey,
			}

			tt.fields.userRepo.On("UpdateUser", toModelsUser(tt.args.r.Filtr), toModelsUser(tt.args.r.Upd)).Return(toModelsUser(tt.want), nil)
			tt.fields.tokenRepo.On("RevokeToken", tt.args.r.Token).Return(nil)
			got, err := s.Update(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthServer.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthServer.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
