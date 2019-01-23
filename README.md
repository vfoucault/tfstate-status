# Terraform TFState parser

The main objective of this tool is to list and report TFStates storaes on a remote storage (S3 Only so far)

## Install

```
go get "github.com/vfoucault/tfstate-status"
```

## Usage

```
Usage of ./tfstate-state:
  -container string
        The ContainerName
  -filter-workspace string
        Filter workspace by name
  -list-empty
        List empty states only  
  -prefix string
        Prefix
  -provider string
        s3 so far
  -threads int
        Number of threads. Default to cores count
  -verbose
        Be Verbose

```

## Example

let's say that my bucket that store all my tfstates is `s3://tfstates`

```
~# AWS_REGION=eu-west-1 tfstate-state -provider aws -container tfstates 
+--------------------------------------+----------------------+-------------------------------+-----------+
|                 FILE                 |      WORKSPACE       |         LAST MODIFIED         | RESOURCES |
+--------------------------------------+----------------------+-------------------------------+-----------+
| awsconfig.tfstate                    | demo                 | 2019-01-18 14:40:43 +0000 UTC |        31 |
| emr-demo-env.tfstate                 | staging              | 2018-12-18 12:42:45 +0000 UTC |        38 |
| emr-demo-env.tfstate                 | production           | 2018-12-19 17:21:43 +0000 UTC |        38 |
| kafka-production.tfstate             | production           | 2018-12-17 09:19:37 +0000 UTC |        42 |
+--------------------------------------+----------------------+-------------------------------+-----------+
```