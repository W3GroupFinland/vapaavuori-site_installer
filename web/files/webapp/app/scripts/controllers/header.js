'use strict';

/**
 * @ngdoc function
 * @name webappApp.controller:HeaderCtrl
 * @description
 * # HeaderCtrl
 * Controller of the webappApp
 */
angular.module('webappApp')
	.controller('HeaderCtl', ['$scope', 'HostmasterService', function ($scope, HostmasterService) {

	$scope.user = null;

	// Get platform listing.
	HostmasterService.getUser().then(function (result) {
		$scope.user = result; // Set the result.
		console.log($scope.user);
	});  
}]);
