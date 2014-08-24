'use strict';

/**
 * @ngdoc service
 * @name webappApp.statusService
 * @description
 * # statusService
 * Factory in the webappApp.
 */

angular.module('webappApp')
	.factory('RefreshService', ['$rootScope', function($rootScope) {
		// Initialize status messages.
		$rootScope.statusMessages = [];
		
		// We return this object to anything injecting our service
		var Service = {};

		var Data = {};

		Service.setData = function(data) {
			Data[data.Type] = data.Data;
			console.log(Data);
		};

		return Service;
}]);
