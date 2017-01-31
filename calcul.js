app.controller('ProgramIndexController', ['$scope', function($scope) {
    $scope.programIndex = Object();
    for (var p of $scope.programs) {
        $scope.programIndex[p.name] = Object();
        $scope.programIndex[p.name]["name"] = p.name;
        $scope.programIndex[p.name]["totalWorkout"] = p.workouts.length;
        $scope.programIndex[p.name]["comment"] = p.comment;
        $scope.programIndex[p.name]["totalLength"] = 0;

        for (var w of p.workouts) {
            $scope.programIndex[p.name]["totalLength"] += (w.laps * w.length);
        }
        $scope.programIndex[p.name]["totalTrackLap"] = $scope.programIndex[p.name]["totalLength"] / 400;
    }
    console.debug($scope.programIndex);
}]);

app.controller('CalculController', ['$scope', function($scope) {
    var trackLength = 400;

    function calcVMADistance(vma, meters, percentage) {
        var result = "";
        var vmaMs = vma * 1000 / 3600;
        var vma100 = 100 / vmaMs;
        var calcul = vma100 / percentage * meters;
	    var stemps = calcul;
	    var minute = ((stemps - (stemps)%60) / 60);
	    var second = ((stemps % 60) * 10) / 10;

        if (minute > 0) {
            result = minute;
            result += "'";
        }
        result += Number((second).toFixed(0));
        if (minute == 0) {
            result += "s";
        }
        return result;
    }

    function calcVMASPeed(vma, percentage) {
        return (vma * percentage) / 100;
    }

    function calcPace(vitesse) {
        var ret = "";
        var e = 1 / vitesse * 60;
        var t = Math.floor(e / 60);
        var n = Math.floor(e - t * 60);
        var r = Math.round(60 * (e - t * 60 - n));
        r == 60 && (n += 1, r = 0);

	    if (n == 0 && r != 0) {
		    return r;
	    }

        ret += n + "'";
        if (r == 0) {
            return ret;
        } else if (r < 10) {
            ret += "0";
        }
        ret += r;

        return ret;
    }

    function calc(time, meters, percentage, vmas) {
        var res = new Object();
        var trackLaps = meters / trackLength; // Todo the fancy half stuff

        for (var vmaTarget of vmas) {
            res[vmaTarget] = Object();
            res[vmaTarget]['vma'] = vmaTarget; // Hack
            res[vmaTarget]['totalTime'] = calcVMADistance(vmaTarget, meters, percentage);
            if (meters >= trackLength) {
                res[vmaTarget]['timeLap'] = calcVMADistance(vmaTarget, trackLength, percentage);
            } else {
                res[vmaTarget]['timeLap'] = "NA";
            }
            res[vmaTarget]['speed'] = calcVMASPeed(vmaTarget, percentage);
            res[vmaTarget]['pace'] = calcPace(res[vmaTarget]['speed']);
        }
        return res;
    }

    $scope.test = function(r, m, p) {
        var res = calc(r, m, p, $scope.vmaWanted);
        return res;
    };

    for (var program of $scope.programs) {
        if (program.name != $scope.programWanted) continue;
        $scope.program = program;
    }
}]);
