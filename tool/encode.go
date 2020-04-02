package tool

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 每次启动时生成一个secret 避免在repo中存放明文。
// 使用hmac算法加密 secret必须是字节流

func SecretGen() []byte {
	//rand.Seed(time.Now().Unix())
	//return []byte(Md5(fmt.Sprintf("%d",rand.Int())))

	// 开发时使用固定的秘钥 方便debug
	return []byte("scarlet")
}
