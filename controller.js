var app = angular.module("Frack", ["ngRoute"]);

app.config(function($routeProvider) {
    $routeProvider
        .when("/workout/:name", {controller: "FrackController", templateUrl: "partials/frack.html"})
        .when("/workout/:name/vma/:vma", {controller: "FrackController", templateUrl: "partials/frack.html"})
        .otherwise({controller: "FrackController", templateUrl: "partials/frack.html"});
});


function range(min, max) {
    var input = Array();
    min = parseInt(min);
    max = parseInt(max);
    for (var i=min; i<=max; i++)
        input.push(i);
    return input;
};


app.filter('range', function() {
    return range;
});

app.controller("FrackController", ['$scope', '$location', '$routeParams', '$http', function($scope, $location, $routeParams, $http) {
    $scope.programWanted = '';
    $scope.vmaWanted = [];
    $scope.allVMAS = range(12, 18);
    $scope.selectedVMA = 18;

    if ($routeParams.name) {
        $scope.programWanted = $routeParams.name;
        $scope.selectedProgram = $routeParams.name;
    }

    // By default use 12-18, if we have a 12:18 in the router config parse it to
    // range int 12.13.14....18...etc. if only one number then pass it straight
    if ($routeParams.vma) {
        if ($routeParams.vma.indexOf(":") == -1) {
            $scope.selectedVMA = $routeParams.vma;
            $scope.vmaWanted = [parseInt($routeParams.vma)];
        } else {
            var sp = $routeParams.vma.split(':');
            $scope.vmaWanted = range(sp[0], sp[1]);
        }
    } else {
        $scope.vmaWanted = range(12, 18); //default
    }

    if (!$scope.programs) {
        var res = $http.get('/frack.yaml');
	    res.success(function(data, status, headers, config) {
            $scope.programs = jsyaml.load(data);;
            $scope.programNames = Array();
            for (var p of $scope.programs)
                $scope.programNames.push(p.name);
	    });
    }

    $scope.submit = function() {
        var t = "", p ="";
        if ( $scope.selectedVMA) {
            t = $scope.selectedVMA;
        } else if ($scope.vmaWanted.length > 1) {
            t = $scope.vmaWanted[0] + ":" + $scope.vmaWanted[$scope.vmaWanted.length - 1];
        } else {
            t = "12:18"; //default
        }

        p = $scope.selectedProgram;
        $location.path("/workout/" + p + "/vma/" + t);
    };

}]);
