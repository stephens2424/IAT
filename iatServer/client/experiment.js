(function () {

  var experiment,
      experimentRequest = $.getJSON("expdata/sample");

  var KEYS = {
    "Left": 37,
    "Right": 39
  }

  var Experiment = function (frames, experimentNode) {
    this.node = experimentNode;
    this.currentFrame = 0;
    this.frames = frames;
    this.advanceTimes = [];
  }

  Experiment.prototype.advance = function () {
    this.currentFrame += 1;
    if (this.currentFrame < this.frames.length) {
      this.node.innerHTML = this.frames[this.currentFrame].HTML;
      this.advanceTimes.push(performance.now());
    } else {
      var $form = $("<form action='/postResults' method='POST'>");
      var $input = $("<input type='hidden' name='times[]'>");
      $.each(this.advanceTimes, function (i, time) {
        $form.append($input.clone().attr('value', time));
      });
      $form.submit();
    }
  }

  $('body').on("keydown", function (e) {
    switch (e.which) {
      case KEYS.Left:
        experiment.advance();
        break;
      case KEYS.Right:
        experiment.advance();
        break;
    }
  });

  $(function () {
    var experimentNode = document.getElementById("experiment");
    experimentRequest.done(function (experimentData) {
      experiment = new Experiment(experimentData.Frames, experimentNode);
      experiment.advance();
    });
  });
})();
