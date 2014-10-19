angular
  .module('Wanderlust')
  .controller('UserInfoController', [
    '$scope',
    'UserInfoResource',
    function ($scope, UserInfoResource) {
      /**
       * Gets the collection of users and displays them to see who has an account
       */
      $scope.userCollection = [];
      $scope.isLoading = !$scope.session.isLoggedIn();
      UserInfoResource.get(function (response) {
        $scope.userCollection = response.Users || {};
      });
    }
  ]);
