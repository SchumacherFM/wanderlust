angular
  .module('Wanderlust')
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
