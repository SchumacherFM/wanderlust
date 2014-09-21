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
