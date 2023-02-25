package middleware

import "golang.org/x/crypto/bcrypt"

/*
password is slow: yes
----------------
password := "..."
hash, err := PasswordHash(password)
if err != nil {
	return err
}
// match = true
// match = false
if !CheckHashPassword(hash, password) {
	return ...
}
then faster: no
-----------
password := "..."
hash := sha256.New()
hash.Write([]byte(password))
sha_hash := hex.EncodeToString(hash.Sum(nil))

fmt.Println("Password -> ", password)
fmt.Println("Hash -> ", sha_hash)
*/

// PasswordHash: hash for password to string
func PasswordHash(password string) (string, error) {
	// GenerateFromPassword(..., bcrypt.DefaultCost{=10})
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckHashPassword: hash for password to bool
func CheckHashPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
