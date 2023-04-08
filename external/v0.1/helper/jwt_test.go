package helper

import (
	"strings"
	"testing"
)

// Test to Validate and create JWT
func TestJWTCreate(t *testing.T) {
	var err error
	var token string
	testInput := []map[string]string{
		{
			"username": "maximo",
			"password": "1234",
		}, {
			"username": "juan",
			"password": "5678",
		}, {
			"username": "maximo",
			"password": "9101112",
		}, {},
	}
	jwt := NewJWT()
	for _, v := range testInput {
		t.Run("Create token", func(t *testing.T) {
			token, err = jwt.Create(v)
			if err != nil {
				t.Errorf("error creating token: %v", err)
			}
			if len(strings.Split(token, ".")) != 3 {
				t.Errorf("mal formed token: %v", token)
			}
		})
		t.Run("Validate token", func(t *testing.T) {
			if token == "" {
				t.Error("Empty token")
			}
			values, err := jwt.Validate(token)
			if err != nil {
				t.Errorf("error validating token: %v", token)
			}
			v2, ok := values.(map[string]interface{})
			if !ok {
				t.Errorf("error asserting the tipe of payload: %v, %v", v2, values)
			}
			for k, value := range v2 {
				if value != v[k] {
					t.Errorf("Values are different Want %s got %s.", v[k], value)
				}
			}
		})
	}
}

// test to Meter the time spent creating tokens aprox 1/ms
// cant reduce allocations to improve performance
func BenchmarkToken(b *testing.B) {
	var jwt *JWT
	b.Run("generate interface", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			jwt = NewJWT()
		}
	})
	b.ResetTimer()
	b.Run("Create Token", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			jwt.Create(nil)
		}
	})

}
