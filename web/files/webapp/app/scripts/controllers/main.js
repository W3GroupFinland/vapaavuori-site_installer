'use strict';

/**
 * @ngdoc function
 * @name webappApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the webappApp
 */
angular.module('webappApp')
  .controller('MainCtrl', 
  	['$scope', 
  	'$rootScope', 
  	'HostmasterService', 
  	'ModalService', 
  	'Templates', 
  	'$route', 
  	'$routeParams',
  	function ($scope, $rootScope, HostmasterService, ModalService, Templates, $route, $routeParams) {

    var setSelectedPlatform = function(selected, data) {
    	return {  		
    		Selected: selected,
  			Data: data,
  		};
    };

  	// Initialize values.
  	$scope.platforms = {};
  	
  	// Set selected platform default values.
  	$scope.selectedPlatform = setSelectedPlatform(false, {});

  	$scope.showNewSiteForm = false;

  	// Set templates to scope.
  	$scope.templates = Templates.newTemplates();
  	$scope.templates.add('siteList', 'views/partials/selected.platform.html');
  	$scope.templates.add('siteDetails', 'views/partials/sitedetails.html');
  	$scope.templates.set('siteList');

  	// Initialize site details.
  	$scope.siteDetails = {};

    $scope.showSiteDetails = function(site) {
    	$scope.siteDetails = site;
    };
	
	// Get platform listing.
	HostmasterService.getPlatforms().then(function (result) {
		$scope.platforms = platformsByName(result.Data); // Set the result.
	  
	  	getSelectedPlatform();
	});

	// Refresh platform data.
	$scope.$on('PLATFORMS', function(_, args) {
	    $scope.platforms = platformsByName(args.Data);
		
	    getSelectedPlatform();
      	$scope.$apply();
	});

	var platformsByName = function(platforms) {
		var data = {};
		for (var i = 0; i < platforms.length; i++) {
			data[platforms[i].Name] = platforms[i];
		}

		return data;
	};

	var getSelectedPlatform = function() {
		if ($routeParams.name !== undefined) {
			$scope.selectedPlatform = setSelectedPlatform(true, $scope.platforms[$routeParams.name]);
		}		
	};

	$scope.platformSelected = function(platform) {
		if ($scope.selectPlatform.length === 0) {
			return false;
		}

		if (platform.Name === $scope.selectedPlatform.Data.Name) {
			return true;
		}

		return false;
	};

	$scope.selectPlatform = function(platform) {
		if (platform.Registered === false) {
			$scope.registerPlatformModal(platform);
		} else {
			$scope.selectedPlatform = setSelectedPlatform(true, platform);
		}

		$scope.templates.set('siteList');
	};

    $scope.registerPlatformModal = function(platform) {
        ModalService.showModal({
            templateUrl: 'views/partials/register.platform.html',
            controller: 'RegisterPlatformCtrl'
        }).then(function(modal) {
            modal.element.modal();
            modal.close.then(function(result) {
            	if (result === 1) {
            		HostmasterService.registerPlatform(platform);
            		$scope.selectedPlatform = setSelectedPlatform(true, platform);
            	}
            });
        });
    };

    $scope.platformRegistered = function(platform) {
    	return platform.Registered;
    };

  }]);
