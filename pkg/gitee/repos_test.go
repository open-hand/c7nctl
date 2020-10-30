package gitee

import (
	"context"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestRepositoriesService_ListTags(t *testing.T) {
	client := NewClient(nil)
	tags, resp, err := client.Repositories.ListTags(context.Background(), "open-hand", "devops-service",
		&ListOptions{AccessToken: "b6fed6c7bad34220693300b936f7c8c4"})
	if err != nil {
		log.Error(resp)
		log.Panic(err)
	}
	for _, tag := range tags {
		println(*tag.Name)
	}
}
