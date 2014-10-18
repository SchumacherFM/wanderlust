/**
 * ErrorInterceptor will be applied in the routes.js file
 */
angular
  .module('Wanderlust')

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
  });
