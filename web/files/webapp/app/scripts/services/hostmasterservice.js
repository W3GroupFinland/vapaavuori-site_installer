'use strict';

/**
 * @ngdoc service
 * @name webappApp.HostmasterService
 * @description
 * # HostmasterService
 * Factory in the webappApp.
 */

angular.module('webappApp')
  .factory('HostmasterService', ['$q', '$rootScope', 'StatusService', function($q, $rootScope, StatusService) {

    // We return this object to anything injecting our service
    var Service = {};

    // Keep all pending requests here until they get responses
    var callbacks = {};

    // Create a unique callback ID to map requests to responses
    var currentCallbackId = 0;

    var ws = {};
    try {
      // Create our websocket object with the address to the websocket
      ws = new WebSocket('ws://localhost:8888/app/ws');
    } catch (err) {
      console.log(err);
      return err;
    }

    // Create boolean to check if connection is open.
    $rootScope.wsConnected = false;

    // On websocket connection open.
    ws.onopen = function(){
        $rootScope.wsConnected = true;
        StatusService.setMessage('Connected to hostmaster service.');
    };

    // On websocket connection close.
    ws.onclose = function() {
      $rootScope.wsConnected = false;
      StatusService.setMessage('Hostmaster service connection closed.');
    };

    // On websocket message.
    ws.onmessage = function(message) {
        var data = angular.fromJson(message.data);

        if (data.Refresh) {
          
          var obj = {
            Type: data.Type,
            Data: data.Data,
          };

          console.log(data);

          // TODO: Find some more elegant way to do this.
          $rootScope.$broadcast(data.Type, obj);
          return;
        }

        listener(data);
    };

    // Websocket message sender. 
    function sendRequest(request) {
      var defer = $q.defer();
      var callbackId = getCallbackId();
      callbacks[callbackId] = {
        time: new Date(),
        cb:defer
      };
      // Add callback id to request.
      request.CallbackId = callbackId;
      // Wait for connection. If connection is in ready state,
      // request will be sent immediately.
      waitSocketConn(ws, function() {
        ws.send(angular.toJson(request));
      }, 1, 1000*30);

      return defer.promise;
    }

    function listener(data) {
      var messageObj = data;
      // If an object exists with CallbackId in our callbacks object, resolve it.
      if(callbacks.hasOwnProperty(messageObj.CallbackId)) {

        var obj = {
          Type: messageObj.Type,
          Data: messageObj.Data,
        };

        $rootScope.$apply(callbacks[messageObj.CallbackId].cb.resolve(obj));
        delete callbacks[messageObj.callbackID];
      }
    }

    // This creates a new callback ID for a request
    function getCallbackId() {
      currentCallbackId += 1;
      if(currentCallbackId > 10000) {
        currentCallbackId = 0;
      }

      return currentCallbackId;
    }

    // Websocket connection waiter function.
    function waitSocketConn(socket, callback, frequency, timeout) {
      // Recursive wait function which runs callback when connected.
      function waitConn(counter) {
        setTimeout(
            function () {
              if ((frequency * counter) > timeout) {
                console.log('Websocket connection timeout.');
                return;
              }

              if (socket.readyState === 1) {
                  callback();
                  return;
              } else if (socket.readyState === 3) {
                console.log('Websocket connection closed or failed.');
                // If websocket connection is closed, return from wait loop.
                return;
              } else {
                  waitConn((counter+1));
              }

            }, frequency); // wait for x milliseconds for the connection...  
      }

      if (socket.readyState !== 1) {
        // Wait for socket connection before sending.
        waitConn(0);
      } else {
        // If socket is in ready state send message.
        callback();
      }
    }

    // Get method for platforms
    Service.getPlatforms = function() {
      var request = {
        Type: 'GET_PLATFORMS'
      };

      // Storing in a variable for clarity on what sendRequest returns
      var promise = sendRequest(request); 
      return promise;
    };

    // Get method for platforms
    Service.registerPlatform = function(platform) {
      var request = {
        Type: 'REGISTER_PLATFORM',
        Data: {
          Name: platform.Name,
        },
      };

      // Storing in a variable for clarity on what sendRequest returns
      var promise = sendRequest(request); 
      return promise;
    };

    // Get method for platforms
    Service.registerFullSite = function(installTemplate) {
      var request = {
        Type: 'REGISTER_FULL_SITE',
        Data: installTemplate,
      };

      // Storing in a variable for clarity on what sendRequest returns
      var promise = sendRequest(request); 
      return promise;
    };    

    // Get method for platforms
    Service.getSiteTemplates = function() {
      var request = {
        Type: 'GET_SITE_TEMPLATES',
      };

      // Storing in a variable for clarity on what sendRequest returns
      var promise = sendRequest(request); 
      return promise;
    };

    // Get method for platforms
    Service.getServerTemplates = function() {
      var request = {
        Type: 'GET_SERVER_TEMPLATES',
      };

      // Storing in a variable for clarity on what sendRequest returns
      var promise = sendRequest(request); 
      return promise;
    };

    // Get method for platforms
    Service.getServerCertificates = function() {
      var request = {
        Type: 'GET_SERVER_CERTIFICATES',
      };

      // Storing in a variable for clarity on what sendRequest returns
      var promise = sendRequest(request); 
      return promise;
    };    

     // Get method for user info
    Service.getUser = function() {
      var request = {
        Type: 'GET_USER'
      };

      // Storing in a variable for clarity on what sendRequest returns
      var promise = sendRequest(request); 
      return promise;
    };   

    return Service;
}]);
