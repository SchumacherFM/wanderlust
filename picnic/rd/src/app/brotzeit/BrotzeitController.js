angular
  .module('Wanderlust')
  .controller('BrotzeitController', [
    '$scope',
    '$modal',
    function ($scope, $modal) {

      $scope.showCronForm = function () {
        alert('Use a popover');
      };

      $scope.isCollapsed = true;

      $scope.openCronHelp = function () {
        $modal.open({
          templateUrl: 'partials/brotzeit/tpl/cronHelp.html',
          controller: 'CronHelpController',
          size: 'lg'
        });
      };

    }])
  // Please note that $modalInstance represents a modal window (instance) dependency.
  // It is not the same as the $modal service used above.
  .controller('CronHelpController', function ($scope, $modalInstance) {

    $scope.close = function () {
      $modalInstance.dismiss('cancel');
    };
  });
