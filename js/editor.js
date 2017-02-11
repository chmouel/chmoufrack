app.controller("Editor", ['$scope', function($scope, $http) {
    $scope.Excercise = [
        {
            name: "",
            comment: "",
            steps: [
                {
                    "type": "warmup",
                    "effort_type": "distance",
                    "effort_unit": "km",
                    "effort":  ""
                },
                {
                    "type": "interval",
                    "laps": "",
                    "length": "",
                    "percentage": "",
                    "rest": "",
                    "effort_type": "distance",
                    "effort_unit": "km"
                },
                {
                    "type": "warmdown",
                    "effort_type": "distance",
                    "effort_unit": "km",
                    "effort":  ""
                }
            ]
        }
    ];
    $scope.excercise = $scope.Excercise[0];

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
        console.log($scope.excercise);
        console.log($scope.excercise.steps[0]);
        console.log($scope.excercise.steps[1]);
        console.log($scope.excercise.steps[2]);
    };

    $scope.addNewWarmupWarmdown = function() {
        $scope.excercise.steps.push({
            "type": "warmup",
            "effort_type": "distance",
            "effort_unit": "km",
            "effort":  ""
        });
    };

    $scope.addNewIntervals = function() {
        $scope.excercise.steps.push({
            "type": "interval",
            "laps": "",
            "length": "",
            "percentage": "",
            "rest": "",
            "effort_type": "distance",
            "effort_unit": "km"
        });
    };

    $scope.removeStep = function(step) {
        $scope.excercise.steps.splice(step, 1);
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
