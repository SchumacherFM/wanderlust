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

