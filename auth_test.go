package hamster

import (
	"encoding/base64"
	//"fmt"
	"github.com/kr/fernet"
	"testing"
	"time"
)

func TestFernet(t *testing.T) {
	k := fernet.MustDecodeKeys("YI1ZYdopn6usnQ/5gMAHg8+pNh6D0DdaJkytdoLWUj0=")
	tok, err := fernet.EncryptAndSign([]byte("mysharedtoken"), k[0])
	if err != nil {
		t.Fatalf("fernet encryption failed %v\n", err)
	}
	stok := base64.URLEncoding.EncodeToString(tok)

	btok, err := base64.URLEncoding.DecodeString(stok)
	//fmt.Println(btok)

	if err != nil {
		t.Fatalf("fernet key decryption failed %v\n", err)
	}

	msg := fernet.VerifyAndDecrypt(btok, 60*time.Second, k)
	if string(msg) != "mysharedtoken" {
		t.Fatalf("verification failed!\n")
	}

}

func Testbcrypt(t *testing.T) {

	password := "password"

	//encrypt, get hash and salt
	hash, salt, err := encryptPassword(password)

	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	//decrypt and match

	if matched := matchPassword(password, hash, salt); !matched {
		t.Fatalf("match failed: %v", err)
	}
}

func Testbase64(t *testing.T) {
	token := "518b65cdcde9e8116e000001"

	//encode it
	encoded_token := encodeBase64Token(token)

	//decode and match
	if decoded_token := decodeToken(encoded_token); token != decoded_token {

		t.Fatalf("decoding match failed: %v")

	}

}
