var app = angular.module('app', ['ngResource']);

function YouskiCtrl($scope, $resource, $youtube) {
  var Entry   = $resource("/search");
  var Related = $resource("/related");
  var related;
  
  var nextRelatedVideo = (function() {
    var index = 0;
    var helper = function() {
      var next = related[index];
      index += 1;
      return next;
    };
    return helper;
  }());
  
  function dequeue(video) {
    var i = $scope.queue.indexOf(video);
    $scope.queue.remove(i);
  }
  
  $scope.paused = false;
  
  $scope.dequeueAndPlay = function(queued) {
    dequeue(queued);
    $youtube.start(queued);
  }
  
  $scope.pause = function() {
    $youtube.pause()
    $scope.paused = true;
  }
  
  $scope.pick = function(entry) {
    $youtube.start(entry, $scope);
    
    related = Related.query({videoId:entry.VideoId}, function() {
      $scope.queue = related.slice(0, 20);
    });
  }
  
  $scope.play = function() {
    $youtube.play()
    $scope.paused = false;
  }
  
  $scope.next = function() {
    var next = nextRelatedVideo();
    
    dequeue(next);
    $youtube.start(next);
    $scope.paused = false;
  }
  
  $scope.replay = function() {
    $youtube.replay()
  }
  
  $scope.search = function(query) {
    $scope.entries = Entry.query({query:query});
  };
};
