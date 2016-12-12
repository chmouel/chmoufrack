var app = angular.module("Frack", ["ngRoute"]);
var default_vma = '14:19'

app.config(function($routeProvider) {
    $routeProvider
        .when("/workout/:name", {controller: "FrackController", templateUrl: "partials/frack.html"})
        .when("/workout/:name/vma/:vma", {controller: "FrackController", templateUrl: "partials/frack.html"})
        .otherwise({controller: "FrackController", templateUrl: "partials/frack.html"});
});


app.controller("FrackController", ['$scope', '$location', '$routeParams', '$http', function($scope, $location, $routeParams, $http) {
    var workout_name = '';
    $scope.workoutdetail = {'ProgramName': "", "Workout": []}
    $scope.selectedVMA = '';
    vma = default_vma;

    if ($routeParams.name) {
        workout_name = $routeParams.name;
    }

    if ($routeParams.vma) {
        vma = $routeParams.vma;
        $scope.selectedVMA = $routeParams.vma;
    }

    res = $http.get('/rest/program/' + workout_name + "/" + vma + "/full");
	res.success(function(data, status, headers, config) {
		$scope.workoutdetail = data;
        $scope.selectedVMA = $scope.workoutdetail.TargetVMA;
	});

    res = $http.get('/rest/programs');
	res.success(function(data, status, headers, config) {
		$scope.programs = data;

        if ($scope.workoutdetail.ProgramName != '') {
            $scope.programs.forEach(function (program, i) {
                if (program.Name == $scope.workoutdetail.ProgramName) {
                    $scope.selectedWorkout = $scope.programs[i];
                }
            });
        }
	});

    $scope.submit = function() {
        $location.path("/workout/" + $scope.selectedWorkout.Name + "/vma/" + $scope.selectedVMA)
    }

}])
