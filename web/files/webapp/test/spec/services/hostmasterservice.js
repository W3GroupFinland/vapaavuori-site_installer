'use strict';

describe('Service: HostmasterService', function () {

  // load the service's module
  beforeEach(module('webappApp'));

  // instantiate service
  var HostmasterService;
  beforeEach(inject(function (_HostmasterService_) {
    HostmasterService = _HostmasterService_;
  }));

  it('should do something', function () {
    expect(!!HostmasterService.getPlatforms()).toBe(true);
  });

});
