package object

import (
	"io"
	"time"
)

type Container struct {
	Name         string `json:"name"`
	Count        int64  `json:"count"`
	Bytes        int64  `json:"bytes"`
	LastModified string `json:"last_modified,omitempty"`
}

type ListContainersOutput struct {
	Containers []Container
}

type Object struct {
	Name         string    `json:"name"`
	Hash         string    `json:"hash"`
	Bytes        int64     `json:"bytes"`
	ContentType  string    `json:"content_type"`
	LastModified time.Time `json:"last_modified"`
}

type ListObjectsInput struct {
	Prefix    string
	Delimiter string
	Marker    string
	Limit     int
}

type ListObjectsOutput struct {
	Objects []Object
}

type ObjectMetadata struct {
	ContentType        string
	ContentDisposition string
	ContentEncoding    string
	Metadata           map[string]string
}

type PutObjectInput struct {
	Container   string
	ObjectName  string
	Body        io.Reader
	ContentType string
	Metadata    map[string]string
}

type PutObjectOutput struct {
	ETag         string
	LastModified time.Time
}

type GetObjectInput struct {
	Container  string
	ObjectName string
}

type GetObjectOutput struct {
	Body          io.ReadCloser
	ContentType   string
	ContentLength int64
	ETag          string
	LastModified  time.Time
	Metadata      map[string]string
}

type CopyObjectInput struct {
	SourceContainer       string
	SourceObjectName      string
	DestinationContainer  string
	DestinationObjectName string
}

type DeleteObjectInput struct {
	Container  string
	ObjectName string
}

type ContainerMetadata struct {
	ReadACL          string
	WriteACL         string
	SyncTo           string
	SyncKey          string
	VersionsLocation string
	ObjectCount      int64
	BytesUsed        int64
	CustomMetadata   map[string]string
}

type UpdateContainerInput struct {
	Name     string
	ReadACL  string
	WriteACL string
	Metadata map[string]string
}

type CreateContainerInput struct {
	Name     string
	ReadACL  string
	WriteACL string
	Metadata map[string]string
}
