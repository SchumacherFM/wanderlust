angular
  .module('Dashboard', [
    'ui.bootstrap',
    'ui.router',
    'ngCookies',
    'ngResource',
    'ui.gravatar'
  ])
  .constant('picnicUrls', {
    auth: '/auth/',
    messages: '/api/messages'
  })
  .constant('AUTH_TOKEN_HEADER', 'X-Auth-Token')
  .constant('AUTH_TOKEN_STORAGE_KEY', 'WL_authToken');

angular.module('ui.gravatar').config(
  [
    'gravatarServiceProvider',
    function (gravatarServiceProvider) {
      gravatarServiceProvider.defaults = {
        size: 40,
        "default": 'monsterid'
      };

      // Use https endpoint
      gravatarServiceProvider.secure = true;
    }
  ]
);
