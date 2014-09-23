/**
 * Master Controller
 */
angular
  .module('Dashboard')
  .controller(
  'MasterCtrl',
  [
    '$scope',
    '$state',
    '$cookieStore',
    '$timeout',
    '$analytics',
    'Session',
    'AuthResource',
    'Alert',
    function ($scope, $state, $cookieStore, $timeout, $analytics, Session, AuthResource, Alert) {

      //<Alerts>
      $scope.alert = Alert;
      $scope.$watchCollection('alert.messages', function (newValue, oldValue) {
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
          $state.go("/");
        });
      };

      $scope.login = function () {
        Session.setLastLoginUrl();
        $state.go("login");
      };
      //</Sessions>

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
