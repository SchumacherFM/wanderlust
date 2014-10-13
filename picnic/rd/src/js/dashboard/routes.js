/**
 * Route configuration for the Dashboard module.
 */
angular.module('Dashboard')
  .config([
    '$stateProvider',
    '$urlRouterProvider',
    '$httpProvider',
    function ($stateProvider, $urlRouterProvider, $httpProvider) {

      // For unmatched routes
      $urlRouterProvider.otherwise('/');

      // Application routes
      $stateProvider
        .state('index', {
          url: '/',
          templateUrl: 'partials/dashboard.html'
        })
        .state('login', {
          url: '/login',
          templateUrl: 'partials/login.html',
          controller: 'LoginCtrl'
        })
        .state('tables', {
          url: '/tables',
          templateUrl: 'partials/tables.html'
        })
        .state('privacy', {
          url: '/privacy',
          templateUrl: 'partials/privacy.html'
        })
        .state('provisioners', {
          url: '/provisioners/:type',
          templateUrl: function ($stateParams) {
            // 404 errors can occur when a template not exists
            var type = $stateParams.type || 'textarea';
            return 'partials/provisioners/' + type + '.html';
          }
        })
        .state('shop', {
          url: '/shop',
          templateUrl: 'partials/shop.html'
        });

      $httpProvider.interceptors.push('AuthInterceptor');
      $httpProvider.interceptors.push('ErrorInterceptor');

    }]);
