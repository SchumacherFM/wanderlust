angular
  .module('Dashboard')
  .controller('ShopCtrl', [
    '$scope',
    'Alert',
    function ($scope, Alert) {
      /**
       * this is only the temporary implementation. All this will be outsources into magento
       * and we will talk directly to Magento REST API. Concept is that GoLang talks to Magento
       * and then provides a route for ng for faster displaying.
       */
      var i = 0;
      var ps = [
        {
          id: i++,
          t: "Google Analytics",
          d: "Get your data automatically from the Google Analytics API.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: true // available
        },
        {
          id: i++,
          t: "Piwik",
          d: "Get your data automatically from the Piwik API.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: false
        },
        {
          id: i++,
          t: "KISSMetrics",
          d: "Get your data automatically from the KISSMetrics API.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: false
        },
        {
          id: i++,
          t: "Clicky",
          d: "Get your data automatically from the Clicky API.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: false
        },
        {
          id: i++,
          t: "Your API",
          d: "Get your data automatically from your custom API.",
          i: "img/800x500.gif",
          p: "€149.00/hourly",
          a: false
        },
        {
          id: i++,
          t: "Webhook",
          d: "Trigger Wanderlust whenever you do a deployment.",
          i: "img/800x500.gif",
          p: "€49.00",
          a: false
        },
        {
          id: i++,
          t: "Concurrency",
          d: "Unlimited number of concurrent request. Default is one request.",
          i: "img/800x500.gif",
          p: "€2.00",
          a: false
        },
        {
          id: i++,
          t: "Automatic updates",
          d: "Alerts you for new versions with auto download.",
          i: "img/800x500.gif",
          p: "€9.00/monthly",
          a: false
        },
        {
          id: i++,
          t: "Magento CE",
          d: "Import data directly from Magento.",
          i: "img/800x500.gif",
          p: "€29.00",
          a: false
        },
        {
          id: i++,
          t: "Magento EE",
          d: "Import data directly from Magento Enterprise.",
          i: "img/800x500.gif",
          p: "€99.00",
          a: false
        },
        {
          id: i++,
          t: "TYPO3",
          d: "Import data directly from TYPO3.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: false
        },
        {
          id: i++,
          t: "Drupal",
          d: "Import data directly from Drupal.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: false
        },
        {
          id: i++,
          t: "Wordpress",
          d: "Import data directly from Wordpress.",
          i: "img/800x500.gif",
          p: "€9.00",
          a: false
        },
        {
          id: i++,
          t: "WooCommerce",
          d: "Import data directly from WooCommerce.",
          i: "img/800x500.gif",
          p: "€49.00",
          a: false
        },
        {
          id: i++,
          t: "Shopware",
          d: "Import data directly from Shopware.",
          i: "img/800x500.gif",
          p: "€49.00",
          a: false
        },
        {
          id: i++,
          t: "RSS/Atom Feed",
          d: "Import data from a RSS/Atom Feed",
          i: "img/800x500.gif",
          p: "€19.00",
          a: false
        }
      ];
      $scope.products = ps;
      $scope.helloWorld = function () {
        Alert.danger('Hello World! This feature is WIP');
      }
    }
  ]);