/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/404LifeFound/cert-manager-alidns-webhook/config"
	"github.com/404LifeFound/cert-manager-alidns-webhook/internal/alidns"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := config.LoadGlobalConfig(); err != nil {
		panic(err)
	}
	log.Info().Msg("starting cert-manager-alidns-webhook")
	cmd.RunWebhookServer(config.GlobalConfig.GroupName, &alidns.AliDNSProviderSolver{})
}
