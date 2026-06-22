package bcrypt

import "testing"

func TestHashPassword_Success(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Errorf("HashPassword() error = %v, want nil", err)
	}
	if hash == "" {
		t.Error("HashPassword() 返回空字符串")
	}
	// bcrypt hash 以 $2a$ 开头
	if len(hash) < 60 {
		t.Errorf("HashPassword() 返回的 hash 长度 = %v, want >= 60", len(hash))
	}
}

func TestHashPassword_DifferentInputs(t *testing.T) {
	hash1, _ := HashPassword("password1")
	hash2, _ := HashPassword("password2")
	if hash1 == hash2 {
		t.Error("不同密码应生成不同的 hash")
	}
}

func TestHashPassword_SameInputDifferentHash(t *testing.T) {
	// bcrypt 每次生成的 salt 不同，相同密码的 hash 也不同
	hash1, _ := HashPassword("samepassword")
	hash2, _ := HashPassword("samepassword")
	if hash1 == hash2 {
		t.Error("相同密码的 hash 应该不同（salt 随机）")
	}
}

func TestCheckPassword_Success(t *testing.T) {
	hash, _ := HashPassword("correctpassword")
	if !CheckPassword("correctpassword", hash) {
		t.Error("CheckPassword() 对正确密码应返回 true")
	}
}

func TestCheckPassword_WrongPassword(t *testing.T) {
	hash, _ := HashPassword("correctpassword")
	if CheckPassword("wrongpassword", hash) {
		t.Error("CheckPassword() 对错误密码应返回 false")
	}
}

func TestCheckPassword_EmptyPassword(t *testing.T) {
	hash, _ := HashPassword("")
	if !CheckPassword("", hash) {
		t.Error("CheckPassword() 对空密码应返回 true")
	}
}

func TestCheckPassword_EmptyHash(t *testing.T) {
	if CheckPassword("password", "") {
		t.Error("CheckPassword() 对空 hash 应返回 false")
	}
}
