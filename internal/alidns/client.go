package alidns

import (
	"fmt"

	"github.com/404LifeFound/cert-manager-alidns-webhook/config"
	"github.com/404LifeFound/cert-manager-alidns-webhook/internal/utils"
	alidns "github.com/alibabacloud-go/alidns-20150109/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	tea_utils "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/rs/zerolog/log"
)

type AliDNS struct {
	Client *alidns.Client
}

func NewAliDNSClient(alidns_config *config.Config) (*AliDNS, error) {
	// default credential chain
	credential, err := credentials.NewCredential(nil)
	if err != nil {
		log.Error().Err(err).Msg("New alicloud credential failed")
		return nil, err
	}

	config := &openapi.Config{
		Credential: credential,
		RegionId:   &alidns_config.AliDNS.Region,
	}

	client, err := alidns.NewClient(config)
	if err != nil {
		log.Error().Err(err).Msg("New alicloud dns client failed")
		return nil, err
	}

	return &AliDNS{
		Client: client,
	}, nil
}

func StringPtr(s string) *string { return &s }
func Int64Ptr(i int64) *int64    { return &i }

func (a *AliDNS) GetTxTRecord(domain, rr string) (*alidns.DescribeDomainRecordsResponseBodyDomainRecordsRecord, error) {
	request := &alidns.DescribeDomainRecordsRequest{
		DomainName: StringPtr(domain),
		RRKeyWord:  StringPtr(rr),
		Type:       StringPtr("TXT"),
	}

	response, err := a.Client.DescribeDomainRecordsWithOptions(request, &tea_utils.RuntimeOptions{})
	if err != nil {
		log.Error().Err(err).Msgf("failed to describe record: %s.%s", rr, domain)
		return nil, err
	}

	records := response.Body.GetDomainRecords()
	for _, r := range records.Record {
		if *r.RR == rr {
			return r, nil
		}
	}

	return nil, fmt.Errorf("txt record does not exist: %s.%s", rr, domain)
}

func (a *AliDNS) AddTxTRecord(domain, rr, value string) error {
	request := &alidns.AddDomainRecordRequest{
		DomainName: StringPtr(domain),
		RR:         StringPtr(rr),
		Type:       StringPtr("TXT"),
		Value:      StringPtr(value),
		TTL:        Int64Ptr(600),
	}

	_, err := a.Client.AddDomainRecordWithOptions(request, &tea_utils.RuntimeOptions{})
	if err != nil {
		log.Error().Err(err).Msgf("failed to create TXT record: %s.%s", rr, domain)
		return err
	}
	return nil
}

func (a *AliDNS) DeleteTxTRecord(id string) error {
	request := &alidns.DeleteDomainRecordRequest{RecordId: StringPtr(id)}

	_, err := a.Client.DeleteDomainRecordWithOptions(request, &tea_utils.RuntimeOptions{})
	if err != nil {
		log.Error().Err(err).Msgf("failed to delete txt record which id is: %s", id)
		return err
	}

	return nil
}

func (a *AliDNS) GetHostedDomain(domain string) (string, error) {
	request := &alidns.DescribeDomainsRequest{
		KeyWord:    StringPtr(utils.UnFqdn(domain)),
		SearchMode: StringPtr("EXACT"),
	}

	response, err := a.Client.DescribeDomainsWithOptions(request, &tea_utils.RuntimeOptions{})
	if err != nil {
		log.Error().Err(err).Msgf("failed to describe domain: %s", domain)
		return "", err
	}

	domains := response.Body.GetDomains()

	if len(domains.Domain) == 0 {
		log.Error().Msgf("domain %s does not exist", domain)
	}

	return *domains.Domain[0].DomainName, nil
}
