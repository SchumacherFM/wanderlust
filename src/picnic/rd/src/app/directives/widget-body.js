angular
  .module('Wanderlust')
  .directive('rdWidgetBody', function () {
    return {
      restrict: 'E',
      requires: '^rdWidget',
      scope: {
        loading: '=',
        bodyclass: '@'
      },
      transclude: true,
      template: '<div class="widget-body" data-ng-class="bodyclass">' +
      '<rd-loading data-ng-show="loading"></rd-loading>' +
      '<div data-ng-hide="loading" class="widget-content" data-ng-transclude></div></div>'
    };
  });
