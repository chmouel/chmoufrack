var app = angular.module("Frack", ["ngRoute"]);

app.config(function($routeProvider) {
    $routeProvider.when("/:name", {controller: "FrackController", templateUrl: "frack.html"});
});


app.controller("FrackController", ['$scope', '$routeParams', '$http', function($scope, $routeParams, $http) {
    var workout_name = 'HelloMoto';
    var vma = '17'
    $scope.programs = {'ProgramName': "", "Workout": []}

    if ($routeParams.name) {
         workout_name = $routeParams.name
    } else {
        //TODO: list all workout to choose from
        return;
    }

    res = $http.get('/rest/program/' + workout_name + "/" + vma + "/full");
	res.success(function(data, status, headers, config) {
		$scope.programs = data;
	});
}])
