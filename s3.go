package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

type S3Broker struct {
	aws.Config
}

func (b *S3Broker) Provision(instanceid, serviceid, planid string) error {
	bucket := bucketNameFromInstanceId(instanceid)
	svc := s3.New(session.New(&b.Config))
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: b.Config.Region,
		},
	})
	if err != nil {
		return errors.Wrapf(err, "couldn't create s3 bucket: %q", instanceid)
	}

	if err := svc.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: aws.String(bucket)}); err != nil {
		return errors.Wrapf(err, "failed to wait for bucket to exist %q", bucket)
	}

	// TODO(dfc) tag bucket with service data

	params := &iam.CreateGroupInput{
		GroupName: aws.String(envOr("GROUP_PATH", "/cloud-foundry/s3/")),
		Path:      aws.String(envOr("USER_PATH", "/cloud-foundry/s3/")),
	}
	resp, err := svc.CreateGroup(params)

	fmt.Printf("Creating service instance %s for service %s plan %s\n", instanceid, serviceid, planid)
	return nil
}

func bucketNameFromInstanceId(instanceid string) string {
	return envOr("BUCKET_NAME_PREFIX", "cloud-foundry-") + instanceid
}

func (b *S3Broker) Deprovision(instanceid, serviceid, planid string) error {
	bucket := bucketNameFromInstanceId(instanceid)
	svc := s3.New(session.New(&b.Config))
	params := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}
	_, err := svc.DeleteBucket(params)
	if err != nil {
		return errors.Wrapf(err, "failed to remove bucket %q", bucket)
	}
	fmt.Printf("Deleting service instance %s for service %s plan %s\n", instanceid, serviceid, planid)
	return nil
}

func (b *S3Broker) Bind(instanceid, bindingid, serviceid, planid string) error {
	fmt.Printf("Creating service binding %s for service %s plan %s instance %s\n",
		bindingid, serviceid, planid, instanceid)

	return nil
}

func (b *S3Broker) Unbind(instanceid, bindingid, serviceid, planid string) error {
	fmt.Printf("Delete service binding %s for service %s plan %s instance %s\n",
		bindingid, serviceid, planid, instanceid)
	return nil
}

func envOr(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
