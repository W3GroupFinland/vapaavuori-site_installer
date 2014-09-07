package controllers

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"log"
)

func (c *HostmasterWS) UpdatePlatforms() error {
	// Populate platform list.
	platforms, err := c.System.GetDrupalPlatforms()
	if err != nil {
		log.Fatal("Error getting drupal platforms in initialization.")
	}

	for _, platform := range platforms.List {
		err := c.populatePlatformSites(platform)
		if err != nil {
			return err
		}
	}

	c.Platforms = platforms

	return nil
}

func (c *HostmasterWS) populatePlatformSites(platform *models.PlatformInfo) error {
	sites, err := c.HostMasterDB.GetPlatformSites(platform.PlatformId)
	if err != nil {
		return err
	}

	for _, site := range sites {
		// Get the configs from database.
		configs, err := c.HostMasterDB.GetSiteServerConfigs(site.Id)
		if err != nil {
			return err
		}

		// Append them to site config slice.
		site.AddConfigs(configs...)

		// Get site domains from database.
		domains, err := c.HostMasterDB.GetSiteDomains(site.Id)
		if err != nil {
			return err
		}

		// Append them to site domain slice.
		site.AddDomains(domains...)

		// Map site domains to site configs.
		site.MapDomainsToConfigs(site.Domains)

		// Finally add site to platform.
		platform.AddSite(site)
	}

	return nil
}
