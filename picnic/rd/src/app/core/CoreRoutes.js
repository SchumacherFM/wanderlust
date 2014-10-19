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
          url: '/provisioners/{type:[a-z0-9]{3,20}}',
          controller: 'ProvisionerController',
          templateUrl: 'partials/provisioner/tpl/form.html',
          data: {
            ncyBreadcrumbLabel: 'Provisioner / {{typeName}}'
          }
        });

      $httpProvider.interceptors.push('AuthInterceptor');
      $httpProvider.interceptors.push('ErrorInterceptor');

    }]);
