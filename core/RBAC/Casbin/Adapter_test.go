package casbin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAdapter_FilePolicy(t *testing.T) {
	modelPath := "./test_model.conf"
	policyPath := "./test_policy.csv"

	// 检查模型和策略文件是否存在
	_, err := os.Stat(modelPath)
	assert.NoError(t, err)
	_, err = os.Stat(policyPath)
	assert.NoError(t, err)

	enforcer, err := NewAdapter(nil, modelPath, policyPath)
	assert.NoError(t, err)
	assert.NotNil(t, enforcer)
	assert.NotNil(t, enforcer.Enforcer)

	// alice 属于 admin，可以 GET /api/user
	ok, err := enforcer.VerifyAuth("alice", "/api/user", "GET")
	assert.NoError(t, err)
	assert.True(t, ok)

	// bob 属于 editor，可以 POST /api/article
	ok, err = enforcer.VerifyAuth("bob", "/api/article", "POST")
	assert.NoError(t, err)
	assert.True(t, ok)

	// alice 不能 POST /api/article
	ok, err = enforcer.VerifyAuth("alice", "/api/article", "POST")
	assert.NoError(t, err)
	assert.False(t, ok)

	// bob 不能 GET /api/user
	ok, err = enforcer.VerifyAuth("bob", "/api/user", "GET")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestUpdateEnforcer(t *testing.T) {
	modelPath := "./test_model.conf"
	policyPath := "./test_policy.csv"

	enforcer, err := NewAdapter(nil, modelPath, policyPath)
	assert.NoError(t, err)
	assert.NotNil(t, enforcer)
	// 热更新策略（此处只是调用，实际可结合动态策略文件测试）
	err = enforcer.UpdateEnforcer()
	assert.NoError(t, err)
}
