'use strict';
angular
  .module('Wanderlust', [
    'ui.bootstrap',
    'ui.router',
    'LocalStorageModule',
    'ngResource',
    'ui.gravatar',
    'ui.bootstrap',
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
    function ($resourceProvider, gravatarServiceProvider) {
      // Don't strip trailing slashes from calculated URLs
      $resourceProvider.defaults.stripTrailingSlashes = false;
      gravatarServiceProvider.defaults = {
        size: 40,
        "default": 'monsterid'
      };

      // Use https endpoint
      gravatarServiceProvider.secure = true;
    }]);