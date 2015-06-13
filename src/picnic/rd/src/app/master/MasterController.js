/**
 * Master Controller
 */
angular
  .module('Wanderlust')
  .controller('MasterCtrl', [
    '$scope',
    '$state',
    'localStorageService',
    'Session',
    'AuthResource',
    function ($scope, $state, localStorageService, Session, AuthResource) {
      'use strict';
      var LS_TOGGLE_KEY = 'wlToggle';

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
  .controller('NavigationController', [
    '$scope',
    'ProvisionerResource',
    function ($scope, ProvisionerResource) {
      'use strict';

      function loadProv() {
        ProvisionerResource.get({prov: ''}, function (result) {
          $scope.provisioners = [];
          $scope.provisioners = result.Collection || [];
        }, function (result) {
          $scope.provisioners = [];
          $scope.provisioners.push({
            Name: result.data,
            Url: "/",
            Icon: "fa-exclamation-circle"
          });
        });
      }

      $scope.$watch(
        function () {
          return $scope.session.isLoggedIn(); // from parent scope
        },
        function (newValue, oldValue) {
          if (true === newValue || (false === newValue && true === oldValue)) {
            loadProv();
          }
        }
      );
      $scope.initProvNav = loadProv;
    }
  ]);
