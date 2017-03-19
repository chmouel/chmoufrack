app.controller("EditController", function($scope, $http, $routeParams, utils, $location, $window) {
  $scope.forms = {};

  $scope.exercise = {};
  $scope.exercise.steps = [];

  var applyProgram = function() {
    angular.forEach(utils.programs, function(program, noop) {
      if (program.name == window.encodeURIComponent($routeParams.name)) {
        $scope.error = '';
        $scope.exercise = program;
        $window.document.title = "ChmouFrack: " + program.name;
      }
    });
  };

  $scope.$on('Facebook:statusChange', function(ev, response) {
    if (response.status == 'connected') {
      utils.FBdoLogged(response).then(function(data) {
      }).then(function(data) {
        utils.getExercises(true).then(function(data) {
          applyProgram();
        });
      });
    } else {
      console.log("no login allowed");
    }
  });

  if (angular.isUndefined($scope.exercise.length) &&
      angular.isDefined(utils.facebookInfo.email) &&
      utils.programs.length !== 0)
    applyProgram();

  $scope.effortDistanceUnits = [
    {
      name: 'Mètres',
      value: 'm'
    },
    {
      name: 'Kilomètre',
      value: 'km'
    }
  ];

  $scope.effortTypeOptions = [
    {
      name: 'Distance',
      value: 'distance'
    },
    {
      name: 'Temps',
      value: 'time'
    }
  ];

  $scope.stepTypeOptions = [
    {
      name: 'Warmup',
      value: 'warmup'
    },
    {
      name: 'Warmdown',
      value: 'warmdown'
    }
  ];

  var showError = function(msg) {
    $scope.error = msg;
    $scope.success = false;
    $window.scrollTo(0, 0);
  };

  $scope.submit = function() {
    var err = null;
    angular.forEach($scope.forms, function(p, noop) {
      if (Object.keys(p.$error).length === 0)
        return;
      angular.forEach(p.$error, function(e, noop) {
        angular.forEach(e, function(a, noop) {
          // NOTE: if we chose time and we have a distance not defined then it's fine
          if (a.$name == 'effortd' && a.$$parentForm.unit.$viewValue == 'time') {
            return;
          }
          // NOTE: same as before
          if (a.$name == 'effortt' && a.$$parentForm.unit.$viewValue == 'distance') {
            return;
          }
          err = a;
        });
      });
    });

    if (err) {
      showError("Vous avez une erreur a corrigez");
      return;
    }

    if (!$scope.facebook.loggedIn) {
      showError("Vous devez être connecter a facebook pour pouvoir enregistrer un programme.");
      return;
    }

    // NOTE(chmou): If a rename delete the old one
    if (angular.isDefined($routeParams.name) && $scope.exercise.name != $routeParams.name) {
      var exist = false;
      angular.forEach($scope.programs, function(p, noop) {
        if ($scope.exercise.name == p.name)
          exist = true;
      });
      if (exist) {
        showError('Vous ne pouvez pas renommer un programme vers ' +
                  'le nom <b>' + $scope.exercise.name +
                  '</b> qui existe, <a href="/#!/edit/' +
                  $scope.exercise.name +'">supprimez</a> le d\'abord.');
        return;
      }
      $scope.delete($routeParams.name);
    }


    utils.submitExercise($scope.exercise).then(function(result) {
      $scope.error = false;
      $scope.success = 'Votre programme a était sauvegarder, le <a href="/#!/workout/' + $scope.exercise.name + '"><strong>voir ici</strong></a>';
      $window.scrollTo(0, 0);
    }).catch(function(error) {
      $scope.success = false;
      $scope.error = "Error while saving: (" + error.status + "): " + error.data.error;
      $window.scrollTo(0, 0);
    });
  };

  $scope.delete = function(t, r) {
    if (!t) {
      t = $scope.exercise.name;
    }

    utils.deleteExercise(t).then(function(result) {
      $location.path('/');
    }).catch(function(error) {
      $scope.success = false;
      $scope.error = error;
      $window.scrollTo(0, 0);
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
    if (!arr.steps) {
      arr.steps = [];
    }
    arr.steps.splice(index, 1);
  };

  $scope.addNewWarmupWarmdown = function(t, arr) {
    if (!arr.steps) {
      arr.steps = [];
    }
    arr.steps.push({
      "type": t,
      "effort_type": "distance",
      "effort": ""
    });
  };

  $scope.addNewRepeat = function(arr) {
    if (!arr.steps) {
      arr.steps = [];
    }

    arr.steps.push({
      "type": "repeat",
      "repeat": {
        steps: [],
        repeat: 0
      }
    });
  };

  $scope.addNewIntervals = function(arr) {
    if (!arr.steps) {
      arr.steps = [];
    }

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

app.directive('selectOnClick', function() {
  return {
    restrict: 'A',
    link: function(scope, element) {
      var focusedElement;
      element.on('click', function() {
        if (focusedElement != this) {
          this.select();
          focusedElement = this;
        }
      });
      element.on('blur', function() {
        focusedElement = null;
      });
    }
  };
});
