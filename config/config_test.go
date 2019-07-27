package config

import (
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestConfig(t *testing.T) {
	v := &Config{
		DBDriver: "postgres",
		DBConn:   "postgres://postgres:postgres@localhost/postgres?sslmode=disable",
		Port:     8000,
		Secret:   "HMAC jwt secret, use env vars instead of this",
		Host:     "http://localhost:8000",
		Minio: &Minio{
			Endpoint:     "/minio/endpoin",
			AccessKey:    "minio",
			AccessSecret: "secret",
		},
		NSQ: &NSQ{
			LookupD: "/nsqlookupd/endpoint",
		},
	}
	b, _ := yaml.Marshal(v)
	ioutil.WriteFile("c.yaml", b, 0600)
}
