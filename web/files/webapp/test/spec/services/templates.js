'use strict';

describe('Service: Templates', function () {

  // load the service's module
  beforeEach(module('webappApp'));

  // instantiate service
  var Templates;
  beforeEach(inject(function (_Templates_) {
    Templates = _Templates_;
  }));

  it('should do something', function () {
    expect(!!Templates).toBe(true);
  });

});
