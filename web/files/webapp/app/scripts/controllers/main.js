'use strict';

/**
 * @ngdoc function
 * @name webappApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the webappApp
 */
angular.module('webappApp')
  .controller('MainCtrl', ['$scope', '$rootScope', 'HostmasterService', 'ModalService', '$route', '$routeParams',
  	function ($scope, $rootScope, HostmasterService, ModalService, $route, $routeParams) {
  	// Initialize values.
  	$scope.platforms = {};
  	$scope.selectedPlatform = {};
  	$scope.showNewSiteForm = false;
	
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
			$scope.selectedPlatform = $scope.platforms[$routeParams.name];
		}		
	};

	$scope.platformSelected = function(platform) {
		if ($scope.selectPlatform.length === 0) {
			return false;
		}

		if (platform.Name === $scope.selectedPlatform.Name) {
			return true;
		}

		return false;
	};

	$scope.selectPlatform = function(platform) {
		if (platform.Registered === false) {
			$scope.registerPlatformModal(platform);
		} else {
			$scope.selectedPlatform = platform;
		}
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
            		$scope.selectedPlatform = platform;
            	}
            });
        });
    };

    $scope.platformRegistered = function(platform) {
    	return platform.Registered;
    };

    $scope.addNewSite = function(platform) {
    	console.log(platform);
    	$scope.showNewSiteForm = true;
    };

  }]);
