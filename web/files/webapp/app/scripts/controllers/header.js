'use strict';

/**
 * @ngdoc function
 * @name webappApp.controller:HeaderCtrl
 * @description
 * # HeaderCtrl
 * Controller of the webappApp
 */
angular.module('webappApp')
	.controller('HeaderCtl', ['$scope', '$location', 'HostmasterService', 
		function ($scope, $location, HostmasterService) {

	$scope.user = null;

	// Get platform listing.
	HostmasterService.getUser().then(function (result) {
		$scope.user = result; // Set the result.
	});

    $scope.navClass = function (page) {
        var currentRoute = $location.path().substring(1) || 'platforms';
        return page === currentRoute ? 'active' : '';
    };
}]);
