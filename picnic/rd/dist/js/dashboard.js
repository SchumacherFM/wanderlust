(function(){ 
angular.module('Dashboard', ['ui.bootstrap', 'ui.router', 'ngCookies']);
'use strict';

/**
 * Route configuration for the Dashboard module.
 */
angular.module('Dashboard').config(['$stateProvider', '$urlRouterProvider',
    function ($stateProvider, $urlRouterProvider) {

        // For unmatched routes
        $urlRouterProvider.otherwise('/');

        // Application routes
        $stateProvider
            .state('index', {
                url: '/',
                templateUrl: 'dashboard.html'
            })
            .state('tables', {
                url: '/tables',
                templateUrl: 'tables.html'
            });
    }]);

/**
 * Master Controller
 */
angular
  .module('Dashboard')
  .controller(
  'MasterCtrl',
  [
    '$scope',
    '$cookieStore',
    function ($scope, $cookieStore) {
      /**
       * Sidebar Toggle & Cookie Control
       */
      var mobileView = 992;

      $scope.getWidth = function () {
        return window.innerWidth;
      };

      $scope.$watch($scope.getWidth, function (newValue, oldValue) {
        if (newValue >= mobileView) {
          if (angular.isDefined($cookieStore.get('toggle'))) {
            $scope.toggle = !$cookieStore.get('toggle');
          }
          else {
            $scope.toggle = true;
          }
        }
        else {
          $scope.toggle = false;
        }

      });

      $scope.toggleSidebar = function () {
        $scope.toggle = !$scope.toggle;

        $cookieStore.put('toggle', $scope.toggle);
      };

      window.onresize = function () {
        $scope.$apply();
      };
    }
  ]
);

/**
 * Alerts Controller
 */
angular
  .module('Dashboard')
  .controller('AlertsCtrl', [
    '$scope',
    function ($scope) {
      $scope.alerts = [
        {type: 'success', msg: 'Thanks for visiting! Feel free to create pull requests to improve the dashboard!'},
        {type: 'danger', msg: 'Found a bug? Create an issue with as many details as you can.'}
      ];

      $scope.addAlert = function () {
        $scope.alerts.push({msg: 'Another alert!'});
      };

      $scope.closeAlert = function (index) {
        $scope.alerts.splice(index, 1);
      };
    }
  ]);

/**
 * Loading Directive
 * @see http://tobiasahlin.com/spinkit/
 */
angular
  .module('Dashboard')
  .directive('rdLoading', function () {
    return {
      restrict: 'AE',
      template: '<div class="loading"><div class="double-bounce1"></div><div class="double-bounce2"></div></div>'
    };
  }
);


})();