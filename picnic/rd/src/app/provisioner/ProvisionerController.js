angular
  .module('Wanderlust')
  .controller('ProvisionerController', [
    '$scope',
    '$stateParams',
    'ProvisionerResource',
    function ($scope, $stateParams, ProvisionerResource) {
      var type = $stateParams.type || 'textarea';
      $scope.typeName = type;

      ProvisionerResource.get({prov: type}).$promise.then(
        function success(data) {
          "use strict";
          console.log('data success', data);
        },
        function err(data) {
          "use strict";
          console.log('data err', data);
        }
      );

      //var formContent = [
      //  {
      //    label: "URL to sitemap.xml",
      //    input: {
      //      type: "text",
      //      name: "sitemap1",
      //      model: "sitemap",
      //      pattern: /^http.+\.xml$/gi,
      //      placeholder: "http://my-server.com/sitemap.xml"
      //    },
      //    helpBlock: "Please enter a valid sitemap URL.",
      //    infoBlock: "Only the first 10 URLs will be parsed."
      //  }
      //];

      $scope.sitemap = '';

      $scope.provSave = function (e) {
        console.log('save', $scope.sitemap, e);
      };
      $scope.validSitemap = /^http.+\.xml$/gi;

    }
  ]);
