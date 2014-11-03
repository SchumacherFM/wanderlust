angular
  .module('Wanderlust')

  // handles all the provisioners
  .factory('ProvisionerResource', [
    '$resource',
    'picnicUrls',
    function ($resource, picnicUrls) {
      return $resource(picnicUrls.provisioners + ':prov', {prov: '@prov'});
    }
  ])
  .factory('ProvisionerForm', [
    '$timeout',
    'ProvisionerResource',
    'Alert',
    function ($timeout, ProvisionerResource, Alert) {
      'use strict';

      return {
        _type: '',
        setType: function (t) {
          this._type = t;
          return this;
        },
        _scope: {},
        setScope: function (s) {
          this._scope = s;
          return this;
        },
        _timeout: null,
        _saveUpdates: function (inputFieldName) {
          var $that = this;
          return function () {
            //console.log($that._scope.provForm);
            if ($that._scope.provForm.$valid) {
              Alert.info("Saved " + inputFieldName);
              //console.log("Saving updates to item #", $that._scope[inputFieldName]);
              ProvisionerResource.save({
                prov: $that._type,
                key: inputFieldName,
                value: $that._scope[inputFieldName]
              }, function (data) {
                console.log('savesuccess', data);
              });
              //  .$promise.then(
              //  function saveSuccess(data) {
              //    console.log('savesuccess', data);
              //  },
              //  function saveError(data) {
              //    console.log('saveerrror', data);
              //  }
              //);
            }
            // invalid input data will be indicated via form input error class
            //Alert.warning("Data is not valid for: " + inputFieldName);

          };
        },
        _debounceUpdate: function (inputFieldName) {
          var $that = this;
          return function (newVal, oldVal) {
            if (newVal !== oldVal) {
              if ($that._timeout) {
                $timeout.cancel($that._timeout);
              }
              $that._timeout = $timeout($that._saveUpdates(inputFieldName), 1.1 * 1000);
            }
          };
        },
        init: function () {
          var $that = this;
          ProvisionerResource.get({prov: $that._type}).$promise.then(
            function success(response) {
              if (!response.data) {
                Alert.warning("Error in retrieving provisioner success data. See console.log for more info.");
                return console.error('Provisioner success error', response);
              }
              angular.forEach(response.data, function (value, inputFieldName) {
                if (!this[inputFieldName]) {
                  this[inputFieldName] = value;
                  this.$watch(inputFieldName, $that._debounceUpdate(inputFieldName));
                }
              }, $that._scope);
            },
            function err(data) {
              Alert.warning("Error in retrieving provisioner data. See console.log for more info.");
              console.error('Provisioner:', data.data || data);
            }
          );
        }
      };

    }
  ]);
