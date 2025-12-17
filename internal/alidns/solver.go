package alidns

import (
	"errors"

	"github.com/404LifeFound/cert-manager-alidns-webhook/config"
	"github.com/404LifeFound/cert-manager-alidns-webhook/internal/utils"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
)

type AliDNSProviderSolver struct {
	AliDNS *AliDNS
}

type AliDNSProviderConfig struct{}

func (a *AliDNSProviderSolver) Name() string {
	return "alidns"
}

func (a *AliDNSProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	log.Info().Msgf("presenting txt record: %v %v", ch.ResolvedFQDN, ch.ResolvedZone)
	domainName, err := a.AliDNS.GetHostedDomain(ch.ResolvedZone)
	if err != nil {
		log.Error().Err(err).Msgf("domain %s not exist", ch.ResolvedZone)
		return err
	}

	rr := utils.ExtractRR(ch.ResolvedFQDN, domainName)

	if err := a.AliDNS.AddTxTRecord(domainName, rr, ch.Key); err != nil {
		log.Error().Err(err).Msgf("faild to add txt record: %v", ch.ResolvedFQDN)
		return err
	}

	log.Info().Msgf("present txt record %v sccess", ch.ResolvedFQDN)
	return nil
}

func (a *AliDNSProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	log.Info().Msgf("clean up txt record: %v %v", ch.ResolvedFQDN, ch.ResolvedZone)
	domainName, err := a.AliDNS.GetHostedDomain(ch.ResolvedZone)
	if err != nil {
		log.Error().Err(err).Msgf("domain %s not exist", ch.ResolvedZone)
		return err
	}

	rr := utils.ExtractRR(ch.ResolvedFQDN, domainName)
	record, err := a.AliDNS.GetTxTRecord(domainName, rr)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get text record %v.%v", rr, domainName)
		return err
	}

	if *record.Value != ch.Key {
		log.Error().Msgf("records value does not match: %v", ch.ResolvedFQDN)
		return errors.New("record value does not match")
	}

	if err := a.AliDNS.DeleteTxTRecord(*record.RecordId); err != nil {
		log.Error().Err(err).Msgf("failed to delete domain record %v", ch.ResolvedFQDN)
		return err
	}

	log.Info().Msgf("clean up txt record success: %v %v", ch.ResolvedFQDN, ch.ResolvedZone)
	return nil
}

func (a *AliDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	alidns_client, err := NewAliDNSClient(&config.GlobalConfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to setup alicloud dns client")
		return err
	}
	a.AliDNS = alidns_client
	return nil
}
