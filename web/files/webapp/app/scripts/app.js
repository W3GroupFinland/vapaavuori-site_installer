'use strict';

/**
 * @ngdoc overview
 * @name webappApp
 * @description
 * # webappApp
 *
 * Main module of the application.
 */
angular
	.module('webappApp', [
		'ngAnimate',
		'ngCookies',
		'ngResource',
		'ngRoute',
		'ngSanitize',
		'ngTouch',
		'angularModalService',
		'services.config',
	])
	.config(function ($routeProvider) {
		$routeProvider
		.when('/', {
			templateUrl: 'views/main.html',
			controller: 'MainCtrl',
		})
		.when('/platform/:name', {
			templateUrl: 'views/main.html',
			controller: 'MainCtrl',
		})
		.when('/platform/:name/site/:index', {
			templateUrl: 'views/main.html',
			controller: 'MainCtrl',
		})				
		.when('/platform/new-site/:name', {
			templateUrl: 'views/add-new-site.html',
			controller: 'MainCtrl',
		})
		.when('/about', {
			templateUrl: 'views/about.html',
			controller: 'AboutCtrl',
		})
		.when('/contact', {
			templateUrl: 'views/contact.html',
			controller: 'ContactCtrl',
		})		
		.otherwise({
			redirectTo: '/'
		});
	});