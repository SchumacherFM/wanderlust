/**
 * Route configuration for the Dashboard module.
 */
angular.module('Wanderlust')
  .config([
    '$stateProvider',
    '$urlRouterProvider',
    '$httpProvider',
    'picnicUrls',
    function ($stateProvider, $urlRouterProvider, $httpProvider, picnicUrls) {

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
        .state('brotzeit', {
          url: '/brotzeit',
          templateUrl: 'partials/brotzeit/tpl/bz.html',
          controller: 'BrotzeitController',
          data: {
            ncyBreadcrumbLabel: 'Brotzeit - The URL Cache'
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
          url: picnicUrls.provisioners + '{type:[a-z0-9]{3,20}}',
          controller: 'ProvisionerController',
          templateUrl: function (matchedParts) {
            return 'partials' + picnicUrls.provisioners + 'tpl/' + (matchedParts.type || '') + '.html';
          },
          data: {
            ncyBreadcrumbLabel: 'Provisioner / {{typeName}}'
          }
        });

      $httpProvider.interceptors.push('AuthInterceptor');
      $httpProvider.interceptors.push('ErrorInterceptor');

    }]);
