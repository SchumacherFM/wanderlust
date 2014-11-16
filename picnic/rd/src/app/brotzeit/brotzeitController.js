angular
  .module('Wanderlust')
  .controller('BrotzeitController', [
    '$scope',
    '$modal',
    'BrotzeitResource',
    'growl',
    function ($scope, $modal, BrotzeitResource, growl) {
      $scope.bzConfigs = [];
      BrotzeitResource.get().$promise.then(
        function success(data) {
          $scope.bzConfigs = data.Collection || [];
          $scope.bzConfigs.forEach(function (bzc) {
            bzc.isCollapsed = true;
            bzc.ScheduleIsValid = bzc.Schedule !== '';
          });
        },
        function err(data) {
          growl.warning('Error in retrieving Brotzeit collection. See console log');
          console.log('BrotzeitResource err', data);
        }
      );

      $scope.saveCronExpression = function (bzModel) {
        BrotzeitResource.save(
          {
            Route: bzModel.Route || '',
            Schedule: bzModel.Schedule || ''
          },
          function success() {
            bzModel.isCollapsed = true;
            bzModel.ScheduleIsValid = true;
            growl.info('Cron Schedule saved!');
          }
        );
      };

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
