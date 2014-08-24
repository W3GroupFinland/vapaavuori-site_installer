'use strict';

/**
 * @ngdoc function
 * @name webappApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the webappApp
 */
angular.module('webappApp')
  .controller('MainCtrl', ['$scope', '$rootScope', 'HostmasterService', 'ModalService', 
  	function ($scope, $rootScope, HostmasterService, ModalService) {
  	// Initialize values.
  	$scope.platforms = [];
  	$scope.selectedPlatform = [];
	
	// Get platform listing.
	HostmasterService.getPlatforms().then(function (result) {
		$scope.platforms = result; // Set the result.
	});

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
		console.log(platform);
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
            	}
            });
        });
    };

  }]);
