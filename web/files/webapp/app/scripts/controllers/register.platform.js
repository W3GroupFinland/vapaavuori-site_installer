'use strict';

/**
 * @ngdoc function
 * @name webappApp.controller:RegisterPlatformCtrl
 * @description
 * # RegisterPlatformCtrl
 * Controller of the webappApp
 */
angular.module('webappApp')
	.controller('RegisterPlatformCtrl', function($scope, close) {
		$scope.close = function(result) {
			close(result, 500); // close, but give 500ms for bootstrap to animate
		};
	});
