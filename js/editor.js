app.controller("Editor", ['$scope', function($scope, $http) {
    $scope.allSteps = [
        {
            "type": "warmup",
            "effort_type": "distance",
            "effort_unit": "km",
            "distance":  "0",
            "time":  "0"
        },
        {
            "type": "interval",
            "laps": 0,
            "length": 0,
            "percentage": 0,
            "rest": "",
            "effort_type": "distance",
            "effort_unit": "km"
        },
        {
            "type": "warmdown",
            "effort_type": "distance",
            "effort_unit": "km",
            "distance":  "0",
            "time":  "0"
        }
    ];

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
        console.table($scope.allSteps);
    };

    $scope.addNewWarmupWarmdown = function() {
        $scope.allSteps.push({
            "type": "warmup",
            "effort_type": "distance",
            "effort_unit": "km",
            "distance":  "0",
            "time":  "0"
        });
    };

    $scope.addNewIntervals = function() {
        $scope.allSteps.push({
            "type": "interval",
            "laps": 0,
            "length": 0,
            "percentage": 0,
            "rest": "",
            "effort_type": "distance",
            "effort_unit": "km"
        });
    }

    $scope.remove = function(step) {
        $scope.allSteps.splice(step, step);
    }

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
