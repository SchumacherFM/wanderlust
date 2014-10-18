/**
 * Route configuration for the Dashboard module.
 */
angular.module('Wanderlust')
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
          templateUrl: 'partials/dashboard/tpl/dashboard.html',
          data: {
            ncyBreadcrumbLabel: 'Dashboard'
          }
        })
        .state('login', {
          url: '/login',
          templateUrl: 'partials/login/tpl/login.html',
          controller: 'LoginController',
          data: {
            ncyBreadcrumbLabel: 'Login'
          }
        })
        .state('tables', {
          url: '/tables',
          templateUrl: 'partials/core/tpl/tables.html',
          data: {
            ncyBreadcrumbLabel: 'Yet another demo table page'
          }
        })
        .state('shop', {
          url: '/shop',
          templateUrl: 'partials/marketplace/tpl/mp.html',
          controller: 'MarketplaceController',
          data: {
            ncyBreadcrumbLabel: 'Shop - Your in-app purchase made easy!'
          }
        })
        .state('privacy', {
          url: '/privacy',
          templateUrl: 'partials/core/tpl/privacy.html',
          data: {
            ncyBreadcrumbLabel: 'Privacy Statement'
          }
        })
        .state('provisioners', {
          url: '/provisioners/:type',
          templateUrl: function ($stateParams) {
            // 404 errors can occur when a template not exists
            var type = $stateParams.type || 'textarea';
            return 'partials/provisioner/tpl/' + type + '.html';
          },
          data: {
            ncyBreadcrumbLabel: 'Provisioner / {{name}}'
          }
        });

      $httpProvider.interceptors.push('AuthInterceptor');
      $httpProvider.interceptors.push('ErrorInterceptor');

    }]);
