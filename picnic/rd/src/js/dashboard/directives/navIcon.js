angular
  .module('Dashboard')
  .directive('rdNavIcon', function () {
    return {
      restrict: 'E',
      scope: {
        icon: '@'
      },
      link: function (scope, element) {
        "use strict";
        var tpl = '';
        if (-1 === scope.icon.indexOf('fa-')) { // img
          tpl = '<span class="menu-icon"><img src="' + scope.icon + '" height="30"/></span>';
        } else { // fa-icon
          tpl = '<span class="menu-icon fa ' + scope.icon + '"></span>';
        }
        element.html(tpl);
      }
    };
  }
);
