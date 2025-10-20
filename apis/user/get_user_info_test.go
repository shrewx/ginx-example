package user

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"

	"ginx-example/constants/status_error"
)

// 注意：涉及数据库的测试需要在集成测试中进行，
// 这里只测试不依赖数据库的参数校验逻辑

func TestGetUserInfo_NegativeUserID(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	// 创建接口实例
	getUserInfo := &GetUserInfo{ID: "-1"}

	// 创建测试上下文
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// 执行接口
	result, err := getUserInfo.Output(ctx)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, status_error.InvalidUserID, err)
}

func TestGetUserInfo_InvalidUserID(t *testing.T) {
	// 设置测试环境
	gin.SetMode(gin.TestMode)

	// 创建接口实例
	getUserInfo := &GetUserInfo{ID: "abc"}

	// 创建测试上下文
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// 执行接口
	result, err := getUserInfo.Output(ctx)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, status_error.InvalidUserID, err)
}
