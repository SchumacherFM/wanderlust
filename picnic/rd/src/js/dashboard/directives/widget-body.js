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
