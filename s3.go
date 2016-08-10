package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

type S3Broker struct {
	aws.Config
}

func (b *S3Broker) Provision(instanceid, serviceid, planid string) error {
	bucket := bucketNameFromInstanceId(instanceid)
	if err := b.createBucket(bucket); err != nil {
		return err
	}
	fmt.Printf("Creating service instance %s for service %s plan %s\n", instanceid, serviceid, planid)
	return nil
}

func (b *S3Broker) createBucket(name string) error {
	svc := s3.New(session.New(&b.Config))
	params := &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: b.Config.Region,
		},
	}

	if _, err := svc.CreateBucket(params); err != nil {
		return errors.Wrapf(err, "couldn't create s3 bucket: %q", name)
	}

	if err := svc.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: aws.String(name)}); err != nil {
		return errors.Wrapf(err, "failed to wait for bucket to exist %q", name)
	}

	// TODO(dfc) tag bucket with service data
	return nil
}

func bucketNameFromInstanceId(instanceid string) string {
	return envOr("BUCKET_NAME_PREFIX", "cloud-foundry-") + instanceid
}

func (b *S3Broker) Deprovision(instanceid, serviceid, planid string) error {
	bucket := bucketNameFromInstanceId(instanceid)
	if err := b.destroyBucket(bucket); err != nil {
		return err
	}
	fmt.Printf("Deleting service instance %s for service %s plan %s\n", instanceid, serviceid, planid)
	return nil
}

func (b *S3Broker) destroyBucket(name string) error {
	svc := s3.New(session.New(&b.Config))
	params := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}
	_, err := svc.DeleteBucket(params)
	return errors.Wrapf(err, "failed to remove bucket %q", name)
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
