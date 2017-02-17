app.controller('ProgramIndexController', ['$scope', function($scope) {
    $scope.programIndex = Object();
    console.log($scope.programs);

    for (var p of $scope.programs) {
        if (p.name == "") {
            console.log("Invalid workout");
            continue;
        }
        $scope.programIndex[p.name] = Object();
        $scope.programIndex[p.name]["name"] = p.name;
        if (p.steps) {
            $scope.programIndex[p.name]["totalWorkout"] = p.steps.length;
        }
        $scope.programIndex[p.name]["comment"] = p.comment;
        $scope.programIndex[p.name]["id"] = p.id;
    }
}]);

app.controller('CalculController', ['$scope', function($scope) {
    var trackLength = 400;

    function calculDistanceForSeconds(vma, seconds, percentage) {
        var vma_ms = vma * 1000 / 3600;
        var SDist = vma_ms * seconds * percentage / 100;
	    SDist = ( Math.round(SDist * 10) ) / 10;
	    return Math.round(SDist);
    }

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

        if (second != 0) {
            if (second < 10) {
                result += 0;
            }
            result += Number((second).toFixed(0));
        } else {
            result += "00";
        }

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
            ret += "00";
            return ret;
        } else if (r < 10) {
            ret += "0";
        }
        ret += r;

        return ret;
    }

    function calc(time, meters, seconds, percentage, vmas) {
        var res = new Object();

        for (var vmaTarget of vmas) {
            res[vmaTarget] = Object();
            res[vmaTarget]['vma'] = vmaTarget; // Hack

            if (seconds > 0) {
                meters = calculDistanceForSeconds(vmaTarget, seconds, percentage);
                res[vmaTarget]['totalTime'] = meters + "m";
            } else {
                var length = calcVMADistance(vmaTarget, meters, percentage);
                res[vmaTarget]['totalTime'] = length;
            }

            var trackLaps = meters / trackLength; // Todo the fancy half stuff

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

    $scope.doCalc = function(r, m, s, p) {
        var res = calc(r, m, s, p, $scope.vmaWanted);
        return res;
    };

    for (var program of $scope.programs) {
        if (program.name != $scope.selectedProgram) continue;
        $scope.program = program;
    }
}]);
