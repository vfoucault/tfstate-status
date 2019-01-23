package providers

import (
	"context"
	"errors"
	"github.com/vfoucault/tfstate-status/models"
	"io"
	"time"
)

type Provider interface {
	ListFiles() ([]ObjectFile, error)
	GetFile(string) (io.ReadCloser, error)
	SetPrefix(string)
	SetKind(string)
	GetKind() string
	ProcessState(fileObject ObjectFile) (*models.TfState, error)
}

type GenericProvider struct {
	Kind    string
	Context context.Context
	Prefix  string
}

func (c *GenericProvider) SetPrefix(prefix string) {
	c.Prefix = prefix
}

func (c *GenericProvider) SetKind(Kind string) {
	c.Kind = Kind
}

func (c *GenericProvider) GetKind() string {
	return c.Kind
}

func NewProviderFactory(providerStr string, containerName string, prefix string) (Provider, error) {
	switch providerStr {
	case "aws":
		p, err := NewAwsProvider(containerName)
		if err != nil {
			return nil, err
		}
		p.SetKind("aws")
		p.SetPrefix(prefix)
		return p, err
		//case "azure":
		//	err, provider := NewAzureProvider(os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_KEY"), containerName)
		//	if err != nil {
		//		return nil, err
		//	}
		//	provider.SetKind("azure")
		//	provider.SetPrefix(prefix)
		//	return provider, err
		//case "gcp":
		//	err, provider := NewGcpProvider(os.Getenv("GOOGLE_CLOUD_PROJECT"), containerName, prefix)
		//	if err != nil {
		//		return nil, err
		//	}
		//	provider.SetKind("gcp")
		//	provider.SetPrefix(prefix)
		//	return provider, nil
	}
	return nil, errors.New("unable to process ConfigFactory")
}

type ObjectFile struct {
	Key          string
	LastModified time.Time
}
