package nhncloud

import (
	"sync"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/compute"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncr"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncs"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/nks"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/iam"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/floatingip"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/loadbalancer"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/securitygroup"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/vpc"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/mariadb"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/mysql"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/postgresql"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/block"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/object"
)

type Client struct {
	config *Config
	mu     sync.Mutex

	iam           *iam.Client
	compute       *compute.Client
	mysqlClient   *mysql.Client
	mariadbClient *mariadb.Client
	pgClient      *postgresql.Client
	vpcClient     *vpc.Client
	sgClient      *securitygroup.Client
	fipClient     *floatingip.Client
	lbClient      *loadbalancer.Client
	blockClient   *block.Client
	objectClient  *object.Client
	nksClient     *nks.Client
	ncrClient     *ncr.Client
	ncsClient     *ncs.Client
}

func New(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, ErrCredentialsRequired
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &Client{config: cfg}, nil
}

func (c *Client) IAM() *iam.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.iam == nil {
		c.iam = iam.NewClient(c.config.Region, c.config.Credentials, c.config.httpClient(), c.config.Debug)
	}
	return c.iam
}

func (c *Client) Compute() *compute.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.compute == nil {
		c.compute = compute.NewClient(c.config.Region, c.config.IdentityCredentials, c.config.httpClient(), c.config.Debug)
	}
	return c.compute
}

func (c *Client) MySQL() *mysql.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.mysqlClient == nil {
		appKey := c.config.AppKeys["rds-mysql"]
		c.mysqlClient = mysql.NewClient(c.config.Region, appKey, c.config.Credentials, c.config.Debug)
	}
	return c.mysqlClient
}

func (c *Client) MariaDB() *mariadb.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.mariadbClient == nil {
		appKey := c.config.AppKeys["rds-mariadb"]
		c.mariadbClient = mariadb.NewClient(c.config.Region, appKey, c.config.Credentials, c.config.Debug)
	}
	return c.mariadbClient
}

func (c *Client) PostgreSQL() *postgresql.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.pgClient == nil {
		appKey := c.config.AppKeys["rds-postgresql"]
		c.pgClient = postgresql.NewClient(c.config.Region, appKey, c.config.Credentials, c.config.Debug)
	}
	return c.pgClient
}

func (c *Client) VPC() *vpc.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.vpcClient == nil {
		c.vpcClient = vpc.NewClient(c.config.Region, c.config.IdentityCredentials, c.config.httpClient(), c.config.Debug)
	}
	return c.vpcClient
}

func (c *Client) SecurityGroup() *securitygroup.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sgClient == nil {
		c.sgClient = securitygroup.NewClient(c.config.Region, c.config.IdentityCredentials, c.config.httpClient(), c.config.Debug)
	}
	return c.sgClient
}

func (c *Client) FloatingIP() *floatingip.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.fipClient == nil {
		c.fipClient = floatingip.NewClient(c.config.Region, c.config.IdentityCredentials, c.config.httpClient(), c.config.Debug)
	}
	return c.fipClient
}

func (c *Client) LoadBalancer() *loadbalancer.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lbClient == nil {
		c.lbClient = loadbalancer.NewClient(c.config.Region, c.config.IdentityCredentials, c.config.httpClient(), c.config.Debug)
	}
	return c.lbClient
}

func (c *Client) BlockStorage() *block.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.blockClient == nil {
		c.blockClient = block.NewClient(c.config.Region, c.config.IdentityCredentials, c.config.httpClient(), c.config.Debug)
	}
	return c.blockClient
}

func (c *Client) ObjectStorage() *object.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.objectClient == nil {
		c.objectClient = object.NewClient(c.config.Region, c.config.IdentityCredentials, c.config.httpClient(), c.config.Debug)
	}
	return c.objectClient
}

func (c *Client) NKS() *nks.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.nksClient == nil {
		c.nksClient = nks.NewClient(c.config.Region, c.config.IdentityCredentials, c.config.httpClient(), c.config.Debug)
	}
	return c.nksClient
}

func (c *Client) NCR() *ncr.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ncrClient == nil {
		appKey := c.config.AppKeys["ncr"]
		c.ncrClient = ncr.NewClient(c.config.Region, appKey, c.config.Credentials, c.config.httpClient(), c.config.Debug)
	}
	return c.ncrClient
}

func (c *Client) NCS() *ncs.Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ncsClient == nil {
		appKey := c.config.AppKeys["ncs"]
		c.ncsClient = ncs.NewClient(c.config.Region, appKey, c.config.Credentials, c.config.httpClient(), c.config.Debug)
	}
	return c.ncsClient
}
