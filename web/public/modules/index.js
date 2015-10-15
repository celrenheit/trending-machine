angular.module('tm',[
  'ng',
	'ngAria',
	'ngAnimate',
	'ngMaterial',

	'ui.router',
  'cfp.hotkeys',
  'mdPickers',
	'socialLinks',

  'tm.common'
])


.config(($stateProvider, $urlRouterProvider, $locationProvider) => {
	$stateProvider
	.state('home', {
    url: "/snapshots?date&l",
    templateUrl: 'modules/index.tpl.html',
    controller: "SnapshotsCtrl",
    resolve: {
      snapshot: (SnapshotsService, $stateParams, ToastService) => SnapshotsService.get($stateParams.date ? moment($stateParams.date, "DD-MM-YYYY") : new Date())
                                                                                  .then((s) => s)
                                                                                  .catch((res) => res)

    }
  })

	$urlRouterProvider.otherwise('/snapshots');
	$locationProvider.hashPrefix('!')
})

.run(($rootScope, $mdMedia) => {
  $rootScope.$mdMedia = $mdMedia
})

.controller('AppCtrl', () => {

})


.controller('SnapshotsCtrl', (SnapshotsService, $scope, snapshot, $mdDatePicker, ToastService, $state, $stateParams, $anchorScroll, $location, $timeout, SocialSharingService) => {
  $scope.showPicker = (ev) => {
    if(is.mobile())
      return;
  	$mdDatePicker(ev, $scope.selectedDate).then((selectedDate) => {
      $scope.selectedDate = selectedDate;
    });
  }

  $scope.$watch('selectedDate', (newDate, oldDate) => {
    if(newDate === oldDate)
      return;

    $state.go('home', {
      date: moment(newDate).format("DD-MM-YYYY")
    })
  })

  if(snapshot.status && snapshot.status != 200) {
    if(snapshot.status === 404) {
      ToastService.simple('Nothing has been found for this date')
    } else {
      ToastService.simple(snapshot.data.error)
    }
    return;
  }
  $scope.snapshot = snapshot
  $scope.selectedDate = moment($stateParams.date, "DD-MM-YYYY").toDate() || new Date()

  var sortLanguages = (langs) => {
    return Object.keys(langs).sort((a,b) => a[0] === "_" || b[0] === "_" || a > b)
  }
  $scope.languages = sortLanguages($scope.snapshot.languages)
  $scope.selectedLang = $stateParams.l && $scope.languages.indexOf($stateParams.l) != -1 ? $stateParams.l : $scope.languages[0]

  $scope.$watch('selectedLang', (newLang, oldLang) => {
    // $state.transitionTo('home', {l: newLang});
    $stateParams.l = newLang
  })

  $scope.share = ($ev, pos) => {
    let url = window.location.origin + "/" +$state.href("home", $stateParams) + "#pos-" + pos
    SocialSharingService.showModal($ev, url)
  }

  $scope.goHome = () =>$state.go('home', {date: moment().format("DD-MM-YYYY")})
  $scope.gotoAnchor = function(x) {
    var newHash = 'pos-' + x;
    if ($location.hash() !== newHash) {
      // set the $location.hash to `newHash` and
      // $anchorScroll will automatically scroll to it
      $location.hash('pos-' + x);
    } else {
      // call $anchorScroll() explicitly,
      // since $location.hash hasn't changed
      $anchorScroll();
    }
  };
  
  $scope.$watch('$location.hash', function () {
      _.defer($anchorScroll);
  });
})

.factory('SnapshotsService', ($http) => {
  return {
    get(date = new Date()) {
      let formatted = moment(date).format("DD-MM-YYYY");
      return new Promise((resolve, reject) => {
        $http.get('/api/snapshots', {
          params:{
            date: formatted
          }
        })
            .then((res) => resolve(res.data),
                  (err) => reject(err))
      })
    }
  }
})

.factory('ToastService', function($mdToast) {
	var toastPosition = {
	    bottom: false,
	    top: true,
	    left: false,
	    right: true
	  };
	  var getToastPosition = function() {
	    return Object.keys(toastPosition)
	      .filter(function(pos) { return toastPosition[pos]; })
	      .join(' ');
	  };
	return {
		simple: function(content, delay) {
			console.log('notif', content, getToastPosition());
			return $mdToast.show(
					$mdToast.simple()
			        .content(content)
			        .position(getToastPosition())
			        .hideDelay(delay || 3000)
			        )
		}
	}
})
