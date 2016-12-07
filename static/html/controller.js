var app = angular.module("Frack", ["ngRoute"]);

app.config(function($routeProvider) {
    $routeProvider.when("/:name", {controller: "FrackController"});
});

function FrackController($scope, $routeParams, $http) {
    var workout_name = 'HelloMoto';
    var vma = '17'
    $scope.programs = {'ProgramName': "", "Workout": []}
    console.debug($routeParams);

    res = $http.get('/rest/program/' + workout_name + "/" + vma + "/full");
	res.success(function(data, status, headers, config) {
		$scope.programs = data;
	});
}
app.controller("FrackController", ['$scope', '$routeParams', '$http', function($scope, $routeParams, $http) {
    var workout_name = 'HelloMoto';
    var vma = '17'
    $scope.programs = {'ProgramName': "", "Workout": []}
    console.debug($routeParams.name);
    return;

    res = $http.get('/rest/program/' + workout_name + "/" + vma + "/full");
	res.success(function(data, status, headers, config) {
		$scope.programs = data;
	});
}])
