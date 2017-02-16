var app = angular.module("Frack", ["ngRoute", "ngSanitize"]);

app.config(function($routeProvider) {
    $routeProvider
        .when("/edit/:name", {controller: "EditController", templateUrl: "html/edit/editor.html"})
        .when("/workout/:name", {controller: "ViewController", templateUrl: "html/view/view.html"})
        .when("/workout/:name/vma/:vma", {controller: "ViewController", templateUrl: "html/view/view.html"})
        .otherwise({controller: "ViewController", templateUrl: "html/view/selection.html"});
});

function range(min, max) {
    var input = [];
    min = parseInt(min);
    max = parseInt(max);
    for (var i=min; i<=max; i++)
        input.push(i.toString());
    return input;
};

// Disable caching: https://goo.gl/yHW1vE
app.config(['$httpProvider', function($httpProvider) {
    //initialize get if not there
    if (!$httpProvider.defaults.headers.get)
        $httpProvider.defaults.headers.get = {};

    $httpProvider.defaults.headers.get['If-Modified-Since'] = 'Mon, 26 Jul 1997 05:00:00 GMT';
    $httpProvider.defaults.headers.get['Cache-Control'] = 'no-cache';
    $httpProvider.defaults.headers.get['Pragma'] = 'no-cache';
}]);

//http://stackoverflow.com/a/36254259
app.directive('input', [function() {
    return {
        restrict: 'E',
        require: '?ngModel',
        link: function(scope, element, attrs, ngModel) {
            if (
                   'undefined' !== typeof attrs.type
                && 'number' === attrs.type
                && ngModel
            ) {
                ngModel.$formatters.push(function(modelValue) {
                    return Number(modelValue);
                });

                ngModel.$parsers.push(function(viewValue) {
                    return Number(viewValue);
                });
            }
        }
    }
}]);
