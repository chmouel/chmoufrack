var app = angular.module("myapp", []);
app.controller("ListController", ['$scope', '$http', function($scope, $http) {
    $scope.personalDetails = [
        {
            'repetition':'3',
            'meters':'400',
            'percentage':'100',
            'repos': '3 minutes tranquillou poto!',
            'programname': 'ProgramREST4963',
        },
        {
            'repetition':'2',
            'meters':'800',
            'percentage':'95',
            'repos': '3 minutes tranquillou poto!',
            'programname': 'ProgramREST4963',
        },
        {
            'repetition':'1',
            'meters':'1000',
            'percentage':'90',
            'repos': '3 minutes tranquillou poto!',
            'programname': 'ProgramREST4963',
        }];

        $scope.addNew = function(personalDetail){
            $scope.personalDetails.push({
                'repetition': "",
                'meters': "",
                'repos': "",
                'programname': "",
            });
            console.debug($scope);
        };

        $scope.post = function(personalDetail){
            console.debug($scope.personalDetails);
            var res = $http.post('/workouts', $scope.personalDetails);
		    res.success(function(data, status, headers, config) {
			    console.debug(data);
		    });;
        };

        $scope.remove = function(){
            var newDataList=[];
            $scope.selectedAll = false;
            angular.forEach($scope.personalDetails, function(selected){
                if(!selected.selected){
                    newDataList.push(selected);
                }
            });
            $scope.personalDetails = newDataList;
        };

    $scope.checkAll = function () {
        if (!$scope.selectedAll) {
            $scope.selectedAll = true;
        } else {
            $scope.selectedAll = false;
        }
        angular.forEach($scope.personalDetails, function(personalDetail) {
            personalDetail.selected = $scope.selectedAll;
        });
    };


}]);
