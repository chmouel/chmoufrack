app.factory('utils', function($http, $q, Facebook) {
  var utils = {};
  utils.facebookInfo = {};

  utils.FBdoLogged = function(response) {
    var deferred = $q.defer();
    var facebookInfo = {
      'id': response.authResponse.userID,
      'accessToken': response.authResponse.accessToken
    };
    Facebook.api('/me', function(response) {
      angular.extend(facebookInfo, response);
      utils.facebookInfo = facebookInfo;
      deferred.resolve(facebookInfo);
    });
    return deferred.promise;
  };

  utils.range = function(min, max) {
    var input = [];
    min = parseInt(min);
    max = parseInt(max);
    for (var i = min; i <= max; i++) {
      input.push(i.toString());
    }
    return input;
  };

  var counter = 0;
  utils.getExercises = function() {
    // Angular $http() and then() both return promises themselves
    return $http({
      method: "GET",
      url: "/v1/exercises"
    }).then(function(result) {
      if (typeof (result.data) === 'string' &&
        result.data.trim() == "null" && counter < 3) {
        console.log("retry");
        utils.getExercises();
      }
      return result.data;
    });
  };

  utils.fbLogin = function() {
    Facebook.login();
  };

  var fbURLarg = function getURLarg() {
    var req = {
      url: 'fbID=' + utils.facebookInfo.id,
      headers: {
        'Authorization': "Bearer " + utils.facebookInfo.accessToken
      }
    };
    return req;
  };

  utils.submitExercise = function(exercise) {

    angular.forEach(exercise.steps, function(step, noop) {
      if (step.effort_type == 'time' && !angular.isUndefined(step.length))
        delete step.length;
      angular.forEach(step, function(k, v) {
        if (angular.isUndefined(k)) delete step[v];
      });
    });

    var req = fbURLarg();
    req.url = '/v1/exercise?' + req.url;
    req.method = 'POST';
    req.data = exercise;
    return $http(req);
  };

  utils.deleteExercise = function(t) {
    var req = fbURLarg();
    req.url = '/v1/exercise/' + t + '?' + req.url;
    req.method = 'DELETE';
    return $http(req);
  };

  return utils;
});

app.filter('secondsToHms', function() {
    return function(x) {
      var d = Number(x);
      var h = Math.floor(d / 3600);
      var m = Math.floor(d % 3600 / 60);
      var s = Math.floor(d % 3600 % 60);

      var hDisplay = h > 0 ? h + "h" : "";
      var mDisplay = m > 0 ? m + "mn " : "";
      var sDisplay = s > 0 ? s + "s" : "";
      return hDisplay + mDisplay + sDisplay;
    };
});

app.directive('checkEffortTime', function() {
  return {
    require: 'ngModel',
    link: function(scope, element, attr, mCtrl) {
      function myValidation(value) {
        if (value === '') {
          mCtrl.$setValidity('required', false);
          mCtrl.$setValidity('invalidEffortFormat', true);
          return;
        }
        if (value.substr(value.length-1) == "m" ||
            value.substr(value.length-1) == "s" ||
            value.substr(value.length-2) == "mn" ||
            value.substr(value.length-2) == "sec") {
          mCtrl.$setValidity('invalidEffortFormat', true);
        } else {
          mCtrl.$setValidity('required', true);
          mCtrl.$setValidity('invalidEffortFormat', false);
        }
        return value;
      }
      mCtrl.$parsers.push(myValidation);
    }
  };
});
