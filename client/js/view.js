app.controller("ViewController", function($scope, $location, $routeParams, $http, utils, Facebook, $window) {
  $scope.facebook = {};
  $scope.facebook.info = {};
  $scope.facebook.loggedIn = false;
  $scope.facebook.ready = false;

  $scope.vmaWanted = [];
  $scope.allVMAS = utils.range(12, 22);
  $scope.rootUrl = $location.absUrl().replace($location.url(), "");

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

  var refreshExercise = function(logged) {
    utils.getExercises(logged).then(function(data) {
      utils.programs = data;
      $scope.programs = utils.programs;
      $scope.programNames = [];
      angular.forEach(utils.programs, function(p, noop) {
        $scope.programNames.push(p.name);
      });
    });
  };

  if (!$scope.programs) {
    refreshExercise(false);
  }

  $scope.$on('Facebook:xfbmlRender', function(ev, response) {
    $scope.facebook.ready = true;
  });

  $scope.$on('Facebook:statusChange', function(ev, response) {
    if (response.status == 'connected') {
      $scope.facebook.loggedIn = true;
      utils.FBdoLogged(response).then(function(data) {
        $scope.facebook.loggedIn = true;
        $scope.facebook.info = data;
      }).then(function(data) {
        refreshExercise(true);
      });
    } else {
      console.log("no login allowed");
    }
  });

  $scope.fbLogin = function() {
    if ($scope.facebook.loggedIn)
      Facebook.logout();
    else
      Facebook.login();
  };

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
