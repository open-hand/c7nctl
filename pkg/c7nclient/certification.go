package c7nclient

import (
	"errors"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"io"
	"math"
	"net/url"
	"strings"
	"time"
)

func (c *C7NClient) ListGenericCert(out io.Writer, projectId int) {
	if projectId == 0 {
		return
	}

	paras := make(map[string]interface{})
	paras["page"] = 1
	paras["size"] = 10000

	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/certs/page_cert", projectId), paras, nil)
	if err != nil {
		fmt.Println("build request error")
	}
	var genericCerts = model.GenericCerts{}
	_, err = c.do(req, &genericCerts)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}

	genericCertInfos := []model.GenericCertInfo{}
	for _, genericCert := range genericCerts.List {
		genericCertInfo := model.GenericCertInfo{
			Id:     genericCert.Id,
			Name:   genericCert.Name,
			Domain: genericCert.Domain,
		}
		genericCertInfos = append(genericCertInfos, genericCertInfo)
	}
	model.PrintGenericCertInfo(genericCertInfos, out)
}

func (c *C7NClient) ListCert(out io.Writer, projectId int, envId int) {
	if projectId == 0 {
		return
	}

	paras := make(map[string]interface{})
	paras["page"] = 1
	paras["size"] = 10000
	paras["env_id"] = envId

	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/certifications/page_by_options", projectId), paras, nil)
	if err != nil {
		fmt.Println("build request error")
	}
	var certifications = model.Certifications{}
	_, err = c.do(req, &certifications)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}

	certificationInfos := []model.CertificationInfo{}
	now := time.Now()
	for _, certification := range certifications.List {
		expireDay := -1
		if certification.ValidUntil != "" {
			validUntil, _ := time.Parse(baseFormat, certification.ValidUntil)
			expireDay = int(math.Floor(validUntil.Sub(now).Seconds() / 3600 / 24))
		}
		certificationInfo := model.CertificationInfo{
			ID:         certification.ID,
			CertName:   certification.CertName,
			CommonName: certification.CommonName,
			Domains:    strings.Join(certification.Domains, ","),
			ExpireDay:  expireDay,
		}

		certificationInfos = append(certificationInfos, certificationInfo)
	}
	model.PrintCertificationInfo(certificationInfos, out)
}

func (c *C7NClient) CreateCert(out io.Writer, projectId int, data *url.Values) {
	if projectId == 0 {
		return
	}

	req, err := c.newRequestWithFormData("POST", fmt.Sprintf("/devops/v1/projects/%d/certifications", projectId), nil, data)
	if err != nil {
		fmt.Printf("build request error")
	}
	var result string
	_, err = c.doHandleString(req, &result)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	fmt.Printf("Successfully create Certification %s", (*data)["certName"])

}

func (c *C7NClient) GetCert(out io.Writer, projectId int, envId int, name string) (error error, result *model.Certification) {
	if projectId == 0 {
		return errors.New("you do not have the permission of the project!"), nil
	}
	if envId == 0 {
		return errors.New("you do not have the permission of the env!"), nil
	}
	paras := make(map[string]interface{})
	paras["env_id"] = envId
	paras["cert_name"] = name
	req, err := c.newRequest("GET", fmt.Sprintf("devops/v1/projects/%d/certifications//query_by_name", projectId), paras, nil)
	if err != nil {
		fmt.Printf("request build err:%v", err)
		return err, nil
	}
	devopsCertification := model.Certification{}
	_, err = c.do(req, &devopsCertification)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return err, nil
	}
	return err, &devopsCertification
}
