angular
  .module('Wanderlust')
  .controller('ProvisionerController', [
    '$scope',
    '$stateParams',
    'ProvisionerResource',
    'ProvisionerForm',
    function ($scope, $stateParams, ProvisionerResource, ProvisionerForm) {
      var type = $stateParams.type || 'textarea';
      $scope.typeName = type;

      ProvisionerForm.setScope($scope).setType(type).init();
    }
  ]);
