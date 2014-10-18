(function(){ 
'use strict';
angular
  .module('Dashboard', [
    'ui.bootstrap',
    'ui.router',
    'LocalStorageModule',
    'ngResource',
    'ui.gravatar',
    'ui.bootstrap',
    'picnic.services',
    'angulartics',
    'angulartics.piwik',
    'ncy-angular-breadcrumb'
  ])
  .constant('picnicUrls', {
    auth: '/auth/',
    users: '/users/',
    sysinfo: '/sysinfo/',
    provisioners: '/provisioners/',
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

  // handles all the provisioners
  .factory('ProvisionerResource', function ($resource, picnicUrls) {
    return $resource(picnicUrls.provisioners + ':prov', {prov: '@prov'});
  })

  // loads the user collection when the dashboard website is open.
  .factory('UserInfoResource', function ($resource, picnicUrls) {
    return $resource(picnicUrls.users, {});
  })
  .factory('SysInfoResource', function ($resource, picnicUrls) {
    return $resource(picnicUrls.sysinfo, {});
  })
  .factory('SysInfoWidgets', function (Session) {
    var loggedIn = Session.isLoggedIn();
    return {
      Goroutines: {
        "icon": "fa-gears",
        "title": 0,
        "comment": "Workers",
        "loading": !loggedIn,
        iconColor: "green"
      },
      Wanderers: {
        "icon": "fa-globe",
        "title": 0,
        "comment": "Wanderers",
        "loading": !loggedIn,
        iconColor: "orange"
      },
      Brotzeit: {
        "icon": "fa-download",
        "title": 0,
        "comment": "Brotzeit",
        "loading": !loggedIn,
        iconColor: "red"
      },
      SessionExpires: {
        "icon": "fa-clock-o",
        "title": 0,
        "comment": "Log out in",
        "loading": !loggedIn,
        iconColor: "blue"
      }
    };
  })
  .factory('AuthInterceptor', function (localStorageService, TrackUser, AUTH_TOKEN_HEADER, AUTH_TOKEN_STORAGE_KEY) {
    // adds for every request the token
    return {
      request: function (config) {
        config.headers = config.headers || {};
        var token = localStorageService.get(AUTH_TOKEN_STORAGE_KEY);
        if (token && token.length > 20) {
          TrackUser.setToken(token);
          config.headers[AUTH_TOKEN_HEADER] = token;
        }
        return config;
      }
    };

  })
  .factory('ErrorInterceptor', function ($q, /*$location, */ Session, Alert) {
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
        if (412 === status) { // 412 pre condition failed: Waiting for login ...
          return rejection;
        }
        if (403 === status) {
          msg = "Sorry, you're not allowed to do this";
        }
        if (400 === status && response.data.errors) {
          msg = "Sorry, your form contains errors, please try again";
        }

        if (response.data && typeof response.data === 'string') {
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
angular.module('Dashboard')
  .config([
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
          templateUrl: 'partials/dashboard.html',
          data: {
            ncyBreadcrumbLabel: 'Dashboard'
          }
        })
        .state('login', {
          url: '/login',
          templateUrl: 'partials/login.html',
          controller: 'LoginCtrl',
          data: {
            ncyBreadcrumbLabel: 'Login'
          }
        })
        .state('tables', {
          url: '/tables',
          templateUrl: 'partials/tables.html',
          data: {
            ncyBreadcrumbLabel: 'Yet another demo table page'
          }
        })
        .state('shop', {
          url: '/shop',
          templateUrl: 'partials/shop.html',
          controller: 'ShopCtrl',
          data: {
            ncyBreadcrumbLabel: 'Shop - Your in-app purchase made easy!'
          }
        })
        .state('privacy', {
          url: '/privacy',
          templateUrl: 'partials/privacy.html',
          data: {
            ncyBreadcrumbLabel: 'Privacy Statement'
          }
        })
        .state('provisioners', {
          url: '/provisioners/:type',
          templateUrl: function ($stateParams) {
            // 404 errors can occur when a template not exists
            var type = $stateParams.type || 'textarea';
            return 'partials/provisioners/' + type + '.html';
          },
          data: {
            ncyBreadcrumbLabel: 'Provisioner / {{name}}'
          }
        });

      $httpProvider.interceptors.push('AuthInterceptor');
      $httpProvider.interceptors.push('ErrorInterceptor');

    }]);

/* Services */
angular.module('picnic.services', [])
  .service('TrackUser', ['$window', 'md5',
    function ($window, md5) {

      function setVar(index, theVar, user) {
        $window._paq.push(['setCustomVariable', // piwik specific
          index, // Index, the number from 1 to 5 where this custom variable name is stored
          theVar,
          md5(user),
          "visit" // scope
        ]);
      }

      return {
        "setUser": function (user) {
          setVar(1, "username", user);
        },
        "setToken": function (token) {
          setVar(2, "token", token);
        }
      };
    }])
  .service('Session', ['$location', 'localStorageService', '$q', 'AUTH_TOKEN_STORAGE_KEY', 'Alert', 'TrackUser',
    function ($location, localStorageService, $q, AUTH_TOKEN_STORAGE_KEY, Alert, TrackUser) {
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
        this.isAdmin = false;
      };

      Session.prototype.set = function (session) {
        TrackUser.setUser(session.userName)
        this.loggedIn = session.loggedIn;
        this.name = session.name;
        this.userName = session.userName;
        this.email = session.email;
        this.isAdmin = session.isAdmin;
      };

      Session.prototype.isLoggedIn = function () {
        return this.loggedIn === true;
      };

      Session.prototype.login = function (result, token) {
        this.set(result);
        this.$delete = result.$delete;
        if (token) {
          TrackUser.setToken(token);
          localStorageService.set(AUTH_TOKEN_STORAGE_KEY, token);
        }
      };

      Session.prototype.logout = function () {
        var $this = this,
            d = $q.defer();
        $this.$delete(function (result) {
          $this.clear();
          d.resolve(result);
          localStorageService.remove(AUTH_TOKEN_STORAGE_KEY);
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
            Url: "/",
            Icon: "fa-exclamation-circle"
          });
        });
      }

      $scope.$watch(
        function () {
          return Session.isLoggedIn();
        },
        function (newValue) {
          if (true === newValue) {
            loadProv();
          }
        }
      );
      loadProv();
    }
  ]);

/**
 * Dashboard Controller
 */
angular
  .module('Dashboard')
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


angular
  .module('Dashboard')
  .controller('LoginCtrl', [
    '$scope',
    '$location',
    '$window',
    'Session',
    'AuthResource',
    'Alert',
    'AUTH_TOKEN_HEADER',
    function ($scope,
              $location,
              $window,
              Session,
              AuthResource,
              Alert,
              AUTH_TOKEN_HEADER) {

      $scope.formData = new AuthResource();

      $scope.login = function () {
        $scope.formData.$save(function saveLoginPost(result, headers) {
          $scope.formData = new AuthResource();

          if (result.loggedIn) {
            Session.login(result, headers(AUTH_TOKEN_HEADER));
            Alert.success("Welcome back, " + result.name);
            var path = Session.getLastLoginUrl();
            if (path) {
              $location.path(path);
            } else {
              $window.history.back();
            }
          }
        });
      };
    }
  ]);
angular
  .module('Dashboard')
  .controller('ProvisionerCtrl', [
    '$scope',
    'Session',
    function ($scope,
              Session) {

      $scope.name = 'Hello';
      console.log('$scope.name', $scope.name)
    }
  ]);
angular
  .module('Dashboard')
  .controller('ShopCtrl', [
    '$scope',
    'Alert',
    function ($scope, Alert) {
      /**
       * this is only the temporary implementation. All this will be outsources into magento
       * and we will talk directly to Magento REST API. Concept is that GoLang talks to Magento
       * and then provides a route for ng for faster displaying.
       */
      var i = 0;
      var ps = [
        {
          id: i++,
          t: "Google Analytics",
          d: "Get your data automatically from the Google Analytics API.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: true // available
        },
        {
          id: i++,
          t: "Piwik",
          d: "Get your data automatically from the Piwik API.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: false
        },
        {
          id: i++,
          t: "KISSMetrics",
          d: "Get your data automatically from the KISSMetrics API.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: false
        },
        {
          id: i++,
          t: "Clicky",
          d: "Get your data automatically from the Clicky API.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: false
        },
        {
          id: i++,
          t: "Your API",
          d: "Get your data automatically from your custom API.",
          i: "img/800x500.gif",
          p: "€149.00/hourly",
          a: false
        },
        {
          id: i++,
          t: "Webhook",
          d: "Trigger Wanderlust whenever you do a deployment.",
          i: "img/800x500.gif",
          p: "€49.00",
          a: false
        },
        {
          id: i++,
          t: "Concurrency",
          d: "Unlimited number of concurrent request. Default is one request.",
          i: "img/800x500.gif",
          p: "€2.00",
          a: false
        },
        {
          id: i++,
          t: "Automatic updates",
          d: "Alerts you for new versions with auto download.",
          i: "img/800x500.gif",
          p: "€9.00/monthly",
          a: false
        },
        {
          id: i++,
          t: "Magento CE",
          d: "Import data directly from Magento.",
          i: "img/800x500.gif",
          p: "€29.00",
          a: false
        },
        {
          id: i++,
          t: "Magento EE",
          d: "Import data directly from Magento Enterprise.",
          i: "img/800x500.gif",
          p: "€99.00",
          a: false
        },
        {
          id: i++,
          t: "TYPO3",
          d: "Import data directly from TYPO3.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: false
        },
        {
          id: i++,
          t: "Drupal",
          d: "Import data directly from Drupal.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: false
        },
        {
          id: i++,
          t: "Wordpress",
          d: "Import data directly from Wordpress.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: false
        },
        {
          id: i++,
          t: "WooCommerce",
          d: "Import data directly from WooCommerce.",
          i: "img/800x500.gif",
          p: "€49.00",
          a: false
        },
        {
          id: i++,
          t: "Shopware",
          d: "Import data directly from Shopware.",
          i: "img/800x500.gif",
          p: "€49.00",
          a: false
        },
        {
          id: i++,
          t: "RSS/Atom Feed",
          d: "Import data from a RSS/Atom Feed",
          i: "img/800x500.gif",
          p: "€19.00",
          a: false
        }
      ];
      $scope.products = ps;
      $scope.helloWorld = function () {
        Alert.danger('Hello World! This feature is WIP');
      }
    }
  ]);
angular
  .module('Dashboard')
  .directive('rdCheck', function () {
    return {
      restrict: 'AE',
      scope: {
        checked: '@'
      },
      template: '<i class="fa fa-check" data-ng-show="checked"></i><i class="fa fa-times" data-ng-hide="checked"></i>'
    };
  }
);


angular
  .module('Dashboard')
  .directive('rdNavLi', function () {
    return {
      restrict: 'E',
      template: '<li data-ng-model="p.Name" class="sidebar-list">' +
        '<a href="#{{p.Url}}" data-analytics-on="click" data-analytics-category="navigation">{{p.Name}}' +
        '<rd-nav-icon icon="{{p.Icon}}"></rd-nav-icon></a></li>',
      scope: {
        p: '='
      }
    };
  })
  .directive('rdNavIcon', function () {
    return {
      restrict: 'E',
      scope: {
        icon: '@'
      },
      link: function (scope, element) {
        'use strict';
        var tpl = '';
        if (-1 === scope.icon.indexOf('fa-')) { // img
          tpl = '<span class="menu-icon"><img src="' + scope.icon + '" height="30"/></span>';
        } else { // fa-icon
          tpl = '<span class="menu-icon fa ' + scope.icon + '"></span>';
        }
        element.html(tpl);
      }
    };
  });

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


angular
  .module('Dashboard')
  .directive('rdWidget', function () {
    return {
      transclude: true,
      template: '<div class="widget" ng-transclude></div>',
      restrict: 'EA'
    };
  }
);


angular
  .module('Dashboard')
  .directive('rdWidgetHeader', function () {
    return {
      requires: '^rdWidget',
      scope: {
        title: '@',
        icon: '@'
      },
      transclude: true,
      template: '<div class="widget-header"> <i class="fa" ng-class="icon"></i> {{title}} <div class="pull-right" ng-transclude></div></div>',
      restrict: 'E'
    };
  });

angular
  .module('Dashboard')
  .directive('rdWidgetBody', function () {
    return {
      requires: '^rdWidget',
      scope: {
        loading: '@?',
        bodyclass: '@'
      },
      transclude: true,
      template: '<div class="widget-body" ng-class="bodyclass">' +
        '<rd-loading ng-show="loading"></rd-loading>' +
        '<div ng-hide="loading" class="widget-content" ng-transclude></div></div>',
      restrict: 'E'
    };
  }
);

})();