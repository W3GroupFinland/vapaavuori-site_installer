'use strict';

/**
 * @ngdoc function
 * @name webappApp.controller:NewSiteCtrl
 * @description
 * # NewSiteCtrl
 * Controller of the webappApp
 */
angular.module('webappApp')
	.controller('NewSiteCtrl', ['$scope', 'HostmasterService', '$routeParams', function ($scope, HostmasterService, $routeParams) {

		console.log($routeParams);
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
			$scope.templates = result;
		});

		// Get site template listing.
		HostmasterService.getServerTemplates().then(function (result) {
			$scope.httpTemplates = result;
			$scope.sslTemplates = result;
		});

		// Get site template listing.
		HostmasterService.getServerCertificates().then(function (result) {
			$scope.certificates = result;
			$scope.certKeys = result;
		});		

		$scope.http = {};
		$scope.http.aliases = [];
		
		$scope.ssl = {};
		$scope.ssl.aliases = [];	

		$scope.submit = function(site) {
			// Copy aliases to their right places.
			site.HttpServer.DomainAliases = $scope.http.aliases;
			site.SSLServer.DomainAliases = $scope.ssl.aliases;
			
			site.HttpServer.Port = parseInt(site.HttpServer.Port);
			site.SSLServer.Port = parseInt(site.SSLServer.Port);
			site.InstallInfo.PlatformId = parseInt(site.InstallInfo.PlatformId);
			site.InstallInfo.PlatformName = $scope.platformName;
			
			$scope.master = angular.copy(site);
			console.log('SITE', site);

			HostmasterService.registerFullSite(site).then(function (result) {
				console.log('RESULT', result);
			});
		};

		$scope.reset = function() {
		$scope.site = angular.copy($scope.master);
		};

		$scope.reset();		
		
		console.log($scope, HostmasterService);
}]);

