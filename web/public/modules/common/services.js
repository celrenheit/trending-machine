angular.module('tm.common', [])


.factory('SocialSharingService', function ($mdDialog) {
	return {
		showModal: (ev, url) => {
			$mdDialog.show({
				controller: ($scope) => {
					$scope.url = url
					$scope.hide = function() {
				    $mdDialog.hide();
				  };
				  $scope.cancel = function() {
				    $mdDialog.cancel();
				  };
				  $scope.answer = function(answer) {
				    $mdDialog.hide(answer);
				  };
				},
				templateUrl: 'modules/common/social-sharing.modal.tpl.html',
				parent: angular.element(document.body),
				targetEvent: ev,
				locals: {
					url: url
				}
			})
			.then(function(answer) {
			}, function() {
			});
		},
	}
})
