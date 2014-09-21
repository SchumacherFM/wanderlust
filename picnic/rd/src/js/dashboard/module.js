angular
  .module('Dashboard', [
    'ui.bootstrap',
    'ui.router',
    'ngCookies',
    'ngResource',
    'ui.gravatar',
    'ui.bootstrap',
    'picnic.services'
  ])
  .constant('picnicUrls', {
    auth: '/auth/',
    users: '/users/',
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
