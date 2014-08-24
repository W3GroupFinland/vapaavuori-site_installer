'use strict';

describe('Service: useStrict', function () {

  // load the service's module
  beforeEach(module('webappApp'));

  // instantiate service
  var useStrict;
  beforeEach(inject(function (_useStrict_) {
    useStrict = _useStrict_;
  }));

  it('should do something', function () {
    expect(!!useStrict).toBe(true);
  });

});
