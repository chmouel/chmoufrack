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

    $http({
        method: 'GET',
        url: '/rest/program/' + $routeParams.name + '/workouts'
    }).then(function successCallback(response) {
		$scope.programDetails = response.data;
    }, function errorCallback(response) {
        $scope.AddNewProgram = true;
        $scope.addNew();
    });

    $("#wrapper").addClass("toggled");

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

    $scope.submit = function(programDetail){
        if ($scope.AddNewProgram) {
            $http({
                method: 'POST',
                url: '/rest/program/' + $scope.programName
            }).then(function successCallback(response) {
            }, function errorCallback(response) {
                console.debug("Failed to create new program");
            });
            $scope.AddNewProgram = false;
        } else {
            $http({
                method: 'DELETE',
                url: '/rest/program/' + $scope.programName + '/purge'
            }).then(function successCallback(response) {
                console.debug("Success purging program");
            }, function errorCallback(response) {
                console.debug("Failed to delete program");
            });
        }
        $http.post('/rest/program/' + $scope.programName + '/workouts',
                   $scope.programDetails);

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
    var res = $http.get('/rest/programs');
	res.success(function(data, status, headers, config) {
		$scope.workoutS = data;
	});

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
