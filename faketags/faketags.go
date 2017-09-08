package faketags

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"strings"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

type Faketags struct {
	namespace string
	sess *session.Session
}

const DefaultNamespace = "faketag"

func New(sess *session.Session) Faketags {
	return NewWithNamespace(sess, DefaultNamespace)
}

func NewWithNamespace(sess *session.Session, namespace string) Faketags {
	return Faketags{
		namespace: namespace,
		sess: sess,
	}
}

func (f *Faketags) TagsForId(id string) (map[string]string, error) {
	ssmApi := ssm.New(f.sess)

	paramName := f.paramName(id)

	resp, err := ssmApi.ListTagsForResource(&ssm.ListTagsForResourceInput{
		ResourceId: aws.String(paramName),
		ResourceType: aws.String(ssm.ResourceTypeForTaggingParameter),
	})
	if err != nil { return nil, err }

	tagMap := map[string]string{}
	for _, tag := range resp.TagList {
		tagMap[*tag.Key] = *tag.Value
	}

	return tagMap, nil
}

func (f *Faketags) paramName(name string) string {
	return fmt.Sprintf("/%s/%s", f.namespace, name) // TODO: make name SSM-safe
}

func (f *Faketags) IdsForTags(tags map[string]string) (map[string]map[string]string, error) {
	tagApi := resourcegroupstaggingapi.New(f.sess)

	filters := []*resourcegroupstaggingapi.TagFilter{}

	for key, val := range tags {
		filters = append(filters, &resourcegroupstaggingapi.TagFilter{
			Key: &key,
			Values: aws.StringSlice([]string{val}),
		})
	}

	input := &resourcegroupstaggingapi.GetResourcesInput{
		TagFilters: filters,
		ResourceTypeFilters: aws.StringSlice([]string{"ssm:parameter"}),
	}

	results := map[string]map[string]string{}

	tagApi.GetResourcesPages(input, func(page *resourcegroupstaggingapi.GetResourcesOutput, lastPage bool) bool {
		for _, resMap := range page.ResourceTagMappingList {
			resTags := map[string]string{}
			for _, tag := range resMap.Tags {
				resTags[*tag.Key] = *tag.Value
			}

			arn := *resMap.ResourceARN

			parts := strings.SplitN(arn, fmt.Sprintf("parameter/%s/", f.namespace), 2)
			id := parts[1]
			results[id] = resTags
		}
		//return !lastPage
		return false // TODO there appears to be a bug in the aws sdk fixme
	})

	return results, nil
}

func (f *Faketags) putParameter(name string) (string, error) {
	ssmApi := ssm.New(f.sess)

	paramName := f.paramName(name)

	_, err := ssmApi.PutParameter(&ssm.PutParameterInput{
		Name: &paramName,
		Type: aws.String(ssm.ParameterTypeString),
		Value: aws.String("faketag"),
		Description: aws.String("faketag"),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == ssm.ErrCodeParameterAlreadyExists {
				return paramName, nil
			}
		}
		return "", err
	}

	return paramName, nil
}

func (f *Faketags) PutTagsForId(id string, tags map[string]string) error {
	ssmApi := ssm.New(f.sess)

	paramName, err := f.putParameter(id)
	if err != nil { return err }

	ssmTags := []*ssm.Tag{}
	for key, val := range tags {
		ssmTags = append(ssmTags, &ssm.Tag{
			Key: aws.String(key),
			Value: aws.String(val),
		})
	}

	_, err = ssmApi.AddTagsToResource(&ssm.AddTagsToResourceInput{
		ResourceType: aws.String(ssm.ResourceTypeForTaggingParameter),
		ResourceId: &paramName,
		Tags: ssmTags,
	})

	return err
}
