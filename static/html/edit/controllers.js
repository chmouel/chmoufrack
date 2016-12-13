var app = angular.module("EditWorkout", ["ngRoute"]);

app.config(function($routeProvider) {
    $routeProvider.
        when("/:name", {
            templateUrl: "partials/editworkout.html",
            controller: DetailController
        })
        .otherwise({
            templateUrl: "partials/default.html"
        })

});

function DetailController($scope, $routeParams, $http) {
    $scope.message = $routeParams.id;
    $scope.AddNewProgram = false;
    $scope.programName = $routeParams.name;
    $("#wrapper").addClass("toggled");

    $http({
        method: 'GET',
        url: '/rest/program/' + $routeParams.name + '/workouts'
    }).then(function successCallback(response) {
		$scope.programDetails = response.data;
    }, function errorCallback(response) {
        $scope.AddNewProgram = true;
        $scope.addNew();
    });

    $scope.addNew = function(programDetail){
        if ($scope.programDetails == "null\n" || typeof($scope.programDetails) === "undefined") { //wtf
            $scope.programDetails = [];
        }

        $scope.programDetails.push({
            'repetition': "",
            'meters': "",
            'repos': "",
            'programname': "",
        });
    };

    $scope.save = function(programDetail){
        if ($scope.AddNewProgram) {
            var firstStep = $http({method: 'POST', url: '/rest/program/' + $scope.programName})
            firstStep.message = "create program"
            firstStep.thenFunction = function() {
                //$scope.workoutS.push({'Comment': '', 'Date': '', 'ID': '', 'Name': $scope.programName});
                $scope.refreshPrograms();
            }
            $scope.AddNewProgram = false;
        } else {
            var firstStep = $http({method: 'DELETE', url: '/rest/program/' + $scope.programName + '/purge'})
            firstStep.message = "purging program"
            firstStep.thenFunction = function() {};
        }

        firstStep.then(function successCallback(response) {
            console.debug($scope.workoutS);
            $http.post('/rest/program/' + $scope.programName + '/workouts', $scope.programDetails);
            firstStep.thenFunction();
        }, function errorCallback(response) {
            console.debug("Failing " + firstStep.message);
        });
    };

    $scope.saveshow = function(){
        $scope.save();
        $(location).attr('href', '/html/program/' + $scope.programName);
    }

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
    $scope.refreshPrograms = function() {
        var res = $http.get('/rest/programs');
	    res.success(function(data, status, headers, config) {
		    $scope.workoutS = data;
	    });
    }
    $scope.refreshPrograms();

    $scope.addNewProgram = function() {
      bootbox.prompt("Add a new Program", (function (program) {
          if (program === null || program == "") {
              return 1;
          }
        $(location).attr('href', '#' + program);
      }))
     }

    $scope.toggleSidebar = function() {
        $("#wrapper").toggleClass("toggled");
    }

}]);
