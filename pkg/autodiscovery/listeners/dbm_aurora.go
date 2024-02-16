// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2020-present Datadog, Inc.

package listeners

import (
	"context"
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/autodiscovery/integration"
	"github.com/DataDog/datadog-agent/pkg/databasemonitoring/aws"
	"github.com/DataDog/datadog-agent/pkg/databasemonitoring/integrations"
	"github.com/DataDog/datadog-agent/pkg/util/containers"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"strconv"
	"sync"
	"time"
)

const dbmAdIdentifier = "database_monitoring_aurora"

func init() {
	Register(dbmAdIdentifier, NewDBMAuroraListener)
}

// DBMAuroraListener implements database-monitoring aurora discovery
type DBMAuroraListener struct {
	sync.RWMutex
	newService       chan<- Service
	delService       chan<- Service
	stop             chan bool
	services         map[string]Service
	config           integrations.AutodiscoverClustersConfig
	awsClients       map[string]aws.RDSClient // cached clients by region
	previousServices map[string]struct{}
}

var _ Service = &DBMAuroraService{}

// DBMAuroraService implements and store results from the Service interface for the DBMAuroraListener
type DBMAuroraService struct {
	adIdentifier string
	entityID     string
	checkName    string
	clusterId    string
	instance     *aws.Instance
}

func NewDBMAuroraListener(Config) (ServiceListener, error) {
	config, err := integrations.NewAutodiscoverClustersConfig()
	if err != nil {
		return nil, err
	}
	l := &DBMAuroraListener{
		stop:       make(chan bool),
		config:     config,
		services:   make(map[string]Service),
		awsClients: make(map[string]aws.RDSClient),
	}
	return l, nil
}

func newDBMAuroraListener(cfg Config, awsClients map[string]aws.RDSClient) (ServiceListener, error) {
	config, err := integrations.NewAutodiscoverClustersConfig()
	if err != nil {
		return nil, err
	}
	l := &DBMAuroraListener{
		config:     config,
		services:   make(map[string]Service),
		awsClients: awsClients,
	}
	return l, nil
}

// Listen listens for new and deleted aurora endpoints
func (l *DBMAuroraListener) Listen(newSvc, delSvc chan<- Service) {
	l.newService = newSvc
	l.delService = delSvc
	go l.discoverAuroraClusters()
}

// Stop stops the listener
func (l *DBMAuroraListener) Stop() {
	l.stop <- true
}

// discoverAuroraClusters discovers aurora clusters according to the configuration
func (l *DBMAuroraListener) discoverAuroraClusters() {
	discoveryTicker := time.NewTicker(time.Duration(l.config.DiscoveryInterval) * time.Second)
	for {
		for _, cluster := range l.config.Clusters {
			ids := make([]string, 0)
			ids = append(ids, cluster.ClusterIds...)
			if _, ok := l.awsClients[cluster.Region]; !ok {
				c, err := aws.NewRDSClient(cluster.Region, l.config.RoleArn)
				if err != nil {
					log.Errorf("error creating aws client for region %s: %s", cluster.Region, err)
					continue
				}
				l.awsClients[cluster.Region] = c
			}
			auroraCluster, err := l.awsClients[cluster.Region].GetAuroraClusterEndpoints(ids)
			if err != nil {
				log.Errorf("error discovering aurora cluster, skipping: %s", err)
				continue
			}
			discoveredServices := make(map[string]struct{})
			for id, c := range auroraCluster {
				for _, instance := range c.Instances {
					if instance == nil {
						continue
					}
					entityID := instance.Digest(string(cluster.Type), id)
					discoveredServices[entityID] = struct{}{}
					l.createService(entityID, string(cluster.Type), id, instance)
				}
			}
			// TODO: should we wait a certain number of run iterations before we remove instances?
			deletedServices := findDeletedServices(l.previousServices, discoveredServices)
			l.deleteServices(deletedServices)
			l.previousServices = discoveredServices
			select {
			case <-l.stop:
				return
			case <-discoveryTicker.C:
			}
		}
	}
}

func (l *DBMAuroraListener) createService(entityID, checkName, clusterId string, instance *aws.Instance) {
	l.Lock()
	defer l.Unlock()
	if _, present := l.services[entityID]; present {
		return
	}
	svc := &DBMAuroraService{
		adIdentifier: dbmAdIdentifier,
		entityID:     entityID,
		checkName:    checkName,
		instance:     instance,
		clusterId:    clusterId,
	}
	l.services[entityID] = svc
	l.newService <- svc
}

func (l *DBMAuroraListener) deleteServices(entityIDs []string) {
	l.Lock()
	defer l.Unlock()
	for _, entityID := range entityIDs {
		if svc, present := l.services[entityID]; present {
			l.delService <- svc
			delete(l.services, entityID)
		}
	}
}

func findDeletedServices(previousServices, discoveredServices map[string]struct{}) []string {
	deletedServices := make([]string, 0)

	for svc := range previousServices {
		if _, exists := discoveredServices[svc]; !exists {
			deletedServices = append(deletedServices, svc)
		}
	}

	return deletedServices
}

// GetServiceID returns the unique entity name linked to that service
func (d *DBMAuroraService) GetServiceID() string {
	return d.entityID
}

// GetTaggerEntity returns the tagger entity
func (d *DBMAuroraService) GetTaggerEntity() string {
	return d.entityID
}

// GetADIdentifiers return the single AD identifier for a static config service
func (d *DBMAuroraService) GetADIdentifiers(ctx context.Context) ([]string, error) {
	return []string{d.adIdentifier}, nil
}

// GetHosts returns the host for the aurora endpoint
func (d *DBMAuroraService) GetHosts(ctx context.Context) (map[string]string, error) {
	return map[string]string{"": d.instance.Endpoint}, nil
}

// GetPorts returns the port for the aurora endpoint
func (d *DBMAuroraService) GetPorts(ctx context.Context) ([]ContainerPort, error) {
	port := int(d.instance.Port)
	return []ContainerPort{{port, fmt.Sprintf("p%d", port)}}, nil
}

// GetTags returns the list of container tags - currently always empty
func (d *DBMAuroraService) GetTags() ([]string, error) {
	return []string{}, nil
}

// GetPid returns nil and an error because pids are currently not supported
func (d *DBMAuroraService) GetPid(ctx context.Context) (int, error) {
	return -1, ErrNotSupported
}

// GetHostname returns nothing - not supported
func (d *DBMAuroraService) GetHostname(ctx context.Context) (string, error) {
	return "", ErrNotSupported
}

func (d *DBMAuroraService) IsReady(ctx context.Context) bool {
	return true
}

// GetCheckNames returns nil
func (d *DBMAuroraService) GetCheckNames(context.Context) []string {
	return []string{d.checkName}
}

// HasFilter returns false on SNMP
//
//nolint:revive
func (d *DBMAuroraService) HasFilter(filter containers.FilterType) bool {
	return false
}

// GetExtraConfig parses the template variables with the extra_ prefix and returns the value
func (d *DBMAuroraService) GetExtraConfig(key string) (string, error) {
	switch key {
	case "region":
		return d.instance.Region, nil
	case "managed_authentication_enabled":
		return strconv.FormatBool(d.instance.IamEnabled), nil
	case "dbclusteridentifier":
		return d.clusterId, nil
	}
	return "", ErrNotSupported
}

// FilterTemplates does nothing.
//
//nolint:revive
func (d *DBMAuroraService) FilterTemplates(m map[string]integration.Config) {
}
