var app = angular.module("Frack", ["ngRoute"]);
var vma = '14:19'

app.config(function($routeProvider) {
    $routeProvider
        .when("/workout/:name", {controller: "FrackController", templateUrl: "partials/frack.html"})
        .when("/workout/:name/vma/:vma", {controller: "FrackController", templateUrl: "partials/frack.html"});
});


app.controller("FrackController", ['$scope', '$routeParams', '$http', function($scope, $routeParams, $http) {
    var workout_name = 'HelloMoto';
    $scope.programs = {'ProgramName': "", "Workout": []}

    if ($routeParams.name) {
        workout_name = $routeParams.name;
    } else {
        //TODO: list all workout to choose from
        return;
    }

    if ($routeParams.vma) {
        vma = $routeParams.vma;
    }

    res = $http.get('/rest/program/' + workout_name + "/" + vma + "/full");
	res.success(function(data, status, headers, config) {
		$scope.programs = data;
	});
}])
