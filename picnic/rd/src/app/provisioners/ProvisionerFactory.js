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
        _timeoutSave: null,
        _saveUpdates: function (inputFieldName) {
          var $that = this;
          return function () {
            if ($that._scope.provForm.$valid) {
              ProvisionerResource.save({
                prov: $that._type,
                key: inputFieldName,
                value: $that._scope[inputFieldName]
              }, function () {
                $that._scope[inputFieldName + 'Saved'] = true;
                // remove the green tick that it has successful saved
                if ($that._timeoutSave) {
                  $timeout.cancel($that._timeoutSave);
                }
                $that._timeoutSave = $timeout(function removeGreenTick() {
                  $that._scope[inputFieldName + 'Saved'] = false;
                }, 2300);

              });
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
              if (!response.data || !Array.isArray(response.data)) {
                Alert.warning("Error in retrieving provisioner success data. See console.log for more info.");
                return console.error('Provisioner success error', response);
              }

              // iterating over the slice from GoLang
              var inputName = '', inputValue = '', i = 0, dl = response.data.length;
              for (i = 0; i < dl; i = i + 2) {
                inputName = response.data[i];
                inputValue = response.data[i + 1];
                if (!$that._scope[inputName]) {
                  $that._scope[inputName] = inputValue;
                  $that._scope[inputName + 'Saved'] = false;
                  $that._scope.$watch(inputName, $that._debounceUpdate(inputName));
                }
              }

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
