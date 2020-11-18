package utils

import (
	"context"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/common/consts"
	"github.com/choerodon/c7nctl/pkg/gitee"
	std_errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"regexp"
)

func GetReleaseTag(app, version string) (targetVersion string, err error) {
	client := gitee.NewClient(nil)

	tags, resp, err := client.Repositories.ListTags(context.Background(), "open-hand", app, &gitee.ListOptions{AccessToken: consts.DefaultGiteeAccessToken})
	if err != nil {
		log.Debug(resp)
		return "", std_errors.WithMessage(err, fmt.Sprintf("Get Relesea %s version failed", app))
	}

	reg := regexp.MustCompile("^" + version + ".\\d+$")
	for _, tag := range tags {
		tagName := *tag.Name
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
