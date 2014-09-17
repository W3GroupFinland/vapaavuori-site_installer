'use strict';

angular.module('services.config', [])
  .constant('configuration', {
    wsServer: 'wss://tivia-hostmaster.dyndns.org:8443/app/ws'
  });
