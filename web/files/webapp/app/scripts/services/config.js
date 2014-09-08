'use strict';

angular.module('services.config', [])
  .constant('configuration', {
    wsServer: 'wss://local.hostmaster.fi:8443/app/ws'
  });
