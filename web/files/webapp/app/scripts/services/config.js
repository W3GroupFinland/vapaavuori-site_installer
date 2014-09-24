'use strict';

angular.module('services.config', [])
  .constant('configuration', {
    wsServer: 'wss://hostmaster.tivia.fi:8443/app/ws'
  });
