package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/stretchr/testify/suite"
)

type UserAPITestSuite struct {
	suite.Suite
	baseURL     string
	accessToken string
	userId      int64
}

func (s *UserAPITestSuite) SetupSuite() {
	s.baseURL = "http://localhost:9000" // Update with your API port
}

func (s *UserAPITestSuite) request(method, path string, body interface{}, token bool) (*http.Response, error) {
	var bodyReader *bytes.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, s.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token && s.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.accessToken)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	return client.Do(req)
}

func (s *UserAPITestSuite) TestA_Register() {
	req := types.RegisterReq{
		Username: fmt.Sprintf("testuser_%d", time.Now().Unix()),
		Password: "password123",
		Phone:    "13800138000",
		Email:    "test@example.com",
	}

	resp, err := s.request("POST", "/api/user/register", req, false)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var result types.RegisterResp
	err = json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)
	s.NotZero(result.UserId)
	s.userId = result.UserId
}

func (s *UserAPITestSuite) TestB_Login() {
	req := types.LoginReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := s.request("POST", "/api/user/login", req, false)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var result types.TokenResp
	err = json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)
	s.NotEmpty(result.AccessToken)
	s.accessToken = result.AccessToken
}

func (s *UserAPITestSuite) TestC_GetProfile() {
	resp, err := s.request("GET", "/api/user/profile", nil, true)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var result types.UserInfo
	err = json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)
	s.NotZero(result.Id)
}

func (s *UserAPITestSuite) TestD_UpdateProfile() {
	req := types.UpdateProfileReq{
		Nickname: "Updated Name",
		Gender:   "male",
	}

	resp, err := s.request("PUT", "/api/user/profile", req, true)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
}

func (s *UserAPITestSuite) TestE_AddressOperations() {
	// Add address
	addReq := types.AddressReq{
		ReceiverName:  "John Doe",
		ReceiverPhone: "13800138000",
		Province:      "TestProvince",
		City:          "TestCity",
		District:      "TestDistrict",
		DetailAddress: "Test Address 123",
		IsDefault:     true,
	}

	resp, err := s.request("POST", "/api/user/addresses", addReq, true)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var address types.Address
	err = json.NewDecoder(resp.Body).Decode(&address)
	s.Require().NoError(err)
	s.NotZero(address.Id)

	// List addresses
	resp, err = s.request("GET", "/api/user/addresses", nil, true)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var addresses []types.Address
	err = json.NewDecoder(resp.Body).Decode(&addresses)
	s.Require().NoError(err)
	s.NotEmpty(addresses)
}

func (s *UserAPITestSuite) TestF_WalletOperations() {
	// Get wallet
	resp, err := s.request("GET", "/api/user/wallet", nil, true)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var wallet types.WalletDetail
	err = json.NewDecoder(resp.Body).Decode(&wallet)
	s.Require().NoError(err)

	// Recharge wallet
	rechargeReq := types.RechargeReq{
		Amount:      100.00,
		PaymentType: 1,
	}

	resp, err = s.request("POST", "/api/user/wallet/recharge", rechargeReq, true)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	// Get transactions
	resp, err = s.request("GET", "/api/user/wallet/transactions?page=1&pageSize=10", nil, true)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	var transactions types.TransactionListResp
	err = json.NewDecoder(resp.Body).Decode(&transactions)
	s.Require().NoError(err)
}

func (s *UserAPITestSuite) TestG_PasswordOperations() {
	// Change password
	changeReq := types.ChangePasswordReq{
		OldPassword: "password123",
		NewPassword: "newpassword123",
	}

	resp, err := s.request("PUT", "/api/user/password/change", changeReq, true)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
}

func TestUserAPI(t *testing.T) {
	suite.Run(t, new(UserAPITestSuite))
}
