app.controller("ViewController", function($scope, $location, $routeParams, $http, rest, $facebook, userInfo, $window) {
    $scope.vmaWanted = [];
    $scope.allVMAS = range(12, 22);
    $scope.rootUrl = $location.absUrl().replace($location.url(), "");

    var item = quoteSource[Math.floor(Math.random()*quoteSource.length)];
    $scope.quote = item.quote;
    $scope.quote_author = item.name;


    $scope.fbLogged = $window.sessionStorage.getItem('fbUserInfo');
    if ($scope.fbLogged) {
        $scope.fbLogged = JSON.parse($scope.fbLogged);
    }

    $scope.fbLogin = function() {
        $facebook.login();
    };

    $scope.fbLogout = function() {
        $facebook.logout();
        $window.sessionStorage.removeItem("fbUserInfo");
        $scope.fbLogged = null;
    };

    $scope.$on('fb.auth.login', function(event, userDetails) {
        userInfo.get().then(function(u) {
            $scope.fbLogged=u;
            $window.sessionStorage.setItem('fbUserInfo', JSON.stringify(u));
        });
    });

    if ($routeParams.name) {
        $scope.selectedProgram = $routeParams.name;
    }

    // By default use 12-18, if we have a 12:18 in the router config parse it to
    // range int 12.13.14....18...etc. if only one number then pass it straight
    if ($routeParams.vma) {
        if ($routeParams.vma.indexOf(":") == -1) {
            $scope.selectedVMA = $routeParams.vma;
            $scope.vmaWanted = [parseInt($routeParams.vma)];
        } else {
            var sp = $routeParams.vma.split(':');
            $scope.vmaWanted = range(sp[0], sp[1]);
        }
    } else {
        $scope.vmaWanted = $scope.allVMAS;
    }

    if (!$scope.programs) {
        var myPromise = rest.getExercises();
        myPromise.then(function(data) {
            $scope.programs = data;
            $scope.programNames = [];
            for (var p of $scope.programs)
                $scope.programNames.push(p.name);
        });
    };

    $scope.submit = function(tourl) {
        var t = "", p ="";
        if ( $scope.selectedVMA) {
            t = $scope.selectedVMA;
        } else if ($scope.vmaWanted.length > 1) {
            t = $scope.vmaWanted[0] + ":" + $scope.vmaWanted[$scope.vmaWanted.length - 1];
        } else {
            t = $scope.allVMAS[0] + ":" + $scope.allVMAS[$scope.allVMAS.length - 1];; //default
        }

        p = tourl ? tourl : $scope.selectedProgram;
        if (angular.isUndefined(p))
            return false;

        $location.path("/workout/" + p + "/vma/" + t);

        return true; //make emacs happy, this remind me of perl modules :-[]
    };

});
