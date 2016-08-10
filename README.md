# dto-s3-broker

A simple ser

# Configuration

## Plans

The broker provides one plan, `basic`, which provides a single s3 bucket shared amongst all applications bound to it.

## Mandatory Configuration

The broker requires the following environment variables to be provided.

- `AUTH_USER`, the basic auth username that protects this broker's instance.
- `AUTH_PASS`, the basic auth password that protects this broker's instance.
- `AWS_ACCESS_KEY`, the access key this broker instance uses to request services from AWS.
- `AWS_SECRET_KEY`, the secret key this broker instance uses to request services from AWS.
- `AWS_REGION`, the region in which s3 buckets will be created.


