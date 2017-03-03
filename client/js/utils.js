app.factory('utils', function($http, $q, Facebook) {
    var utils = {};
    utils.facebookInfo = {};

    utils.FBdoLogged = function(response){
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
        for (var i=min; i<=max; i++)
            input.push(i.toString());
        return input;
    };

    var counter = 0;
    utils.getExercises = function() {
        // Angular $http() and then() both return promises themselves
        return $http({method:"GET", url:"/v1/exercises"}).then(function(result){
            if (typeof(result.data) === 'string' &&
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
            }};
        return req;
    };

    utils.submitExercise = function(exercise) {
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