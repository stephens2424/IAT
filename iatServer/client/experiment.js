(function () {

  var experiment,
      experimentRequest = $.getJSON("expdata/sample");

  var KEYS = {
    "Left": 37,
    "Right": 39
  }

  var Experiment = function (frames, subjectID, experimentNode) {
    this.node = experimentNode;
    this.currentFrame = 0;
    this.frames = frames;
    this.advanceTimes = [];
    this.responses = [];
    this.subjectID = subjectID;
  }

  Experiment.prototype.advance = function () {
    this.currentFrame += 1;
    if (this.currentFrame < this.frames.length) {
      this.node.innerHTML = this.frames[this.currentFrame].HTML;
      this.advanceTimes.push(performance.now());
    } else {
      var $form = $("<form action='/postResults' method='POST'>");
      $form.append("<input type='hidden' name='subjectID' value='"+ this.subjectID + "'>");
      var $timeInput = $("<input type='hidden' name='times[]'>");
      var $respInput = $("<input type='hidden' name='resp[]'>");
      $.each(this.advanceTimes, function (i, time) {
        $form.append($timeInput.clone().attr('value', time));
      });
      $.each(this.responses, function (i, resp) {
        $form.append($respInput.clone().attr('value', resp));
      });
      $form.submit();
    }
  }

  $('body').on("keydown", function (e) {
    switch (e.which) {
      case KEYS.Left:
        experiment.responses.push("l");
        experiment.advance();
        break;
      case KEYS.Right:
        experiment.responses.push("r");
        experiment.advance();
        break;
    }
  });

  $(function () {
    var experimentNode = document.getElementById("experiment");
    experimentRequest.done(function (experimentData) {
      experiment = new Experiment(experimentData.Frames, experimentData.SubjectID, experimentNode);
      experiment.advance();
    });
  });
})();
