'use strict';

/**
 * @ngdoc function
 * @name webappApp.controller:NewSiteCtrl
 * @description
 * # NewSiteCtrl
 * Controller of the webappApp
 */
angular.module('webappApp')
	.controller('NewSiteCtrl', ['$scope', 'HostmasterService', '$routeParams', 'ModalService', 
		function ($scope, HostmasterService, $routeParams, ModalService) {

		$scope.platformName = '';
		if ($routeParams.name !== undefined) {
			$scope.platformName = $routeParams.name;
		}

		$scope.master = {};
		$scope.templates = [];
		$scope.httpTemplates = [];
		$scope.sslTemplates = [];
		$scope.certificates = [];
		$scope.certKeys = [];

		$scope.site = {
			'InstallInfo' : {
				'SiteName' : '',
				'SubDirectory' : '',
				'TemplatePath' : '',
				'PlatformName' : '',
				'PlatformId' : null,
			},
			'HttpServer' : {
				'Template' : '',
				'Port' : null,
				'DomainInfo' : {
					'DomainName' : '',
					'Host' : '',
				},
				'DomainAliases' : [],
			},
			'SSLServer' : {
				'Template' : '',
				'Port' : null,
				'DomainInfo' : {
					'DomainName' : '',
					'Host' : '',
				},
				'DomainAliases' : [],
				'Certificate' : '',
				'Key' : '',
			},
		};

		// Get site template listing.
		HostmasterService.getSiteTemplates().then(function (result) {
			$scope.templates = result.Data;
		});

		// Get site template listing.
		HostmasterService.getServerTemplates().then(function (result) {
			$scope.httpTemplates = result.Data;
			$scope.sslTemplates = result.Data;
		});

		// Get site template listing.
		HostmasterService.getServerCertificates().then(function (result) {
			$scope.certificates = result.Data;
			$scope.certKeys = result.Data;
		});		

		$scope.http = {};
		$scope.http.aliases = [];
		
		$scope.ssl = {};
		$scope.ssl.aliases = [];
		$scope.errors = [];

		$scope.submit = function(site) {
			// Copy aliases to their right places.
			if (site.HttpServer !== undefined) {
				site.HttpServer.DomainAliases = $scope.http.aliases;
				site.HttpServer.Port = parseInt(site.HttpServer.Port);
			}
			if (site.SSLServer !== undefined) {
				site.SSLServer.DomainAliases = $scope.ssl.aliases;
				site.SSLServer.Port = parseInt(site.SSLServer.Port);
			}

			if (site.InstallInfo === undefined) {
				var error = {
					Key: 'Basic info',
					Error: 'Undefined',
				};
				$scope.errors = [error];
				return;
			}
			
			site.InstallInfo.PlatformId = parseInt(site.InstallInfo.PlatformId);
			site.InstallInfo.PlatformName = $scope.platformName;
			
			$scope.master = angular.copy(site);

			HostmasterService.registerFullSite(site).then(function (result) {
				if (result.Type === 'FORM_ERROR') {
					$scope.errors = result.Data;
					return;
				}

				if (result.Type === 'PROCESS_STARTING') {
					$scope.siteInstallModal();
				}
			});
		};

		$scope.siteInstallModal = function() {
	        ModalService.showModal({
	            templateUrl: 'views/partials/site.install.html',
	            controller: 'SiteInstallProcessCtrl'
	        }).then(function(modal) {
	            modal.element.modal();
	            modal.close.then(function() {
	            });
	        });			
		};

		$scope.reset = function() {
			$scope.site = angular.copy($scope.master);
		};

		$scope.reset();		
		
}]);

