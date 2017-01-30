var app = angular.module("Frack", ["ngRoute"]);

app.config(function($routeProvider) {
    $routeProvider
        .when("/workout/:name", {controller: "FrackController", templateUrl: "partials/frack.html"})
        .when("/workout/:name/vma/:vma", {controller: "FrackController", templateUrl: "partials/frack.html"})
        .otherwise({controller: "FrackController", templateUrl: "partials/frack.html"});
});


function range(input, min, max) {
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


    if ($routeParams.name) {
        $scope.programWanted = $routeParams.name;
    }

    if ($routeParams.vma) {
        $scope.vmaWanted = [parseInt($routeParams.vma)];
    } else {
        $scope.vmaWanted = range([], 12, 18);
    }

    var res = $http.get('/rest/programs');
	res.success(function(data, status, headers, config) {
		$scope.programs = data;
	});

    $scope.submit = function() {
        $location.path("/workout/" + $scope.selectedWorkout.Name);
    };

}])
