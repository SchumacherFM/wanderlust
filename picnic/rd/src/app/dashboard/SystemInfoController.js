/**
 * Dashboard Controller
 */
angular
  .module('Wanderlust')
  .controller('SystemInfoController', [
    '$scope',
    '$timeout',
    'SysInfoResource',
    'SysInfoWidgets',
    function ($scope, $timeout, SysInfoResource, SysInfoWidgets) {
      var timeoutSecs = 3000,
          timeoutPromise;

      function tick() { // @todo should be websocket
        SysInfoResource.get().$promise.then(function success(data) {
          angular.forEach(data, function (v, k) {
            if (SysInfoWidgets[k]) {
              SysInfoWidgets[k].title = parseInt(v, 10); // fight against all evil ;-)
            }
          });

          if (SysInfoWidgets.SessionExpires) {
            var s = SysInfoWidgets.SessionExpires.title,
                m = parseInt(s / 60, 10);
            s = s - (m * 60);
            SysInfoWidgets.SessionExpires.title = m + 'm ' + s + 's';
          }
          $scope.sysInfoWidgets = SysInfoWidgets;
          timeoutPromise = $timeout(tick, timeoutSecs);
        }, function error() {
          // this interval cancels itself when the user logs out
          $scope.isLoading = !$scope.session.isLoggedIn();
          $scope.sysInfoWidgets = SysInfoWidgets;
        });
      }

      if (true === $scope.session.isLoggedIn()) {
        tick();
      }
      $scope.isLoading = !$scope.session.isLoggedIn();
      $scope.sysInfoWidgets = SysInfoWidgets;

      // Cancel interval on page changes
      $scope.$on('$destroy', function () {
        if (angular.isDefined(timeoutPromise)) {
          $timeout.cancel(timeoutPromise);
          timeoutPromise = undefined;
        }
      });
    }
  ]);
