package convert

// See https://github.com/catapult-project/catapult/blob/master/tracing/docs/embedding-trace-viewer.md
// This is almost verbatim copy of:
// https://github.com/catapult-project/catapult/blob/master/tracing/bin/index.html
// on revision 5f9e4c3eaa555bdef18218a89f38c768303b7b6e.
var templTrace = `
<html>
<head>
<link href="/trace_viewer_html" rel="import">
<style type="text/css">
  html, body {
    box-sizing: border-box;
    overflow: hidden;
    margin: 0px;
    padding: 0;
    width: 100%;
    height: 100%;
  }
  #trace-viewer {
    width: 100%;
    height: 100%;
  }
  #trace-viewer:focus {
    outline: none;
  }
</style>
<script>
'use strict';
(function() {
  var viewer;
  var url;
  var model;

  function load() {
    var req = new XMLHttpRequest();
    var is_binary = /[.]gz$/.test(url) || /[.]zip$/.test(url);
    req.overrideMimeType('text/plain; charset=x-user-defined');
    req.open('GET', url, true);
    if (is_binary)
      req.responseType = 'arraybuffer';

    req.onreadystatechange = function(event) {
      if (req.readyState !== 4)
        return;

      window.setTimeout(function() {
        if (req.status === 200)
          onResult(is_binary ? req.response : req.responseText);
        else
          onResultFail(req.status);
      }, 0);
    };
    req.send(null);
  }

  function onResultFail(err) {
    var overlay = new tr.ui.b.Overlay();
    overlay.textContent = err + ': ' + url + ' could not be loaded';
    overlay.title = 'Failed to fetch data';
    overlay.visible = true;
  }

  function onResult(result) {
    model = new tr.Model();
    var opts = new tr.importer.ImportOptions();
    opts.shiftWorldToZero = false;
    var i = new tr.importer.Import(model, opts);
    var p = i.importTracesWithProgressDialog([result]);
    p.then(onModelLoaded, onImportFail);
  }

  function onModelLoaded() {
    viewer.model = model;
    viewer.viewTitle = "trace";

    if (!model || model.bounds.isEmpty)
      return;
    var sel = window.location.hash.substr(1);
    if (sel === '')
      return;
    var parts = sel.split(':');
    var range = new (tr.b.Range || tr.b.math.Range)();
    range.addValue(parseFloat(parts[0]));
    range.addValue(parseFloat(parts[1]));
    viewer.trackView.viewport.interestRange.set(range);
  }

  function onImportFail(err) {
    var overlay = new tr.ui.b.Overlay();
    overlay.textContent = tr.b.normalizeException(err).message;
    overlay.title = 'Import error';
    overlay.visible = true;
  }

  document.addEventListener('DOMContentLoaded', function() {
    var container = document.createElement('track-view-container');
    container.id = 'track_view_container';

    viewer = document.createElement('tr-ui-timeline-view');
    viewer.track_view_container = container;
    viewer.appendChild(container);

    viewer.id = 'trace-viewer';
    viewer.globalMode = true;
    document.body.appendChild(viewer);

    url = '/jsontrace?{{PARAMS}}';
    load();
  });
}());
</script>
</head>
<body>
</body>
</html>
`
