package utils

import (
	"context"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/common/consts"
	std_errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/yidaqiang/go-chartmuseum/chartmuseum"
	helm_repo "helm.sh/helm/v3/pkg/repo"
	"regexp"
	"strings"
)

var client *chartmuseum.Client

func GetReleaseTag(repo, app, version string) (targetVersion string, err error) {
	if repo == "" {
		repo = consts.DefaultRepoUrl
	}
	url, path := matchChartRepo(repo)
	if client == nil {
		if client, err = chartmuseum.NewClient(url, nil); err != nil {
			return "", err
		}
	}

	charts := new(helm_repo.ChartVersions)
	if resp, err := client.ChartService.ListChartVersion(context.Background(), path, app, charts); err != nil {
		log.Debug(resp)
		return "", std_errors.WithMessage(err, fmt.Sprintf("Get Relesea %s version failed", app))
	}

	reg := regexp.MustCompile("^" + version + ".\\d+$")
	for _, c := range *charts {
		tagName := c.Version
		if reg.MatchString(tagName) {
			if targetVersion == "" {
				targetVersion = tagName
			}
			log.Debugf("%s version %s", app, targetVersion)
			if VersionOrdinal(targetVersion) < VersionOrdinal(tagName) {
				targetVersion = tagName
			}
		}
	}
	return targetVersion, nil
}

func matchChartRepo(repo string) (string, string) {
	spaceReg, _ := regexp.Compile(`/([a-zA-Z0-9]+)/`)

	idx := spaceReg.FindStringSubmatchIndex(repo)

	return repo[:idx[2]], strings.Trim(repo[idx[2]:], "/")
}

func VersionOrdinal(version string) string {
	// ISO/IEC 14651:2011
	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1
	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			panic("VersionOrdinal: invalid version")
		}
		vo = append(vo, b)
		vo[j]++
	}
	return string(vo)
}
