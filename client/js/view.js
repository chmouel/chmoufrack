app.controller("ViewController", function($scope, $location, $routeParams, $http, utils, Facebook, $window) {
  $scope.facebook = {};
  $scope.facebook.info = {};
  $scope.facebook.loggedIn = false;
  $scope.facebook.ready = false;

  $scope.vmaWanted = [];
  $scope.allVMAS = utils.range(12, 22);
  $scope.rootUrl = $location.absUrl().replace($location.url(), "");

  $scope.$on('Facebook:xfbmlRender', function(ev, response) {
    $scope.facebook.ready = true;
  });

  $scope.$on('Facebook:statusChange', function(ev, response) {
    if (response.status == 'connected') {
      $scope.facebook.loggedIn = true;
      utils.FBdoLogged(response).then(function(data) {
        $scope.facebook.loggedIn = true;
        $scope.facebook.info = data;
      });
    } else {
      console.log("no login allowed");
    }
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
      $scope.vmaWanted = utils.range(sp[0], sp[1]);
    }
  } else {
    $scope.vmaWanted = $scope.allVMAS;
  }

  if (!$scope.programs) {
    var myPromise = utils.getExercises();
    myPromise.then(function(data) {
      $scope.programs = data;
      $scope.programNames = [];
      angular.forEach(data, function(p, noop) {
        $scope.programNames.push(p.name);
      });
    });
  }

  $scope.fbLogin = function() {
    utils.fbLogin();
  }

  $scope.add = function() {
    $location.path('/add');
  };

  $scope.switchtoProgram = function(tourl) {
    var t = "";
    var p = "";
    if ($scope.selectedVMA) {
      t = $scope.selectedVMA;
    } else if ($scope.vmaWanted.length > 1) {
      t = $scope.vmaWanted[0] + ":" + $scope.vmaWanted[$scope.vmaWanted.length - 1];
    } else {
      t = $scope.allVMAS[0] + ":" + $scope.allVMAS[$scope.allVMAS.length - 1]; //default
    }

    p = tourl ? tourl : $scope.selectedProgram;
    if (angular.isUndefined(p)) {
      return false;
    }

    $location.path("/workout/" + p + "/vma/" + t);

    return true; //make emacs happy, this remind me of perl modules :-[]
  };

});

app.controller('ProgramIndexController', function($scope, utils) {
    $scope.programIndex = {};
    var myPromise = utils.getExercises();
    myPromise.then(function(data) {
        angular.forEach(data, function(p, noop) {
            if (p.name !== "") {
                $scope.programIndex[p.name] = {};
                $scope.programIndex[p.name].name = p.name;
                if (p.steps) {
                    $scope.programIndex[p.name].totalWorkout = p.steps.length;
                }
                $scope.programIndex[p.name].comment = p.comment;
                $scope.programIndex[p.name].id = p.id;
                $scope.programIndex[p.name].fb = p.fb;
            }
        });
    });
});
