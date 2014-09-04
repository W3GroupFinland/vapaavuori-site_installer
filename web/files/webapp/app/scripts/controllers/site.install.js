'use strict';

/**
 * @ngdoc function
 * @name webappApp.controller:SiteInstallProcessCtrl
 * @description
 * # SiteInstallProcessCtrl
 * Controller of the webappApp
 */
angular.module('webappApp')
	.controller('SiteInstallProcessCtrl', ['$scope', '$rootScope', 'close', function($scope, $rootScope, close) {
		$scope.processMsg = [];

		$rootScope.$on('PROCESS_MESSAGE', function(_, args) {
			$scope.processMsg.push(args.Data.Message);
			$scope.$apply();
		});		

		$scope.close = function(result) {
			close(result, 500); // close, but give 500ms for bootstrap to animate
		};
	}]);
