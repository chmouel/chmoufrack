app.controller("EditController", ['$scope', '$http', '$routeParams', function($scope, $http, $routeParams) {
    $scope.exercise = Object();
    $scope.exercise.steps = [];

    if ($routeParams.name) {
        var res = $http.get('/v1/exercise/' + $routeParams.name );
	    res.then(function(response) {
            $scope.exercise = response.data;
	    }, function errorCallBack(response) {
            $scope.exercise = new Object();
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
        console.log($scope.exercise);
        var res = $http.post('/v1/exercise', $scope.exercise);
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
        arr.steps.splice(index, 1);
    };

    $scope.addNewWarmupWarmdown = function(t, arr) {
        arr.steps.push({
            "type": t,
            "effort_type": "distance",
            "effort":  ""
        });
    };

    $scope.addNewRepeat = function(arr) {
        arr.steps.push({
            "type": "repeat",
            "repeat": {steps: [], repeat: 0}
        });
    };


    $scope.addNewIntervals = function(arr) {
        arr.steps.push({
            "type": "interval",
            "laps": "",
            "length": "",
            "percentage": "",
            "rest": "",
            "effort_type": "distance"
        });
    };


}]);

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
