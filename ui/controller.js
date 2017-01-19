var app = angular.module("Frack", ["ngRoute"]);
var default_vma = '14:19'

app.config(function($routeProvider) {
    $routeProvider
        .when("/workout/:name", {controller: "FrackController", templateUrl: "partials/frack.html"})
        .when("/workout/:name/vma/:vma", {controller: "FrackController", templateUrl: "partials/frack.html"})
        .otherwise({controller: "FrackController", templateUrl: "partials/frack.html"});
});

app.filter('range', function() {
  return function(input, min, max) {
    min = parseInt(min);
    max = parseInt(max);
    for (var i=min; i<=max; i++)
      input.push(i);
    return input;
  };
});

app.controller("FrackController", ['$scope', '$location', '$routeParams', '$http', function($scope, $location, $routeParams, $http) {
    var workout_name = '';
    $scope.workoutdetail = {'Name': "", "Workout": []}
    $scope.selectedVMA = '';
    vma = default_vma;

    if ($routeParams.name) {
        workout_name = $routeParams.name;
    }

    if ($routeParams.vma == "undefined") {
        $location.path("/workout/" + workout_name + "/vma/" + default_vma)
        return;
    }

    if ($routeParams.vma) {
        vma = $routeParams.vma;
        $scope.selectedVMA = vma;
    }

    res = $http.get('/rest/program/' + workout_name + "/" + vma);
	res.success(function(data, status, headers, config) {
		$scope.workoutdetail = data;
        $scope.selectedVMA = $scope.workoutdetail.TargetVMA;
	});

    res = $http.get('/rest/programs');
	res.success(function(data, status, headers, config) {
		$scope.programs = data;

        if ($scope.workoutdetail.Name != '') {
            $scope.programs.forEach(function (program, i) {
                if (program.Name == $scope.workoutdetail.Name) {
                    $scope.selectedWorkout = $scope.programs[i];
                }
            });
        }
	});

    $scope.submit = function() {
        $location.path("/workout/" + $scope.selectedWorkout.Name + "/vma/" + $scope.selectedVMA)
    }

}])
