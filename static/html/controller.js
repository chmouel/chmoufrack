var app = angular.module("Frack", []);

app.config(['$locationProvider', function($locationProvider){
    $locationProvider.html5Mode({
        enabled: true,
        requireBase: false
    });
}]);

app.controller("FrackController", ['$scope', '$http', '$location', function($scope, $http, $location) {
    var workout_name = $location.search().name;
    var vma = $location.search().vma;
    var res = '';
    if (!workout_name) {
        workout_name = 'HelloMoto';
    }

    if (vma) {
        res = $http.get('/rest/program/' + workout_name + "/" + vma + "/full");
    } else {
        res = $http.get('/rest/program/' + workout_name + "/full");
    }

    $scope.programs = {'ProgramName': "", "Workout": []}
	res.success(function(data, status, headers, config) {
		$scope.programs = data;
	});
}]);
