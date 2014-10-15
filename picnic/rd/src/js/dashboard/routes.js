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
          templateUrl: 'partials/dashboard.html',
          data: {
            ncyBreadcrumbLabel: 'Dashboard'
          }
        })
        .state('login', {
          url: '/login',
          templateUrl: 'partials/login.html',
          controller: 'LoginCtrl',
          data: {
            ncyBreadcrumbLabel: 'Login'
          }
        })
        .state('tables', {
          url: '/tables',
          templateUrl: 'partials/tables.html',
          data: {
            ncyBreadcrumbLabel: 'Yet another demo table page'
          }
        })
        .state('privacy', {
          url: '/privacy',
          templateUrl: 'partials/privacy.html',
          data: {
            ncyBreadcrumbLabel: 'Privacy Statement'
          }
        })
        .state('provisioners', {
          url: '/provisioners/:type',
          templateUrl: function ($stateParams) {
            // 404 errors can occur when a template not exists
            var type = $stateParams.type || 'textarea';
            return 'partials/provisioners/' + type + '.html';
          },
          data: {
            ncyBreadcrumbLabel: 'Provisioner / {{name}}'
          }
        })
        .state('shop', {
          url: '/shop',
          templateUrl: 'partials/shop.html'
        });

      $httpProvider.interceptors.push('AuthInterceptor');
      $httpProvider.interceptors.push('ErrorInterceptor');

    }]);
