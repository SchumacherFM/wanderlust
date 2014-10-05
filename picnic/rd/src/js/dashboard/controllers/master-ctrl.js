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
  .controller('systemInfo', [
    '$scope',
    '$timeout',
    'SysInfoResource',
    'SysInfoWidgets',
    function ($scope, $timeout, SysInfoResource, SysInfoWidgets) {
      var loggedIn = $scope.session.isLoggedIn(),
          timeoutSecs = 3000,
          timeoutPromise;

      function tick() { // @todo should be websocket
        SysInfoResource.get().$promise.then(function success(data) {
          angular.forEach(data, function (v, k) {
            if (SysInfoWidgets[k]) {
              SysInfoWidgets[k].title = parseInt(v, 10); // fight against all evil ;-)
              SysInfoWidgets[k].loading = !loggedIn;
            }
          });
          $scope.sysInfoWidgets = SysInfoWidgets;
          timeoutPromise = $timeout(tick, timeoutSecs);
        }, function error() {
          // this interval cancels itself when the user logs out
          angular.forEach(SysInfoWidgets, function (obj) {
            obj.loading = true;
          });
          $scope.sysInfoWidgets = SysInfoWidgets;
        });
      }

      if (true === loggedIn) {
        tick();
      }
      $scope.sysInfoWidgets = SysInfoWidgets;

      // Cancel interval on page changes
      $scope.$on('$destroy', function () {
        if (angular.isDefined(timeoutPromise)) {
          $timeout.cancel(timeoutPromise);
          timeoutPromise = undefined;
        }
      });
    }
  ])
  .controller('userInfo', [
    '$scope',
    'UserInfoResource',
    function ($scope, UserInfoResource) {
      $scope.userCollection = [];
      $scope.isLoading = !$scope.session.isLoggedIn();
      UserInfoResource.get(function (response) {
        $scope.userCollection = response.Users || {};
      });
    }
  ]);
