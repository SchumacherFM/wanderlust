/**
 * Master Controller
 */
angular
  .module('Dashboard')
  .controller('MasterCtrl', [
    '$scope',
    '$state',
    'localStorageService',
    '$timeout',
    'Session',
    'AuthResource',
    'Alert',
    function ($scope, $state, localStorageService, $timeout, Session, AuthResource, Alert) {
      'use strict';
      var LS_TOGGLE_KEY = 'wlToggle';
      //<Alerts>
      $scope.alert = Alert;
      $scope.$watchCollection('alert.messages', function () {
        $timeout(function () {
          //  @todo        $analytics.eventTrack('alert.messages', {  category: 'category' });
          Alert.dismissLast();
        }, 3000);
      });
      //</Alerts>

      //<Sessions>
      Session.init(AuthResource);
      $scope.session = Session;

      $scope.logout = function () {
        Session.logout().then(function () {
          $state.go('index');
        });
      };

      $scope.login = function () {
        Session.setLastLoginUrl();
        $state.go('login');
      };
      //</Sessions>

      /**
       * Sidebar Toggle & localStorageService Control
       */
      $scope.toggle = localStorageService.get(LS_TOGGLE_KEY) !== 'false';
      var mobileView = 992;
      $scope.getWidth = function () {
        return window.innerWidth;
      };
      $scope.$watch($scope.getWidth, function (newValue) {
        $scope.toggle = false;
        if (newValue >= mobileView) {
          $scope.toggle = true;
          if (localStorageService.get(LS_TOGGLE_KEY)) {
            $scope.toggle = localStorageService.get(LS_TOGGLE_KEY) !== 'false';
          }
        }
      });
      $scope.toggleSidebar = function () {
        $scope.toggle = !$scope.toggle;
        localStorageService.set(LS_TOGGLE_KEY, $scope.toggle);
      };
      window.onresize = function () {
        $scope.$apply();
      };
    }
  ])
  .controller('navigation', [
    '$scope',
    'ProvisionerResource',
    'Session',
    function ($scope, ProvisionerResource, Session) {
      'use strict';

      function loadProv() {
        ProvisionerResource.get({prov: ''}, function (result) {
          $scope.provisioners = [];
          $scope.provisioners = result.Collection || [];
        }, function (result) {
          $scope.provisioners = [];
          $scope.provisioners.push({
            Name: result.data,
            Url: "",
            Icon: "fa-exclamation-circle"
          });
        });
      }

      $scope.$watch(
        function () {
          return Session.isLoggedIn();
        },
        function (newValue, oldValue) {
          if (newValue === oldValue) {
            return;
          }
          loadProv();
        }
      );
      // when user logged in this may fire twice but with the $on listener, this will fire
      // only when the user is logged out. Is there a better solution?
      $scope.$on('$viewContentLoaded', function () {
        if (false === Session.isLoggedIn()) {
          loadProv();
        }
      });
    }
  ]);
