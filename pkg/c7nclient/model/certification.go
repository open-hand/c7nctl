package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
	"strconv"
	"time"
)

const baseFormat = "2006-01-02 15:04:05"

type jsonTime time.Time

func (this jsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("%s", time.Time(this).Format(baseFormat))
	return []byte(stamp), nil
}

func (this *jsonTime) UnmarshalJSON(b []byte) error {
	now, err := time.Parse(fmt.Sprintf("\"%s\"", baseFormat), string(b))
	*this = jsonTime(now)
	return err
}

type Certifications struct {
	Pages int             `json:"pages"`
	Size  int             `json:"size"`
	Total int             `json:"total"`
	List  []Certification `json:"list"`
}

type Certification struct {
	ID                         int      `json:"id"`
	OrganizationID             int      `json:"organizationId"`
	CertName                   string   `json:"certName"`
	CommonName                 string   `json:"commonName"`
	Domains                    []string `json:"domains"`
	Type                       string   `json:"type"`
	Status                     string   `json:"status"`
	ValidFrom                  string   `json:"validFrom"`
	ValidUntil                 string   `json:"validUntil"`
	EnvID                      int      `json:"envId"`
	EnvName                    string   `json:"envName"`
	EnvConnected               bool     `json:"envConnected"`
	CommandType                string   `json:"commandType"`
	CommandStatus              string   `json:"commandStatus"`
	Error                      string   `json:"error"`
	SkipCheckProjectPermission bool     `json:"skipCheckProjectPermission"`
}

type CertificationInfo struct {
	ID         int
	CertName   string
	CommonName string
	Domains    string
	ExpireDay  int
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

type GenericCerts struct {
	Pages int           `json:"pages"`
	Size  int           `json:"size"`
	Total int           `json:"total"`
	List  []GenericCert `json:"list"`
}

type GenericCert struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	KeyValue  string `json:"keyValue"`
	CertValue string `json:"certValue"`
	Domain    string `json:"domain"`
	Type      string `json:"type"`
}

type GenericCertInfo struct {
	Id     int
	Name   string
	Domain string
}

func PrintCertificationInfo(certificationInfo []CertificationInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Id", "GetName", "Domain", "ExpireDay")
	for _, r := range certificationInfo {
		table.AddRow(r.ID, r.CertName, r.Domains, strconv.Itoa(r.ExpireDay)+"天")
	}
	fmt.Fprintf(out, table.String())
}

func PrintGenericCertInfo(genericCertInfos []GenericCertInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Id", "GetName", "Domain")
	for _, r := range genericCertInfos {
		table.AddRow(r.Id, r.Name, r.Domain)
	}
	fmt.Fprintf(out, table.String())
}
