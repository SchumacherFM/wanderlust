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
          setVar(1, "username", user)
        },
        "setToken": function (token) {
          setVar(2, "token", token)
        }
      }
    }])
  .service('Session', ['$location', '$window', '$q', 'AUTH_TOKEN_STORAGE_KEY', 'Alert', 'TrackUser',
    function ($location, $window, $q, AUTH_TOKEN_STORAGE_KEY, Alert, TrackUser) {
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
