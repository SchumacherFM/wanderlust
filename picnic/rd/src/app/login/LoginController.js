angular
  .module('Wanderlust')
  .controller('LoginController', [
    '$scope',
    '$location',
    '$window',
    'Session',
    'AuthResource',
    'Alert',
    'AUTH_TOKEN_HEADER',
    function ($scope,
              $location,
              $window,
              Session,
              AuthResource,
              Alert,
              AUTH_TOKEN_HEADER) {

      $scope.formData = new AuthResource();

      $scope.login = function () {
        $scope.formData.$save(function saveLoginPost(result, headers) {
          $scope.formData = new AuthResource();

          if (result.loggedIn) {
            Session.login(result, headers(AUTH_TOKEN_HEADER));
            Alert.success("Welcome back, " + result.name);
            var path = Session.getLastLoginUrl();
            if (path) {
              $location.path(path);
            } else {
              $window.history.back();
            }
          }
        });
      };
    }
  ]);