app.service('$youtube', function($window) {
  var urlParams = "?" +
    "rel=0," +
    "&iv_load_policy=3" +
    "&autohide=1" +
    "&modestbranding=1" +
    "&autoplay=1" +
    "&enablejsapi=1" +
    "&playerapiid=\"ytplayer\"" +
    "&version=3" +
    "&autoplay=1"
  
  function ytPlayer() {
    return document.getElementById("myplayer");
  }
  
  this.pause = function() {
    ytPlayer().pauseVideo();
  }
  
  this.play = function() {
    ytPlayer().playVideo();
  }
  
  this.start = function(video, scope) {
    if(!ytPlayer()) {
      swfobject.embedSWF(
        "http://www.youtube.com/v/" + video.VideoId + urlParams,
        "apiplayer", "640", "360", "8", null, null,
        { allowScriptAccess: "always" },
        { id: "myplayer" }
      );
    
      $window.onYouTubePlayerReady = function(playerId) {
        ytPlayer().addEventListener("onStateChange", "stateChanged")
      }
      scope.loaded = true;
    } else {
      ytPlayer().loadVideoById(video.VideoId);
    }
  }
});
