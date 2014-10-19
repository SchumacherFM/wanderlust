angular
  .module('Wanderlust')

  // handles all the provisioners
  .factory('ProvisionerResource', function ($resource, picnicUrls) {
    return $resource(picnicUrls.provisioners + ':prov', {prov: '@prov'});
  });
