var app = angular.module("Frack", ["ngRoute", "ngSanitize", "ngFacebook"]);
app.config(['$facebookProvider', function($facebookProvider) {
    $facebookProvider.setAppId('3518596602').setPermissions(['email']);
}]);

app.run(['$rootScope', '$window', function($rootScope, $window) {
    (function(d, s, id) {
      var js, fjs = d.getElementsByTagName(s)[0];
      if (d.getElementById(id)) return;
      js = d.createElement(s); js.id = id;
      js.src = "//connect.facebook.net/en_US/sdk.js";
      fjs.parentNode.insertBefore(js, fjs);
    }(document, 'script', 'facebook-jssdk'));
    $rootScope.$on('fb.load', function() {
      $window.dispatchEvent(new Event('fb.load'));
    });
}]);

app.config(function($routeProvider) {
    $routeProvider
        .when("/add", {controller: "EditController", templateUrl: "html/edit/editor.html"})
        .when("/edit/:name", {controller: "EditController", templateUrl: "html/edit/editor.html"})
        .when("/workout/:name", {controller: "ViewController", templateUrl: "html/view/view.html"})
        .when("/workout/:name/vma/:vma", {controller: "ViewController", templateUrl: "html/view/view.html"})
        .otherwise({controller: "ViewController", templateUrl: "html/view/selection.html"});
});

function range(min, max) {
    var input = [];
    min = parseInt(min);
    max = parseInt(max);
    for (var i=min; i<=max; i++)
        input.push(i.toString());
    return input;
};

// Disable caching: https://goo.gl/yHW1vE
app.config(['$httpProvider', function($httpProvider) {
    //initialize get if not there
    if (!$httpProvider.defaults.headers.get)
        $httpProvider.defaults.headers.get = {};

    $httpProvider.defaults.headers.get['If-Modified-Since'] = 'Mon, 26 Jul 1997 05:00:00 GMT';
    $httpProvider.defaults.headers.get['Cache-Control'] = 'no-cache';
    $httpProvider.defaults.headers.get['Pragma'] = 'no-cache';
}]);

//http://stackoverflow.com/a/36254259
app.directive('input', [function() {
    return {
        restrict: 'E',
        require: '?ngModel',
        link: function(scope, element, attrs, ngModel) {
            if (
                   'undefined' !== typeof attrs.type
                && 'number' === attrs.type
                && ngModel
            ) {
                ngModel.$formatters.push(function(modelValue) {
                    return Number(modelValue);
                });

                ngModel.$parsers.push(function(viewValue) {
                    return Number(viewValue);
                });
            }
        }
    };
}]);

app.factory('userInfo', function($facebook, $q) {
    function UserInfoService() {
        var self = this;

        self.getURLarg = function() {
            var deferred = $q.defer();
            self.get().then(function(u) {
                    var req = {
                        url: 'fbID=' + u.id,
                        headers: {
                            'Authorization': "Bearer " + u.auth.accessToken
                        }
                    };
                deferred.resolve(req);
            });
            return deferred.promise;
        };

        self.get = function() {
            var userInfo = new Object();
            var deferred = $q.defer();
            $facebook.cachedApi('/me').then(function(user) {
                var auth = $facebook.getAuthResponse();
                userInfo = user;
                userInfo['auth'] = auth;
                deferred.resolve(userInfo);
            });

            return deferred.promise;
        };
    };
    return new UserInfoService();
});

app.factory('rest', function($http, userInfo) {
    var deleteExercise = function(t) {
        return userInfo.getURLarg().then(
            function(req) {
                req['url'] = '/v1/exercise/' + t + '?' + req.url;
                req['method'] = 'DELETE';
                return $http(req);
            });
    };

    var counter = 0;
    var getExercises = function() {
        // Angular $http() and then() both return promises themselves
        return $http({method:"GET", url:"/v1/exercises"}).then(function(result){
            if (typeof(result.data) === 'string' &&
                result.data.trim() == "null" && counter < 3) {
                console.log("retry");
                getExercises();
            };
            return result.data;
        });
    };

    var submitExercise = function(exercise) {
        return userInfo.getURLarg().then(
            function(req) {
                req['url'] = '/v1/exercise?' + req.url;
                req['method'] = 'POST';
                req['data'] = exercise;
                return $http(req);
            });
    };
    return {
        getExercises: getExercises,
        deleteExercise: deleteExercise,
        submitExercise: submitExercise
    };
});
