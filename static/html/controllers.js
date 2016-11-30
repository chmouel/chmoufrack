var app = angular.module("EditWorkout", ["ngRoute"]);

app.config(function($routeProvider) {
    $routeProvider.
        when("/:name", {
            templateUrl : "editworkout.html",
            controller: DetailController
        })
});

function DetailController($scope, $routeParams, $http) {
    $scope.message = $routeParams.id;

    var res = $http.get('/program/' + $routeParams.name + "/workouts");
	res.success(function(data, status, headers, config) {
		$scope.programDetails = data;
        $scope.programName = $routeParams.name;
	});;

    $("#wrapper").addClass("toggled");

    $scope.addNew = function(programDetail){
        console.debug($scope.programDetails);
        if ($scope.programDetails == "null\n") { //wtf
            $scope.programDetails = [];
        }

        $scope.programDetails.push({
            'repetition': "",
            'meters': "",
            'repos': "",
            'programname': "",
        });
    };

    $scope.post = function(programDetail){
        var res = $http.delete('/program/' + $scope.programName + '/purge')
		res.success(function(data, status, headers, config) {
			console.debug(data);
		});;

        var res = $http.post('/program/' + $scope.programName + '/workouts', $scope.programDetails);
		res.success(function(data, status, headers, config) {
			console.debug(data);
		});;
    };

    $scope.remove = function(){
        var newDataList=[];
        $scope.selectedAll = false;
        angular.forEach($scope.programDetails, function(selected){
            if(!selected.selected){
                newDataList.push(selected);
            }
        });
        $scope.programDetails = newDataList;
    };

    $scope.checkAll = function () {
        if (!$scope.selectedAll) {
            $scope.selectedAll = true;
        } else {
            $scope.selectedAll = false;
        }
        angular.forEach($scope.programDetails, function(programDetail) {
            programDetail.selected = $scope.selectedAll;
        });
    };

}

app.controller("ListController", ['$scope', '$http', function($scope, $http) {
    var res = $http.get('/programs');
	res.success(function(data, status, headers, config) {
		$scope.workoutS = data;
	});;

}]);
