var app = angular.module("Frack", ["ngRoute", "ngSanitize"]);

app.config(function($routeProvider) {
    $routeProvider
        .when("/edit", {controller: "Editor", templateUrl: "partials/editor.html"})
        .when("/workout/:name", {controller: "FrackController", templateUrl: "partials/frack.html"})
        .when("/workout/:name/vma/:vma", {controller: "FrackController", templateUrl: "partials/frack.html"})
        .otherwise({controller: "FrackController", templateUrl: "partials/selection.html"});
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
