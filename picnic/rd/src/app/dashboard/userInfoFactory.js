angular
  .module('Wanderlust')

  // loads the user collection when the dashboard website is open.
  .factory('UserInfoResource', function ($resource, picnicUrls) {
    return $resource(picnicUrls.users, {});
  });
