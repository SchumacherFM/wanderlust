angular
  .module('Wanderlust')
  .factory('SysInfoResource', function ($resource, picnicUrls) {
    return $resource(picnicUrls.sysinfo, {});
  })
  .factory('SysInfoWidgets', function () {
    return {
      Goroutines: {
        "icon": "fa-gears",
        "title": 0,
        "comment": "Workers",
        iconColor: "green"
      },
      Wanderers: {
        "icon": "fa-globe",
        "title": 0,
        "comment": "Wanderers",
        iconColor: "orange"
      },
      Brotzeit: {
        "icon": "fa-download",
        "title": 0,
        "comment": "Brotzeit",
        iconColor: "red"
      },
      SessionExpires: {
        "icon": "fa-clock-o",
        "title": 0,
        "comment": "Log out in",
        iconColor: "blue"
      }
    };
  });
