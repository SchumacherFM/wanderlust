(function(){ 
angular
  .module('Dashboard', [
    'ui.bootstrap',
    'ui.router',
    'ngCookies',
    'ngResource',
    'ui.gravatar'
  ])
  .constant('picnicUrls', {
    auth: '/auth/',
    messages: '/api/messages'
  })
  .constant('AUTH_TOKEN_HEADER', 'X-Auth-Token')
  .constant('AUTH_TOKEN_STORAGE_KEY', 'WL_authToken');

angular.module('ui.gravatar').config(
  [
    'gravatarServiceProvider',
    function (gravatarServiceProvider) {
      gravatarServiceProvider.defaults = {
        size: 40,
        "default": 'monsterid'
      };

      // Use https endpoint
      gravatarServiceProvider.secure = true;
    }
  ]
);

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
    }]
);

/* Services */

(function () {
  'use strict';
  angular.module('picnic.services', [])
    .service('Session', ['$location', '$window', '$q', 'AUTH_TOKEN_STORAGE_KEY', 'Alert',
      function ($location, $window, $q, AUTH_TOKEN_STORAGE_KEY, Alert) {
        var noRedirectUrls = [
          "/login",
          "/changepass",
          "/recoverpass",
          "/signup"
        ];

        function isNoRedirectFromLogin(url) {
          var result = false;
          angular.forEach(noRedirectUrls, function (value) {
            if (value == url) {
              result = true;
            }
          });
          return result;
        }

        function Session() {
          this.clear();
          this.lastLoginUrl = null;
        }

        Session.prototype.init = function (authResource) {
          this.resource = authResource;
          this.sync();
        };

        Session.prototype.sync = function () {
          var $this = this,
              d = $q.defer();
          $this.resource.get({}, function (result) {
            $this.login(result);
            d.resolve(result);
          });
          return d.promise;
        };

        Session.prototype.redirectToLogin = function () {
          this.clear();
          this.setLastLoginUrl();
          Alert.danger("You must be logged in");
          $location.path("/login");
        };

        Session.prototype.check = function () {
          var $this = this;
          $this.sync().then(function () {
            if (!$this.loggedIn) {
              $this.redirectToLogin();
            }
          });
        };

        Session.prototype.setLastLoginUrl = function () {
          this.lastLoginUrl = $location.path();
        };

        Session.prototype.getLastLoginUrl = function () {
          var url = this.lastLoginUrl;
          if (isNoRedirectFromLogin(url)) {
            url = null;
          }
          this.lastLoginUrl = null;
          return url;
        };

        Session.prototype.clear = function () {
          this.loggedIn = false;
          this.name = null;
          this.email = null;
          this.id = null;
          this.isAdmin = false;
        };

        Session.prototype.set = function (session) {
          this.loggedIn = session.loggedIn;
          this.name = session.name;
          this.email = session.email;
          this.id = session.id;
          this.isAdmin = session.isAdmin;
        };

        Session.prototype.login = function (result, token) {
          this.set(result);
          this.$delete = result.$delete;
          if (token) {
            $window.localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, token);
          }
        };

        Session.prototype.logout = function () {
          var $this = this,
              d = $q.defer();
          $this.$delete(function (result) {
            $this.clear();
            d.resolve(result);
            $window.localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY);
          });
          return d.promise;
        };

        return new Session();

      }
    ])
    .service('Auth', [
      '$resource',
      'urls',
      function ($resource, urls) {
        return $resource(urls.auth, {}, {
          'signup': {
            method: 'POST',
            url: urls.auth + 'signup'
          },
          'recoverPassword': {
            method: 'PUT',
            url: urls.auth + 'recoverpass'
          },
          'changePassword': {
            method: 'PUT',
            url: urls.auth + 'changepass'
          }
        });
      }
    ])
    .service('Alert', [
      function () {
        function Alert() {
          var $this = this;
          $this.messages = [];

          var addMessage = function (type, message) {
            $this.messages.push({
              message: message,
              type: type
            });
          };

          $this.dismiss = function (index) {
            $this.messages.splice(index, 1);
          };

          $this.dismissLast = function () {
            $this.messages.pop();
          };

          $this.success = addMessage.bind(null, "success");
          $this.info = addMessage.bind(null, "info");
          $this.warning = addMessage.bind(null, "warning");
          $this.danger = addMessage.bind(null, "danger");
        }

        return new Alert();

      }
    ]);
})();

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