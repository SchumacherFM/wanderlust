(function(){ 
angular
  .module('Dashboard', [
    'ui.bootstrap',
    'ui.router',
    'ngCookies',
    'ngResource',
    'ui.gravatar',
    'ui.bootstrap',
    'picnic.services',
    'angulartics',
    'angulartics.piwik'
  ])
  .constant('picnicUrls', {
    auth: '/auth/',
    users: '/users/',
    messages: '/api/messages'
  })
  .constant('AUTH_TOKEN_HEADER', 'X-Auth-Token')
  .constant('AUTH_TOKEN_STORAGE_KEY', 'WL_authToken')
  .config([
    '$resourceProvider',
    'gravatarServiceProvider',
    function ($resourceProvider, gravatarServiceProvider) {
      // Don't strip trailing slashes from calculated URLs
      $resourceProvider.defaults.stripTrailingSlashes = false;
      gravatarServiceProvider.defaults = {
        size: 40,
        "default": 'monsterid'
      };

      // Use https endpoint
      gravatarServiceProvider.secure = true;
    }]);

/**
 * ErrorInterceptor will be applied in the routes.js file
 */
angular
  .module('Dashboard')
  .factory('ErrorInterceptor', function ($q, $location, Session, Alert) {
    return {

      response: function (response) {
        return response;
      },

      responseError: function (response) {
        var rejection = $q.reject(response),
            status = parseInt(response.status, 10),
            msg = 'Sorry, an error has occurred';

        if (401 === status) {
          Session.redirectToLogin();
          return;
        }
        if (404 === status) {
          // handle locally
          return;
        }
        if (403 === status) {
          msg = "Sorry, you're not allowed to do this";
        }
        if (400 === status && response.data.errors) {
          msg = "Sorry, your form contains errors, please try again";
        }

        if (response.data && typeof(response.data) === 'string') {
          msg = response.data;
        }
        console.log('msg', msg);
        if (msg.length > 0) {
          Alert.danger(msg);
        }
        return rejection;
      }
    };
  });

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

/* Services */
angular.module('picnic.services', [])
  .service('Session', ['$location', '$window', '$q', 'AUTH_TOKEN_STORAGE_KEY', 'Alert',
    function ($location, $window, $q, AUTH_TOKEN_STORAGE_KEY, Alert) {
      var noRedirectUrls = {
        "/login": true,
        "/changepass": true,
        "/recoverpass": true
      };

      function isNoRedirectFromLogin(url) {
        return noRedirectUrls[url] || false;
      }

      function Session() {
        this.clear();
        this.lastLoginUrl = null;
      }

      Session.prototype.init = function (authResource) {
        this.authResource = authResource;
        this.sync();
      };

      Session.prototype.sync = function () {
        var $this = this,
            d = $q.defer();
        $this.authResource.get({}, function (result) {
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
        if (true === isNoRedirectFromLogin(url)) {
          url = null;
        }
        this.lastLoginUrl = null;
        return url;
      };

      Session.prototype.clear = function () {
        this.loggedIn = false;
        this.name = null;
        this.userName = null;
        this.email = null;
        this.id = null;
        this.isAdmin = false;
      };

      Session.prototype.set = function (session) {
        this.loggedIn = session.loggedIn;
        this.name = session.name;
        this.userName = session.userName;
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
  .service('AuthResource', [
    '$resource',
    'picnicUrls',
    function ($resource, picnicUrls) {
      return $resource(picnicUrls.auth, {}, {
        'recoverPassword': {
          method: 'PUT',
          url: picnicUrls.auth + 'recoverpass'
        },
        'changePassword': {
          method: 'PUT',
          url: picnicUrls.auth + 'changepass'
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
    'Session',
    'AuthResource',
    'Alert',
    function ($scope, $state, $cookieStore, $timeout, Session, AuthResource, Alert) {

      //<Alerts>
      $scope.alert = Alert;
      $scope.$watchCollection('alert.messages', function (newValue, oldValue) {
        $timeout(function () {
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

angular
  .module('Dashboard')
  .controller('LoginCtrl', [
    '$scope',
    '$location',
    '$window',
    '$http',
    'Session',
    'AuthResource',
    'Alert',
    'AUTH_TOKEN_HEADER',
    function ($scope,
              $location,
              $window,
              $http,
              Session,
              AuthResource,
              Alert,
              AUTH_TOKEN_HEADER) {

      $scope.formData = new AuthResource();

      $scope.login = function () {
        $scope.formData.$save(function saveLoginPost(result, headers) {
          $scope.formData = new AuthResource();
          console.log(result, headers(AUTH_TOKEN_HEADER));
          if (result.loggedIn) {
            Session.login(result, headers(AUTH_TOKEN_HEADER));
            Alert.success("Welcome back, " + result.name);
            var path = Session.getLastLoginUrl();
            console.log(path);
            //if (path) {
            //  $location.path(path);
            //} else {
            //  $window.history.back();
            //}
          }
        });
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