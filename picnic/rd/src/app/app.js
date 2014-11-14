'use strict';
angular
  .module('Wanderlust', [
    'ui.bootstrap',
    'ui.router',
    'LocalStorageModule',
    'ngResource',
    'ui.gravatar',
    'angular-growl',
    'picnic.services',
    'angulartics',
    'angulartics.piwik',
    'ncy-angular-breadcrumb'
  ])
  .constant('picnicUrls', {
    auth: '/auth/',
    users: '/users/',
    sysinfo: '/sysinfo/',
    provisioners: '/provisioners/',
    messages: '/api/messages'
  })
  .constant('AUTH_TOKEN_HEADER', 'X-Auth-Token')
  .constant('AUTH_TOKEN_STORAGE_KEY', 'WL_authToken')
  .config([
    '$resourceProvider',
    'gravatarServiceProvider',
    'growlProvider',
    function ($resourceProvider, gravatarServiceProvider, growlProvider) {
      // Don't strip trailing slashes from calculated URLs
      $resourceProvider.defaults.stripTrailingSlashes = false;
      gravatarServiceProvider.defaults = {
        size: 40,
        "default": 'monsterid'
      };

      // Use https endpoint
      gravatarServiceProvider.secure = true;

      growlProvider.globalTimeToLive({success: 1000, error: 2000, warning: 3000, info: 4000});

    }]);

if (!Array.isArray) {
  Array.isArray = function (arg) {
    return Object.prototype.toString.call(arg) === '[object Array]';
  };
}
