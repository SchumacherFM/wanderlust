/**
 * ErrorInterceptor will be applied in the routes.js file
 */
angular
  .module('Dashboard')
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
      Provisioners: {
        "icon": "fa-database",
        "title": 0,
        "comment": "Provisioners",
        "loading": !loggedIn,
        iconColor: "blue"
      }
    };
  })
  .factory('SysInfoResource', function ($resource, picnicUrls) {
    return $resource(picnicUrls.sysinfo, {});
  })
  .factory('AuthInterceptor', function ($window, TrackUser, AUTH_TOKEN_HEADER, AUTH_TOKEN_STORAGE_KEY) {
    // adds for every request the token
    return {
      request: function (config) {
        config.headers = config.headers || {};
        var token = $window.localStorage.getItem(AUTH_TOKEN_STORAGE_KEY);
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
