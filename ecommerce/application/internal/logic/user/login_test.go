package user

import (
    "context"
    "errors"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
    "github.com/stretchr/testify/assert"
)

func TestLoginLogic_Login(t *testing.T) {
    tests := []struct {
        name    string
        req     *types.LoginReq
        mock    func(mockUser *User)
        want    *types.TokenResp
        wantErr error
    }{
        {
            name: "successful login",
            req: &types.LoginReq{
                Username: "test_user",
                Password: "password123",
            },
            mock: func(mockUser *User) {
                mockUser.EXPECT().Login(
                    context.Background(),
                    &user.LoginRequest{
                        Username: "test_user",
                        Password: "password123",
                    },
                ).Return(&user.LoginResponse{
                    AccessToken:  "test_access_token",
                    RefreshToken: "test_refresh_token",
                    ExpiresIn:    3600,
                }, nil)
            },
            want: &types.TokenResp{
                AccessToken:  "test_access_token",
                RefreshToken: "test_refresh_token",
                ExpiresIn:    3600,
            },
            wantErr: nil,
        },
        {
            name: "failed login",
            req: &types.LoginReq{
                Username: "test_user",
                Password: "wrong_password",
            },
            mock: func(mockUser *User) {
                mockUser.EXPECT().Login(
                    context.Background(),
                    &user.LoginRequest{
                        Username: "test_user",
                        Password: "wrong_password",
                    },
                ).Return(nil, errors.New("invalid credentials"))
            },
            want:    nil,
            wantErr: errors.New("invalid credentials"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mock user client
            mockUser := NewUser(t)
            
            // Setup mock expectations
            tt.mock(mockUser)

            // Create service context with mock
            svcCtx := &svc.ServiceContext{
                UserRpc: mockUser,
            }

            // Create login logic instance
            logic := NewLoginLogic(context.Background(), svcCtx)

            // Execute login
            got, err := logic.Login(tt.req)

            // Assert results
            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr.Error(), err.Error())
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}