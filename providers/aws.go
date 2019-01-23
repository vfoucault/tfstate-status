package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vfoucault/tfstate-status/models"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type ProviderAws struct {
	GenericProvider
	Session    *session.Session
	S3Service  *s3.S3
	BucketName string
	Downloader *s3manager.Downloader
	Delimiter  string
}

func (c *ProviderAws) ListFiles() (fileList []ObjectFile, err error) {
	// TODO: Should be cleaned up, there is not error that could be examined.
	params := s3.ListObjectsInput{
		Bucket:    aws.String(c.BucketName),
		Prefix:    &c.Prefix,
		Delimiter: &c.Delimiter,
		MaxKeys:   aws.Int64(1000),
	}

	p := request.Pagination{
		NewRequest: func() (*request.Request, error) {
			req, _ := c.S3Service.ListObjectsRequest(&params)
			req.SetContext(c.Context)
			return req, nil
		},
	}
	for p.Next() {
		page := p.Page().(*s3.ListObjectsOutput)
		for _, obj := range page.Contents {
			if !strings.HasSuffix(*obj.Key, "/") {
				fileList = append(fileList, ObjectFile{Key: *obj.Key, LastModified: *obj.LastModified})
			}
		}
	}
	return
}

func (c *ProviderAws) GetFile(fileName string) (rc io.ReadCloser, err error) {
	result, err := c.S3Service.GetObjectWithContext(c.Context,
		&s3.GetObjectInput{
			Bucket: aws.String(c.BucketName),
			Key:    aws.String(fileName),
		}, func(r *request.Request) {
			// Comment out to have transport decode the object.
			//r.HTTPRequest.Header.Add("Accept-Encoding", "gzip")
		})
	if err != nil {
		return
	}
	rc = result.Body
	return
}

func NewAwsProvider(bucketName string) (Provider, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	provider := ProviderAws{}
	provider.Context = context.Background()
	provider.Session = sess
	provider.S3Service = svc
	provider.BucketName = bucketName
	provider.Downloader = s3manager.NewDownloader(sess)
	// Delimiter is very useful
	provider.Delimiter = ""
	return &provider, nil
}

func (c *ProviderAws) newWorkspace(filePath string) (*models.TerraformState, error) {
	var state = new(models.TerraformState)
	stream, err := c.GetFile(filePath)
	if err != nil {
		return new(models.TerraformState), err
	}
	jsonData, _ := ioutil.ReadAll(stream)
	json.Unmarshal(jsonData, &state)

	return state, nil
}

func (c *ProviderAws) ProcessState(fileObject ObjectFile) (*models.TfState, error) {
	var err error
	if filepath.Ext(fileObject.Key) != ".tfstate" {
		err := errors.New(fmt.Sprintf("prefix %v is not a state file. It Should have the .tfstate extention.", fileObject.Key))
		return new(models.TfState), err
	}
	paths := strings.Split(fileObject.Key, "/")
	fileName := paths[len(paths)-1]

	state := new(models.TfState)
	if paths[0] != "env:" {
		state.Workspace = "default"
	} else {
		state.Workspace = paths[1]
	}
	state.Name = strings.Split(fileName, ".")[0]
	state.LastModified = fileObject.LastModified
	state.FileName = fileName

	state.State, err = c.newWorkspace(fileObject.Key)
	return state, err
}
