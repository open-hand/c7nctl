package model

import "time"

type Certifications struct {
	TotalPages       int             `json:"totalPages"`
	TotalElements    int             `json:"totalElements"`
	NumberOfElements int             `json:"numberOfElements"`
	Size             int             `json:"size"`
	Number           int             `json:"number"`
	Content          []Certification `json:"content"`
	Empty            bool            `json:"empty"`
}

type Certification struct {
	ID                         int       `json:"id"`
	OrganizationID             int       `json:"organizationId"`
	CertName                   string    `json:"certName"`
	CommonName                 string    `json:"commonName"`
	Domains                    []string  `json:"domains"`
	Type                       string    `json:"type"`
	Status                     string    `json:"status"`
	ValidFrom                  time.Time `json:"validFrom"`
	ValidUntil                 time.Time `json:"validUntil"`
	EnvID                      int       `json:"envId"`
	EnvName                    string    `json:"envName"`
	EnvConnected               bool      `json:"envConnected"`
	CommandType                string    `json:"commandType"`
	CommandStatus              string    `json:"commandStatus"`
	Error                      string    `json:"error"`
	SkipCheckProjectPermission bool      `json:"skipCheckProjectPermission"`
}

type CertificationPostInfo struct {
	Domains   []string `json:"domains"`
	EnvID     int      `json:"envId"`
	CertName  string   `json:"certName"`
	Type      string   `json:"type"`
	CertValue string   `json:"certValue"`
	KeyValue  string   `json:"keyValue"`
}

type Certificate struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       CertSpec `json:"spec"`
}

type CertSpec struct {
	CommonName string            `json:"commonName"`
	DnsNames   []string          `json:"dnsNames"`
	IssuerRef  map[string]string `json:"issuerRef"`
	Acme       CertAcme          `json:"acme"`
	ExistCert  ExistCert         `json:"existCert"`
}

type CertAcme struct {
	Config []CertConfig `json:"config"`
}

type ExistCert struct {
	Key  string `json:"key"`
	Cert string `json:"cert"`
}

type CertConfig struct {
	Http01  map[string]string `json:"http01"`
	Domains []string          `json:"domains"`
}
