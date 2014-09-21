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
    'Session',
    'Auth',
    function ($scope, $cookieStore,Session,Auth) {

      $scope.session = Session;
      Session.init(Auth);

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

      $scope.login = function () {
        alert("@todo login")
      };
      $scope.logout = function () {
        alert("@todo logout")
      };

    }
  ]
);
