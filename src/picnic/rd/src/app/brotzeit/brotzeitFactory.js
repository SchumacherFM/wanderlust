angular
  .module('Wanderlust')
  .factory('BrotzeitResource', [
    '$resource',
    'picnicUrls',
    function ($resource, picnicUrls) {
      return $resource(picnicUrls.brotzeit);
    }
  ]);
