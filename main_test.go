package main

import (
	"testing"
)

var ca = `-----BEGIN CERTIFICATE-----
MIIDBjCCAe6gAwIBAgIBATANBgkqhkiG9w0BAQsFADA8MTowOAYDVQQDDDFNeVNR
TF9TZXJ2ZXJfNS43LjMzX0F1dG9fR2VuZXJhdGVkX0NBX0NlcnRpZmljYXRlMB4X
DTIxMDYyNDEzNTM1OVoXDTMxMDYyMjEzNTM1OVowPDE6MDgGA1UEAwwxTXlTUUxf
U2VydmVyXzUuNy4zM19BdXRvX0dlbmVyYXRlZF9DQV9DZXJ0aWZpY2F0ZTCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANoxVz08XCRnA6r44LoqNnW2ODE+
F+WrFhsXvUm8eyIBzTcbhc8jrR4bBjx5JmtVTf4aSySB4YYqswqTkIDBvSyITCAq
4nfTn/EgBH6W0+qpUcR5QSQs4filRlYz8cpN/ggj9ZCklWv1uVLBkXYCG8uDpQo/
pjMbP1Qfw9LixRULLCrUjBRi0xNHsJguD9NmCK1KhZ8AQP2lBjEQ4ttFFsdvNIlr
QMOILq6bi4uFGIh/uFdfnsAfk8T6SLoCMLPUM8Z+S2bvkcWECX1YfvDfED7hxjEb
Vyj1nQUjWPk+OuRSlDzntKj2DqZP+/igGgaNm3UDJrjgqclPZas8JpP/ZiMCAwEA
AaMTMBEwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAlNldjvJg
VkMjLaRHrlcWmQR2rWQfVMAp7kDeYtmylhv1BIFBLT4gB8gXwPGjjfpMlj0lh8Rr
8f0OdQ4udRaaSrcAd93LC+snVTbiC/x72G99bJ7RKOfmtjRV43KN8V8XVNudthYW
tCDnTjwGF8630g/RKYgBdaKOcjNLRxOB5wCD8Ft2h9IRu+lX/qux1qv+jqo/9WKB
EntjpUC7pUbFq/Z6pMFS9j3mlyjRdb6IpebbKtgHWlr5U6+yft5w8+XWuY5+Jf17
I/NNIET1wmcgvW9JQ9GBkyJwybDqd8RBYSYxeRs/j7hzi5T00HIVQpLzCloKT138
h3Su3dsC5Uu2PQ==
-----END CERTIFICATE-----
`

var clientCert = `-----BEGIN CERTIFICATE-----
MIIDBzCCAe+gAwIBAgIBAzANBgkqhkiG9w0BAQsFADA8MTowOAYDVQQDDDFNeVNR
TF9TZXJ2ZXJfNS43LjMzX0F1dG9fR2VuZXJhdGVkX0NBX0NlcnRpZmljYXRlMB4X
DTIxMDYyNDEzNTQwMFoXDTMxMDYyMjEzNTQwMFowQDE+MDwGA1UEAww1TXlTUUxf
U2VydmVyXzUuNy4zM19BdXRvX0dlbmVyYXRlZF9DbGllbnRfQ2VydGlmaWNhdGUw
ggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDKl9enXlq4ZnaJgmo/t0fK
vwwgoY2s0t5Dm0j4bsGGU6Qohq9LlQoxz/umBtPhlXURYnDZelJ+YF+HUnbydeMn
X2XQMCqA+R9abX7sMKKQWAzhqZnPp/jUo3xevLKOlekNUYbzdi9jquttvJiEQ/Pt
UsFQfmL2LSbNuIufwbRo5hWg1v6fHSZ8l6xxteY8GvHn0gANZVbsJlMvR+SKYGSo
uiHKrLJOMvOMd4i4hq6VtccIMIuCoYRym1k3bok+l94CXr5l5oxXHretRuOQQavY
UbidRumuik31JuCnVyV6bJI4qvxc0nxvEO60NSDhZDmxzUR+7potZ7q4ZkwEPWY/
AgMBAAGjEDAOMAwGA1UdEwEB/wQCMAAwDQYJKoZIhvcNAQELBQADggEBAAWJfJqG
T/EI+2Y+M1sn2VpGStR5sdghMSJ/WOZTBU4YBfcRY9GrVfPkDbbATwdMs5BLVbyf
AnFC5HUmGwSdZGwtZfFXArCtY1jfsPUy/ryiADz5Ilsb5CqjkyupaHo0DBdt/5LF
FY+6019nCQ64WTddxlauKMig6Li6voBskF+wvs7QDzsJ+eT8nR4+iT6e0ztCUxA9
u6tslVj24x2r/1B74lRmpBalEe5a9T13gYYjFCV3Alpzzxo9JTwDg0vsTscMDAK4
BRgdYFiEDC4BsCjGg5zDWOzmRrIyZbTk8LgJlOZUgwHCkZjWLEdR1yz0svEO0KoD
RmYNkvuKS+LqhgM=
-----END CERTIFICATE-----
`

var clientKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAypfXp15auGZ2iYJqP7dHyr8MIKGNrNLeQ5tI+G7BhlOkKIav
S5UKMc/7pgbT4ZV1EWJw2XpSfmBfh1J28nXjJ19l0DAqgPkfWm1+7DCikFgM4amZ
z6f41KN8XryyjpXpDVGG83YvY6rrbbyYhEPz7VLBUH5i9i0mzbiLn8G0aOYVoNb+
nx0mfJescbXmPBrx59IADWVW7CZTL0fkimBkqLohyqyyTjLzjHeIuIaulbXHCDCL
gqGEcptZN26JPpfeAl6+ZeaMVx63rUbjkEGr2FG4nUbpropN9Sbgp1clemySOKr8
XNJ8bxDutDUg4WQ5sc1Efu6aLWe6uGZMBD1mPwIDAQABAoIBAHlo8/E2Z3G43fLc
eCYFBX4obfEkQA5BNLgZWGqhnVKNz4QF6wmFoN1nQe0mjFXVRpoLf4uldciMotMr
BWOaaG4RCqfwJZizMynzEqwjHVkZm2alSCsdA88SxpbTyMIQbuDvact/mTjGll5m
NpSVUJp6DKXbiBz1xBpnLS2pXR+7Kid7PFFxcI2svMOMnDtk/Yeib0tddJ8b9LLd
w4rrJ4XIEiOoxHSXoo4sfQd63Gfe1E9foIiJbLwqPPwdWmx0YXNFbThGqNVqmwtW
UCcWsn3p9CcQN17Yvn2n9WffWeLEPtTvA3efcfKh2ei25L2TFvuY1cUL9kXlprmN
ouI06FECgYEA7a4mODJYCaHco7VtIhwieG4c3saKJ2N8AMfqJxpfc1YfbhE6zKue
ETl/EcCn4fmSEwlGpOXigPls1N129PMggwXJvoX+tX3o1juYuI469L3QWs4RO+Y8
cPt+Nm3kjMQWjaiBxesqF9D2sLFxcADq+d9CkVaK9ywEdo4thSOzuKkCgYEA2jVc
r5S1TLq7LLa4yc8FNiHm4igKQNX5JaNqL92VBoN5NLGwXbrAENhnNOTgjap9qlJO
o6nUDWc2HBJoo7LLB4fXyabzbYPLO/zpeto/OvLR6s2flTBnloz1mswMenehcR12
VtdtFKwNeeb97XQST6neqgUlubCBL99kI6NGcKcCgYBbl1ulFkl01Mo40AZRObRS
4mP/uVSt3xl/F4r1LKWRxaNw//S/wHa2PojoJ2zKmSkgatetXeVOPFAjK4DW3gYu
V4GqCOrht7aNMAQnQrAXdjofc/+SLKQoCdJWWTzUvg/O4Ru5UGk1KlKWPprvRXtH
dTCq30XIAE9r/FgwUVTFmQKBgQC2mwlhdhTyjSGuHYsudyxpTJ2mjNNnxqilMsfe
HpRjywpXaoupJGsurapQvWidpRtcuVxN4gUu4jyJ6W3f7/Ov7aJLgoccErq1DMVm
pdVcQ3AgCsuNUeCeQIAYdsHKiupzerBZpBYIwqLUbNX7LLNWL1XC+mPDrU2u/kcS
+8Qq6QKBgBysxdqOzgYoiuWqFDn7AxbWBZiqUTTDir5G8j/i+6cteywpurv9cOpa
h+6Z4j6yNRt9iMiFeKlAr65YrF68hIic5ZnwsDJY9uT6r8PKQxoPOWcIBvddCH1O
8Na4Pqe1NQetPs3bHCCt/iDa7kUVzI5E6X1/QtWHopVW1b9Ih6r1
-----END RSA PRIVATE KEY-----
`

func TestGenerateDSNFromConfig(t *testing.T) {
	cfg := &ClientConfig{
		FullDsnEmpty: true,
		User:         "user",
		Pass:         "pass",
		Host:         "localhost",
		Port:         1234,
		Database:     "db",
	}

	want := "user:pass@tcp(localhost:1234)/db?multiStatements=true"
	got, err := cfg.MySQLConfig()

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("generated dsn unexpected, got: %q want: %q", got, want)
	}
}

func TestGenerateDSNFromConfigNoPass(t *testing.T) {
	cfg := &ClientConfig{
		FullDsnEmpty: true,
		User:         "user",
		Host:         "localhost",
		Port:         1234,
		Database:     "db",
	}

	want := "user@tcp(localhost:1234)/db?multiStatements=true"
	got, err := cfg.MySQLConfig()

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("generated dsn unexpected, got: %q want: %q", got, want)
	}
}

func TestGenerateDSNWithTLSFromConfig(t *testing.T) {
	cfg := &ClientConfig{
		CA:           ca,
		ClientKey:    clientKey,
		ClientCert:   clientCert,
		FullDsnEmpty: true,
		User:         "user",
		Pass:         "pass",
		Host:         "localhost",
		Port:         1234,
		Database:     "db",
	}

	want := "user:pass@tcp(localhost:1234)/db?multiStatements=true&tls=custom"
	got, err := cfg.MySQLConfig()

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("generated dsn unexpected, got: %q want: %q", got, want)
	}
}

func TestGenerateDSNWithDSN(t *testing.T) {
	cfg := &ClientConfig{
		FullDsn: "user:pass@tcp(localhost:1234)/db?multiStatements=true",
	}

	want := "user:pass@tcp(localhost:1234)/db?multiStatements=true"
	got, err := cfg.MySQLConfig()

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("generated dsn unexpected, got: %q want: %q", got, want)
	}
}

func TestGenerateDSNWithDSNWithTLS(t *testing.T) {
	cfg := &ClientConfig{
		FullDsn:    "user:pass@tcp(localhost:1234)/db?multiStatements=true&tls=custom",
		CA:         ca,
		ClientKey:  clientKey,
		ClientCert: clientCert,
		PortEmpty:  true,
	}

	want := "user:pass@tcp(localhost:1234)/db?multiStatements=true&tls=custom"
	got, err := cfg.MySQLConfig()

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("generated dsn unexpected, got: %q want: %q", got, want)
	}
}

func TestGenerateDSNWithUDS(t *testing.T) {
	cfg := &ClientConfig{
		FullDsnEmpty: true,
		Database:     "db",
		PortEmpty:    true,
		User:         "user",
		Pass:         "pass",
		SocketPath:   "/var/run/socket.sock",
	}

	want := "user:pass@unix(/var/run/socket.sock)/db?multiStatements=true"
	got, err := cfg.MySQLConfig()

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("generated dsn unexpected, got: %q want: %q", got, want)
	}
}

func TestGenerateDSNWithUDSWithTLS(t *testing.T) {
	cfg := &ClientConfig{
		FullDsnEmpty: true,
		Database:     "db",
		User:         "user",
		Pass:         "pass",
		PortEmpty:    true,
		SocketPath:   "/var/run/socket.sock",
		CA:           ca,
		ClientKey:    clientKey,
		ClientCert:   clientCert,
	}

	want := "user:pass@unix(/var/run/socket.sock)/db?multiStatements=true&tls=custom"
	got, err := cfg.MySQLConfig()

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("generated dsn unexpected, got: %q want: %q", got, want)
	}
}

func TestGenerateDSNWithUDSNoPass(t *testing.T) {
	cfg := &ClientConfig{
		FullDsnEmpty: true,
		Database:     "db",
		PortEmpty:    true,
		User:         "user",
		SocketPath:   "/var/run/socket.sock",
	}

	want := "user@unix(/var/run/socket.sock)/db?multiStatements=true"
	got, err := cfg.MySQLConfig()

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("generated dsn unexpected, got: %q; want: %q", got, want)
	}
}

func TestGenerateDSNWithUDSDefaults(t *testing.T) {
	cfg := &ClientConfig{
		FullDsnEmpty: true,
		Database:     "db",
		PortEmpty:    true,
		User:         "user",
	}

	want := "user@unix(/var/run/mysqld/mysqld.sock)/db?multiStatements=true"
	got, err := cfg.MySQLConfig()

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("generated dsn unexpected, got: %q; want: %q", got, want)
	}
}
