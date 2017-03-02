app.controller("EditController", function($scope, $http, $routeParams, utils, userInfo, $location) {
    $scope.exercise = {};
    $scope.exercise.steps = [];

    if ($routeParams.name) {
        var res = $http.get('/v1/exercise/' + $routeParams.name );
	    res.then(function(response) {
            $scope.exercise = response.data;
	    }, function errorCallBack(response) {
            $scope.exercise = {};
            $scope.exercise.steps = {};
            $scope.NotFound = $routeParams.name;
        });
    }

    $scope.effortDistanceUnits = [
        { name: 'Mètres', value: 'm'},
        { name: 'Kilomètre', value: 'km'}
    ];

    $scope.effortTypeOptions = [
        { name: 'Distance', value: 'distance' },
        { name: 'Temps', value: 'time' }
    ];

    $scope.stepTypeOptions = [
        { name: 'Warmup', value: 'warmup' },
        { name: 'Warmdown', value: 'warmdown' }
    ];

    $scope.submit = function() {
        //TODO(chmou): proper error showing (just in case someone wonder we do
        //proper check in API server too)
        if (!$scope.fbLogged) {
            console.error("You should have login to be able to do this here?");
            return;
        }

        // NOTE(chmou): If a rename delete the old one
        if ($scope.exercise.name != $routeParams.name) {
            $scope.delete($routeParams.name);
        }

        utils.submitExercise($scope.exercise).then(function (result) {
            $location.path('/workout/' + $scope.exercise.name);
        });
    };

    $scope.delete = function(t, r) {
        if (!t) t = $scope.exercise.name;

        utils.deleteExercise(t).then(function (result) {
            $location.path('/');
        });
    };

    $scope.swapUp = function(index, arr) {
        var item = arr.steps.splice(index, 1);
        arr.steps.splice(index - 1, 0, item[0]);
    };

    $scope.swapDown = function(index, arr) {
        var item = arr.steps.splice(index, 1);
        arr.steps.splice(index + 1, 0, item[0]);
    };

    $scope.removeStep = function(index, arr) {
        if (!arr.steps)
            arr.steps = [];
        arr.steps.splice(index, 1);
    };

    $scope.addNewWarmupWarmdown = function(t, arr) {
        if (!arr.steps)
            arr.steps = [];

        arr.steps.push({
            "type": t,
            "effort_type": "distance",
            "effort":  ""
        });
    };

    $scope.addNewRepeat = function(arr) {
        if (!arr.steps)
            arr.steps = [];

        arr.steps.push({
            "type": "repeat",
            "repeat": {steps: [], repeat: 0}
        });
    };

    $scope.addNewIntervals = function(arr) {
        if (!arr.steps)
            arr.steps = [];

        arr.steps.push({
            "type": "interval",
            "laps": "",
            "length": "",
            "percentage": "",
            "rest": "",
            "effort_type": "distance"
        });
    };


});

app.directive('selectOnClick', function () {
    return {
        restrict: 'A',
        link: function (scope, element) {
            var focusedElement;
            element.on('click', function () {
                if (focusedElement != this) {
                    this.select();
                    focusedElement = this;
                }
            });
            element.on('blur', function () {
                focusedElement = null;
            });
        }
    };
});
