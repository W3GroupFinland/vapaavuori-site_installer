'use strict';

/**
 * @ngdoc service
 * @name webappApp.statusService
 * @description
 * # statusService
 * Factory in the webappApp.
 */

angular.module('webappApp')
	.factory('StatusService', ['$rootScope', function($rootScope) {
		// Initialize status messages.
		$rootScope.statusMessages = [];
		
		// We return this object to anything injecting our service
		var Service = {};

		// Refresh platform data.
		$rootScope.$on('STATUS_MESSAGE', function(_, args) {
		      $rootScope.statusMessages.push(args.Message);
		      $rootScope.$apply();
		});		
		
		// Get method for platforms
		Service.setMessage = function(msg) {
			$rootScope.statusMessages.push(msg);
		};

		return Service;
}]);
