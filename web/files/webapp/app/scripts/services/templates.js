'use strict';

/**
 * @ngdoc service
 * @name webappApp.Templates
 * @description
 * # Templates
 * Service in the webappApp.
 */
angular.module('webappApp')
  .service('Templates', function Templates() {
    
    var Service = {};

    function NewTemplates() {}

    NewTemplates.prototype = {
    	Selected: '',
    	Templates: {},
    	add: function(name, url) {
    		this.Templates[name] = {Name: name, Url: url};
    	},
    	set: function(name) {
    		this.Selected = name;
    	},
    	selected: function() {
    		if (this.Templates[this.Selected] === undefined) {
    			return '';
    		}

    		return this.Templates[this.Selected].Url;
    	}
    };

    Service.newTemplates = function() {
    	return new NewTemplates();
    };

    return Service;
  });
