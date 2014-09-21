/**
 * Route configuration for the Dashboard module.
 */

angular.module('Dashboard').config([
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
      });

    $httpProvider.interceptors.push('ErrorInterceptor');

  }]);
