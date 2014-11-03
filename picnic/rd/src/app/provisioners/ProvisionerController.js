angular
  .module('Wanderlust')
  .controller('ProvisionerController', [
    '$scope',
    '$stateParams',
    'ProvisionerResource',
    'Alert',
    function ($scope, $stateParams, ProvisionerResource, Alert) {
      var type = $stateParams.type || 'textarea';
      $scope.typeName = type;

      ProvisionerResource.get({prov: type}).$promise.then(
        function success(data) {
          angular.forEach(data, function (value, key) {
            // make sure we're not adding $resolved and $promise
            if (!this[key] && (angular.isString(value) || angular.isNumber(value))) {
              this[key] = value;
              // addWatcher(key)
            }
          }, $scope);
        },
        function err(data) {
          Alert.warning("Error in retrieving provisioner data. See console.log for more info.");
          console.log('data err', data);
        }
      );
      // @todo http://adamalbrecht.com/2013/10/30/auto-save-your-model-in-angular-js-with-watch-and-debounce/
    }
  ]);
