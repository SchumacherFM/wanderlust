/**
 * Master Controller
 */
angular
  .module('Dashboard')
  .controller('MasterCtrl', [
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
          $state.go('index');
        });
      };

      $scope.login = function () {
        Session.setLastLoginUrl();
        $state.go('login');
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
  ])
  .controller('systemInfo', [
    '$scope',
    '$timeout',
    'SysInfoResource',
    function ($scope, $timeout, SysInfoResource) {
      var loggedOut = !$scope.session.loggedIn;

      (function tick() { // @todo should be websocket
        $scope.xdata = SysInfoResource.get(function () {
          $timeout(tick, 1000);
        });

        console.log($scope.xdata)

      })();

      $scope.sysInfoWidgets = [
        {
          "icon": "fa-gears",
          "title": 80,
          "comment": "Workers",
          "loading": loggedOut,
          iconColor: "green"
        },
        {
          "icon": "fa-globe",
          "title": 136,
          "comment": "Wanderers",
          "loading": loggedOut,
          iconColor: "orange"
        },
        {
          "icon": "fa-download",
          "title": 16,
          "comment": "Brotzeit",
          "loading": loggedOut,
          iconColor: "red"
        },
        {
          "icon": "fa-database",
          "title": 3,
          "comment": "Provisioners",
          "loading": loggedOut,
          iconColor: "blue"
        }
      ];
    }
  ]);
