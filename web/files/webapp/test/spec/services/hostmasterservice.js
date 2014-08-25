'use strict';

describe('Service: hostmasterService', function () {

  // load the service's module
  beforeEach(module('webappApp'));

  // instantiate service
  var hostmasterService;
  beforeEach(inject(function (_hostmasterService_) {
    hostmasterService = _hostmasterService_;
  }));

  it('should do something', function () {
    expect(!!hostmasterService).toBe(true);
  });

});
