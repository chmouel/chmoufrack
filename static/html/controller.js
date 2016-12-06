var workout_name = 'HelloMoto';

var app = angular.module("Frack", []);

app.config(['$locationProvider', function($locationProvider){
    $locationProvider.html5Mode({
        enabled: true,
        requireBase: false
    });
}]);

app.controller("FrackController", ['$scope', '$http', '$location', function($scope, $http, $location) {
    var workout_name = $location.search().name;

    $scope.programs = {'ProgramName': "", "Workout": []}
    var res = $http.get('/rest/program/' + workout_name + "/full");
	res.success(function(data, status, headers, config) {
		$scope.programs = data;
	});
}]);
